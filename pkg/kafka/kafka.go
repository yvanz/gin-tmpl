package kafka

import (
	"context"
	"fmt"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/yvanz/gin-tmpl/pkg/logger"
)

var (
	// 仅用于单例模式
	_defaultKafka *CliCfg
)

type Cli interface {
	Address() string                       // Address 返回kafka地址
	NewConsumer() (*ConsumerClient, error) // NewConsumer 新建消费者，kafka连接失败将返回error
	NewAsyncProducerClient() (AsyncProducer, error)
	NewSyncProducerClient() (sarama.SyncProducer, error)
}

type CliCfg struct {
	ctx      context.Context
	kafkaCfg *sarama.Config // 生产者配置
	config   *Config        // kafka配置
	addr     []string       // kafka集群地址
}

// 仅用于单例模式
func Default() *CliCfg {
	return _defaultKafka
}

func (k *CliCfg) Address() string {
	return strings.Join(k.addr, ",")
}

func (k *CliCfg) checkCli() error {
	if k == nil {
		return fmt.Errorf("kafka client is not initialized yet")
	}

	return nil
}

func (k *CliCfg) NewAsyncProducerClient() (AsyncProducer, error) {
	if err := k.checkCli(); err != nil {
		return nil, err
	}

	logger.Debug("waiting for creating kafka producer")
	ctx, cancel := context.WithCancel(k.ctx)

	go func() {
		<-k.ctx.Done()
		logger.Infof("srv stopped, stop producer together")
		cancel()
	}()

	cli := &AsyncProducerClient{
		messageChan: make(chan *sendMessage, k.config.QueueLength),
		asyncError:  make(chan *sarama.ProducerError, k.config.QueueLength),
		errLength:   k.config.QueueLength,
		ctx:         ctx,
	}

	// 异步生产者不建议把 Errors 和 Successes 都开启，一般开启 Errors 就行
	k.kafkaCfg.Producer.Return.Successes = false
	producer, err := sarama.NewAsyncProducer(k.addr, k.kafkaCfg)
	if err != nil {
		return nil, fmt.Errorf("create async producer failed: %s", err.Error())
	}

	cli.asyncProducer = producer
	return cli, nil
}

func (k *CliCfg) NewSyncProducerClient() (sarama.SyncProducer, error) {
	if err := k.checkCli(); err != nil {
		return nil, err
	}

	return sarama.NewSyncProducer(k.addr, k.kafkaCfg)
}

func (k *CliCfg) NewConsumer() (*ConsumerClient, error) {
	if err := k.checkCli(); err != nil {
		return nil, err
	}

	return &ConsumerClient{
		kafkaOptions: k,
		group:        make(map[string]sarama.ConsumerGroup),
	}, nil
}
