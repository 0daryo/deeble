package pubsub

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/pubsub"
	"github.com/0daryo/deeble/converter/producer"
)

type Poller struct {
	sub *pubsub.Subscription
	producer.Producer
	ResultChan chan *producer.Message
	ErrChan    chan error
}

func NewPoller(ctx context.Context, projectID string, subscriptionName string, p producer.Producer) (*Poller, error) {
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	sub := client.Subscription(subscriptionName)
	log.Println("Subscribing to " + subscriptionName)
	return &Poller{
		sub:        sub,
		Producer:   p,
		ResultChan: make(chan *producer.Message),
		ErrChan:    make(chan error),
	}, nil
}

func (p *Poller) Close() {
	close(p.ResultChan)
	close(p.ErrChan)
}

func (p *Poller) Poll(ctx context.Context) error {
	if err := p.sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
		message, err := p.Producer.Produce(m.Data)
		if err != nil {
			p.ErrChan <- fmt.Errorf("failed to unmarshal message: %w", err)
			return
		}
		p.ResultChan <- message
		m.Ack()
	}); err != nil {
		return fmt.Errorf("failed to subscribe: %w", err)
	}
	return nil
}
