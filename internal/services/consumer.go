package services

import (
	"fmt"
	"kafctl/internal/config"
	"kafctl/internal/logger"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type IConsumer interface {
	GetOffset(topic string, partition int, timeoutMs int) error
	GetTopicOffsets(topic *string) error
	ConsumeMessages(topic string) error
	ConsumeMessage(topic string) ([]*kafka.Message, error)
	ConsumeMessagesInFile() error
	Close() error
}

type Consumer struct {
	consumer IRdConsumer
}

func NewConsumer() (IConsumer, error) {
	// Create a new Kafka consumer
	consumer, err := CreateRdConsumer()
	if err != nil {
		return &Consumer{}, err
	}

	return &Consumer{consumer: consumer}, nil
}

func (c *Consumer) GetTopicOffsets(topic *string) error {

	metaData, err := c.consumer.GetMetadata(topic, false, 1000)
	if err != nil {
		logger.Error("Error getting consumer info", "error", err)
		return err
	}

	for _, partition := range metaData.Topics[*topic].Partitions {

		low, high, err := c.consumer.QueryWatermarkOffsets(*topic, partition.ID, 1000)
		if err != nil {
			logger.Error("Error getting offsets", "error", err)
			return err
		}

		logger.Info("Offset", "Id", partition.ID, "low", low, "high", high)
	}

	return nil
}

func (c *Consumer) GetOffset(topic string, partition int, timeoutMs int) error {

	// Query the watermark offsets
	low, high, err := c.consumer.QueryWatermarkOffsets(topic, int32(partition), 1000)
	if err != nil {
		logger.Error("Error getting offsets", "error", err)
		return err
	}

	logger.Info("First offset ", "firstOffset", low)
	logger.Info("Last offset", "lastOffset", high)
	return nil
}

func (c *Consumer) ConsumeMessage(topic string) ([]*kafka.Message, error) {

	logger.Info("Consuming from topic", "topic", topic)

	// Subscribe to the topic
	err := c.consumer.Subscribe(topic, nil)
	if err != nil {
		logger.Error("Error subscribing topic", "error", err)
		return nil, err
	}

	// Get the list of partitions for the topic
	partitions, err := c.consumer.Assignment()
	if err != nil {
		logger.Error("Failed to get partitions", "error", err)
	}

	// Seek to the beginning of each partition
	for _, partition := range partitions {
		err = c.consumer.Seek(kafka.TopicPartition{Topic: &topic, Partition: partition.Partition, Offset: kafka.OffsetBeginning}, -1)
		if err != nil {
			logger.Error("Failed to get partitions", "error", err)
		}
	}
	var msgRes []*kafka.Message
	// Consume messages
	for {
		msg, err := c.consumer.ReadMessage(1000 * time.Millisecond)
		if err == nil {
			msgRes = append(msgRes, msg)
			headers := ""
			for _, header := range msg.Headers {
				headers += fmt.Sprintf("%s: %s, ", header.Key, string(header.Value))
			}
			logger.Info(fmt.Sprintf("Offset=%d, Key=%s, TimeStamp=%s",
				msg.TopicPartition.Offset, string(msg.Key), msg.Timestamp))

			logger.Info(fmt.Sprintf("Offset=%d, Key=%s, Headers=%s,Message=%s",
				msg.TopicPartition.Offset, string(msg.Key), headers, string(msg.Value)))

		} else if err.(kafka.Error).Code() == kafka.ErrTimedOut {
			// Timeout error, no more messages to read
			logger.Info("No more messages to read, exiting.")
			break
		} else {
			logger.Error("Consumer error", "error", err, "msg", msg)
			return nil, err
		}
	}
	return msgRes, nil
}

func (c *Consumer) ConsumeMessages(topic string) error {

	logger.Info("Consuming from topic", "topic", topic)

	// Subscribe to the topic
	err := c.consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		logger.Error("Error subscribing topic", "error", err)
		return err
	}
	// Consume messages
	for {
		msg, err := c.consumer.ReadMessage(1000 * time.Millisecond)
		if err == nil {
			headers := ""
			for _, header := range msg.Headers {
				headers += fmt.Sprintf("%s: %s, ", header.Key, string(header.Value))
			}
			logger.Info(fmt.Sprintf("Offset=%d, Key=%s, TimeStamp=%s",
				msg.TopicPartition.Offset, string(msg.Key), msg.Timestamp))

			logger.Info(fmt.Sprintf("Offset=%d, Key=%s, Headers=%s,Message=%s",
				msg.TopicPartition.Offset, string(msg.Key), headers, string(msg.Value)))

		} else if err.(kafka.Error).Code() == kafka.ErrTimedOut {
			// Timeout error, no more messages to read
			logger.Info("No more messages to read, exiting.")
			break
		} else {
			logger.Error("Consumer error", "error", err, "msg", msg)
			return err
		}
	}
	return nil
}

func (c *Consumer) ConsumeMessagesInFile() error {

	// Open a file to write the messages
	file, err := os.Create(config.OutputFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Subscribe to the topic
	c.consumer.SubscribeTopics([]string{config.Topic}, nil)

	fmt.Printf("Consuming from topic: %s, output file: %s\n", config.Topic, config.OutputFile)

	// Consume messages
	for {
		msg, err := c.consumer.ReadMessage(-1)
		if err == nil {
			// Write message details to the file
			headers := ""
			for _, header := range msg.Headers {
				headers += fmt.Sprintf("%s: %s, ", header.Key, string(header.Value))
			}
			fmt.Printf("Offset=%d, Key=%s, TimeStamp=%s \n",
				msg.TopicPartition.Offset, string(msg.Key), msg.Timestamp)

			_, err := file.WriteString(fmt.Sprintf("Offset=%d, Key=%s, Headers=%s,\nMessage=%s \n\n",
				msg.TopicPartition.Offset, string(msg.Key), headers, string(msg.Value)))
			if err != nil {
				panic(err)
			}
		} else {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}

}

func (c *Consumer) Close() error {
	err := c.consumer.Close()
	if err != nil {
		return err
	}

	return nil
}
