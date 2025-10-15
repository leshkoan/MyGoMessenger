package main

import (
	"time"

	"github.com/jmoiron/sqlx"
)

// Message представляет сообщение в системе.
// Соответствует таблице 'messages' в базе данных.
type Message struct {
	ID         string    `json:"id" db:"id"`
	FromUserID string    `json:"from_user_id" db:"from_user_id"`
	ToUserID   string    `json:"to_user_id" db:"to_user_id"`
	Text       string    `json:"text" db:"text"`
	SentAt     time.Time `json:"sent_at" db:"sent_at"`
}

// CreateMessage вставляет новое сообщение в базу данных и возвращает созданное сообщение.
func CreateMessage(db *sqlx.DB, fromUserID, toUserID, text string) (*Message, error) {
	msg := &Message{}
	err := db.QueryRowx("INSERT INTO messages (from_user_id, to_user_id, text) VALUES ($1, $2, $3) RETURNING id, from_user_id, to_user_id, text, sent_at", fromUserID, toUserID, text).StructScan(msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

// GetMessageHistory извлекает историю сообщений между двумя пользователями.
func GetMessageHistory(db *sqlx.DB, user1ID, user2ID string) ([]Message, error) {
	var messages []Message
	err := db.Select(&messages, "SELECT id, from_user_id, to_user_id, text, sent_at FROM messages WHERE (from_user_id = $1 AND to_user_id = $2) OR (from_user_id = $2 AND to_user_id = $1) ORDER BY sent_at ASC", user1ID, user2ID)
	if err != nil {
		return nil, err
	}
	return messages, nil
}
