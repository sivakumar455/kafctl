package services

import (
	"kafctl/internal/logger"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type IRdConsumer interface {
	Close() (err error)

	ReadMessage(timeout time.Duration) (*kafka.Message, error)
	SubscribeTopics(topics []string, rebalanceCb kafka.RebalanceCb) (err error)
	Subscribe(topic string, rebalanceCb kafka.RebalanceCb) (err error)

	GetMetadata(topic *string, allTopics bool, timeoutMs int) (*kafka.Metadata, error)

	QueryWatermarkOffsets(topic string, partition int32, timeoutMs int) (low int64, high int64, err error)

	Assignment() (partitions []kafka.TopicPartition, err error)
	Seek(partition kafka.TopicPartition, ignoredTimeoutMs int) error
}

type RdConsumer struct {
	rdc IRdConsumer
}

// func CreateRdConsumer() (*kafka.Consumer, error) {

// 	// Create a new Kafka consumer with SSL configuration
// 	consumerCfg, err := CreateConsumerConfig()
// 	if err != nil {
// 		return nil, err
// 	}
// 	c, err := kafka.NewConsumer(consumerCfg)
// 	if err != nil {
// 		logger.Error("Error creating consumer ", "error", err)
// 		return nil, err
// 	}
// 	return c, nil
// }

func (c *RdConsumer) Close() error {

	err := c.rdc.Close()
	if err != nil {
		return err
	}

	return nil
}

func (c *RdConsumer) ReadMessage(timeout time.Duration) (*kafka.Message, error) {
	msg, err := c.rdc.ReadMessage(timeout)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (c *RdConsumer) SubscribeTopics(topics []string, rebalanceCb kafka.RebalanceCb) (err error) {

	err = c.rdc.SubscribeTopics(topics, rebalanceCb)
	if err != nil {
		logger.Error("Error subscribing topic", "error", err)
		return err
	}

	return nil
}

func (c *RdConsumer) Subscribe(topic string, rebalanceCb kafka.RebalanceCb) (err error) {
	err = c.rdc.Subscribe(topic, rebalanceCb)
	if err != nil {
		return err
	}
	return nil
}

func (c *RdConsumer) QueryWatermarkOffsets(topic string, partition int32, timeoutMs int) (low int64, high int64, err error) {

	low, high, err = c.rdc.QueryWatermarkOffsets(topic, partition, timeoutMs)
	if err != nil {
		return 0, 0, err
	}
	return low, high, nil
}

func (c *RdConsumer) GetMetadata(topic *string, allTopics bool, timeoutMs int) (*kafka.Metadata, error) {

	metadata, err := c.rdc.GetMetadata(topic, allTopics, timeoutMs)
	if err != nil {
		return nil, err
	}
	return metadata, nil
}

func (c *RdConsumer) Assignment() (partitions []kafka.TopicPartition, err error) {

	partitions, err = c.rdc.Assignment()
	if err != nil {
		return nil, err
	}
	return partitions, nil
}

func (c *RdConsumer) Seek(partition kafka.TopicPartition, ignoredTimeoutMs int) error {
	err := c.rdc.Seek(partition, ignoredTimeoutMs)
	if err != nil {
		return err
	}
	return nil
}
