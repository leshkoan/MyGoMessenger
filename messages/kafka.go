package main

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
)

const (
	kafkaBroker = "kafka:9092"
	topic       = "messages"
)

// KafkaProducer является оболочкой для kafka.Writer.
type KafkaProducer struct {
	writer *kafka.Writer
}

// NewKafkaProducer создает нового продюсера Kafka.
func NewKafkaProducer() *KafkaProducer {
	w := &kafka.Writer{
		Addr:     kafka.TCP(kafkaBroker),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	return &KafkaProducer{writer: w}
}

// ProduceMessage отправляет сообщение в Kafka.
func (p *KafkaProducer) ProduceMessage(ctx context.Context, msg *Message) error {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx,
		kafka.Message{
			Key:   []byte(msg.ID),
			Value: msgBytes,
		},
	)
}

// Close закрывает писатель Kafka.
func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}