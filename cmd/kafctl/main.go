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
	"os"
)

var (
	topic, kafkaBroker, outputFile, groupId, configFile string
	enableSSL                                           bool
	view                                                bool
)

func main() {
	// Define flags
	//var topic, kafkaBroker, outputFile, groupId, configFile string
	//var enableSSL bool
	flag.StringVar(&topic, "topic", "", "Kafka topic to consume from (mandatory)")
	flag.StringVar(&topic, "t", "", "Kafka topic to consume from (mandatory, shorthand)")

	flag.StringVar(&kafkaBroker, "kafkaBroker", "", "Kafka broker address (mandatory)")
	flag.StringVar(&kafkaBroker, "b", "", "Kafka broker address (mandatory, shorthand)")

	flag.StringVar(&outputFile, "outputFile", "consumer_out.txt", "File to write output to")
	flag.StringVar(&outputFile, "o", "consumer_out.txt", "File to write output to (shorthand)")

	flag.StringVar(&groupId, "groupId", "x-consumer-test", "Kafka consumer group ID")
	flag.StringVar(&groupId, "g", "x-consumer-test", "Kafka consumer group ID (shorthand)")

	flag.BoolVar(&enableSSL, "enableSSL", false, "Enable SSL configuration")
	flag.BoolVar(&enableSSL, "s", false, "Enable SSL configuration (shorthand)")

	flag.StringVar(&configFile, "config", "ssl_properties.json", "Path to SSL configuration file")
	flag.StringVar(&configFile, "f", "ssl_properties.json", "Path to SSL configuration file (shorthand)")

	flag.BoolVar(&view, "view", false, "dash oard for kafctl")
	flag.BoolVar(&view, "v", false, "dash oard for kafctl")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	// Parse flagss
	flag.Parse()

	if kafkaBroker == "" {
		logger.Error("Error: topic and kafkaBroker flags are mandatory")
		flag.Usage()
		os.Exit(1)
	}

	config.InitConfig(kafkaBroker, groupId, configFile, topic, outputFile, enableSSL)
	//err := services.ConsumeMessagesInFile(kafkaBroker, groupId, configFile, topic, outputFile, enableSSL)

	if view {
		logger.Info("Running kafctl dashboard on localhost:8989")
		handlers.HanldlerInit()
	}

}
