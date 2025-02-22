package agent

import (
	"net/http"
)

type Expression struct {
	id     string `json:"id"`
	status string `sjon:"status"`
	result string `json:"result"`
}

type Expressions struct {
	list []Expression `json:"expressions"`
}

func Agent(w http.ResponseWriter, r *http.Request) {

}
