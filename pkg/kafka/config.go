package kafka

import (
	"context"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/yvanz/gin-tmpl/pkg/gadget"
	"github.com/yvanz/gin-tmpl/pkg/logger"
)

type (
	Config struct {
		Addr         string `yaml:"addr" env:"KafkaAddr" env-description:"address of kafka cluster" json:"addr,omitempty"`
		QueueLength  int    `yaml:"queue_length" json:"queue_length,omitempty"`
		KafkaVersion string `yaml:"kafka_version" env:"KafkaVersion" env-description:"version of kafka cluster" json:"kafka_version,omitempty"`
		EnableLog    bool   `yaml:"enable_log" env:"KafkaEnableLog" env-description:"enable kafka log or not" json:"enable_log,omitempty"`
		LogLevel     string `yaml:"log_level" env:"KafkaLogLevel" env-description:"record logs in which level, only support debug/info" json:"log_level,omitempty"`
	}
)

// BuildKafka 创建Kafka客户端实例
func (c *Config) BuildKafka(ctx context.Context) (*CliCfg, error) {
	logger.Debug("build kafka client")

	if _defaultKafka != nil {
		return _defaultKafka, nil
	}

	version, err := sarama.ParseKafkaVersion(c.KafkaVersion)
	if err != nil {
		logger.Warnf("parse kafka version failed: %s, use version %v instead", err.Error(), sarama.V0_10_2_2)
		version = sarama.V0_10_2_2
	}

	if c.EnableLog {
		switch c.LogLevel {
		case LogDebug, "info":
		default:
			c.LogLevel = LogDebug
		}

		sarama.Logger = newKafkaLog(c.LogLevel)
	}

	addr := strings.Split(strings.TrimSpace(c.Addr), ",")
	kCli := &CliCfg{
		addr:   addr,
		config: c,
		ctx:    ctx,
	}

	if c.QueueLength < 1024 {
		c.QueueLength = 1024
	}

	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Net.DialTimeout = 3 * time.Second
	kafkaConfig.Metadata.Retry.Max = 1

	kafkaConfig.Producer.RequiredAcks = sarama.WaitForAll
	kafkaConfig.Producer.Return.Successes = true
	kafkaConfig.Producer.Compression = sarama.CompressionSnappy

	kafkaConfig.Consumer.Return.Errors = true
	kafkaConfig.Consumer.Offsets.Initial = sarama.OffsetNewest
	kafkaConfig.ClientID = gadget.UUID()
	kafkaConfig.Version = version

	kCli.kafkaCfg = kafkaConfig
	_defaultKafka = kCli

	return kCli, nil
}
