package kafka

import (
	"log"

	"github.com/IBM/sarama"
)

// Producer represents a Kafka producer
type Producer struct {
	producer sarama.SyncProducer
	topic    string
}

// NewProducer creates a new Kafka producer
func NewProducer(brokers []string, topic string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &Producer{
		producer: producer,
		topic:    topic,
	}, nil
}

// Publish publishes a message to the Kafka topic
func (p *Producer) Publish(message []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.ByteEncoder(message),
	}

	_, _, err := p.producer.SendMessage(msg)
	if err != nil {
		log.Printf("Failed to publish message to Kafka: %v", err)
		return err
	}

	return nil
}

// Close shuts down the Kafka producer
func (p *Producer) Close() error {
	return p.producer.Close()
}
