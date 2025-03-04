package main

import (
	orchestrator "calc-service/internal/orchestrator"
	"log"
)

func main() {
	// TODO: инициализировать объект конфига

	// TODO: инициализировать логгер

	// Инициализация оркестратора
	orchestrator := orchestrator.New()

	// Старт сервера приложения
	if err := orchestrator.Run(); err != nil {
		log.Fatal(err)
	}
}
