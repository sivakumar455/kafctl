package services

import (
	"fmt"
	"kafctl/internal/config"
	"kafctl/internal/models"
	"reflect"
	"testing"
)

func TestGetClusterDetails(t *testing.T) {
	config.KafkaBroker = "localhost:9092"
	want := models.Broker{Id: 1, Host: "localhost", Port: 9092}
	clusterInfo, err := GetClusterDetails()
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
	err := CreateTopic(topic, 2, 1)
	if err != nil {
		t.Error("Error creating topic")
	}

	res, err := DescribeTopic(topic)
	if err != nil {
		t.Error("Error getting topic details")
	}

	if res[0].Name != topic {
		t.Error("Error creating topic name", topic)
	}

	err = DeleteTopic(topic)

	if err != nil {
		t.Error("Error deleting topic name:", topic)
	}

	_, err = DescribeTopic(topic)
	if err != nil {
		t.Error("Error getting topic details: ", topic)
	}

}
