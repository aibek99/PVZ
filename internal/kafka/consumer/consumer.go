package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/IBM/sarama"

	"Homework-1/internal/kafka"
	kafkaModel "Homework-1/internal/model/kafka"
)

var _ kafka.Consumer = (*Consumer)(nil)

// Consumer is
type Consumer struct {
	consumer sarama.Consumer
	messages io.Writer
}

// NewConsumer is
func NewConsumer(brokers []string, message io.Writer) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("sarama.NewConsumer: %w", err)
	}

	return &Consumer{
		consumer: consumer,
		messages: message,
	}, nil
}

// ReadMessages is
func (c *Consumer) ReadMessages(ctx context.Context, topic string) error {
	partitions, err := c.consumer.Partitions(topic)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	for _, partition := range partitions {
		wg.Add(1)

		go func(partition int32) {
			defer wg.Done()

			partitionConsumer, err := c.consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
			if err != nil {
				return
			}

			for {
				select {
				case err = <-partitionConsumer.Errors():
					{
						log.Printf("[kafka][consumer][ReadMessages] error: %v", err)
					}
				case msg := <-partitionConsumer.Messages():
					{
						var message kafkaModel.Message
						if err = json.Unmarshal(msg.Value, &message); err != nil {
							log.Printf("[kafka][consumer][ReadMessages] Failed to unmarshal Kafka message: %v", err)
							continue
						}
						messageString := fmt.Sprintf(
							"Received Kafka Message: Method: %s, Request: %s, Timestamp: %v\n",
							message.Method,
							message.Request,
							message.Timestamp,
						)

						_, err = c.messages.Write([]byte(messageString))
						if err != nil {
							log.Printf("[kafka][consumer][ReadMessages] Failed to Write Kafka message: %v", err)
							continue
						}
					}
				case <-ctx.Done():
					{
						log.Println("[kafka][consumer][ReadMessages] Shutting down consumer for partition", partition)
						return
					}
				}
			}

		}(partition)
	}

	wg.Wait()

	return fmt.Errorf("consumer canceled")
}

// Close is
func (c *Consumer) Close() error {
	err := c.consumer.Close()
	if err != nil {
		return fmt.Errorf("c.consumer.Close: %w", err)
	}

	return nil
}
