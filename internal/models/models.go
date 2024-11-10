package models

import "github.com/confluentinc/confluent-kafka-go/v2/kafka"

type TopicDetails struct {
	Name                 string
	TopicID              kafka.UUID
	IsInternal           bool
	AuthorizedOperations []kafka.ACLOperation
	Partitions           []kafka.TopicPartitionInfo
}

type BrokerInfo struct {
	Brokers []kafka.BrokerMetadata
	Status  string
	Topics  map[string]kafka.TopicMetadata
}
