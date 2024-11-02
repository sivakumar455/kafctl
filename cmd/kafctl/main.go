/*
Copyright © 2024 Siva Kumar <EMAIL ADDRESS>
*/
package main

import (
	"flag"
	"fmt"
	"kafctl/internal/config"
	"kafctl/internal/handlers"
	"kafctl/internal/logger"
	"kafctl/internal/services"
	"net/http"
	"os"
)

func main() {
	// Define flags
	var topic, kafkaBroker, outputFile, groupId, configFile string
	var enableSSL, view bool
	flag.StringVar(&topic, "topic", "", "Kafka topic to consume from (mandatory)")
	flag.StringVar(&topic, "t", "", "Kafka topic to consume from (mandatory, shorthand)")

	flag.StringVar(&kafkaBroker, "kafkaBroker", "", "Kafka broker address (mandatory)")
	flag.StringVar(&kafkaBroker, "b", "", "Kafka broker address (mandatory, shorthand)")

	flag.StringVar(&outputFile, "outputFile", "", "File to write output to")
	flag.StringVar(&outputFile, "o", "", "File to write output to (shorthand)")

	flag.StringVar(&groupId, "groupId", "", "Kafka consumer group ID")
	flag.StringVar(&groupId, "g", "", "Kafka consumer group ID (shorthand)")

	flag.BoolVar(&enableSSL, "enableSSL", false, "Enable SSL configuration")
	flag.BoolVar(&enableSSL, "s", false, "Enable SSL configuration (shorthand)")

	flag.StringVar(&configFile, "config", "", "Path to SSL configuration file")
	flag.StringVar(&configFile, "f", "", "Path to SSL configuration file (shorthand)")

	flag.BoolVar(&view, "view", false, "dash oard for kafctl")
	flag.BoolVar(&view, "v", false, "dash oard for kafctl")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	// Parse flagss
	flag.Parse()

	// initializing with configs
	err := config.InitConfig(kafkaBroker, groupId, configFile, topic, outputFile, enableSSL, view)
	if err != nil {
		logger.Error("Error initializing kafka config", "error", err)
		os.Exit(1)
	}

	// running kafView if enabled
	if config.KafView {

		app := handlers.Application{}
		admin, err := services.NewKafAdmin()
		if err != nil {
			logger.Error("Error initializing kafka admin", "error", err)
			os.Exit(1)
		}
		defer admin.Close()

		mux, err := app.Routes(admin)
		if err != nil {
			logger.Error("Error creating routes", "error", err)
			os.Exit(1)
		}

		logger.Info("Running kafView dashboard on: ", "url", "http://"+config.KafViewUrl)
		err = http.ListenAndServe(config.KafViewUrl, mux)
		if err != nil {
			logger.Error("Error opening kafView", "error", err)
			os.Exit(1)
		}
	}
}
