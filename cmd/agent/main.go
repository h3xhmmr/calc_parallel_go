package main

import (
	application "calc_parallel/internal/application/app"
	"os"
)

func main() {
	os.Setenv("COMPUTING_POWER", "10")
	application.Agent_start()
}
