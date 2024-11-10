package services

import (
	"kafctl/internal/services/mocks"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setup(t *testing.T) (IKafAdmin, *mocks.MockAdminClient, func(), error) {
	mockAdmin := new(mocks.MockAdminClient)
	NewKafAdmin = func() (IKafAdmin, error) {
		return &KafAdmin{mockAdmin}, nil
	}
	resetAdmin := func() { NewKafAdmin = CreateKafAdmin } // Reset after test

	kafAdmin, err := NewKafAdmin()
	if err != nil {
		t.Error("Error creating KafAdmin")
		return nil, nil, nil, err
	}
	return kafAdmin, mockAdmin, resetAdmin, nil
}

func Test_NewAdminClient_Fail(t *testing.T) {
	// Mock implementation
	NewAdminClient = func(cfg *kafka.ConfigMap) (*kafka.AdminClient, error) {
		return nil, errors.New("mock error")
	}
	defer func() { NewAdminClient = kafka.NewAdminClient }() // Reset after test

	// Now call your function that uses CreateAdminClient
	_, err := CreateAdminClient()
	expectedErr := "mock error"

	if err == nil || err.Error() != expectedErr {
		t.Errorf("Expected %s, got %v", expectedErr, err)
	}
}

func Test_NewAdminClient_Success(t *testing.T) {
	// Mock implementation
	NewAdminClient = func(cfg *kafka.ConfigMap) (*kafka.AdminClient, error) {
		return &kafka.AdminClient{}, nil
	}
	defer func() { NewAdminClient = kafka.NewAdminClient }() // Reset after test

	// Now call your function that uses CreateAdminClient
	_, err := CreateAdminClient()
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
}

func Test_GetClusterDetails(t *testing.T) {

	kafAdmin, mockAdmin, resetAdmin, err := setup(t)
	if err != nil {
		t.Error("Error creating Admin")
	}

	defer resetAdmin()

	expectedBrokers := []kafka.BrokerMetadata{
		{ID: 1, Host: "localhost", Port: 9092},
	}

	mockedBrokers := []kafka.BrokerMetadata{
		{ID: 1, Host: "localhost", Port: 9092},
	}
	mockedMetadata := &kafka.Metadata{
		Brokers: mockedBrokers,
	}

	mockAdmin.On("GetMetadata", mock.Anything, true, mock.Anything).Return(mockedMetadata, nil)

	brokers, err := kafAdmin.GetClusterDetails()
	assert.NoError(t, err)
	assert.Equal(t, expectedBrokers, brokers)

	mockAdmin.AssertExpectations(t)
}

func TestCreateTopic(t *testing.T) {

	topic := "Aatest101"

	kafAdmin, mockAdmin, resetAdmin, err := setup(t)
	assert.NoError(t, err)
	defer resetAdmin()

	topRes := []kafka.TopicResult{{Topic: topic, Error: kafka.Error{}}}
	mockAdmin.On("CreateTopics", mock.Anything, mock.Anything, mock.Anything).Return(topRes, nil)

	topicDesc := []kafka.TopicDescription{{topic, kafka.UUID{}, kafka.Error{}, false, nil, nil}}
	descTopic := kafka.DescribeTopicsResult{topicDesc}
	mockAdmin.On("DescribeTopics", mock.Anything, mock.Anything, mock.Anything).Return(descTopic, nil)

	err = kafAdmin.CreateTopic(topic, 2, 1)
	assert.NoError(t, err)

	res, err := kafAdmin.DescribeTopic(topic)
	assert.NoError(t, err)
	assert.Equal(t, topic, res.TopicDescriptions[0].Name)

}

func Test_DeleteTopic(t *testing.T) {

	topic := "Aatest123"

	kafAdmin, mockAdmin, resetAdmin, err := setup(t)
	assert.NoError(t, err)

	defer resetAdmin()

	createTopicRes := []kafka.TopicResult{{Topic: topic, Error: kafka.Error{}}}
	mockAdmin.On("CreateTopics", mock.Anything, mock.Anything, mock.Anything).Return(createTopicRes, nil)

	delTopicRes := []kafka.TopicResult{{Topic: topic, Error: kafka.Error{}}}
	mockAdmin.On("DeleteTopics", mock.Anything, mock.Anything, mock.Anything).Return(delTopicRes, nil)

	err = kafAdmin.CreateTopic(topic, 2, 1)
	assert.NoError(t, err)

	err = kafAdmin.DeleteTopic(topic)
	assert.NoError(t, err)

}

func Test_ListOffsets(t *testing.T) {

	kafAdmin, mockAdmin, resetAdmin, err := setup(t)
	assert.NoError(t, err)

	defer resetAdmin()

	listOffsetRes := kafka.ListOffsetsResult{}
	mockAdmin.On("ListOffsets", mock.Anything, mock.Anything, mock.Anything).Return(listOffsetRes, nil)

	topic := "Aatest11"
	partition := 0
	err = kafAdmin.GetListOffsets(topic, partition)
	assert.NoError(t, err)
}
