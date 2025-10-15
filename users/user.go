package main

import (
	"time"

	"github.com/jmoiron/sqlx"
)

// User представляет пользователя в системе.
// Соответствует таблице 'users' в базе данных.
// Теги json используются для управления выводом, когда объект User преобразуется в JSON.
// Теги db используются sqlx для сопоставления полей структуры столбцам базы данных.
type User struct {
	ID        string    `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// CreateUser вставляет нового пользователя в базу данных и возвращает созданного пользователя.
func CreateUser(db *sqlx.DB, username string) (*User, error) {
	user := &User{}
	err := db.QueryRowx("INSERT INTO users (username) VALUES ($1) RETURNING id, username, created_at", username).StructScan(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// FindUserByUsername извлекает пользователя из базы данных по его имени пользователя.
func FindUserByUsername(db *sqlx.DB, username string) (*User, error) {
	user := &User{}
	err := db.Get(user, "SELECT id, username, created_at FROM users WHERE username=$1", username)
	if err != nil {
		return nil, err
	}
	return user, nil
}
