package handlers

import (
	"fmt"
	"html/template"
	"kafctl/internal/services"
	"log"
	"net/http"
	"strings"
)

func createTopicFormHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{"./web/ui/home.html", "./web/ui/createTopicForm.html"}

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

func createTopicHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		topicName := r.FormValue("topicName")
		// Process the topic name (e.g., save to database)
		numPartitions := 1
		numReplicas := 1

		err := services.CreateTopic(topicName, numPartitions, numReplicas)
		if err != nil {
			fmt.Fprintf(w, "%v", "Error creating topic")
			return
		}

		fmt.Fprintf(w, "Topic '%s' created successfully!", topicName)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}

}

func deleteTopicHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodDelete {

		path := r.URL.Path
		parts := strings.Split(path, "/")
		if len(parts) < 3 {
			http.Error(w, "Invalid URL", http.StatusBadRequest)
			return
		}
		topicName := parts[2]

		err := services.DeleteTopic(topicName)
		if err != nil {
			fmt.Fprintf(w, "%v", "Error deleting topic")
			return
		}

		fmt.Fprintf(w, "Topic '%s' deleted successfully!", topicName)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}

}

func describeTopicHandler(w http.ResponseWriter, r *http.Request) {
	// Handle topic creation logic here
	// For now, just redirect to home

	topicName := r.URL.Query().Get("name")
	if topicName == "" {
		http.Error(w, "Missing topic name", http.StatusBadRequest)
		return
	}

	topics, err := services.DescribeTopic(topicName)
	if err != nil {
		fmt.Fprintf(w, "%v", "Error creating topic")
		return
	}

	// fmt.Fprintf(w, "Topic '%s' created successfully!", topics)

	files := []string{"./web/ui/home.html", "./web/ui/topicdetails.html"}

	funcMap := template.FuncMap{
		"countPartitions": countPartitions,
		"countReplicas":   countReplicas,
		"countIsrs":       countIsrs}

	tmpl := template.Must(template.New("topicdetail").Funcs(funcMap).ParseFiles(files...))

	err = tmpl.ExecuteTemplate(w, "describeTopic", topics[0])
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", 500)
		return
	}

}
