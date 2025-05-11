package application

import (
	"errors"
)

func Calc(a float64, op string, b float64) (float64, error) {
	var res float64
	switch op {
	case "*":
		res = a * b
	case "/":
		if b == 0 {
			return 0, errors.New("Error: division by zero")
		}
		res = a / b
	case "+":
		res = a + b
	case "-":
		res = a - b
	}
	return res, nil
}
