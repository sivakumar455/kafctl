package services

import (
	"fmt"
	"kafctl/internal/config"
	"strings"

	"math/rand"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func NewProducer() (*kafka.Producer, error) {

	cfg, err := CreateProducerConfig()
	if err != nil {
		fmt.Printf("Failed to create producer: %s", err)
		return nil, err
	}

	producer, err := kafka.NewProducer(cfg)

	if err != nil {
		fmt.Printf("Failed to create producer: %s", err)
		return nil, err
	}

	return producer, nil
}

func ProduceMessage(topic string, key, data []byte, headerMap string) error {

	producer, err := NewProducer()
	if err != nil {
		return err
	}

	headers := []kafka.Header{}

	if headerMap != "" {
		headerList := strings.Split(headerMap, ",")
		for _, val := range headerList {
			res := strings.Split(val, "=")
			header := kafka.Header{Key: res[0], Value: []byte(res[1])}
			headers = append(headers, header)
		}
	}

	kafkaMsg := func(topic string, key, data []byte, headers []kafka.Header) *kafka.Message {
		return &kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Key:            key,
			Value:          data,
			Headers:        headers,
		}
	}

	err = producer.Produce(kafkaMsg(topic, []byte(key), []byte(data), headers), nil)
	if err != nil {
		return err
	}

	// Wait for all messages to be delivered
	producer.Flush(2 * 100)
	producer.Close()
	return nil
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
