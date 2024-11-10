package config

import (
	"encoding/json"
	"kafctl/internal/logger"
	"os"
)

var (
	Topic, KafkaBroker, OutputFile, GroupId, ConfigFile, KafViewUrl string
	EnableSSL, KafView                                              bool
)

type AppConfig struct {
	KafkaBroker   string
	Topic         string
	GroupId       string
	SslConfigFile string
	OutputFile    string
	EnableSSL     bool
	KafView       bool
	KafViewUrl    string
}

const APP_CFG_PATH string = "./app_config.json"

func InitConfig(kafkaBroker, groupId, sslConfigFile, topic, outputFile string, enableSSL bool, kafView bool) error {

	appCfgFile, err := os.Open(APP_CFG_PATH)
	if err != nil {
		logger.Error("Error opening app config file..")
		return err
	}

	var appConfig AppConfig

	err = json.NewDecoder(appCfgFile).Decode(&appConfig)
	if err != nil {
		logger.Error("Error reading app config file")
		return err
	}

	logger.Debug("App config: ", "config", appConfig)

	KafkaBroker = appConfig.KafkaBroker
	Topic = appConfig.Topic
	GroupId = appConfig.GroupId
	ConfigFile = appConfig.SslConfigFile
	OutputFile = appConfig.OutputFile
	EnableSSL = appConfig.EnableSSL
	KafView = appConfig.KafView
	KafViewUrl = appConfig.KafViewUrl

	if kafkaBroker != "" {
		KafkaBroker = kafkaBroker
	}
	if topic != "" {
		Topic = topic
	}
	if groupId != "" {
		GroupId = groupId
	}
	if sslConfigFile != "" {
		ConfigFile = sslConfigFile
	}
	if outputFile != "" {
		OutputFile = outputFile
	}
	if enableSSL {
		EnableSSL = enableSSL
	}
	if kafView {
		KafView = kafView
	}

	return nil
}
