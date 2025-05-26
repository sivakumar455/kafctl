package services

import (
	"encoding/json"
	"kafctl/internal/config"
	"kafctl/internal/logger"
	"os"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/uuid"
)

type SSLConfig struct {
	SecurityProtocol string `json:"securityProtocol"`
	SslCaLocation    string `json:"sslCaLocation"`
	SslCertLocation  string `json:"sslCertLocation"`
	SslKeyLocation   string `json:"sslKeyLocation"`
	SslKeyPassword   string `json:"sslKeyPassword"`
}

var NewAdminClient = kafka.NewAdminClient

func SetSSLConfig(cfgMap *kafka.ConfigMap, sslConfigFile string) error {
	// Read SSL configuration from file
	sslCfgFile, err := os.Open(sslConfigFile)
	if err != nil {
		return err
	}
	defer sslCfgFile.Close()

	var sslConfig SSLConfig
	if err := json.NewDecoder(sslCfgFile).Decode(&sslConfig); err != nil {
		return err
	}

	cfgMap.SetKey("security.protocol", sslConfig.SecurityProtocol)
	if sslConfig.SecurityProtocol == "SSL" || sslConfig.SecurityProtocol == "SASL_SSL" {
		cfgMap.SetKey("ssl.ca.location", sslConfig.SslCaLocation)
		cfgMap.SetKey("ssl.certificate.location", sslConfig.SslCertLocation)
		cfgMap.SetKey("ssl.key.location", sslConfig.SslKeyLocation)
		cfgMap.SetKey("ssl.key.password", sslConfig.SslKeyPassword)
	}
	return nil
}

func CreateConsumerConfig() (*kafka.ConfigMap, error) {
	consumerCfg := &kafka.ConfigMap{
		"bootstrap.servers": config.KafkaBroker,
		// "group.id":           config.GroupId,
		"enable.auto.commit": false,
		// "max.poll.records":   10,
		"auto.offset.reset":    "earliest",
		"group.id":             uuid.New().String(),
		"enable.partition.eof": "true", // To get EOF events, useful for knowing when we've read to the end
	}

	if config.EnableSSL {
		err := SetSSLConfig(consumerCfg, config.ConfigFile)
		if err != nil {
			return nil, err
		}
	}
	return consumerCfg, nil
}

func CreateProducerConfig() (*kafka.ConfigMap, error) {
	producerCfg := &kafka.ConfigMap{
		"bootstrap.servers": config.KafkaBroker,
		"acks":              "all",
	}

	if config.EnableSSL {
		err := SetSSLConfig(producerCfg, config.ConfigFile)
		if err != nil {
			return nil, err
		}
	}
	return producerCfg, nil
}

func NewAdminConfig() (*kafka.ConfigMap, error) {
	adminCfg := &kafka.ConfigMap{
		"bootstrap.servers": config.KafkaBroker,
		"debug":             "broker,protocol",
	}

	if config.EnableSSL {
		err := SetSSLConfig(adminCfg, config.ConfigFile)
		if err != nil {
			return nil, err
		}
	}
	return adminCfg, nil
}

func CreateAdminClient() (IRdAdminClient, error) {
	cfg, err := NewAdminConfig()
	if err != nil {
		logger.Error("Failed to create admin Config", "error", err)
		return nil, err
	}

	admin, err := NewAdminClient(cfg)
	if err != nil {
		logger.Error("Failed to create admin Client", "error", err)
		return nil, err
	}
	return admin, nil
}
