// Package kafka ...
//
//go:generate minimock -g -i Producer -o ./mock/producer_mock.go -n ProducerMock
package kafka

import kafkaModel "Homework-1/internal/model/kafka"

// Producer is
type Producer interface {
	SendMessage(topic string, message kafkaModel.Message) error
	Close() error
	Topic() string
}
