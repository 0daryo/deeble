package spanner

import (
	"context"
	"fmt"

	"cloud.google.com/go/spanner"
	"github.com/0daryo/deeble/converter/producer"
)

var (
	ErrUnkownEventType = fmt.Errorf("unknown event type")
)

type Consumer struct {
	cli spanner.Client
}

func (c *Consumer) Consume(ctx context.Context, message producer.Message) error {
	switch message.EventType {
	case producer.Insert:
		return c.insert(ctx, message)
	}
	return ErrUnkownEventType
}

func (c *Consumer) insert(ctx context.Context, message producer.Message) error {
	_, err := c.cli.Apply(ctx, []*spanner.Mutation{
		spanner.Insert(message.TableName, message.TargetKeys(), message.TargetValues()),
	})
	return err
}
