package main

import (
	"context"

	messagesv1 "MyGoMessenger/gen/go/messages"

	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// server — это реализация интерфейса MessageServiceServer.
type server struct {
	messagesv1.UnimplementedMessageServiceServer // для прямой совместимости
	db                                           *sqlx.DB
	producer                                     *KafkaProducer
}

// NewServer создает новый экземпляр сервера.
func NewServer(db *sqlx.DB, producer *KafkaProducer) *server {
	return &server{db: db, producer: producer}
}

// SendMessage отправляет новое сообщение.
func (s *server) SendMessage(ctx context.Context, req *messagesv1.SendMessageRequest) (*messagesv1.SendMessageResponse, error) {
	msg, err := CreateMessage(s.db, req.GetFromUserId(), req.GetToUserId(), req.GetText())
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to create message")
	}

	if err := s.producer.ProduceMessage(ctx, msg); err != nil {
		// Логируем ошибку, но не прерываем RPC, так как сообщение уже в БД
		// log.Printf("Не удалось отправить сообщение в Kafka: %v", err)
	}

	return &messagesv1.SendMessageResponse{
		Message: &messagesv1.Message{
			Id:         msg.ID,
			FromUserId: msg.FromUserID,
			ToUserId:   msg.ToUserID,
			Text:       msg.Text,
			SentAt:     timestamppb.New(msg.SentAt),
		},
	}, nil
}

// GetMessageHistory получает историю сообщений между двумя пользователями.
func (s *server) GetMessageHistory(ctx context.Context, req *messagesv1.GetMessageHistoryRequest) (*messagesv1.GetMessageHistoryResponse, error) {
	messages, err := GetMessageHistory(s.db, req.GetUser1Id(), req.GetUser2Id())
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to get message history")
	}

	var grpcMessages []*messagesv1.Message
	for _, msg := range messages {
		grpcMessages = append(grpcMessages, &messagesv1.Message{
			Id:         msg.ID,
			FromUserId: msg.FromUserID,
			ToUserId:   msg.ToUserID,
			Text:       msg.Text,
			SentAt:     timestamppb.New(msg.SentAt),
		})
	}

	return &messagesv1.GetMessageHistoryResponse{Messages: grpcMessages}, nil
}
