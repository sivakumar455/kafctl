package handlers

import (
	"fmt"
	"html/template"
	"kafctl/internal/logger"
	"kafctl/internal/services"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type IKafConsumerHandlers interface {
	ViewTopic(w http.ResponseWriter, r *http.Request)
	ViewMessages(w http.ResponseWriter, r *http.Request)
}

type KafConsumerHandlers struct {
	//kch services.IConsumer
}

// func setup() services.IConsumer {
// 	consumer, err := services.NewConsumer()
// 	if err != nil {
// 		logger.Error("Error creating consumer")
// 	}
// 	defer consumer.Close()

// 	return consumer
// }

func NewKafConsumerHandlers() IKafConsumerHandlers {
	return &KafConsumerHandlers{}
}

func (kch *KafConsumerHandlers) ViewTopic(w http.ResponseWriter, r *http.Request) {

	topicName := r.URL.Query().Get("topicname")
	if topicName == "" {
		http.Error(w, "Missing topic name", http.StatusBadRequest)
		return
	}

	slog.Info("Getting messages", "topic", topicName)

	consumer, err := services.NewConsumer()
	if err != nil {
		logger.Error("Error creating consumer")
	}
	defer consumer.Close()

	//var messages []*kafka.Message
	//msg, err := consumer.ConsumeMessage(topicName)
	//msg, err := consumer.GetLatestRecords(topicName, 20)
	msg, err := consumer.GetMessagesInfo(topicName)
	if err != nil {
		fmt.Fprintf(w, "%v", "Error viewing messages")
		return
	}

	logger.Info("Messages are fetched")

	dataTemplate := map[string]any{
		"Message":   msg,
		"TopicName": topicName,
	}

	funcMap := template.FuncMap{
		"countPartitions": countPartitions,
		"countReplicas":   countReplicas,
		"countIsrs":       countIsrs,
		"inc":             incrementer}

	files := []string{BASE_TEMPL_PATH, TOPIC_DETAILS_TEMPL_PATH, VIEW_TOPIC_TEMPL_PATH}

	tmpl := template.Must(template.New("topicViewer").Funcs(funcMap).ParseFiles(files...))
	logger.Info("Rendering messages template")
	err = tmpl.ExecuteTemplate(w, "topicViewer", dataTemplate)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", 500)
		return
	}
	logger.Info("Rendered messages template for viewTopic")

	//fmt.Fprintf(w, "%s", "hello messages")
}

func (kch *KafConsumerHandlers) ViewMessages(w http.ResponseWriter, r *http.Request) {

	topicName := r.URL.Query().Get("topicname")
	if r.Method == http.MethodPost || r.Method == http.MethodGet {
		// Try form values first (POST), then query params (GET)
		if topicName == "" {
			topicName = r.FormValue("topicname")
		}
	}
	if topicName == "" {
		http.Error(w, "Missing topic name", http.StatusBadRequest)
		return
	}

	// Get partition filter (optional) - check both form and query params
	partitionStr := r.FormValue("partition")
	if partitionStr == "" {
		partitionStr = r.URL.Query().Get("partition")
	}

	var selectedPartition *int32
	if partitionStr != "" && partitionStr != "all" {
		if parsed, err := strconv.Atoi(partitionStr); err == nil && parsed >= 0 {
			partID := int32(parsed)
			selectedPartition = &partID
		}
	}

	// Get count (messages per partition) - check both form and query params
	countStr := r.FormValue("numMessages")
	if countStr == "" {
		countStr = r.URL.Query().Get("numMessages")
	}

	countPerPartition := 20 // default
	if countStr != "" {
		if parsed, err := strconv.Atoi(countStr); err == nil && parsed > 0 && parsed <= 1000 {
			countPerPartition = parsed
		}
	}

	slog.Info("Getting messages", "topic", topicName, "partition", selectedPartition, "count", countPerPartition)

	consumer, err := services.NewConsumer()
	if err != nil {
		logger.Error("Error creating consumer")
		http.Error(w, "Error creating consumer", http.StatusInternalServerError)
		return
	}
	defer consumer.Close()

	// Get messages from all partitions (or filtered by partition)
	msg, err := consumer.GetLatestRecords(topicName, countPerPartition)
	if err != nil {
		fmt.Fprintf(w, "Error viewing messages: %v", err)
		return
	}

	// Filter by partition if specified
	if selectedPartition != nil {
		filtered := make([]*kafka.Message, 0)
		for _, m := range msg {
			if m != nil && m.TopicPartition.Partition == *selectedPartition {
				filtered = append(filtered, m)
			}
		}
		msg = filtered
		logger.Info("Filtered messages by partition", "partition", *selectedPartition, "count", len(msg))
	}

	logger.Info("Messages are fetched", "count", len(msg))

	dataTemplate := map[string]any{
		"Message":   msg,
		"TopicName": topicName,
	}

	funcMap := template.FuncMap{
		"countPartitions": countPartitions,
		"countReplicas":   countReplicas,
		"countIsrs":       countIsrs,
		"inc":             incrementer}

	files := []string{BASE_TEMPL_PATH, TOPIC_DETAILS_TEMPL_PATH, VIEW_TOPIC_TEMPL_PATH, MESSAGE_TEMPL_PATH}

	tmpl := template.Must(template.New("viewMessages").Funcs(funcMap).ParseFiles(files...))
	logger.Info("Rendering messages template")
	err = tmpl.ExecuteTemplate(w, "viewMessages", dataTemplate)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", 500)
		return
	}
	logger.Info("Rendered messages template for viewMessages")
}
