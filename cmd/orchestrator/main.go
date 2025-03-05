package main

import (
	"log"
	"sprint2-final-task/internal/orchestrator"
)

func main() {
	// Инициализация оркестратора
	orchestrator := orchestrator.New()

	// Старт сервера приложения
	if err := orchestrator.Run(); err != nil {
		log.Fatal(err)
	}
}
