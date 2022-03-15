package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
)

func main() {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, os.Getenv("GCP_PROJECT"))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	sub := client.Subscription("tutorial.inventory.customers-sub")
	log.Println("Subscribing to tutorial.inventory.customers-sub")
	err = sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
		fmt.Printf("Got message: %q\n", string(m.Data))
		m.Ack() // Acknowledge that we've consumed the message.
	})
	if err != nil {
		log.Fatalf("Failed to receive message: %v", err)
	}
}
