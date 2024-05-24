package kafka

import "time"

// Message is a struct to represent the message to be sent to Kafka.
type Message struct {
	Method    string    `json:"method"`
	Request   string    `json:"request"`
	Timestamp time.Time `json:"timestamp"`
}
