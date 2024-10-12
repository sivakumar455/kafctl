package config

var (
	Topic, KafkaBroker, OutputFile, GroupId, ConfigFile string
	EnableSSL                                           bool
)

func InitConfig(kafkaBroker, groupId, configFile, topic, outputFile string, enableSSL bool) {

	KafkaBroker = kafkaBroker
	Topic = topic
	GroupId = groupId
	ConfigFile = configFile
	OutputFile = outputFile
	EnableSSL = enableSSL
}
