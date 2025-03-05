package orchestrator

import (
	"log"
	"net/http"
	"sprint2-final-task/internal/config"
	"sprint2-final-task/pkg/models"
	"sync"
)

type Orchestrator struct {
	taskQueue   chan models.Task
	expressions map[string]models.Expression
	tasks       map[string]models.Task
	mu          sync.RWMutex
	config      *config.Config
	server      *http.Server
}

func New() *Orchestrator {
	return &Orchestrator{
		taskQueue:   make(chan models.Task, 100),
		expressions: make(map[string]models.Expression),
		tasks:       make(map[string]models.Task),
		config:      config.LoadConfig(),
	}
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
