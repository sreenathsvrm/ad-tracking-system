package kafka

import (
	"log"
	"sync"

	"github.com/IBM/sarama"
)

// Consumer represents a Kafka consumer
type Consumer struct {
	consumer sarama.Consumer
	handler  func(message *sarama.ConsumerMessage)
	wg       sync.WaitGroup
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(brokers []string, handler func(message *sarama.ConsumerMessage)) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		consumer: consumer,
		handler:  handler,
	}, nil
}

// Consume starts consuming messages from the specified topic
func (c *Consumer) Consume(topic string) {
	partitions, err := c.consumer.Partitions(topic)
	if err != nil {
		log.Printf("Failed to get partitions: %v", err)
		return
	}

	for _, partition := range partitions {
		c.wg.Add(1)
		go c.consumePartition(topic, partition)
	}

	log.Printf("Started consuming messages from topic: %s", topic)
}

// consumePartition consumes messages from a specific partition
func (c *Consumer) consumePartition(topic string, partition int32) {
	defer c.wg.Done()

	partitionConsumer, err := c.consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
	if err != nil {
		log.Printf("Failed to consume partition: %v", err)
		return
	}
	defer partitionConsumer.Close()

	for msg := range partitionConsumer.Messages() {
		c.handler(msg)
	}
}

// Close stops the consumer and releases resources
func (c *Consumer) Close() error {
	c.wg.Wait()
	return c.consumer.Close()
}
