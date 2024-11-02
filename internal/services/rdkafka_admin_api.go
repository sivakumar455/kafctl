package services

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type IRdAdminClient interface {
	GetMetadata(topic *string, allTopics bool, timeoutMs int) (*kafka.Metadata, error)
	Close()
	CreateTopics(ctx context.Context, topics []kafka.TopicSpecification,
		options ...kafka.CreateTopicsAdminOption) (result []kafka.TopicResult, err error)
	DeleteTopics(ctx context.Context, topics []string,
		options ...kafka.DeleteTopicsAdminOption) (result []kafka.TopicResult, err error)
	DescribeTopics(ctx context.Context, topics kafka.TopicCollection,
		options ...kafka.DescribeTopicsAdminOption) (result kafka.DescribeTopicsResult, err error)
	ListOffsets(ctx context.Context, topicPartitionOffsets map[kafka.TopicPartition]kafka.OffsetSpec,
		options ...kafka.ListOffsetsAdminOption) (result kafka.ListOffsetsResult, err error)
}

type RdKafkaAdmin struct {
	admin IRdAdminClient
}

func (ka *RdKafkaAdmin) GetMetadata(topic *string, allTopics bool, timeoutMs int) (*kafka.Metadata, error) {
	return ka.admin.GetMetadata(topic, allTopics, timeoutMs)
}

func (ka *RdKafkaAdmin) Close() {
	ka.admin.Close()
}

func (ka *RdKafkaAdmin) CreateTopics(ctx context.Context, topics []kafka.TopicSpecification,
	options ...kafka.CreateTopicsAdminOption) (result []kafka.TopicResult, err error) {
	return ka.admin.CreateTopics(ctx, topics, options...)
}

func (ka *RdKafkaAdmin) DeleteTopics(ctx context.Context, topics []string,
	options ...kafka.DeleteTopicsAdminOption) (result []kafka.TopicResult, err error) {
	return ka.admin.DeleteTopics(ctx, topics, options...)
}

func (ka *RdKafkaAdmin) DescribeTopics(ctx context.Context, topics kafka.TopicCollection,
	options ...kafka.DescribeTopicsAdminOption) (result kafka.DescribeTopicsResult, err error) {
	return ka.admin.DescribeTopics(ctx, topics, options...)
}

func (ka *RdKafkaAdmin) ListOffsets(ctx context.Context, topicPartitionOffsets map[kafka.TopicPartition]kafka.OffsetSpec,
	options ...kafka.ListOffsetsAdminOption) (result kafka.ListOffsetsResult, err error) {
	return ka.admin.ListOffsets(ctx, topicPartitionOffsets, options...)
}
