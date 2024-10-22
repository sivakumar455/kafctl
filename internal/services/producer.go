package services

import (
	"fmt"
	"kafctl/internal/config"
	"os"

	"math/rand"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func NewProducer() *kafka.Producer {

	cfg, err := CreateProducerConfig()
	if err != nil {
		fmt.Printf("Failed to create producer: %s", err)
		return nil
	}

	producer, err := kafka.NewProducer(cfg)

	if err != nil {
		fmt.Printf("Failed to create producer: %s", err)
		os.Exit(1)
	}

	return producer
}

func ProduceMessage() {

	producer := NewProducer()
	kafkaMsg := func(topic string, key, data []byte) *kafka.Message {
		return &kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Key:            key,
			Value:          data,
		}
	}

	users := [...]string{"eabara", "jsmith", "sgarcia", "jbernard", "htanaka", "awalther"}
	items := [...]string{"book", "alarm clock", "t-shirts", "gift card", "batteries"}
	//topic := "purchases"

	for n := 0; n < 10; n++ {
		key := users[rand.Intn(len(users))]
		data := items[rand.Intn(len(items))]

		producer.Produce(kafkaMsg(config.Topic, []byte(key), []byte(data)), nil)
	}

	// Wait for all messages to be delivered
	producer.Flush(15 * 1000)
	producer.Close()
}
