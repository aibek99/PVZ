package kafka

import "context"

// Consumer is
type Consumer interface {
	ReadMessages(ctx context.Context, topic string) error
	Close() error
}
