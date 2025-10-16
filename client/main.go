package main

import (
	"context"
	"fmt"
	"log"

	messagesv1 "mygomessenger.com/mygomessenger/gen/go/messages"
	usersv1 "mygomessenger.com/mygomessenger/gen/go/users"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	userSvcAddr    = "localhost:8081"
	messageSvcAddr = "localhost:8082"
)

func main() {
	// Устанавливаем соединение с сервисом пользователей.
	userConn, err := grpc.NewClient(userSvcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect to user service: %v", err)
	}
	defer userConn.Close()
	userClient := usersv1.NewUserServiceClient(userConn)

	// Устанавливаем соединение с сервисом сообщений.
	msgConn, err := grpc.NewClient(messageSvcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect to message service: %v", err)
	}
	defer msgConn.Close()
	messageClient := messagesv1.NewMessageServiceClient(msgConn)

	// --- Создание пользователей ---
	log.Println("--- Creating Users ---")
	userA, err := createUser(userClient, "user_a")
	if err != nil {
		log.Fatalf("Failed to create user_a: %v", err)
	}

	userB, err := createUser(userClient, "user_b")
	if err != nil {
		log.Fatalf("Failed to create user_b: %v", err)
	}

	// --- Отправка сообщений ---
	log.Println("--- Sending Messages ---")
	// 3 сообщения от A к B
	for i := 1; i <= 3; i++ {
		text := fmt.Sprintf("Hello from user_a, message #%d", i)
		sendMessage(messageClient, userA.Id, userB.Id, text)
	}

	// 3 сообщения от B к A
	for i := 1; i <= 3; i++ {
		text := fmt.Sprintf("Hello from user_b, message #%d", i)
		sendMessage(messageClient, userB.Id, userA.Id, text)
	}

	log.Println("--- Test Complete ---")
}

func createUser(client usersv1.UserServiceClient, username string) (*usersv1.User, error) {
	r, err := client.RegisterUser(context.Background(), &usersv1.RegisterUserRequest{Username: username})
	if err != nil {
		log.Printf("could not create user %s: %v", username, err)
		// Пытаемся получить пользователя, если он уже существует
		return nil, err // В этом тесте мы завершаемся с ошибкой, если не можем создать пользователя
	}
	log.Printf("User Created: ID=%s, Username=%s", r.GetUser().GetId(), r.GetUser().GetUsername())
	return r.GetUser(), nil
}

func sendMessage(client messagesv1.MessageServiceClient, fromID, toID, text string) {
	r, err := client.SendMessage(context.Background(), &messagesv1.SendMessageRequest{FromUserId: fromID, ToUserId: toID, Text: text})
	if err != nil {
		log.Printf("could not send message: %v", err)
		return
	}
	log.Printf("Message Sent: ID=%s, From=%s, To=%s", r.GetMessage().GetId(), r.GetMessage().GetFromUserId(), r.GetMessage().GetToUserId())
}
