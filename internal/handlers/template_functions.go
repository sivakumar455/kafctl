package handlers

import "github.com/confluentinc/confluent-kafka-go/v2/kafka"

func countPartitions(partitions []kafka.PartitionMetadata) int {
	return len(partitions)
}

func countIsrs(partitions []kafka.PartitionMetadata) int {
	total := 0
	for _, partition := range partitions {
		total += len(partition.Isrs)
	}
	return total
}

func countReplicas(partitions []kafka.PartitionMetadata) int {
	total := 0
	for _, partition := range partitions {
		total += len(partition.Replicas)
	}
	return total
}

func incrementer() func() int {
	i := 0
	return func() int {
		i++
		return i
	}
}
