package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	// Импортируем сгенерированные пакеты напрямую из локальных файлов
	pb_messages "github.com/leshkoan/MyGoMessenger/gen/go/messages"
	pb_users "github.com/leshkoan/MyGoMessenger/gen/go/users"
)

const (
	userSvcAddr    = "localhost:8081"
	messageSvcAddr = "localhost:8082"
)

func main() {
	// Устанавливаем соединение с сервисом пользователей
	userConn, err := grpc.NewClient(userSvcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Не удалось подключиться к сервису пользователей: %v", err)
	}
	defer userConn.Close()
	userClient := pb_users.NewUserServiceClient(userConn)

	// Устанавливаем соединение с сервисом сообщений
	msgConn, err := grpc.NewClient(messageSvcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Не удалось подключиться к сервису сообщений: %v", err)
	}
	defer msgConn.Close()
	messageClient := pb_messages.NewMessageServiceClient(msgConn)

	log.Println("=== Тестирование мессенджера с тремя пользователями ===")

	// Создание пользователей
	log.Println("Создание пользователей...")
	userA, err := createUser(userClient, "user_a")
	if err != nil {
		log.Fatalf("Ошибка создания пользователя A: %v", err)
	}
	log.Printf("Пользователь A создан: ID=%s, Username=%s", userA.Id, userA.Username)

	userB, err := createUser(userClient, "user_b")
	if err != nil {
		log.Fatalf("Ошибка создания пользователя B: %v", err)
	}
	log.Printf("Пользователь B создан: ID=%s, Username=%s", userB.Id, userB.Username)

	userC, err := createUser(userClient, "user_c")
	if err != nil {
		log.Fatalf("Ошибка создания пользователя C: %v", err)
	}
	log.Printf("Пользователь C создан: ID=%s, Username=%s", userC.Id, userC.Username)

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

func createUser(client pb_users.UserServiceClient, username string) (*pb_users.User, error) {
	r, err := client.RegisterUser(context.Background(), &pb_users.RegisterUserRequest{Username: username})
	if err != nil {
		log.Printf("Не удалось создать пользователя %s: %v", username, err)
		return nil, err
	}
	return r.GetUser(), nil
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
