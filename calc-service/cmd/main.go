package main

import (
	application "calc-service/internal/app"
	"log"
)

func main() {
	// TODO: инициализировать объект конфига

	// TODO: инициализировать логгер

	// TODO: инициализировать приложение (app)
	app := application.New()
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}

	// TODO: запустить сервер приложения
}
