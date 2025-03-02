package producer

import (
	"ad-tracking-system/internal/domain/models"
	"encoding/json"

	"github.com/IBM/sarama"
)

// KafkaProducer represents a Kafka producer for ad click events
type KafkaProducer struct {
	producer sarama.SyncProducer
	topic    string
}

// NewKafkaProducer creates a new Kafka producer
func NewKafkaProducer(brokers []string, topic string) (*KafkaProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &KafkaProducer{
		producer: producer,
		topic:    topic,
	}, nil
}

// PublishClickEvent publishes a click event to Kafka
func (kp *KafkaProducer) PublishClickEvent(click models.ClickEvent) error {
	message, err := json.Marshal(click)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: kp.topic,
		Value: sarama.ByteEncoder(message),
	}

	_, _, err = kp.producer.SendMessage(msg)
	return err
}

// Close shuts down the Kafka producer
func (kp *KafkaProducer) Close() error {
	return kp.producer.Close()
}
