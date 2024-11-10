package mocks

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/stretchr/testify/mock"
)

// MockAdminClient is a mock implementation of the IRdKafkaAdminClient interface
type MockAdminClient struct {
	mock.Mock
}

func (m *MockAdminClient) GetMetadata(topic *string, allTopics bool, timeoutMs int) (*kafka.Metadata, error) {
	args := m.Called(topic, allTopics, timeoutMs)
	return args.Get(0).(*kafka.Metadata), args.Error(1)
}

func (m *MockAdminClient) Close() {
	m.Called()
}

func (m *MockAdminClient) CreateTopics(ctx context.Context, topics []kafka.TopicSpecification,
	options ...kafka.CreateTopicsAdminOption) (result []kafka.TopicResult, err error) {
	args := m.Called(ctx, topics, options)
	return args.Get(0).([]kafka.TopicResult), args.Error(1)
}

func (m *MockAdminClient) DeleteTopics(ctx context.Context, topics []string,
	options ...kafka.DeleteTopicsAdminOption) (result []kafka.TopicResult, err error) {
	args := m.Called(ctx, topics, options)
	return args.Get(0).([]kafka.TopicResult), args.Error(1)
}

func (m *MockAdminClient) DescribeTopics(ctx context.Context, topics kafka.TopicCollection,
	options ...kafka.DescribeTopicsAdminOption) (result kafka.DescribeTopicsResult, err error) {
	args := m.Called(ctx, topics, options)
	return args.Get(0).(kafka.DescribeTopicsResult), args.Error(1)
}

func (m *MockAdminClient) ListOffsets(ctx context.Context, topicPartitionOffsets map[kafka.TopicPartition]kafka.OffsetSpec,
	options ...kafka.ListOffsetsAdminOption) (result kafka.ListOffsetsResult, err error) {
	args := m.Called(ctx, topicPartitionOffsets, options)
	return args.Get(0).(kafka.ListOffsetsResult), args.Error(1)
}
