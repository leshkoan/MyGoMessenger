package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
)

// RegistrationRequest определяет структуру тела запроса для регистрации пользователя.
type RegistrationRequest struct {
	Username string `json:"username"`
}

// RegisterHandler создает новый http.HandlerFunc для регистрации пользователя.
// Требует подключения к базе данных.
func RegisterHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var req RegistrationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.Username == "" {
			http.Error(w, "Username cannot be empty", http.StatusBadRequest)
			return
		}

		// Проверяем, существует ли уже пользователь
		_, err := FindUserByUsername(db, req.Username)
		if err == nil {
			http.Error(w, "Username already exists", http.StatusConflict)
			return
		}

		user, err := CreateUser(db, req.Username)
		if err != nil {
			log.Printf("Failed to create user: %v", err) // Логируем фактическую ошибку
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	}
}
