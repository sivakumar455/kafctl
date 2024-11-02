package services

import (
	"fmt"
	"kafctl/internal/config"
	"testing"
)

func TestsGetOffsets(t *testing.T) {

	config.KafkaBroker = "localhost:9092"
	topic := "Aatest11"
	partition := 0
	err := GetOffsets(topic, partition)
	fmt.Print("Test completed com")
	if err != nil {
		t.Error("Error getting offsets")
	}

	fmt.Print("Test completed")
}
