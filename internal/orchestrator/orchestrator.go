package orchestrator

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
)

type Expression struct {
	id     string `json:"id"`
	status string `json:"status"`
	result string `json:"result"`
}

type Task struct {
	id             string `json:"id"`
	arg1           string `json:"arg1"`
	arg2           string `json:"arg2"`
	operation      string `json:"operation"`
	operation_time int    `json:"operation_time"`
}

type Expressions_resp struct {
	expressions []Expression `json:"expressions"`
}

func Run_Orchestrator() {
	os.Setenv("TIME_ADDITION_MS", "200")
	os.Setenv("TIME_SUBTRACTION_MS", "200")
	os.Setenv("TIME_MULTIPLICATIONS_MS", "200")
	os.Setenv("TIME_DIVISIONS_MS", "200")
}

func Post_From_Agent(w http.ResponseWriter, r *http.Request) {

}

func Get_From_Agent(w http.ResponseWriter, r *http.Request) {

}

func parce(expression string) ([]string, string, error) {
	expression = strings.ReplaceAll(expression, " ", "")
	exp_nums := expression
	exp_signs := expression
	for _, x := range expression {
		if !strings.Contains("1234567890+-/*()", string(x)) {
			return nil, "", ErrInvalidExpression
		}
	}
	for _, i := range "+-/*()" {
		exp_nums = strings.ReplaceAll(exp_nums, string(i), " ")
	}
	exp_nums = strings.TrimRight(exp_nums, " ")
	exp_nums_slice := strings.Split(exp_nums, " ")
	for _, i := range "1234567890" {
		exp_signs = strings.ReplaceAll(exp_signs, string(i), "")
	}
	for j := 0; j < len(expression)-1; j++ {
		if strings.Contains("+-*/(", string(expression[j])) && strings.Contains("+-*/)", string(expression[j+1])) {
			return nil, "", ErrInvalidExpression
		}
	}
	if strings.Count(expression, "(") != strings.Count(expression, ")") {
		return nil, "", ErrInvalidExpression
	}
	return exp_nums_slice, exp_signs, nil
}

func get_expressions(w http.ResponseWriter, r *http.Request) {
	res, err_m := json.Marshal(exp_resp)
	if err_m != nil {
		http.Error(w, `{"error":"error while marshaling"}`, http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(res)
	if err != nil {
		http.Error(w, `{"error":"error while writing response"}`, http.StatusBadRequest)
	}
}

func get_expression(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	res, err_m := json.Marshal(exp_map[id])
	if err_m != nil {
		http.Error(w, `{"error":"error while marshaling"}`, http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(res)
	if err != nil {
		http.Error(w, `{"error":"error while writing response"}`, http.StatusBadRequest)
	}
}

func add_expression(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"error":"invalid Body"}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var expression Expression
	err = json.Unmarshal(body, &expression)
	if err != nil {
		http.Error(w, `{"error":"invalid Body"}`, http.StatusBadRequest)
		return
	}

	id := uuid.NewString()
	expression.id = id
	expression.status = "in process"
	exp_resp.expressions = append(exp_resp.expressions, expression)
	exp_map[id] = expression

	_, err = w.Write([]byte(id))
	if err != nil {
		http.Error(w, `{"error writing response"}`, http.StatusBadRequest)
	}
}
