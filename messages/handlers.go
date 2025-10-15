package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
)

// SendMessageRequest определяет структуру тела запроса для отправки сообщения.
type SendMessageRequest struct {
	FromUserID string `json:"from_user_id"`
	ToUserID   string `json:"to_user_id"`
	Text       string `json:"text"`
}

// SendMessageHandler создает новый http.HandlerFunc для отправки сообщений.
func SendMessageHandler(db *sqlx.DB, producer *KafkaProducer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var req SendMessageRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		msg, err := CreateMessage(db, req.FromUserID, req.ToUserID, req.Text)
		if err != nil {
			log.Printf("Failed to create message: %v", err)
			http.Error(w, "Failed to create message", http.StatusInternalServerError)
			return
		}

		if err := producer.ProduceMessage(context.Background(), msg); err != nil {
			log.Printf("Failed to produce message to Kafka: %v", err)
			// Мы не возвращаем ошибку пользователю, так как сообщение уже сохранено.
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(msg)
	}
}

// GetHistoryHandler создает новый http.HandlerFunc для получения истории сообщений.
func GetHistoryHandler(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		user1 := r.URL.Query().Get("user1")
		user2 := r.URL.Query().Get("user2")

		if user1 == "" || user2 == "" {
			http.Error(w, "user1 and user2 query parameters are required", http.StatusBadRequest)
			return
		}

		messages, err := GetMessageHistory(db, user1, user2)
		if err != nil {
			log.Printf("Failed to get message history: %v", err)
			http.Error(w, "Failed to get message history", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(messages)
	}
}