package services

import (
	"fmt"
	"kafctl/internal/config"
	"kafctl/internal/logger"
	"log"
	"os"
	"sort"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const (
	// Timeout for individual Kafka operations like GetMetadata, QueryWatermarkOffsets
	metadataTimeoutMs = 5000
	// Timeout for each poll call
	pollTimeoutMs = 1000
	// Max duration to wait for all messages
	overallOperationTimeout = 20 * time.Second
)

type MessageRecord struct {
	Topic     string
	Partition int32
	Offset    int64
	Timestamp time.Time
	Key       string
	Value     string
}

type IConsumer interface {
	GetOffset(topic string, partition int, timeoutMs int) error
	GetTopicOffsets(topic *string) error
	ConsumeMessages(topic string) error
	ConsumeMessage(topic string) ([]*kafka.Message, error)
	GetLatestRecords(topic string, countPerPartition int) ([]*kafka.Message, error)
	ConsumeMessagesInFile() error
	Close() error
}

type Consumer struct {
	consumer *kafka.Consumer
}

func CreateConsumer() (*kafka.Consumer, error) {

	// Create a new Kafka consumer with SSL configuration
	consumerCfg, err := CreateConsumerConfig()
	if err != nil {
		return nil, err
	}
	c, err := kafka.NewConsumer(consumerCfg)
	if err != nil {
		logger.Error("Error creating consumer ", "error", err)
		return nil, err
	}
	return c, nil
}

func NewConsumer() (IConsumer, error) {
	// Create a new Kafka consumer
	consumer, err := CreateConsumer()
	if err != nil {
		return &Consumer{}, err
	}

	return &Consumer{consumer: consumer}, nil
}

func (c *Consumer) GetTopicOffsets(topic *string) error {

	metaData, err := c.consumer.GetMetadata(topic, false, 1000)
	if err != nil {
		logger.Error("Error getting consumer info", "error", err)
		return err
	}

	for _, partition := range metaData.Topics[*topic].Partitions {

		low, high, err := c.consumer.QueryWatermarkOffsets(*topic, partition.ID, 1000)
		if err != nil {
			logger.Error("Error getting offsets", "error", err)
			return err
		}

		logger.Info("Offset", "Id", partition.ID, "low", low, "high", high)
	}

	return nil
}

func (c *Consumer) GetOffset(topic string, partition int, timeoutMs int) error {

	// Query the watermark offsets
	low, high, err := c.consumer.QueryWatermarkOffsets(topic, int32(partition), 1000)
	if err != nil {
		logger.Error("Error getting offsets", "error", err)
		return err
	}

	logger.Info("First offset ", "firstOffset", low)
	logger.Info("Last offset", "lastOffset", high)
	return nil
}

func (c *Consumer) ConsumeMessage(topic string) ([]*kafka.Message, error) {

	logger.Info("Consuming from topic", "topic", topic)

	// Subscribe to the topic
	err := c.consumer.Subscribe(topic, nil)
	if err != nil {
		logger.Error("Error subscribing topic", "error", err)
		return nil, err
	}

	// Get the list of partitions for the topic
	partitions, err := c.consumer.Assignment()
	if err != nil {
		logger.Error("Failed to get partitions", "error", err)
	}

	// Seek to the beginning of each partition
	for _, partition := range partitions {
		err = c.consumer.Seek(kafka.TopicPartition{Topic: &topic, Partition: partition.Partition, Offset: kafka.OffsetBeginning}, -1)
		if err != nil {
			logger.Error("Failed to get partitions", "error", err)
		}
	}
	var msgRes []*kafka.Message
	// Consume messages

	count := 0
	for count < 4 {
		msg, err := c.consumer.ReadMessage(5000 * time.Millisecond)
		if err == nil {
			msgRes = append(msgRes, msg)
			headers := ""
			for _, header := range msg.Headers {
				headers += fmt.Sprintf("%s: %s, ", header.Key, string(header.Value))
			}
			logger.Info(fmt.Sprintf("Offset=%d, Key=%s, TimeStamp=%s",
				msg.TopicPartition.Offset, string(msg.Key), msg.Timestamp))

			logger.Info(fmt.Sprintf("Offset=%d, Key=%s, Headers=%s,Message=%s",
				msg.TopicPartition.Offset, string(msg.Key), headers, string(msg.Value)))
			count = count + 1

		} else if err.(kafka.Error).Code() == kafka.ErrTimedOut {
			// Timeout error, no more messages to read
			logger.Info("No more messages to read, exiting.")
			break
		} else {
			logger.Error("Consumer error", "error", err, "msg", msg)
			return nil, err
		}
	}
	return msgRes, nil
}

// GetLatestRecords fetches the latest 'count' records from each partition of the topic.
func (c *Consumer) GetLatestRecords(topic string, countPerPartition int) ([]*kafka.Message, error) {

	consumer := c.consumer
	defer consumer.Close()

	log.Printf("Consumer created for topic %s with countPerPartition %d", topic, countPerPartition)

	// 1. Determine partitions for the topic
	metadata, err := consumer.GetMetadata(&topic, false, metadataTimeoutMs)
	if err != nil {
		return nil, fmt.Errorf("failed to get metadata for topic %s: %w", topic, err)
	}

	topicMetadata, ok := metadata.Topics[topic]
	// 1. Check if the topic was found in the map.
	if !ok {
		// If the topic is not found, there's no topicMetadata.Error to check.
		// topicMetadata would be its zero value here, and its Error field would be kafka.ErrNoError (0).
		return nil, fmt.Errorf("topic '%s' not found in metadata", topic)
	}

	// 2. If the topic was found, check its Error field.
	//    kafka.ErrNoError (which has a value of 0) indicates no error.
	//    Any other value indicates an error.
	if topicMetadata.Error.Code() != kafka.ErrNoError {
		// Create a kafka.Error object from the ErrorCode for better error reporting.
		// The ErrorCode itself is just an integer. kafka.NewError provides the string representation.
		kafkaErr := kafka.NewError(topicMetadata.Error.Code(), "", false) // Description can be empty or more specific if available
		return nil, fmt.Errorf("error in metadata for topic '%s': %w", topic, kafkaErr)
	}

	if len(topicMetadata.Partitions) == 0 {
		log.Printf("Topic %s has no partitions.", topic)
		return nil, nil // Or an error, depending on desired behavior
	}

	log.Printf("Topic %s has %d partitions.", topic, len(topicMetadata.Partitions))

	// 2. For each partition, determine end offsets and seek positions
	assignments := make([]kafka.TopicPartition, 0, len(topicMetadata.Partitions))
	originalHighWatermarks := make(map[int32]int64)
	expectedMessagesPerPartition := make(map[int32]int)
	// New: Track partitions for which we've received an EOF
	partitionsExhausted := make(map[int32]bool)

	for _, pMeta := range topicMetadata.Partitions {
		partitionID := pMeta.ID
		if pMeta.Error.Code() != kafka.ErrNoError {
			log.Printf("Skipping partition %d for topic %s due to metadata error: %v", partitionID, topic, kafka.NewError(pMeta.Error.Code(), "", false))
			continue
		}

		lowWm, highWm, err := consumer.QueryWatermarkOffsets(topic, partitionID, metadataTimeoutMs)
		if err != nil {
			// Log and continue, or return error. For robustness, let's try to process other partitions.
			log.Printf("Warning: failed to query watermark offsets for %s [%d]: %v. Skipping this partition.", topic, partitionID, err)
			continue
		}
		originalHighWatermarks[partitionID] = highWm // Still useful for sanity checks if needed
		log.Printf("Partition %d: LowWM: %d, HighWM: %d", partitionID, lowWm, highWm)

		// We want to read 'countPerPartition' messages ending at 'highWm-1'.
		// So, the ideal start would be 'highWm - countPerPartition'.
		targetStartOffset := highWm - int64(countPerPartition)

		// Actual seek offset must be >= lowWm and >= 0.
		effectiveSeekOffset := targetStartOffset
		if effectiveSeekOffset < lowWm {
			log.Printf("Partition %d: Target start offset %d is less than LowWM %d. Adjusting seek to LowWM.", partitionID, targetStartOffset, lowWm)
			effectiveSeekOffset = lowWm
		}
		if effectiveSeekOffset < 0 { // Should mostly be covered by lowWm, but as a safeguard
			effectiveSeekOffset = 0
		}

		// Number of messages we actually expect to fetch from this partition.
		numExpectedThisPartition := int64(0)
		if highWm > effectiveSeekOffset { // Only if there's a potential range to read
			numExpectedThisPartition = highWm - effectiveSeekOffset
		} else {
			log.Printf("Partition %d: HighWM (%d) is not greater than effectiveSeekOffset (%d). Expecting 0 messages.", partitionID, highWm, effectiveSeekOffset)
		}

		log.Printf("Partition %d: Original HighWM: %d, Calculated Effective Seek Offset: %d, Expected Messages to Fetch: %d",
			partitionID, highWm, effectiveSeekOffset, numExpectedThisPartition)

		if numExpectedThisPartition > 0 {
			assignments = append(assignments, kafka.TopicPartition{
				Topic:     &topic,
				Partition: partitionID,
				Offset:    kafka.Offset(effectiveSeekOffset),
			})
			// Store the *actual number of messages we can get from this range*,
			// capped by countPerPartition if the range was larger (though our seek tries to limit this).
			// If effectiveSeekOffset was adjusted up by lowWm, numExpectedThisPartition might be < countPerPartition.
			// We should only expect what's available.
			expectedMessagesPerPartition[partitionID] = int(numExpectedThisPartition)
			partitionsExhausted[partitionID] = false // Initialize
		} else {
			log.Printf("Partition %d: Not assigning as no messages are expected to be fetched.", partitionID)
		}
	}
	if len(assignments) == 0 {
		log.Println("No partitions to assign.")
		return []*kafka.Message{}, nil
	}

	// 3. Assign partitions to the consumer
	err = consumer.Assign(assignments)
	if err != nil {
		return nil, fmt.Errorf("failed to assign partitions: %w", err)
	}
	log.Printf("Assigned partitions: %+v", assignments)

	// 4. Poll for messages
	var collectedRawRecords []*kafka.Message // Store all messages before final sorting
	collectedPerPartitionCount := make(map[int32]int)
	for _, assign := range assignments {
		collectedPerPartitionCount[assign.Partition] = 0
	}

	allDone := false
	startTime := time.Now()

	// Check if initially allDone (e.g., no assignments were made)
	if len(expectedMessagesPerPartition) == 0 {
		allDone = true
	}

	for !allDone && time.Since(startTime) < overallOperationTimeout {
		ev := consumer.Poll(pollTimeoutMs)
		if ev == nil {
			log.Println("Poll timed out...")
			// Continue to check allDone condition
		}

		switch e := ev.(type) {
		case *kafka.Message:
			log.Printf("Received message from %s [%d] @ %d", *e.TopicPartition.Topic, e.TopicPartition.Partition, e.TopicPartition.Offset)
			partitionID := e.TopicPartition.Partition

			// Check if we still *expect* messages from this partition based on our calculated expected count
			// and if we haven't already marked it as exhausted.
			if expectedCount, ok := expectedMessagesPerPartition[partitionID]; ok {
				if collectedPerPartitionCount[partitionID] < expectedCount {
					// Further check: ensure message offset is within the originally intended range.
					// The effectiveSeekOffset was our start, originalHighWatermarks[partitionID] was our end.
					// This check helps if auto.offset.reset jumped us outside our desired window.
					// However, librdkafka usually handles the "offset out of range" by not delivering the message.
					// The primary gate is `collectedPerPartitionCount[partitionID] < expectedCount`.
					if int64(e.TopicPartition.Offset) < originalHighWatermarks[partitionID] { // Sanity check
						collectedRawRecords = append(collectedRawRecords, e)
						collectedPerPartitionCount[partitionID]++
					} else {
						log.Printf("Partition %d: Skipped message at offset %d as it's >= originalHighWatermark %d.",
							partitionID, e.TopicPartition.Offset, originalHighWatermarks[partitionID])
					}
				} else {
					// Already collected expected for this partition, but got more (should be rare with Assign & EOF)
					log.Printf("Partition %d: Received message at offset %d but already collected expected %d. Ignoring.",
						partitionID, e.TopicPartition.Offset, expectedCount)
				}
			} else {
				// Message from a partition we didn't expect to get messages from (should not happen with Assign)
				log.Printf("Warning: Received message from un-expected partition %d. Ignoring.", partitionID)
			}

		case kafka.PartitionEOF:
			log.Printf("Reached EOF for %s [%d] at offset %v\n",
				*e.Topic, e.Partition, e.Offset)
			// Mark this partition as exhausted for the current fetch attempt.
			// It means we won't get more messages from its current assigned range.
			partitionsExhausted[e.Partition] = true

		case kafka.Error:
			if e.IsFatal() {
				return nil, fmt.Errorf("fatal consumer error: %w", e)
			}
			log.Printf("Non-fatal consumer error: %v", e)
			// If error indicates partition is problematic, could mark as exhausted too, e.g.
			// if e.Code() == kafka.ErrUnknownTopicOrPart || e.Code() == kafka.ErrBrokerNotAvailable {
			//  if tp := e.TopicPartition(); tp != nil {
			//      partitionsExhausted[tp.Partition] = true
			//  }
			//}

		default:
			// log.Printf("Ignored event: %v\n", e)
		}

		// Check if all partitions are done
		allDone = true
		if len(expectedMessagesPerPartition) == 0 { // If no partitions were expected to yield messages
			allDone = true
		} else {
			for pID, expectedCount := range expectedMessagesPerPartition {
				// A partition is "pending" if we haven't collected enough AND we haven't hit EOF for it.
				if collectedPerPartitionCount[pID] < expectedCount && !partitionsExhausted[pID] {
					allDone = false
					break
				}
			}
		}
	}

	if time.Since(startTime) >= overallOperationTimeout && !allDone {
		log.Printf("Warning: Operation timed out after %s. Collected %d raw records.", overallOperationTimeout, len(collectedRawRecords))
		log.Printf("Current collected counts: %+v", collectedPerPartitionCount)
		log.Printf("Expected counts: %+v", expectedMessagesPerPartition)
		log.Printf("Exhausted partitions: %+v", partitionsExhausted)
	}

	log.Printf("Finished polling. Total raw records collected: %d", len(collectedRawRecords))

	sort.SliceStable(collectedRawRecords, func(i, j int) bool {
		if collectedRawRecords[i].TopicPartition.Partition != collectedRawRecords[j].TopicPartition.Partition {
			return collectedRawRecords[i].TopicPartition.Partition > collectedRawRecords[j].TopicPartition.Partition
		}
		return collectedRawRecords[i].TopicPartition.Offset > collectedRawRecords[j].TopicPartition.Offset
	})

	finalRecords := make([]MessageRecord, 0, len(collectedRawRecords))
	for _, rec := range collectedRawRecords {
		finalRecords = append(finalRecords, MessageRecord{
			Topic:     *rec.TopicPartition.Topic,
			Partition: rec.TopicPartition.Partition,
			Offset:    int64(rec.TopicPartition.Offset),
			Timestamp: rec.Timestamp,
			Key:       string(rec.Key),   // Assuming UTF-8 string
			Value:     string(rec.Value), // Assuming UTF-8 string
		})
	}

	log.Printf("Returning %d deserialized records.", len(finalRecords))
	return collectedRawRecords, nil

}

func (c *Consumer) ConsumeMessages(topic string) error {

	logger.Info("Consuming from topic", "topic", topic)

	// Subscribe to the topic
	err := c.consumer.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		logger.Error("Error subscribing topic", "error", err)
		return err
	}
	// Consume messages
	for {
		msg, err := c.consumer.ReadMessage(1000 * time.Millisecond)
		if err == nil {
			headers := ""
			for _, header := range msg.Headers {
				headers += fmt.Sprintf("%s: %s, ", header.Key, string(header.Value))
			}
			logger.Info(fmt.Sprintf("Offset=%d, Key=%s, TimeStamp=%s",
				msg.TopicPartition.Offset, string(msg.Key), msg.Timestamp))

			logger.Info(fmt.Sprintf("Offset=%d, Key=%s, Headers=%s,Message=%s",
				msg.TopicPartition.Offset, string(msg.Key), headers, string(msg.Value)))

		} else if err.(kafka.Error).Code() == kafka.ErrTimedOut {
			// Timeout error, no more messages to read
			logger.Info("No more messages to read, exiting.")
			break
		} else {
			logger.Error("Consumer error", "error", err, "msg", msg)
			return err
		}
	}
	return nil
}

func (c *Consumer) ConsumeMessagesInFile() error {

	// Open a file to write the messages
	file, err := os.Create(config.OutputFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Subscribe to the topic
	c.consumer.SubscribeTopics([]string{config.Topic}, nil)

	fmt.Printf("Consuming from topic: %s, output file: %s\n", config.Topic, config.OutputFile)

	// Consume messages
	for {
		msg, err := c.consumer.ReadMessage(-1)
		if err == nil {
			// Write message details to the file
			headers := ""
			for _, header := range msg.Headers {
				headers += fmt.Sprintf("%s: %s, ", header.Key, string(header.Value))
			}
			fmt.Printf("Offset=%d, Key=%s, TimeStamp=%s \n",
				msg.TopicPartition.Offset, string(msg.Key), msg.Timestamp)

			_, err := file.WriteString(fmt.Sprintf("Offset=%d, Key=%s, Headers=%s,\nMessage=%s \n\n",
				msg.TopicPartition.Offset, string(msg.Key), headers, string(msg.Value)))
			if err != nil {
				panic(err)
			}
		} else {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}

}

func (c *Consumer) Close() error {
	err := c.consumer.Close()
	if err != nil {
		return err
	}

	return nil
}
