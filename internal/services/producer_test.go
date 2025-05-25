package services

import (
	"kafctl/internal/config"
	"math/rand"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func TestPublish(t *testing.T) {

}

func ProduceMessageTest() error {

	producer, err := NewProducer()
	if err != nil {
		return err
	}
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
	return nil
}
