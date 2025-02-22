package agent

import (
	"net/http"
	"time"
)

type Task struct {
	id             string `json:"id"`
	arg1           string `json:"arg1"`
	arg2           string `json:"arg2"`
	operation      string `json:"operation"`
	operation_time int    `json:"operation_time"`
}

func worker(jobs <-chan string, results chan<- string, t int) {
	for j := range jobs {
		for i := 0; i <= t; t++ {
			time.Sleep(time.Millisecond)
		}
		res, _ := Calculate(j)
		results <- res
	}
}

func Agent(w http.ResponseWriter, r *http.Request) {
}
