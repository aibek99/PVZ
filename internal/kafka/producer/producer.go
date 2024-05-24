package producer

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"

	"Homework-1/internal/config"
	"Homework-1/internal/kafka"
	kafkaModel "Homework-1/internal/model/kafka"
)

var _ kafka.Producer = (*Producer)(nil)

// Producer is
type Producer struct {
	syncProducer sarama.SyncProducer
	topic        string
}

// NewProducer is
func NewProducer(cfg config.Kafka) (*Producer, error) {
	saramaConfig := sarama.NewConfig()
	// wants a message from broker leader replica acks = 1
	saramaConfig.Producer.RequiredAcks = sarama.WaitForLocal
	// retying send message max 2 times
	saramaConfig.Producer.Retry.Max = 2
	// waits message delivered
	saramaConfig.Producer.Return.Successes = true
	// Create the Kafka Producer
	syncProducer, err := sarama.NewSyncProducer(cfg.Brokers, saramaConfig)
	if err != nil {
		return nil, err
	}

	return &Producer{
		syncProducer: syncProducer,
		topic:        cfg.Topic,
	}, nil
}

// SendMessage is
func (p *Producer) SendMessage(topic string, message kafkaModel.Message) error {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(messageBytes),
	}

	partition, offset, err := p.syncProducer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("p.syncProducer.SendMessage: %w", err)
	}

	log.Println("[kafka][producer][send_message] Partition: ", partition, " Offset: ", offset, " AnswerID:", message.Method)

	return nil
}

// Close is
func (p *Producer) Close() error {
	err := p.syncProducer.Close()
	if err != nil {
		return fmt.Errorf("p.syncProducer.Close: %w", err)
	}

	return nil
}

// Topic is
func (p *Producer) Topic() string {
	return p.topic
}
