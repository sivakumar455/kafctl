/*
Copyright © 2024 Siva Kumar <EMAIL ADDRESS>
*/
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type SSLConfig struct {
	SecurityProtocol string `json:"securityProtocol"`
	SslCaLocation    string `json:"sslCaLocation"`
	SslCertLocation  string `json:"sslCertLocation"`
	SslKeyLocation   string `json:"sslKeyLocation"`
	SslKeyPassword   string `json:"sslKeyPassword"`
}

func main() {
	// Define flags
	var topic, kafkaBroker, outputFile, groupId, configFile string
	var enableSSL bool
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

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	// Parse flagss
	flag.Parse()

	if topic == "" || kafkaBroker == "" {
		fmt.Println("Error: topic and kafkaBroker flags are mandatory")
		flag.Usage()
		os.Exit(1)
	}

	// Create a new Kafka consumer with SSL configuration
	consumerCfg := &kafka.ConfigMap{
		"bootstrap.servers": kafkaBroker,
		"group.id":          groupId,
		"auto.offset.reset": "earliest",
	}

	// Read SSL configuration from file
	if enableSSL {

		sslCfgFile, err := os.Open(configFile)
		if err != nil {
			panic(err)
		}
		defer sslCfgFile.Close()

		var sslConfig SSLConfig
		if err := json.NewDecoder(sslCfgFile).Decode(&sslConfig); err != nil {
			panic(err)
		}

		consumerCfg.SetKey("security.protocol", sslConfig.SecurityProtocol)
		if sslConfig.SecurityProtocol == "SSL" || sslConfig.SecurityProtocol == "SASL_SSL" {
			consumerCfg.SetKey("ssl.ca.location", sslConfig.SslCaLocation)
			consumerCfg.SetKey("ssl.certificate.location", sslConfig.SslCertLocation)
			consumerCfg.SetKey("ssl.key.location", sslConfig.SslKeyLocation)
			consumerCfg.SetKey("ssl.key.password", sslConfig.SslKeyPassword)
		}
	}

	c, err := kafka.NewConsumer(consumerCfg)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Consuming from topic: %s, output file: %s\n", topic, outputFile)

	// Subscribe to the topic
	c.SubscribeTopics([]string{topic}, nil)

	// Open a file to write the messages
	file, err := os.Create(outputFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	defer c.Close()

	// Consume messages
	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			// Write message details to the file
			headers := ""
			for _, header := range msg.Headers {
				headers += fmt.Sprintf("%s: %s, ", header.Key, string(header.Value))
			}
			fmt.Printf("Offset=%d, Key=%s, TimeStamp=%s \n",
				msg.TopicPartition.Offset, string(msg.Key), msg.Timestamp)

			_, err := file.WriteString(fmt.Sprintf("Offset=%d, Key=%s, Headers=%s,\nMessage=%s \n\n",
				msg.TopicPartition.Offset, string(msg.Key), headers, string(msg.Value)))
			if err != nil {
				panic(err)
			}
		} else {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
}
