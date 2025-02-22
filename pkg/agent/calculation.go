package agent

import (
	"errors"
	"os"
	"strconv"
)

func Calculate(job string) (string, error) {
	os.Setenv("TIME_ADDITION_MS", "200")
	os.Setenv("TIME_SUBTRACTION_MS", "200")
	os.Setenv("TIME_MULTIPLICATIONS_MS", "200")
	os.Setenv("TIME_DIVISIONS_MS", "200")

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
