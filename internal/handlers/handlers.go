package handlers

import (
	"encoding/json"
	"html/template"
	"kafctl/internal/logger"
	"kafctl/internal/models"
	"kafctl/internal/services"
	"log"
	"net/http"
	"time"
)

const BASE_TEMPL_PATH string = "./web/ui/base.html"
const HOME_TEMPL_PATH string = "./web/ui/home.html"
const TOPIC_FORM_TEMPL_PATH string = "./web/ui/topicform.html"
const TOPIC_DETAILS_TEMPL_PATH string = "./web/ui/topicdetails.html"

type Data struct {
	Content string `json:"content"`
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	data := Data{Content: "Page Last updated at " + time.Now().Format(time.RFC3339)}
	json.NewEncoder(w).Encode(data)
}

func home(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	brokerInfo := models.BrokerInfo{}
	brokerInfo.Status = "UP"
	kafAdmin := services.KafAdmin{}
	brokers, err := kafAdmin.GetClusterDetails()
	if err != nil {
		logger.Error("Error getting cluster details: ", "error", err)
		brokerInfo.Status = "DOWN"
	} else {
		brokerInfo.Brokers = brokers
	}
	topics, err := kafAdmin.GetAllTopics()
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

	tmpl := template.Must(template.New("base").Funcs(funcMap).ParseFiles(files...))

	err = tmpl.ExecuteTemplate(w, "base", brokerInfo)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", 500)
		return
	}
}
