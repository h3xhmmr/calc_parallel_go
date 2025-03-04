package main

import (
	"calc_parallel/internal/application"
	"os"
)

func main() {
	os.Setenv("COMPUTING_POWER", "10")
	application.Agent_start()
}
