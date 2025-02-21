package main

import (
	application "calc-service/internal/app"
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	os.Setenv("PORT", "8081")

	errChan := make(chan error, 1)

	go func() {
		app := application.New()
		errChan <- app.Run()
	}()

	select {
	case err := <-errChan:
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	default:
	}
}
