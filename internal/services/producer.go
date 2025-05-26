package services

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

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

func ProduceMessage(topic, key, headerMap string, data []byte) error {

	producer, err := NewProducer()
	if err != nil {
		return err
	}

	defer producer.Close()

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

	deliveryCh := make(chan kafka.Event)

	err = producer.Produce(kafkaMsg(topic, []byte(key), []byte(data), headers), deliveryCh)
	if err != nil {
		slog.Error("Error producing message")
		return err
	}

	defer close(deliveryCh)

	// Wait for all messages to be delivered
	//producer.Flush(2 * 100)
	//producer.Close()

	select {

	case event := <-deliveryCh:
		msg := event.(*kafka.Message)
		if msg.TopicPartition.Error != nil {
			slog.Error("Failed to deliver message ", "error", msg.TopicPartition.Error)
			return err
		}
		slog.Info("Published message to ", "topic", *msg.TopicPartition.Topic, "partition", msg.TopicPartition.Partition, "offset", msg.TopicPartition.Offset)
	case <-time.After(10 * time.Second):
		slog.Error("Message delivery timed out")
		return fmt.Errorf("message delivery timed out")
	}

	slog.Info("Message published successfully")

	return nil
}
