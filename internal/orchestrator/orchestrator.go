package orchestrator

import (
	"encoding/json"
	"go-final-task/internal/config"
	"go-final-task/pkg/database"
	"go-final-task/pkg/models"
	"log"
	"net/http"
	"regexp"
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

var (
	db     *database.SqlDB
	mu     sync.Mutex // Мьютекс для синхронизации доступа к результатам
	ctxKey contextKey = "expression id"
	userID userid     = "user id"
)

func (o *Orchestrator) Run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/calculate", o.Calculate)
	mux.HandleFunc("/api/v1/expressions", o.ListAllExpressions)
	mux.HandleFunc("/api/v1/expressions/{id}", o.GetExpressionByID)
	mux.HandleFunc("/api/v1/internal/task", o.ManageTask)

	registerRoutes(mux)

	handler := corsMiddleware(loggingMiddleware(mux))

	o.server = &http.Server{
		Addr:    o.config.Host + ":" + o.config.Port,
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

func checkCookie(cookie *http.Cookie, err error) bool {
	if err != nil {
		return false
	}

	token := cookie.Value
	return !(len(token) == 0)
}

type ErrorResponse struct {
	Res string `json:"error" example:"Internal server error"`
}

func errorResponse(w http.ResponseWriter, err string, statusCode int) {
	w.WriteHeader(statusCode)
	e := ErrorResponse{Res: err}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(e)
}

func checkId(id string) bool {
	if id == "-1" || id == "" {
		return false
	}

	pattern := "^[0-9]+$"
	r := regexp.MustCompile(pattern)
	return r.MatchString(id)
}
