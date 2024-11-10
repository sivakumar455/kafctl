package handlers

import (
	"kafctl/internal/services"
	"net/http"
)

type Application struct{}

func (app *Application) Routes(admin services.IKafAdmin) (http.Handler, error) {

	mux := http.NewServeMux()

	handlers := NewKafkaHandlers(admin)

	consumerHandler := NewKafConsumerHandlers()

	fileServer := http.FileServer(http.Dir("./web/ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", handlers.home)
	mux.HandleFunc("/data", handlers.dataHandler)

	mux.HandleFunc("/topics", handlers.GetTopicsHandler)
	mux.HandleFunc("/createtopicform", handlers.createTopicFormHandler)
	mux.HandleFunc("/createtopic", handlers.createTopicHandler)

	mux.HandleFunc("/topic-details", handlers.describeTopicHandler)
	mux.HandleFunc("/delete-topic/", handlers.deleteTopicHandler)

	mux.HandleFunc("/view-messages", consumerHandler.GetMessages)

	mux.HandleFunc("/publishform", publishForm)
	mux.HandleFunc("/publishpayload", publishPayload)

	// Register pprof handlers
	// mux.HandleFunc("/debug/pprof/", http.HandlerFunc(pprof.Index))
	// mux.HandleFunc("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	// mux.HandleFunc("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	// mux.HandleFunc("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	// mux.HandleFunc("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))

	return mux, nil

}
