package services

import (
	"context"
	"encoding/json"
	"fmt"
	"kafctl/internal/config"
	"kafctl/internal/logger"
	"kafctl/internal/models"
	"os"
	"strconv"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type SSLConfig struct {
	SecurityProtocol string `json:"securityProtocol"`
	SslCaLocation    string `json:"sslCaLocation"`
	SslCertLocation  string `json:"sslCertLocation"`
	SslKeyLocation   string `json:"sslKeyLocation"`
	SslKeyPassword   string `json:"sslKeyPassword"`
}

func SetSSLConfig(consumerCfg *kafka.ConfigMap, configFile string) error {

	// Read SSL configuration from file
	sslCfgFile, err := os.Open(configFile)
	if err != nil {
		return err
	}
	defer sslCfgFile.Close()

	var sslConfig SSLConfig
	if err := json.NewDecoder(sslCfgFile).Decode(&sslConfig); err != nil {
		return err
	}

	consumerCfg.SetKey("security.protocol", sslConfig.SecurityProtocol)
	if sslConfig.SecurityProtocol == "SSL" || sslConfig.SecurityProtocol == "SASL_SSL" {
		consumerCfg.SetKey("ssl.ca.location", sslConfig.SslCaLocation)
		consumerCfg.SetKey("ssl.certificate.location", sslConfig.SslCertLocation)
		consumerCfg.SetKey("ssl.key.location", sslConfig.SslKeyLocation)
		consumerCfg.SetKey("ssl.key.password", sslConfig.SslKeyPassword)
	}
	return nil
}

func CreateConsumerConfig() (*kafka.ConfigMap, error) {

	// Create a new Kafka consumer with SSL configuration
	consumerCfg := &kafka.ConfigMap{
		"bootstrap.servers": config.KafkaBroker,
		"group.id":          config.GroupId,
		"auto.offset.reset": "earliest",
	}

	// Read SSL configuration from file
	if config.EnableSSL {
		err := SetSSLConfig(consumerCfg, config.ConfigFile)
		if err != nil {
			return nil, err
		}
	}

	return consumerCfg, nil
}

func CreateProducerConfig() (*kafka.ConfigMap, error) {

	// Create a new Kafka consumer with SSL configuration
	producerCfg := &kafka.ConfigMap{
		"bootstrap.servers": config.KafkaBroker,
		"acks":              "all",
	}

	// Read SSL configuration from file
	if config.EnableSSL {
		err := SetSSLConfig(producerCfg, config.ConfigFile)
		if err != nil {
			return nil, err
		}
	}

	return producerCfg, nil
}

func CreateAdminConfig() (*kafka.ConfigMap, error) {

	// Create a new Kafka consumer with SSL configuration
	adminCfg := &kafka.ConfigMap{"bootstrap.servers": config.KafkaBroker}

	// Read SSL configuration from file
	if config.EnableSSL {
		err := SetSSLConfig(adminCfg, config.ConfigFile)
		if err != nil {
			return nil, err
		}
	}

	return adminCfg, nil
}

func CreateAdminClient() (*kafka.AdminClient, error) {

	cfg, err := CreateAdminConfig()
	if err != nil {
		return nil, err
	}

	admin, err := kafka.NewAdminClient(cfg)
	if err != nil {
		logger.Error("Failed to create admin Client", "error", err)
		return nil, err
	}

	return admin, nil
}

func GetClusterDetails() ([]models.Broker, error) {

	adminClient, err := CreateAdminClient()
	if err != nil {
		return nil, err
	}
	defer adminClient.Close()

	// Fetch metadata
	metadata, err := adminClient.GetMetadata(nil, true, 10000)
	if err != nil {
		return nil, err
	}

	// Print broker information
	logger.Info("Brokers Details")
	brokers := []models.Broker{}
	for _, broker := range metadata.Brokers {
		logger.Info(fmt.Sprintf("ID: %d, Host: %s, Port: %d\n", broker.ID, broker.Host, broker.Port))
		brokers = append(brokers, models.Broker{Id: broker.ID, Host: broker.Host, Port: broker.Port})
	}

	return brokers, nil

}

func GetAllTopics() (*models.TopicsMap, error) {

	admin, err := CreateAdminClient()
	if err != nil {
		return nil, err
	}
	logger.Info("Admin created")
	metadata, err := admin.GetMetadata(nil, true, 5000)
	if err != nil {
		return nil, err
	}
	logger.Info("Metadat call was successfull")
	// topics := metadata.Topics

	// logger.Info("Total topics: ", "", len(topics))

	// for k := range topics {
	// 	logger.Debug(fmt.Sprintf("Topics: %v", k))
	// }

	topicLst := models.TopicsMap{}

	// Print topic information
	topicsMap := make(map[string]models.TopicMetadata)

	for topicName, topic := range metadata.Topics {
		logger.Debug(fmt.Sprintf("Topic: %s, Partitions: %d\n", topicName, len(topic.Partitions)))
		partitionMetaData := []models.PartitionMetadata{}
		for partitionID, partition := range topic.Partitions {
			logger.Debug(fmt.Sprintf("  Partition: %d, Leader: %d, Replicas: %v, ISR: %v\n",
				partitionID, partition.Leader, partition.Replicas, partition.Isrs))
			partitionMetaData = append(partitionMetaData, models.PartitionMetadata{ID: partition.ID,
				Leader: partition.Leader, Replicas: partition.Replicas, Isrs: partition.Isrs})
		}

		topicsMap[topicName] = models.TopicMetadata{Topic: topicName, Partitions: partitionMetaData}

	}
	logger.Info("Total topics: ", "", len(topicsMap))
	topicLst.TopicsMap = topicsMap

	return &topicLst, nil

}

func CreateTopic(topic string, numParts, replicationFactor int) error {

	admin, err := CreateAdminClient()
	if err != nil {
		return err
	}
	defer admin.Close()
	logger.Info("Admin created")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create topics on cluster.
	// Set Admin options to wait for the operation to finish (or at most 60s)
	maxDur, err := time.ParseDuration("5s")
	if err != nil {
		return err
	}
	results, err := admin.CreateTopics(
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

	// Print results
	for _, result := range results {
		logger.Info("result: ", "result", result)
	}
	return nil

}

func DeleteTopic(topic string) error {

	admin, err := CreateAdminClient()
	if err != nil {
		return err
	}
	defer admin.Close()
	logger.Info("Admin created")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create topics on cluster.
	// Set Admin options to wait for the operation to finish (or at most 60s)
	maxDur, err := time.ParseDuration("5s")
	if err != nil {
		return err
	}

	topics := []string{topic}
	results, err := admin.DeleteTopics(ctx, topics, kafka.SetAdminOperationTimeout(maxDur))
	if err != nil {
		logger.Error("Failed to delete topics: ", "error", err)
		return err
	}

	res := results[0]
	logger.Info("Topic deleted successfully", "result", res)

	return nil

}

func DescribeTopic(topic string) ([]models.TopicDetails, error) {

	admin, err := CreateAdminClient()
	if err != nil {
		return nil, err
	}
	defer admin.Close()
	logger.Info("Admin created")

	// Call DescribeTopics.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	includeAuthorizedOperations := true

	describeTopicsResult, err := admin.DescribeTopics(
		ctx,
		kafka.NewTopicCollectionOfTopicNames([]string{topic}),
		kafka.SetAdminOptionIncludeAuthorizedOperations(
			includeAuthorizedOperations))
	if err != nil {
		fmt.Printf("Failed to describe topics: %s\n", err)
		return nil, err
	}

	// Print results
	fmt.Printf("A total of %d topic(s) described:\n\n",
		len(describeTopicsResult.TopicDescriptions))

	for _, t := range describeTopicsResult.TopicDescriptions {
		if t.Error.Code() != 0 {
			fmt.Printf("Topic: %s has error: %s\n",
				t.Name, t.Error)
			continue
		}
		fmt.Printf("Topic: %s has succeeded\n", t.Name)
		fmt.Printf("Topic Id: %s\n", t.TopicID)
		if includeAuthorizedOperations {
			fmt.Printf("Allowed operations: %s\n", t.AuthorizedOperations)
		}
		for i := 0; i < len(t.Partitions); i++ {
			fmt.Printf("\tPartition id: %d with leader: %s\n",
				t.Partitions[i].Partition, t.Partitions[i].Leader)
			fmt.Printf("\t\tThe in-sync replica count is: %d, they are: \n\t\t%s\n",
				len(t.Partitions[i].Isr), t.Partitions[i].Isr)
			fmt.Printf("\t\tThe replica count is: %d, they are: \n\t\t%s\n",
				len(t.Partitions[i].Replicas), t.Partitions[i].Replicas)
		}
		fmt.Printf("\n")
	}

	var topics []models.TopicDetails

	for _, t := range describeTopicsResult.TopicDescriptions {
		if t.Error.Code() != 0 {
			fmt.Printf("Topic: %s has error: %s\n", t.Name, t.Error)
			continue
		}

		var partitions []models.PartitionDetails
		for _, p := range t.Partitions {
			partitions = append(partitions, models.PartitionDetails{
				PartitionID: p.Partition,
				Leader:      strconv.Itoa(p.Leader.ID),
			})
		}

		topics = append(topics, models.TopicDetails{
			Name:       t.Name,
			Partitions: partitions,
		})

	}
	return topics, nil

}
