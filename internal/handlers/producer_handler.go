package handlers

import (
	"fmt"
	"html/template"
	"kafctl/internal/logger"
	"kafctl/internal/services"
	"log"
	"net/http"
)

func publishForm(w http.ResponseWriter, r *http.Request) {

	files := []string{PUBLISH_FORM_TEMPL_PATH, BASE_TEMPL_PATH}

	funcMap := template.FuncMap{
		"countPartitions": countPartitions,
		"countReplicas":   countReplicas,
		"countIsrs":       countIsrs}

	tmpl := template.Must(template.New("publishform").Funcs(funcMap).ParseFiles(files...))

	// get topics
	admin, err := services.NewKafAdmin()
	if err != nil {
		logger.Error("Error initializing kafka admin", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	//defer admin.Close()

	topics, err := admin.GetAllTopics()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "publishform", topics)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func publishPayload(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		topicName := r.FormValue("topicName")
		payload := r.FormValue("payload")
		optionalHeaders := r.FormValue("optionalHeaders")

		logger.Info("Publish payload options", "topicName", topicName, "payload", payload, "optionalHeaders", optionalHeaders)

		err := services.ProduceMessage(topicName, []byte(payload), []byte(payload), optionalHeaders)
		if err != nil {
			fmt.Fprintf(w, "Error publishing to topic  %s", topicName)
			return
		}

		fmt.Fprintf(w, "Topic %s published successfully!", topicName)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}
