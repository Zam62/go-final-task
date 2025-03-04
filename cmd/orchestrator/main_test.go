package main

import (
	orchestrator "calc-service/internal/orchestrator"
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	os.Setenv("PORT", "8081")

	errChan := make(chan error, 1)

	go func() {
		orchestrator := orchestrator.New()
		errChan <- orchestrator.Run()
	}()

	select {
	case err := <-errChan:
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	default:
	}
}
