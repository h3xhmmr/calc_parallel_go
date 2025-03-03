package agent

import (
	"errors"
	"strconv"
)

func Calculate(job string) (string, error) {

	if string(job[2]) == "0" && string(job[1]) == "/" {
		return "", errors.New("Error: dividing by zero")
	}

	a, err := strconv.Atoi(string(job[0]))
	if err != nil {
		return "", err
	}
	b, err := strconv.Atoi(string(job[2]))
	if err != nil {
		return "", err
	}

	res := ""
	switch string(job[1]) {
	case "*":
		res = strconv.Itoa(a * b)
	case "/":
		res = strconv.Itoa(a / b)
	case "+":
		res = strconv.Itoa(a + b)
	case "-":
		res = strconv.Itoa(a - b)
	}
	return res, nil
}
