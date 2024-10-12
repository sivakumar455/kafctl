package models

type PartitionMetadata struct {
	ID       int32
	Leader   int32
	Replicas []int32
	Isrs     []int32
}

// TopicMetadata contains per-topic metadata
type TopicMetadata struct {
	Topic      string
	Partitions []PartitionMetadata
}

type ACLOperation int

type UUID struct {
	// Most Significant Bits.
	mostSignificantBits int64
	// Least Significant Bits.
	leastSignificantBits int64
	// Base64 representation
	base64str string
}

type TopicDetails struct {
	Name                 string
	TopicID              UUID
	AuthorizedOperations []ACLOperation
	Partitions           []PartitionDetails
}

type PartitionDetails struct {
	PartitionID int
	Leader      string
	Isr         []string
	Replicas    []string
}

type TopicsMap struct {
	TopicsMap map[string]TopicMetadata
}

type Broker struct {
	Id   int32
	Host string
	Port int
}

type BrokerList struct {
	Brokers []Broker
	Status  string
	Topics  TopicsMap
}
