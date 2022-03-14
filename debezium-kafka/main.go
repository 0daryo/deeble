package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	// r := kafka.NewReader(kafka.ReaderConfig{
	// 	Brokers:   []string{"localhost:9092"},
	// 	Topic:     "dbserver1.inventory.customers",
	// 	Partition: 0,
	// 	MinBytes:  10e3, // 10KB
	// 	MaxBytes:  10e6, // 10MB
	// })
	// r.SetOffset(42)
	// log.Println("Starting to read messages from topic dbhistory.inventory")
	// for {
	// 	m, err := r.ReadMessage(context.Background())
	// 	if err != nil {
	// 		log.Printf("error while reading message: %v", err)
	// 		break
	// 	}
	// 	fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
	// }

	// if err := r.Close(); err != nil {
	// 	log.Fatal("failed to close reader:", err)
	// }
	// log.Println("Finished reading messages from topic dbhistory.inventory")
	// to consume messages
	topic := "dbserver1.inventory.customers"
	partition := 0
	log.Println("Starting to read messages from topic", topic)
	conn, err := kafka.DialLeader(context.Background(), "tcp", "172.28.0.5:9092", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	batch := conn.ReadBatch(10e3, 1e6) // fetch 10KB min, 1MB max

	b := make([]byte, 10e3) // 10KB max per message
	for {
		n, err := batch.Read(b)
		if err != nil {
			break
		}
		fmt.Println(string(b[:n]))
	}

	if err := batch.Close(); err != nil {
		log.Fatal("failed to close batch:", err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close connection:", err)
	}
}
