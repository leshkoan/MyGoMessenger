package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	consumer := NewKafkaConsumer()
	defer consumer.Close()

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		consumer.ConsumeMessages(ctx)
	}()

	log.Println("Notification service started. Waiting for messages...")

	// Wait for termination signal
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	<-sigchan

	log.Println("Shutting down notification service...")
	cancel()
}