package main

import (
	"log"
	"net/http"
)

func main() {
	db, err := NewDBConnection()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	producer := NewKafkaProducer()
	defer producer.Close()

	http.HandleFunc("/messages/send", SendMessageHandler(db, producer))
	http.HandleFunc("/messages/history", GetHistoryHandler(db))

	log.Println("Message service starting on port 8002...")
	if err := http.ListenAndServe(":8002", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}