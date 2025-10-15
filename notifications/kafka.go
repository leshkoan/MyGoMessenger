package main

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

const (
	kafkaBroker = "kafka:9092"
	topic       = "messages"
	groupID     = "notification-group"
)

// KafkaConsumer является оболочкой для kafka.Reader.
type KafkaConsumer struct {
	reader *kafka.Reader
}

// NewKafkaConsumer создает нового потребителя Kafka.
func NewKafkaConsumer() *KafkaConsumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{kafkaBroker},
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	return &KafkaConsumer{reader: r}
}

// ConsumeMessages читает сообщения из Kafka и логирует их.
func (c *KafkaConsumer) ConsumeMessages(ctx context.Context) {
	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Failed to read message: %v", err)
			break
		}
		fmt.Printf("Получено сообщение со смещением %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
	}
}

// Close закрывает читателя Kafka.
func (c *KafkaConsumer) Close() error {
	return c.reader.Close()
}