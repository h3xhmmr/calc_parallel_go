package main

import (
	"calc_parallel/internal/orchestrator"
	"calc_parallel/pkg/agent"
	"os"
)

func main() {
	os.Setenv("TIME_ADDITION_MS", "200")
	os.Setenv("TIME_SUBTRACTION_MS", "200")
	os.Setenv("TIME_MULTIPLICATIONS_MS", "200")
	os.Setenv("TIME_DIVISIONS_MS", "200")
	os.Setenv("COMPUTING_POWER", "10")
	orc := orchestrator.NewOrchestrator()
	orc.Run_Orchestrator()
	agent.Start()
}
