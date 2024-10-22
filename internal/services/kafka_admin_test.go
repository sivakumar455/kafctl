package services

import (
	"fmt"
	"kafctl/internal/config"
	"reflect"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func TestGetClusterDetails(t *testing.T) {
	config.KafkaBroker = "localhost:9092"
	want := kafka.BrokerMetadata{ID: 1, Host: "localhost", Port: 9092}
	app := KafAdmin{}
	clusterInfo, err := app.GetClusterDetails()
	if err != nil {
		t.Errorf("Test failed with err %v", err)
	}

	if !reflect.DeepEqual(want, clusterInfo[0]) {
		t.Error("cluster info assertion error")
	}
	fmt.Println(clusterInfo)

}

func TestCreateTopic(t *testing.T) {
	config.KafkaBroker = "localhost:9092"

	topic := "Aatest101"
	app := KafAdmin{}

	err := app.CreateTopic(topic, 2, 1)
	if err != nil {
		t.Error("Error creating topic")
	}
	res, err := app.DescribeTopic(topic)
	if err != nil {
		t.Error("Error getting topic details")
	}

	if res.TopicDescriptions[0].Name != topic {
		t.Error("Error creating topic name", topic)
	}

	err = app.DeleteTopic(topic)

	if err != nil {
		t.Error("Error deleting topic name:", topic)
	}

	_, err = app.DescribeTopic(topic)
	if err != nil {
		t.Error("Error getting topic details: ", topic)
	}

}

func TestGetTopicDetails(t *testing.T) {
	config.KafkaBroker = "localhost:9092"
	topic := "Aatest123"
	app := KafAdmin{}

	err := app.CreateTopic(topic, 2, 1)
	if err != nil {
		t.Error("Error creating topic", topic)
	}

	res, err := app.DescribeTopic(topic)
	if err != nil {
		t.Error("Error getting topic details")
	}

	if res.TopicDescriptions[0].Name != topic {
		t.Error("Err in getting topic details")
	}

}

func TestListOffsets(t *testing.T) {

	config.KafkaBroker = "localhost:9092"

	topic := "Aatest11"
	partition := 0
	app := KafAdmin{}

	err := app.ListOffsets(topic, partition)
	if err != nil {
		t.Error("Err getting offsetlist")
	}
}
