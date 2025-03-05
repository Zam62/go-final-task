package orchestrator

import (
	"encoding/json"
	"log"
	"net/http"
	calculation "sprint2-final-task/pkg"
	"sprint2-final-task/pkg/models"
	"strings"
	"text/template"
)

type CalculationRequest struct {
	Expression string `json:"expression"`
}

type CalculationResponse struct {
	Result float64 `json:"result"`
	Error  string  `json:"error,omitempty"`
}

func (o *Orchestrator) Calculate(w http.ResponseWriter, r *http.Request) {
	wrapped := w.(*responseWriter)

	if r.Method != http.MethodPost {
		wrapped.error = "метод не разрешен"
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req CalculationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		wrapped.error = "некорректное тело запроса"
		http.Error(w, `{"error": "Expression is not valid"}`, http.StatusUnprocessableEntity)
		return
	}

	result, err := calculation.Calc(req.Expression)
	response := CalculationResponse{}

	if err != nil {
		wrapped.error = err.Error()
		response.Error = "Expression is not valid"
		w.WriteHeader(http.StatusUnprocessableEntity)
	} else {
		response.Result = result
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		wrapped.error = "ошибка сериализации"
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}
}

func (o *Orchestrator) ListAllExpressions(w http.ResponseWriter, r *http.Request) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	expressions := make([]models.Expression, 0, len(o.expressions))
	for _, expr := range o.expressions {
		expressions = append(expressions, expr)
	}

	json.NewEncoder(w).Encode(struct {
		Expressions []models.Expression `json:"expressions"`
	}{Expressions: expressions})
}

func (o *Orchestrator) GetExpressionByID(w http.ResponseWriter, r *http.Request) {
}

func (o *Orchestrator) ManageTask(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		o.ManageTask(w, r)
		return
	}

	var result models.ResultRequest
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	o.mu.Lock()
	task, ok := o.tasks[result.TaskID]
	if !ok {
		o.mu.Unlock()
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	task.Status = result.Status
	task.Result = &result.Result
	// task.UpdatedAt = result.Updated
	o.tasks[result.TaskID] = task

	o.mu.Unlock()
	w.WriteHeader(http.StatusOK)
}

// Обработчик главной страницы.
func IndexPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/home" {
		http.NotFound(w, r)
		return
	}
	// Используем функцию template.ParseFiles() для чтения файла шаблона.
	// Если возникла ошибка, мы запишем детальное сообщение ошибки и
	// используя функцию http.Error() мы отправим пользователю
	// ответ: 500 Internal Server Error (Внутренняя ошибка на сервере)
	ts, err := template.ParseFiles("./web/template/home.page.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	// Затем мы используем метод Execute() для записи содержимого
	// шаблона в тело HTTP ответа. Последний параметр в Execute() предоставляет
	// возможность отправки динамических данных в шаблон.
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}

	path := r.URL.Path
	if strings.HasSuffix(path, "js") {
		w.Header().Set("Content-Type", "text/javascript")
	} else {
		w.Header().Set("Content-Type", "text/css")
	}
}
