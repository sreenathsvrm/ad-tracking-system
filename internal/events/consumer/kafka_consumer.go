package consumer

import (
	"log"
	"sync"

	"github.com/IBM/sarama"
)

// KafkaConsumer represents a Kafka consumer for ad click events
type KafkaConsumer struct {
	consumer sarama.Consumer
	handler  func(message *sarama.ConsumerMessage)
	wg       sync.WaitGroup
}

// NewKafkaConsumer creates a new Kafka consumer
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
	}, nil
}

// Consume starts consuming messages from the specified topic
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

// consumePartition consumes messages from a specific partition
func (kc *KafkaConsumer) consumePartition(topic string, partition int32) {
	defer kc.wg.Done()

	partitionConsumer, err := kc.consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
	if err != nil {
		log.Printf("Failed to consume partition: %v", err)
		return
	}
	defer partitionConsumer.Close()

	for msg := range partitionConsumer.Messages() {
		kc.handler(msg)
	}
}

// Close stops the consumer and releases resources
func (kc *KafkaConsumer) Close() error {
	kc.wg.Wait()
	return kc.consumer.Close()
}
