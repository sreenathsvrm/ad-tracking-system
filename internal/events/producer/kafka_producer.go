package producer

import (
	"ad-tracking-system/internal/domain/models"
	"ad-tracking-system/internal/utils/circuitbreaker"
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
	"github.com/sony/gobreaker"
)

type KafkaProducer struct {
	producer sarama.SyncProducer
	topic    string
	cb       *gobreaker.CircuitBreaker
}

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
		cb:       circuitbreaker.NewCircuitBreaker("kafka-producer"), // Initialize circuit breaker
	}, nil
}

func (kp *KafkaProducer) PublishClickEvent(click models.ClickEvent) error {
	// Wrap Kafka operation with circuit breaker
	_, err := kp.cb.Execute(func() (interface{}, error) {
		message, err := json.Marshal(click)
		if err != nil {
			return nil, err
		}

		msg := &sarama.ProducerMessage{
			Topic: kp.topic,
			Value: sarama.ByteEncoder(message),
		}

		_, _, err = kp.producer.SendMessage(msg)
		return nil, err
	})
	if err != nil {
		log.Printf("Failed to publish click event (circuit breaker): %v", err)
		return err
	}

	return nil
}

func (kp *KafkaProducer) Close() error {
	return kp.producer.Close()
}
