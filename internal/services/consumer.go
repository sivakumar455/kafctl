package services

import (
	"fmt"
	"kafctl/internal/config"
	"os"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func NewConsumer() (*kafka.Consumer, error) {

	// Create a new Kafka consumer with SSL configuration
	consumerCfg, err := CreateConsumerConfig()
	if err != nil {
		return nil, err
	}
	c, err := kafka.NewConsumer(consumerCfg)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Consuming from topic: %s\n", config.Topic)
	return c, nil
}

func ConsumeMessages(kafkaBroker, groupId, configFile, topic string, enableSSL bool) error {

	c, err := NewConsumer()
	if err != nil {
		panic(err)
	}
	defer c.Close()
	fmt.Printf("Consuming from topic: %s\n", topic)

	// Subscribe to the topic
	c.SubscribeTopics([]string{topic}, nil)

	// Consume messages
	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			// Write message details to the file
			headers := ""
			for _, header := range msg.Headers {
				headers += fmt.Sprintf("%s: %s, ", header.Key, string(header.Value))
			}
			fmt.Printf("Offset=%d, Key=%s, TimeStamp=%s \n",
				msg.TopicPartition.Offset, string(msg.Key), msg.Timestamp)

			fmt.Printf("Offset=%d, Key=%s, Headers=%s,\nMessage=%s \n\n",
				msg.TopicPartition.Offset, string(msg.Key), headers, string(msg.Value))

		} else {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}

}

func ConsumeMessagesInFile() error {

	c, err := NewConsumer()
	if err != nil {
		panic(err)
	}
	defer c.Close()

	// Open a file to write the messages
	file, err := os.Create(config.OutputFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Subscribe to the topic
	c.SubscribeTopics([]string{config.Topic}, nil)

	fmt.Printf("Consuming from topic: %s, output file: %s\n", config.Topic, config.OutputFile)

	// Consume messages
	for {
		msg, err := c.ReadMessage(-1)
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
