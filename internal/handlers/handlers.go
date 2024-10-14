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

type Data struct {
	Content string `json:"content"`
}

func countPartitions(partitions []models.PartitionMetadata) int {
	return len(partitions)
}

func countIsrs(partitions []models.PartitionMetadata) int {
	total := 0
	for _, partition := range partitions {
		total += len(partition.Isrs)
	}
	return total
}

func countReplicas(partitions []models.PartitionMetadata) int {
	total := 0
	for _, partition := range partitions {
		total += len(partition.Replicas)
	}
	return total
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	data := Data{Content: "Page Last updated at " + time.Now().Format(time.RFC3339)}
	json.NewEncoder(w).Encode(data)
}

func GetTopicsHandler(w http.ResponseWriter, r *http.Request) {

	brokerInfo := models.BrokerList{}
	brokerInfo.Status = "Kafka is Up and Running"

	topics, err := services.GetAllTopics()
	if err != nil {
		logger.Error("Err getting topics: ", "error", err)
	} else {
		brokerInfo.Topics = *topics
	}

	files := []string{"./web/ui/base.html",
		"./web/ui/home.html"}

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

func home(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	brokerInfo := models.BrokerList{}
	brokerInfo.Status = "Kafka is Up and Running"

	brokers, err := services.GetClusterDetails()
	if err != nil {
		logger.Error("Error getting cluster details: ", "error", err)
		//http.Error(w, "Error getting cluster details", 500)
		brokerInfo.Status = "Kafka is down"
	} else {
		brokerInfo.Brokers = brokers
	}

	topics, err := services.GetAllTopics()
	if err != nil {
		logger.Error("Err getting topics: ", "error", err)
	} else {
		brokerInfo.Topics = *topics
	}

	//logger.Info(fmt.Sprintf("%v", topics))

	files := []string{"./web/ui/base.html",
		"./web/ui/home.html"}

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

func HanldlerInit() {

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./web/ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/data", dataHandler)

	mux.HandleFunc("/topics", GetTopicsHandler)
	mux.HandleFunc("/createtopicform", createTopicFormHandler)
	mux.HandleFunc("/createtopic", createTopicHandler)

	mux.HandleFunc("/topic-details", describeTopicHandler)
	mux.HandleFunc("/delete-topic/", deleteTopicHandler)

	logger.Info("Starting server on port 8989")

	err := http.ListenAndServe("localhost:8989", mux)
	if err != nil {
		logger.Error("Error", "error", err)
		return
	}

}
