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

	http.HandleFunc("/users/register", RegisterHandler(db))

	log.Println("User service starting on port 8001...")
	if err := http.ListenAndServe(":8001", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}