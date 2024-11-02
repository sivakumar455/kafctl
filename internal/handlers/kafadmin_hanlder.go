package handlers

import (
	"fmt"
	"html/template"
	"kafctl/internal/logger"
	"kafctl/internal/models"
	"log"
	"net/http"
	"strings"
)

func (kah *KafAdminHandlers) createTopicFormHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{HOME_TEMPL_PATH, TOPIC_FORM_TEMPL_PATH}

	funcMap := template.FuncMap{
		"countPartitions": countPartitions,
		"countReplicas":   countReplicas,
		"countIsrs":       countIsrs}

	tmpl := template.Must(template.New("createform").Funcs(funcMap).ParseFiles(files...))

	err := tmpl.ExecuteTemplate(w, "createTopicForm", nil)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", 500)
		return
	}

}

func (kah *KafAdminHandlers) createTopicHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		topicName := r.FormValue("topicName")
		numPartitions := 1
		numReplicas := 1

		err := kah.kafAdmin.CreateTopic(topicName, numPartitions, numReplicas)
		if err != nil {
			fmt.Fprintf(w, "%v", "Error creating topic")
			return
		}

		fmt.Fprintf(w, "Topic '%s' created successfully!", topicName)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}

}

func (kah *KafAdminHandlers) deleteTopicHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodDelete {

		path := r.URL.Path
		parts := strings.Split(path, "/")
		if len(parts) < 3 {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}
		topicName := parts[2]
		err := kah.kafAdmin.DeleteTopic(topicName)
		if err != nil {
			fmt.Fprintf(w, "%v", "Error deleting topic")
			return
		}

		fmt.Fprintf(w, "Topic '%s' deleted successfully!", topicName)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}

}

func (kah *KafAdminHandlers) describeTopicHandler(w http.ResponseWriter, r *http.Request) {

	topicName := r.URL.Query().Get("name")
	if topicName == "" {
		http.Error(w, "Missing topic name", http.StatusBadRequest)
		return
	}
	topics, err := kah.kafAdmin.DescribeTopic(topicName)
	if err != nil {
		fmt.Fprintf(w, "%v", "Error creating topic")
		return
	}

	files := []string{HOME_TEMPL_PATH, TOPIC_DETAILS_TEMPL_PATH}

	funcMap := template.FuncMap{
		"countPartitions": countPartitions,
		"countReplicas":   countReplicas,
		"countIsrs":       countIsrs}

	tmpl := template.Must(template.New("topicdetail").Funcs(funcMap).ParseFiles(files...))

	err = tmpl.ExecuteTemplate(w, "describeTopic", topics.TopicDescriptions[0])
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", 500)
		return
	}

}

func (kah *KafAdminHandlers) GetTopicsHandler(w http.ResponseWriter, r *http.Request) {

	brokerInfo := models.BrokerInfo{}
	brokerInfo.Status = "Kafka is Up and Running"
	topics, err := kah.kafAdmin.GetAllTopics()
	if err != nil {
		logger.Error("Err getting topics: ", "error", err)
	} else {
		brokerInfo.Topics = topics
	}

	files := []string{BASE_TEMPL_PATH, HOME_TEMPL_PATH}

	funcMap := template.FuncMap{
		"countPartitions": countPartitions,
		"countReplicas":   countReplicas,
		"countIsrs":       countIsrs}

	tmpl := template.Must(template.New("topicsl").Funcs(funcMap).ParseFiles(files...))

	err = tmpl.ExecuteTemplate(w, "topics", brokerInfo)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", 500)
		return
	}
}
