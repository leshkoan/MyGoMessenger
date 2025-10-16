package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Драйвер PostgreSQL
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	// Импортируем сгенерированные пакеты напрямую из локальных файлов
	pb_messages "mygomessenger.com/mygomessenger/gen/go/messages"
	pb_users "mygomessenger.com/mygomessenger/gen/go/users"
)

const (
	userSvcAddr    = "localhost:8081"
	messageSvcAddr = "localhost:8082"
)

// User представляет пользователя в системе (локальная копия структуры)
type User struct {
	ID        string    `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

func main() {
	// Устанавливаем соединение с сервисом пользователей
	userConn, err := grpc.NewClient(userSvcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Не удалось подключиться к сервису пользователей: %v", err)
	}
	defer userConn.Close()

	// Устанавливаем соединение с сервисом сообщений
	msgConn, err := grpc.NewClient(messageSvcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Не удалось подключиться к сервису сообщений: %v", err)
	}
	defer msgConn.Close()
	messageClient := pb_messages.NewMessageServiceClient(msgConn)

	log.Println("=== Тестирование мессенджера с тремя пользователями ===")

	// Создание пользователей (или получение существующих)
	log.Println("Создание пользователей (или получение существующих)...")
	userA, err := createOrGetUser("user_a")
	if err != nil {
		log.Fatalf("Ошибка создания/получения пользователя A: %v", err)
	}
	log.Printf("Пользователь A: ID=%s, Username=%s", userA.Id, userA.Username)

	userB, err := createOrGetUser("user_b")
	if err != nil {
		log.Fatalf("Ошибка создания/получения пользователя B: %v", err)
	}
	log.Printf("Пользователь B: ID=%s, Username=%s", userB.Id, userB.Username)

	userC, err := createOrGetUser("user_c")
	if err != nil {
		log.Fatalf("Ошибка создания/получения пользователя C: %v", err)
	}
	log.Printf("Пользователь C: ID=%s, Username=%s", userC.Id, userC.Username)

	// Отправка сообщений между всеми пользователями
	log.Println("Отправка сообщений между всеми пользователями...")

	// Пользователь A отправляет сообщения B и C
	log.Println("Пользователь A отправляет сообщения...")
	sendMessage(messageClient, userA.Id, userB.Id, "Привет от A пользователю B!")
	time.Sleep(1 * time.Second)
	sendMessage(messageClient, userA.Id, userC.Id, "Привет от A пользователю C!")
	time.Sleep(1 * time.Second)

	// Пользователь B отправляет сообщения A и C
	log.Println("Пользователь B отправляет сообщения...")
	sendMessage(messageClient, userB.Id, userA.Id, "Привет от B пользователю A!")
	time.Sleep(1 * time.Second)
	sendMessage(messageClient, userB.Id, userC.Id, "Привет от B пользователю C!")
	time.Sleep(1 * time.Second)

	// Пользователь C отправляет сообщения A и B
	log.Println("Пользователь C отправляет сообщения...")
	sendMessage(messageClient, userC.Id, userA.Id, "Привет от C пользователю A!")
	time.Sleep(1 * time.Second)
	sendMessage(messageClient, userC.Id, userB.Id, "Привет от C пользователю B!")

	log.Println("=== Тест завершен успешно! Все пользователи обменялись сообщениями ===")
}

// createOrGetUser создает пользователя или возвращает существующего
func createOrGetUser(username string) (*pb_users.User, error) {
	// Подключаемся к базе данных
	db, err := connectToDB()
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к базе данных: %v", err)
	}
	defer db.Close()

	// Пытаемся найти пользователя
	user, err := findUserByUsername(db, username)
	if err == nil {
		// Пользователь найден, возвращаем его
		log.Printf("Пользователь %s уже существует, используем существующего пользователя", username)
		return &pb_users.User{
			Id:       user.ID,
			Username: user.Username,
		}, nil
	}

	// Пользователь не найден, создаем нового через gRPC сервис
	userConn, err := grpc.NewClient(userSvcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к сервису пользователей: %v", err)
	}
	defer userConn.Close()
	userClient := pb_users.NewUserServiceClient(userConn)

	// Создаем пользователя через gRPC
	r, err := userClient.RegisterUser(context.Background(), &pb_users.RegisterUserRequest{Username: username})
	if err != nil {
		return nil, fmt.Errorf("не удалось создать пользователя %s: %v", username, err)
	}

	log.Printf("Пользователь %s создан: ID=%s, Username=%s", username, r.GetUser().Id, r.GetUser().Username)
	return r.GetUser(), nil
}

func connectToDB() (*sqlx.DB, error) {
	dsn := fmt.Sprintf("host=localhost user=%s password=%s dbname=%s sslmode=disable",
		getEnvOrDefault("POSTGRES_USER", "user"),
		getEnvOrDefault("POSTGRES_PASSWORD", "password"),
		getEnvOrDefault("POSTGRES_DB", "mydatabase"))

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func findUserByUsername(db *sqlx.DB, username string) (*User, error) {
	user := &User{}
	err := db.Get(user, "SELECT id, username, created_at FROM users WHERE username=$1", username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func sendMessage(client pb_messages.MessageServiceClient, fromID, toID, text string) {
	r, err := client.SendMessage(context.Background(), &pb_messages.SendMessageRequest{
		FromUserId: fromID,
		ToUserId:   toID,
		Text:       text,
	})
	if err != nil {
		log.Printf("Не удалось отправить сообщение: %v", err)
		return
	}
	log.Printf("Сообщение отправлено: ID=%s, От=%s, Кому=%s, Текст=%s",
		r.GetMessage().GetId(), r.GetMessage().GetFromUserId(), r.GetMessage().GetToUserId(), r.GetMessage().GetText())
}
