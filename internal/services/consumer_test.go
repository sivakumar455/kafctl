package services

import (
	"fmt"
	"kafctl/internal/config"
	"testing"
)

func Test_GetTopicOffsets(t *testing.T) {

	config.KafkaBroker = "localhost:9092"
	topic := "__consumer_offsets"
	consumer, err := NewConsumer()
	if err != nil {
		t.Error("Error creating consumer")
	}

	err = consumer.GetTopicOffsets(&topic)
	if err != nil {
		t.Error("Error getting offsets")
	}

	fmt.Println("Test completed")

	//consumer.Close()

}

func Test_ConsumeMessages(t *testing.T) {

	config.KafkaBroker = "localhost:9092"
	config.GroupId = "test-grp4"
	topic := "ATest_topic1"

	consumer, err := NewConsumer()
	if err != nil {
		t.Error("Error creating consumer")
	}

	err = consumer.ConsumeMessages(topic)
	if err != nil {
		t.Error("Error consuming msg")
	}
	consumer.Close()
}

func Test_ConMessage(t *testing.T) {

	config.KafkaBroker = "localhost:9092"
	config.GroupId = "test-grp12"
	topic := "ATest_topic1"

	consumer, err := NewConsumer()
	if err != nil {
		t.Error("Error creating consumer")
	}

	msg, err := consumer.ConsumeMessage(topic)
	if err != nil {
		t.Error("Error consuming msg")
	}
	fmt.Println(len(msg))
	consumer.Close()
}
