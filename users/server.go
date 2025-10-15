package main

import (
	"context"

	usersv1 "github.com/leshkoan/MyGoMessenger/gen/go/users"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// server — это реализация интерфейса UserServiceServer.
type server struct {
	usersv1.UnimplementedUserServiceServer // для прямой совместимости
	db *sqlx.DB
}

// NewServer создает новый экземпляр сервера.
func NewServer(db *sqlx.DB) *server {
	return &server{db: db}
}

// RegisterUser регистрирует нового пользователя.
func (s *server) RegisterUser(ctx context.Context, req *usersv1.RegisterUserRequest) (*usersv1.RegisterUserResponse, error) {
	if req.GetUsername() == "" {
		return nil, status.Error(codes.InvalidArgument, "Username cannot be empty")
	}

	// Проверяем, существует ли уже пользователь
	_, err := FindUserByUsername(s.db, req.GetUsername())
	if err == nil {
		return nil, status.Error(codes.AlreadyExists, "Username already exists")
	}

	// Создаем пользователя
	user, err := CreateUser(s.db, req.GetUsername())
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to create user")
	}

	return &usersv1.RegisterUserResponse{
		User: &usersv1.User{
			Id:        user.ID,
			Username:  user.Username,
			CreatedAt: timestamppb.New(user.CreatedAt),
		},
	}, nil
}
