package orchestrator

import (
	"log"
	"net/http"
	"os"
	"sprint2-final-task/pkg/models"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
}

type Orchestrator struct {
	taskQueue   chan models.Task
	expressions map[string]models.Expression
	tasks       map[string]models.Task
	mu          sync.RWMutex
	config      *Config
	server      *http.Server
}

func New() *Orchestrator {
	return &Orchestrator{
		taskQueue:   make(chan models.Task, 100),
		expressions: make(map[string]models.Expression),
		tasks:       make(map[string]models.Task),
		config:      loadConfig(),
	}
}

func loadConfig() *Config {
	config := &Config{}

	log.Println("Загрузка файла конфигурации...")
	err := godotenv.Load()
	if err != nil {
		log.Println("Ошибка загрузки файла конфигурации:", err)
	} else {
		log.Println("Файл конфигурации загружен успешно")
	}

	config.Port = os.Getenv("PORT")

	if config.Port == "" {
		log.Println("PORT не установлен, используется порт по умолчанию 8080")
		config.Port = "8080"
	}

	return config
}

func (o *Orchestrator) Run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/calculate", o.Calculate)
	mux.HandleFunc("/api/v1/expressions", o.ListAllExpressions)
	mux.HandleFunc("/api/v1/expressions/{id}", o.GetExpressionByID)
	mux.HandleFunc("/api/v1/internal/task", o.ManageTask)

	registerRoutes(mux)

	handler := corsMiddleware(loggingMiddleware(mux))

	o.server = &http.Server{
		Addr:    ":" + o.config.Port,
		Handler: handler,
	}

	log.Printf("Запуск сервера на порту %s\n", o.config.Port)
	return o.server.ListenAndServe()
}

func registerRoutes(mux *http.ServeMux) {
	fs := http.FileServer(http.Dir("./web/static"))

	mux.HandleFunc("/", IndexPageHandler)
	mux.Handle("/css/", http.StripPrefix("/css", fs))
	mux.Handle("/js/", http.StripPrefix("/js", fs))

}
