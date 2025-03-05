// internal/config/config.go
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

// type ConfigOld struct {
//     Env            string     `yaml:"env" env-default:"local"`
//     StoragePath    string     `yaml:"storage_path" env-required:"true"`
//     GRPC           GRPCConfig `yaml:"grpc"`
//     MigrationsPath string
//     TokenTTL       time.Duration `yaml:"token_ttl" env-default:"1h"`
// }

// type GRPCConfig struct {
//     Port    int           `yaml:"port"`
//     Timeout time.Duration `yaml:"timeout"`
// }

// func loadConfigOld() Config {
// 	os.Setenv("HOST", "localhost")
// 	os.Setenv("PORT", "8080")
// 	return Config{
// 		Host: os.Getenv("HOST"),
// 		Port: os.Getenv("PORT"),
// 	}
// }

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
		log.Println("HOST не установлен, используется порт по умолчанию localhost")
		config.Host = "localhost"
	}

	return config
}
