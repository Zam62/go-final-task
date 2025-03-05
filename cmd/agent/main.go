package main

import (
	"log"
	"os"
	"sprint2-final-task/internal/agent"
	"strconv"
	"sync"
)

func main() {
	orchestratorURL := os.Getenv("ORCHESTRATOR_URL")
	if orchestratorURL == "" {
		orchestratorURL = "http://localhost:8080"
	}

	computingPower := 4
	if cp := os.Getenv("COMPUTING_POWER"); cp != "" {
		if val, err := strconv.Atoi(cp); err == nil {
			computingPower = val
		}
	}

	agent := agent.NewAgent(orchestratorURL)
	var wg sync.WaitGroup

	log.Printf("Starting agent with %d workers, connecting to %s", computingPower, orchestratorURL)
	for i := 0; i < computingPower; i++ {
		wg.Add(1)
		go agent.Worker(&wg)
	}

	wg.Wait()
}
