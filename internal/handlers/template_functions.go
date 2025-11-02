package handlers

import "github.com/confluentinc/confluent-kafka-go/v2/kafka"

func countPartitions(partitions []kafka.PartitionMetadata) int {
	return len(partitions)
}

func countIsrs(partitions []kafka.PartitionMetadata) int {
	if len(partitions) == 0 {
		return 0
	}
	// Return the ISR count for the first partition (all partitions should have same replication factor)
	// If partitions have different ISR counts, return the minimum (indicating some are out of sync)
	minIsrs := len(partitions[0].Isrs)
	for _, partition := range partitions {
		if len(partition.Isrs) < minIsrs {
			minIsrs = len(partition.Isrs)
		}
	}
	return minIsrs
}

func countReplicas(partitions []kafka.PartitionMetadata) int {
	if len(partitions) == 0 {
		return 0
	}
	// Return the replication factor (replicas per partition)
	// All partitions in a topic should have the same replication factor
	return len(partitions[0].Replicas)
}

func incrementer() func() int {
	i := 0
	return func() int {
		i++
		return i
	}
}
