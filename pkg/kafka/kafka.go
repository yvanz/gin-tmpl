package kafka

import (
	"context"
	"fmt"
	"gin-tmpl/pkg/gadget"
	"gin-tmpl/pkg/logger"
	"strings"
	"time"

	"github.com/Shopify/sarama"
)

var (
	// 仅用于单例模式
	_defaultKafka *CliCfg
)

type Cli interface {
	// Address 返回kafka地址
	Address() string
	// NewProducerClient 新建生产者，kafka连接失败将返回error
	NewAsyncProducerClient() (AsyncProducer, error)
	NewSyncProducerClient() (sarama.SyncProducer, error)
	// NewConsumer 新建消费者，kafka连接失败将返回error
	NewConsumer() (*ConsumerClient, error)
}

type CliCfg struct {
	// kafka集群地址
	addr []string
	// 生产者配置
	kafkaCfg *sarama.Config
	// 上线文
	ctx context.Context
	// kafka配置
	config *Config
}

// 仅用于单例模式
func Default() *CliCfg {
	return _defaultKafka
}

// BuildKafka 创建Kafka客户端实例
func (c *Config) BuildKafka(ctx context.Context) (*CliCfg, error) {
	if len(c.KafkaVersion) == 0 {
		c.KafkaVersion = `0.8.2.0`
	}
	version, err := sarama.ParseKafkaVersion(c.KafkaVersion)
	if err != nil {
		return nil, fmt.Errorf(`kafka's config version is invalid'`)
	}

	if c.EnableLog {
		sarama.Logger = logger.DefaultLog
	}

	kCli := new(CliCfg)
	kCli.ctx = ctx
	kCli.config = c

	addr := strings.Split(strings.TrimSpace(c.Addr), ",")
	kCli.addr = addr

	// kCli.ctx, kCli.cancel = context.WithCancel(context.Background())
	if c.QueueLength < 1024 {
		c.QueueLength = 1024
	}
	// producer config
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForAll
	kafkaConfig.Producer.Return.Successes = true
	kafkaConfig.Producer.Compression = sarama.CompressionSnappy
	kafkaConfig.Version = version

	kafkaConfig.Consumer.Return.Errors = true
	kafkaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	kafkaConfig.ClientID = gadget.UUID()

	kCli.kafkaCfg = kafkaConfig

	if _defaultKafka == nil {
		_defaultKafka = kCli
	}

	return kCli, nil
}

func (k *CliCfg) Address() string {
	return strings.Join(k.addr, ",")
}

func (k *CliCfg) NewAsyncProducerClient() (AsyncProducer, error) {
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
	producer, err := sarama.NewAsyncProducer(k.addr, k.kafkaCfg)
	if err != nil {
		return nil, fmt.Errorf("kafka error, %s", err.Error())
	}
	cli.asyncProducer = producer
	return cli, nil
}

func (k *CliCfg) NewSyncProducerClient() (sarama.SyncProducer, error) {
	return sarama.NewSyncProducer(k.addr, k.kafkaCfg)
}

func (k *CliCfg) NewConsumer() (*ConsumerClient, error) {
	consumer, err := sarama.NewConsumer(k.addr, k.kafkaCfg)
	if err != nil {
		return nil, fmt.Errorf("can't create kafka consumer, address: %s, err: %v", k.addr, err)
	}

	return &ConsumerClient{
		kafkaOptions: k,
		consumer:     consumer,
		group:        make(map[string]sarama.ConsumerGroup),
	}, nil
}

type AsyncProducerClient struct {
	// 异步生产者接口，用于生产者实际操作
	asyncProducer sarama.AsyncProducer
	// 错误消息队列
	asyncError chan *sarama.ProducerError
	// 发送生产消息的队列
	messageChan chan *sendMessage
	ctx         context.Context
	// 错误消息最大长度
	errLength int
	// 生产者线程是否运行
	isRunning bool
}

type sendMessage struct {
	topic string
	value []byte
}

type AsyncProducer interface {
	// 运行异步生产者线程
	RunAsyncProducer()
	// 生产消息
	Produce(topic string, value []byte) error
	// 返回生产者发送消息失败的chan
	ProducerErrors() <-chan *sarama.ProducerError
	// 关闭线程
	CloseProducer()
	// 运行状态
	IsRunning() bool
}

func (client *AsyncProducerClient) RunAsyncProducer() {
	client.isRunning = true
	// 循环判断哪个通道发送过来数据.
	logger.Infof("start kafka producer goroutine...")
	go func(producer sarama.AsyncProducer) {
		for {
			select {
			case suc := <-producer.Successes():
				if suc != nil {
					logger.Debugf("offset:%d, timestamp:%s, partitions:%d", suc.Offset, suc.Timestamp.String(), suc.Partition)
				}
			case fail := <-producer.Errors():
				if fail != nil {
					logger.Error("send message to kafka producer err: ", fail.Error())
					// 写入错误队列，若队列长度已满，则移除第一个元素
					if len(client.asyncError) >= client.errLength {
						<-client.asyncError
					}

					client.asyncError <- fail
				}
			case <-client.ctx.Done():
				logger.Info("stop async producer")
				return
			}
		}
	}(client.asyncProducer)

	go func(producer sarama.AsyncProducer) {
		for {
			select {
			case <-client.ctx.Done():
				logger.Info("stop reporter queue")
				client.isRunning = false
				return

			case m := <-client.messageChan:
				msg := &sarama.ProducerMessage{
					Topic:     m.topic,
					Value:     sarama.ByteEncoder(m.value),
					Timestamp: time.Now(),
				}

				producer.Input() <- msg
				logger.Debugf("sent to kafka, topic:%s, messages_len: %d", msg.Topic, msg.Value.Length())
			}
		}
	}(client.asyncProducer)
}

func (client *AsyncProducerClient) Produce(topic string, value []byte) error {
	if !client.isRunning {
		client.RunAsyncProducer()
	}

	msg := &sendMessage{
		topic: topic,
		value: value,
	}

	client.messageChan <- msg
	return nil
}

func (client *AsyncProducerClient) ProducerErrors() <-chan *sarama.ProducerError {
	return client.asyncError
}

func (client *AsyncProducerClient) CloseProducer() {
	client.isRunning = false
	client.asyncProducer.AsyncClose()
}

func (client *AsyncProducerClient) IsRunning() bool {
	return client.isRunning
}

type Consumer interface {
	// 运行消费者线程
	RunConsumer(group string, topic []string, handler sarama.ConsumerGroupHandler) error
	// 关闭线程
	Close()
	// 获取消费消息
	// Message() <-chan *sarama.ConsumerMessage
	// 运行状态
	IsRunning() bool
}

type ConsumerClient struct {
	kafkaOptions *CliCfg
	consumer     sarama.Consumer
	isRunning    bool
	group        map[string]sarama.ConsumerGroup
}

func (cc *ConsumerClient) RunConsumer(group string, topic []string, handler sarama.ConsumerGroupHandler) error {
	cc.isRunning = true

	consumerGroupClient, err := sarama.NewConsumerGroup(cc.kafkaOptions.addr, group, cc.kafkaOptions.kafkaCfg)
	if err != nil {
		return err
	}
	cc.group[group] = consumerGroupClient

	go func() {
		for {
			err := <-consumerGroupClient.Errors()
			if err != nil {
				logger.Error(err)
			}
		}
	}()

	go func() {
		for {
			err = consumerGroupClient.Consume(cc.kafkaOptions.ctx, topic, handler)
			if err != nil {
				logger.Errorf("Error from consumer: %v", err)
			}
		}
	}()
	return nil
}

func (cc *ConsumerClient) Close() {
	cc.isRunning = false
	cc.consumer.Close()
}

func (cc *ConsumerClient) IsRunning() bool {
	return cc.isRunning
}

type ConsumerGroupHandler struct {
	handler func(*sarama.ConsumerMessage)
}

func NewConsumerGroupHandler(handler func(msg *sarama.ConsumerMessage)) sarama.ConsumerGroupHandler {
	return &ConsumerGroupHandler{handler: handler}
}

func (ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h ConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		logger.Debugf("Message topic:%q partition:%d offset:%d", msg.Topic, msg.Partition, msg.Offset)
		h.handler(msg)
		sess.MarkMessage(msg, "")
	}
	return nil
}
