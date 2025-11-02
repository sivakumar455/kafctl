package services

import (
	"context"
	"fmt"
	"kafctl/internal/logger"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type IKafAdmin interface {
	GetClusterDetails() ([]kafka.BrokerMetadata, error)
	GetAllTopics() (map[string]kafka.TopicMetadata, error)
	CreateTopic(topic string, numParts, replicationFactor int) error
	DeleteTopic(topic string) error
	DescribeTopic(topic string) (kafka.DescribeTopicsResult, error)
	GetListOffsets(topicName string, partition int) error
	Close()
}

type KafAdmin struct {
	admin IRdAdminClient
}

var NewKafAdmin = CreateKafAdmin

var kafkaAdminInstance IKafAdmin

func CreateKafAdmin() (IKafAdmin, error) {

	if kafkaAdminInstance == nil {
		adminClient, err := CreateAdminClient()
		if err != nil {
			return &KafAdmin{}, err
		}
		kafkaAdminInstance = &KafAdmin{admin: adminClient}
	}

	return kafkaAdminInstance, nil
}

func (ka *KafAdmin) GetClusterDetails() ([]kafka.BrokerMetadata, error) {

	logger.Debug("Fetching Brokers Details")
	metadata, err := ka.admin.GetMetadata(nil, true, 10000)
	if err != nil {
		return nil, err
	}

	for _, broker := range metadata.Brokers {
		logger.Debug(fmt.Sprintf("ID: %d, Host: %s, Port: %d\n", broker.ID, broker.Host, broker.Port))
	}

	return metadata.Brokers, nil
}

func (ka *KafAdmin) GetAllTopics() (map[string]kafka.TopicMetadata, error) {

	logger.Info("Fetching all topics")

	metadata, err := ka.admin.GetMetadata(nil, true, 3000)
	if err != nil {
		return nil, err
	}

	for topicName, topic := range metadata.Topics {
		logger.Debug(fmt.Sprintf("Topic: %s, Partitions: %d\n", topicName, len(topic.Partitions)))
		for partitionID, partition := range topic.Partitions {
			logger.Debug(fmt.Sprintf("  Partition: %d, Leader: %d, Replicas: %v, ISR: %v\n",
				partitionID, partition.Leader, partition.Replicas, partition.Isrs))
		}
	}
	logger.Info("Total number of topics: ", "total", len(metadata.Topics))

	return metadata.Topics, nil

}

func (ka *KafAdmin) CreateTopic(topic string, numParts, replicationFactor int) error {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	maxDur, err := time.ParseDuration("5s")
	if err != nil {
		return err
	}
	results, err := ka.admin.CreateTopics(
		ctx,
		// Multiple topics can be created simultaneously
		// by providing more TopicSpecification structs here.
		[]kafka.TopicSpecification{{
			Topic:             topic,
			NumPartitions:     numParts,
			ReplicationFactor: replicationFactor}},
		// Admin options
		kafka.SetAdminOperationTimeout(maxDur))
	if err != nil {
		logger.Error("Failed to create topic: ", "error", err)
		return err
	}

	// Check results for errors
	for _, result := range results {
		if result.Error.Code() != kafka.ErrNoError {
			errMsg := fmt.Sprintf("failed to create topic '%s': %s", result.Topic, result.Error.String())
			logger.Error("Failed to create topic", "topic", result.Topic, "error", errMsg)
			return fmt.Errorf("failed to create topic '%s': %s", result.Topic, result.Error.String())
		}
		logger.Info("Topic created successfully", "topic", result.Topic)
	}
	return nil
}

func (ka *KafAdmin) DeleteTopic(topic string) error {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	maxDur, err := time.ParseDuration("5s")
	if err != nil {
		return err
	}

	topics := []string{topic}
	results, err := ka.admin.DeleteTopics(ctx, topics, kafka.SetAdminOperationTimeout(maxDur))
	if err != nil {
		logger.Error("Failed to delete topics: ", "error", err)
		return err
	}

	res := results[0]
	logger.Info("Topic deleted successfully", "result", res)

	return nil
}

func (ka *KafAdmin) DescribeTopic(topic string) (kafka.DescribeTopicsResult, error) {

	logger.Info("getting details for topic - ", "topic", topic)

	// Call DescribeTopics.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	includeAuthorizedOperations := true

	describeTopicsResult, err := ka.admin.DescribeTopics(
		ctx,
		kafka.NewTopicCollectionOfTopicNames([]string{topic}),
		kafka.SetAdminOptionIncludeAuthorizedOperations(
			includeAuthorizedOperations))
	if err != nil {
		logger.Error("Failed to describe topics:", "error", err)
		return kafka.DescribeTopicsResult{}, err
	}

	// Print results
	logger.Debug(fmt.Sprintf("A total of %d topic(s) described:\n\n", len(describeTopicsResult.TopicDescriptions)))
	logger.Debug("Described for topic", "topic", describeTopicsResult.TopicDescriptions[0].Name)
	printTopicsDescriptionResult(describeTopicsResult, includeAuthorizedOperations)
	return describeTopicsResult, nil
}

func printTopicsDescriptionResult(describeTopicsResult kafka.DescribeTopicsResult, includeAuthorizedOperations bool) {
	for _, t := range describeTopicsResult.TopicDescriptions {
		if t.Error.Code() != 0 {
			logger.Error(fmt.Sprintf("Topic: %s has error: %s\n",
				t.Name, t.Error))
			continue
		}
		logger.Info("Topic has succeeded", "topic", t.Name, "Topic Id", t.TopicID)
		if includeAuthorizedOperations {
			logger.Debug("Allowed operations ", "operations", t.AuthorizedOperations)
		}
		for i := 0; i < len(t.Partitions); i++ {
			logger.Debug(fmt.Sprintf("Partition id: %d with leader: %s",
				t.Partitions[i].Partition, t.Partitions[i].Leader))
			logger.Debug(fmt.Sprintf("The in-sync replica count is: %d, they are: %s",
				len(t.Partitions[i].Isr), t.Partitions[i].Isr))
			logger.Debug(fmt.Sprintf("The replica count is: %d, they are: %s",
				len(t.Partitions[i].Replicas), t.Partitions[i].Replicas))
		}
	}
}

func (ka *KafAdmin) GetListOffsets(topicName string, partition int) error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	topicPartitionOffsets := make(map[kafka.TopicPartition]kafka.OffsetSpec)

	tp := kafka.TopicPartition{Topic: &topicName, Partition: int32(partition)}

	topicPartitionOffsets[tp] = kafka.EarliestOffsetSpec

	results, err := ka.admin.ListOffsets(ctx, topicPartitionOffsets,
		kafka.SetAdminIsolationLevel(kafka.IsolationLevelReadCommitted))
	if err != nil {
		fmt.Printf("Failed to List offsets: %v\n", err)
		return err
	}

	for tp, info := range results.ResultInfos {
		fmt.Printf("Topic: %s Partition: %d\n", *tp.Topic, tp.Partition)
		if info.Error.Code() != kafka.ErrNoError {
			fmt.Printf("	ErrorCode: %d ErrorMessage: %s\n\n", info.Error.Code(), info.Error.String())
		} else {
			fmt.Printf("	Offset: %d Timestamp: %d\n\n", info.Offset, info.Timestamp)
		}
	}

	return nil

}

func (ka *KafAdmin) Close() {
	ka.admin.Close()
}
