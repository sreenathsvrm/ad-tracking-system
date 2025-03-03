package consumer

import (
	"ad-tracking-system/internal/utils/circuitbreaker"
	"log"
	"sync"

	"github.com/IBM/sarama"
	"github.com/sony/gobreaker"
)

type KafkaConsumer struct {
	consumer sarama.Consumer
	handler  func(message *sarama.ConsumerMessage)
	wg       sync.WaitGroup
	cb       *gobreaker.CircuitBreaker
}

func NewKafkaConsumer(brokers []string, handler func(message *sarama.ConsumerMessage)) (*KafkaConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{
		consumer: consumer,
		handler:  handler,
		cb:       circuitbreaker.NewCircuitBreaker("kafka-consumer"), // Initialize circuit breaker
	}, nil
}

func (kc *KafkaConsumer) Consume(topic string) {
	partitions, err := kc.consumer.Partitions(topic)
	if err != nil {
		log.Printf("Failed to get partitions: %v", err)
		return
	}

	for _, partition := range partitions {
		kc.wg.Add(1)
		go kc.consumePartition(topic, partition)
	}

	log.Printf("Started consuming messages from topic: %s", topic)
}

func (kc *KafkaConsumer) consumePartition(topic string, partition int32) {
	defer kc.wg.Done()

	partitionConsumer, err := kc.consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
	if err != nil {
		log.Printf("Failed to consume partition: %v", err)
		return
	}
	defer partitionConsumer.Close()

	for msg := range partitionConsumer.Messages() {
		// Wrap message handling with circuit breaker
		_, err := kc.cb.Execute(func() (interface{}, error) {
			kc.handler(msg)
			return nil, nil
		})
		if err != nil {
			log.Printf("Failed to process message (circuit breaker): %v", err)
		}
	}
}

func (kc *KafkaConsumer) Close() error {
	kc.wg.Wait()
	return kc.consumer.Close()
}
