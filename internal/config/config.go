package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Host string
	Port string
}

// инициализировать объект конфига
func LoadConfig() *Config {
	config := &Config{}

	log.Println("Загрузка файла конфигурации...")
	err := godotenv.Load()
	if err != nil {
		log.Println("Ошибка загрузки файла конфигурации:", err)
	} else {
		log.Println("Файл конфигурации загружен успешно")
	}

	config.Port = os.Getenv("PORT")
	config.Host = os.Getenv("HOST")

	if config.Port == "" {
		log.Println("PORT не установлен, используется порт по умолчанию 8080")
		config.Port = "8080"
	}

	if config.Host == "" {
		log.Println("HOST не установлен, используется host по умолчанию localhost")
		config.Host = "localhost"
	}

	return config
}
