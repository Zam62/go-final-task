package application

import (
	calculation "calc-service/pkg"
	"encoding/json"
	"log"
	"net/http"
	"text/template"
)

type CalculationRequest struct {
	Expression string `json:"expression"`
}

type CalculationResponse struct {
	Result float64 `json:"result"`
	Error  string  `json:"error,omitempty"`
}

func (a *Application) calculateHandler(w http.ResponseWriter, r *http.Request) {
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

// Обработчик главной страницы.
func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/home" {
		http.NotFound(w, r)
		return
	}

	// Используем функцию template.ParseFiles() для чтения файла шаблона.
	// Если возникла ошибка, мы запишем детальное сообщение ошибки и
	// используя функцию http.Error() мы отправим пользователю
	// ответ: 500 Internal Server Error (Внутренняя ошибка на сервере)
	// ts, err := template.ParseFiles("./ui/html/home.page.tmpl")
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
}
