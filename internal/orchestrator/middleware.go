package orchestrator

import (
	"context"
	"encoding/json"
	"go-final-task/pkg/crypto/jwt"
	"go-final-task/pkg/models"
	"log"
	"net/http"
	"strings"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status     int
	error      string
	expression string
}

type (
	ExpressionReq struct {
		Expression string `json:"expression"`
	}

	RespID struct {
		Id int `json:"id"`
	}

	// ErrorResponse struct {
	// 	Res string `json:"error" example:"Internal server error"`
	// }

	Expression struct {
		exp string
		id  int
	}

	contextKey string
	userid     string
)

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &responseWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		log.Printf("Входящий %s запрос на %s", r.Method, r.URL.Path)

		next.ServeHTTP(wrapped, r)

		logMessage := "Завершен %s %s - статус: %d, длительность: %v"
		logArgs := []interface{}{
			r.Method,
			r.URL.Path,
			wrapped.status,
			time.Since(start),
		}

		if wrapped.status != http.StatusOK {
			logMessage += ", ошибка: %s"
			if wrapped.expression != "" {
				logMessage += ", выражение: %s"
				logArgs = append(logArgs, wrapped.error, wrapped.expression)
			} else {
				logArgs = append(logArgs, wrapped.error)
			}
		}

		log.Printf(logMessage, logArgs...)
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func logsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Method: %s, URL: %s", r.Method, r.URL)
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		log.Printf("Method: %s, completion time: %v", r.Method, duration)
	})
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token string
		cookie, err := r.Cookie("jwt")
		if checkCookie(cookie, err) {
			token = cookie.Value
			log.Print("token was taken from cookie")
		} else {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				errorResponse(w, "authorization is required", http.StatusUnauthorized)
				log.Printf("Code: %v, user unauthorized", http.StatusUnauthorized)
				return
			}

			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				errorResponse(w, "invalid token format", http.StatusUnauthorized)
				log.Printf("Code: %v, invalid token format", http.StatusUnauthorized)
				return
			}
			token = tokenParts[1]
			log.Print("token was taken from header")
		}

		claims, id := jwt.Verify(token)
		if !claims {
			errorResponse(w, "invalid token", http.StatusUnauthorized)
			log.Printf("Code: %v, invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userID, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func databaseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			errorResponse(w, "invalid request method", http.StatusMethodNotAllowed)
			log.Printf("Code: %v, Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var body ExpressionReq
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			errorResponse(w, "internal server error", http.StatusInternalServerError)
			log.Printf("Code: %v, error with decoding request body", http.StatusInternalServerError)
			return
		}

		// добавляем выражение в базу данных и получаем id
		log.Printf("Adding expression to database")
		e := &models.Expression{
			Expression: body.Expression,
			Status:     "in process",
			Result:     0.0,
		}
		expID, err := db.InsertExpression(r.Context(), e, r.Context().Value(userID).(int))
		if err != nil {
			errorResponse(w, "internal server error", http.StatusInternalServerError)
			log.Printf("Code: %v, error with database", http.StatusInternalServerError)
			return
		}
		respID := RespID{Id: int(expID)}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(respID)

		// запускаем вычисления в фоновом режиме
		go func() {
			expr := &Expression{
				exp: body.Expression,
				id:  int(expID),
			}
			ctx := context.WithValue(r.Context(), ctxKey, expr)
			next.ServeHTTP(w, r.WithContext(ctx))
		}()
	})
}
