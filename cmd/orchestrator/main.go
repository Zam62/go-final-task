package main

import (
	"go-final-task/internal/orchestrator"
	"log"
)

func main() {
	// Инициализация оркестратора
	orchestrator := orchestrator.New()

	// Старт сервера приложения
	if err := orchestrator.Run(); err != nil {
		log.Fatal(err)
	}
}
