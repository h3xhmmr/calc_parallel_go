package main

import (
	application "calc_parallel/internal/application/app"
	"log"
	"os"
)

func main() {
	os.Setenv("TIME_ADDITION_MS", "200")
	os.Setenv("TIME_SUBTRACTION_MS", "200")
	os.Setenv("TIME_MULTIPLICATIONS_MS", "200")
	os.Setenv("TIME_DIVISIONS_MS", "200")
	orc := application.NewOrchestrator()
	err := orc.Run_Orchestrator()
	if err != nil {
		log.Printf("err whith server, %v", err)
	}
}
