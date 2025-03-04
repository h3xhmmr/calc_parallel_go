package application

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Time_env struct {
	TimeAddition        int
	TimeSubtraction     int
	TimeMultiplications int
	TimeDivisions       int
}

func Time() *Time_env {
	t_add, _ := strconv.Atoi(os.Getenv("TIME_ADDITION_MS"))
	t_sub, _ := strconv.Atoi(os.Getenv("TIME_SUBTRACTION_MS"))
	t_mul, _ := strconv.Atoi(os.Getenv("TIME_MULTIPLICATIONS_MS"))
	t_div, _ := strconv.Atoi(os.Getenv("TIME_DIVISIONS_MS"))

	return &Time_env{
		TimeAddition:        t_add,
		TimeSubtraction:     t_sub,
		TimeMultiplications: t_mul,
		TimeDivisions:       t_div,
	}
}

type Orchestrator struct {
	op_time     *Time_env
	exprStore   map[string]*Expression
	taskStore   map[string]Task
	taskQueue   []Task
	mu          sync.Mutex
	exprCounter int64
	taskCounter int64
}

func NewOrchestrator() *Orchestrator {
	return &Orchestrator{
		op_time:   Time(),
		exprStore: make(map[string]*Expression),
		taskStore: make(map[string]Task),
		taskQueue: make([]Task, 0),
	}
}

type Expression struct {
	ID     string   `json:"id"`
	Expr   string   `json:"expression"`
	Status string   `json:"status"`
	Result *float64 `json:"result,omitempty"`
	AST    *ASTNode `json:"-"`
}

func (o *Orchestrator) CalculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Wrong Method"}`, http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Expression string `json:"expression"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Expression == "" {
		http.Error(w, `{"error":"Invalid Body"}`, http.StatusUnprocessableEntity)
		return
	}
	ast, err := ParseAST(req.Expression)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusUnprocessableEntity)
		return
	}
	o.mu.Lock()
	o.exprCounter++
	exprID := fmt.Sprintf("%d", o.exprCounter)
	expr := &Expression{
		ID:     exprID,
		Expr:   req.Expression,
		Status: "in process",
		AST:    ast,
	}
	o.exprStore[exprID] = expr
	o.scheduleTasks(expr)
	o.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": exprID})
}

func (o *Orchestrator) Handler_expressions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"Wrong Method"}`, http.StatusMethodNotAllowed)
		return
	}
	o.mu.Lock()
	defer o.mu.Unlock()
	exprs := make([]*Expression, 0, len(o.exprStore))
	for _, expr := range o.exprStore {
		if expr.AST != nil && expr.AST.IsLeaf {
			expr.Status = "completed"
			expr.Result = &expr.AST.Value
		}
		exprs = append(exprs, expr)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"expressions": exprs})
}

func (o *Orchestrator) Handler_Id(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"Wrong Method"}`, http.StatusMethodNotAllowed)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/expressions/")
	o.mu.Lock()
	expr, ok := o.exprStore[id]
	o.mu.Unlock()
	if !ok {
		http.Error(w, `{"error":"Expression not found"}`, http.StatusNotFound)
		return
	}
	if expr.AST != nil && expr.AST.IsLeaf {
		expr.Status = "completed"
		expr.Result = &expr.AST.Value
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"expression": expr})
}

func (o *Orchestrator) Handler_Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"Wrong Method"}`, http.StatusMethodNotAllowed)
		return
	}
	o.mu.Lock()
	defer o.mu.Unlock()
	if len(o.taskQueue) == 0 {
		http.Error(w, `{"error":"No task available"}`, http.StatusNotFound)
		return
	}
	task := o.taskQueue[0]
	o.taskQueue = o.taskQueue[1:]
	if expr, exists := o.exprStore[task.ExprID]; exists {
		expr.Status = "in_progress"
	}
	var Task_req struct {
		Task Task `json:"task"`
	}
	Task_req.Task = task
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Task_req)
}

func (o *Orchestrator) Handler_post(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Wrong Method"}`, http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		ID     string  `json:"id"`
		Result float64 `json:"result"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.ID == "" {
		http.Error(w, `{"error":"Invalid Body"}`, http.StatusUnprocessableEntity)
		return
	}
	o.mu.Lock()
	task, ok := o.taskStore[req.ID]
	if !ok {
		o.mu.Unlock()
		http.Error(w, `{"error":"Task not found"}`, http.StatusNotFound)
		return
	}
	task.Node.IsLeaf = true
	task.Node.Value = req.Result
	delete(o.taskStore, req.ID)
	if expr, exists := o.exprStore[task.ExprID]; exists {
		o.scheduleTasks(expr)
		if expr.AST.IsLeaf {
			expr.Status = "completed"
			expr.Result = &expr.AST.Value
		}
	}
	o.mu.Unlock()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"result accepted"}`))
}

func (o *Orchestrator) scheduleTasks(expr *Expression) {
	var traverse func(node *ASTNode)
	traverse = func(node *ASTNode) {
		if node == nil || node.IsLeaf {
			return
		}
		traverse(node.Left)
		traverse(node.Right)
		if node.Left != nil && node.Right != nil && node.Left.IsLeaf && node.Right.IsLeaf {
			if !node.TaskScheduled {
				o.taskCounter++
				taskID := fmt.Sprintf("%d", o.taskCounter)
				var opTime int
				switch node.Operator {
				case "+":
					opTime = o.op_time.TimeAddition
				case "-":
					opTime = o.op_time.TimeSubtraction
				case "*":
					opTime = o.op_time.TimeMultiplications
				case "/":
					opTime = o.op_time.TimeDivisions
				}
				task := Task{
					Id:             taskID,
					ExprID:         expr.ID,
					Arg1:           node.Left.Value,
					Arg2:           node.Right.Value,
					Operation:      node.Operator,
					Operation_time: opTime,
					Node:           node,
				}
				node.TaskScheduled = true
				o.taskStore[taskID] = task
				o.taskQueue = append(o.taskQueue, task)
			}
		}
	}
	traverse(expr.AST)
}

func (o *Orchestrator) Run_Orchestrator() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/calculate", o.CalculateHandler)
	mux.HandleFunc("/api/v1/expressions", o.Handler_expressions)
	mux.HandleFunc("/api/v1/expressions/", o.Handler_Id)
	mux.HandleFunc("/internal/task", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			o.Handler_Get(w, r)
		} else if r.Method == http.MethodPost {
			o.Handler_post(w, r)
		} else {
			http.Error(w, `{"error":"Wrong Method"}`, http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"error":"Not Found"}`, http.StatusNotFound)
	})
	go func() {
		for {
			time.Sleep(8 * time.Second)
			o.mu.Lock()
			if len(o.taskQueue) > 0 {
				log.Printf("tasks in queue: %d", len(o.taskQueue))
			}
			o.mu.Unlock()
		}
	}()
	return http.ListenAndServe(":8080", mux)
}
