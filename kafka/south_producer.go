package kafka

import (
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

// NorthProducer represents the Kafka producer for the North region
type SouthProducer struct {
	producer sarama.SyncProducer
}

// NewSouthProducer initializes and returns a new SouthProducer
func NewSouthProducer(brokers []string) (*SouthProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to start South Kafka producer: %w", err)
	}

	return &SouthProducer{producer: producer}, nil
}

// SendMessage sends a message to a specific topic in the South region
func (p *SouthProducer) SendMessage(topic, message string) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}
	_, _, err := p.producer.SendMessage(msg)
	if err != nil {
		log.Printf("Failed to send message to topic %s: %v", topic, err)
		return err
	}
	log.Printf("Message sent to topic %s in South region", topic)
	return nil
}
