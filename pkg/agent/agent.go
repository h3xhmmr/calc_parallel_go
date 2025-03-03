package agent

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
	id             string `json:"id"`
	arg1           string `json:"arg1"`
	arg2           string `json:"arg2"`
	operation      string `json:"operation"`
	operation_time int    `json:"operation_time"`
}

type Task_resp struct {
	id     string `json:"id"`
	result string `json:"result"`
}

func Start() {
	ComputingPower, _ := strconv.Atoi(os.Getenv("COMPUTING_POWER"))
	for i := 0; i < ComputingPower; i++ {
		go Worker(i)
	}
	select {}
}

func Worker(id int) {
	for {
		resp, err := http.Get("http://localhost:8080/localhost/internal/task")
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

		var taskresp struct {
			Task `json:"task"`
		}

		err = json.NewDecoder(resp.Body).Decode(&taskresp)
		resp.Body.Close()
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		task := taskresp.Task
		time.Sleep(time.Duration(task.operation_time) * time.Millisecond)
		result, err := Calculate(task.arg1 + task.operation + task.arg2)
		if err != nil {
			log.Printf("Worker %d: error %s: %v", id, task.id, err)
			continue
		}

		res := map[string]interface{}{
			"id":     task.id,
			"result": result,
		}

		response, _ := json.Marshal(res)
		respPost, err := http.Post("http://localhost:8080/localhost/internal/task", "application/json", bytes.NewReader(response))

		if err != nil {
			log.Printf("Worker %d: error posting task %s: %v", id, task.id, err)
			continue
		}

		if respPost.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(respPost.Body)
			log.Printf("Worker %d: error posting task %s: %s", id, task.id, string(body))
		} else {
			log.Printf("Worker %d: completed task %s result %s", id, task.id, result)
		}
		respPost.Body.Close()
	}
}
