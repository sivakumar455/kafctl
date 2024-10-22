package handlers

import (
	"net/http"
)

type Application struct{}

func (app *Application) Routes() (http.Handler, error) {

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

	// Register pprof handlers
	// mux.HandleFunc("/debug/pprof/", http.HandlerFunc(pprof.Index))
	// mux.HandleFunc("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	// mux.HandleFunc("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	// mux.HandleFunc("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	// mux.HandleFunc("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))

	return mux, nil

}
