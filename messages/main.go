package main

import (
	"log"
	"net"

	messagesv1 "MyGoMessenger/gen/go/messages"

	"google.golang.org/grpc"
)

const (
	port = ":8002"
)

func main() {
	db, err := NewDBConnection()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	producer := NewKafkaProducer()
	defer producer.Close()

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	messagesv1.RegisterMessageServiceServer(s, NewServer(db, producer))

	log.Printf("Message gRPC service starting on port %s...", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
