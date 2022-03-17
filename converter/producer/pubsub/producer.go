package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"cloud.google.com/go/pubsub"
	"github.com/0daryo/deeble/converter/producer"
)

func Polling(ctx context.Context, projectID string, subscriptionName string, resultChan chan producer.Message, errChan chan error) error {
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return err
	}
	fmt.Printf("projectID %+v\n", projectID)
	fmt.Printf("subscriptionName %+v\n", subscriptionName)
	sub := client.Subscription(subscriptionName)
	log.Println("Subscribing to " + subscriptionName)
	if err := sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
		var msg producer.Message
		if err := json.Unmarshal(m.Data, &msg); err != nil {
			errChan <- fmt.Errorf("failed to unmarshal message: %w", err)
			return
		}
		resultChan <- msg
		m.Ack()
	}); err != nil {
		return fmt.Errorf("failed to subscribe: %w", err)
	}
	return nil
}
