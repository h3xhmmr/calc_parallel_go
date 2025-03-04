package application

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Task struct {
	ExprID         string   `json:"-"`
	Id             string   `json:"id"`
	Arg1           float64  `json:"arg1"`
	Arg2           float64  `json:"arg2"`
	Operation      string   `json:"operation"`
	Operation_time int      `json:"operation_time"`
	Node           *ASTNode `json:"-"`
}

type Task_resp struct {
	Id     string `json:"id"`
	Result string `json:"result"`
}

func Agent_start() {
	ComputingPower, _ := strconv.Atoi(os.Getenv("COMPUTING_POWER"))
	for i := 0; i < ComputingPower; i++ {
		log.Printf("start worker %d", i)
		go Worker(i)
	}
	select {}
}

func Worker(id int) {
	for {
		resp, err := http.Get("http://localhost:8080/internal/task")
		if err != nil {
			log.Printf("worker %d: error with task: %v", id, err)
			time.Sleep(2 * time.Second)
			continue
		}

		if resp.StatusCode == http.StatusNotFound {
			resp.Body.Close()
			time.Sleep(1 * time.Second)
			continue
		}

		var taskResp struct {
			Task struct {
				Id             string  `json:"id"`
				Arg1           float64 `json:"arg1"`
				Arg2           float64 `json:"arg2"`
				Operation      string  `json:"operation"`
				Operation_time int     `json:"operation_time"`
			} `json:"task"`
		}
		err = json.NewDecoder(resp.Body).Decode(&taskResp)
		resp.Body.Close()
		if err != nil {
			log.Printf("Worker %d: error decoding task: %v", id, err)
			time.Sleep(1 * time.Second)
			continue
		}
		task := taskResp.Task
		log.Printf("Worker %d: received task %s: %f %s %f, time %d ms", id, task.Id, task.Arg1, task.Operation, task.Arg2, task.Operation_time)
		time.Sleep(time.Duration(task.Operation_time) * time.Millisecond)
		result, err := Calc(task.Arg1, task.Operation, task.Arg2)
		if err != nil {
			log.Printf("Worker %d: error %s: %v", id, task.Id, err)
			continue
		}

		var res struct {
			ID     string  `json:"id"`
			Result float64 `json:"result"`
		}
		res.ID = task.Id
		res.Result = result

		response, _ := json.Marshal(res)
		respPost, err := http.Post("http://localhost:8080/internal/task", "application/json", bytes.NewReader(response))

		if err != nil {
			log.Printf("Worker %d: error posting task %s: %v", id, task.Id, err)
			continue
		}

		if respPost.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(respPost.Body)
			log.Printf("Worker %d: error posting %s: %s", id, task.Id, string(body))
		} else {
			log.Printf("Worker %d: completed task %s result %f", id, task.Id, result)
		}
		respPost.Body.Close()
	}
}
