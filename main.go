package main

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

func main() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"localhost:9092"},
		Topic:     "dbhistory.inventory",
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
	})
	r.SetOffset(42)
	log.Println("Starting to read messages from topic dbhistory.inventory")
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Printf("error while reading message: %v", err)
			break
		}
		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
	log.Println("Finished reading messages from topic dbhistory.inventory")
}
