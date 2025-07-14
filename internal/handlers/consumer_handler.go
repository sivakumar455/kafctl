package handlers

import (
	"fmt"
	"html/template"
	"kafctl/internal/logger"
	"kafctl/internal/services"
	"log"
	"log/slog"
	"net/http"
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
	msg, err := consumer.GetLatestRecords(topicName, 20)
	//msg, err := consumer.GetMessagesInfo(topicName)
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
