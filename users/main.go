package main

import (
	"log"
	"net"

	usersv1 "github.com/leshkoan/MyGoMessenger/gen/go/users"
	"google.golang.org/grpc"
)

const (
	port = ":8001"
)

func main() {
	db, err := NewDBConnection()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	usersv1.RegisterUserServiceServer(s, NewServer(db))

	log.Printf("User gRPC service starting on port %s...", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
