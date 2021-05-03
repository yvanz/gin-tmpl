package kafka

type Config struct {
	Addr         string `yaml:"addr"`
	QueueLength  int    `yaml:"queue_length"`
	KafkaVersion string `yaml:"kafka_version"`
	EnableLog    bool   `yaml:"enable_log"`
}
