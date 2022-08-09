/*
@Date: 2021/11/9 16:23
@Author: yvanz
@File : producer
*/

package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/Shopify/sarama"
	"github.com/yvanz/gin-tmpl/pkg/logger"
)

type sendMessage struct {
	topic string
	key   string
	value []byte
}

type AsyncProducer interface {
	RunAsyncProducer()                                        // 运行异步生产者线程
	Produce(topic string, value []byte, keys ...string) error // 生产消息
	ProducerErrors() <-chan *sarama.ProducerError             // 返回生产者发送消息失败的chan
	CloseProducer()                                           // 关闭线程
	IsRunning() bool                                          // 运行状态
}

type AsyncProducerClient struct {
	ctx           context.Context
	asyncProducer sarama.AsyncProducer       // 异步生产者接口，用于生产者实际操作
	asyncError    chan *sarama.ProducerError // 错误消息队列
	messageChan   chan *sendMessage          // 发送生产消息的队列
	errLength     int                        // 错误消息最大长度
	isRunning     bool                       // 生产者线程是否运行
}

func (p *AsyncProducerClient) RunAsyncProducer() {
	p.isRunning = true
	// 循环判断哪个通道发送过来数据.
	logger.Infof("start kafka producer goroutine")

	go func(producer sarama.AsyncProducer) {
		// 异步生产者发送后必须把返回值从 Errors 或者 Successes 中读出来，不然会阻塞 sarama 内部处理逻辑，导致只能发出去一条消息
		for {
			select {
			case suc := <-producer.Successes():
				if suc != nil {
					logger.Debugf("produce success, offset: %d, timestamp: %s, partitions: %d", suc.Offset, suc.Timestamp.String(), suc.Partition)
				}
			case fail := <-producer.Errors():
				if fail != nil {
					logger.Errorf("send message to kafka producer err: %s", fail.Error())
					// 写入错误队列，若队列长度已满，则移除第一个元素
					if len(p.asyncError) >= p.errLength {
						<-p.asyncError
					}

					p.asyncError <- fail
				}
			case <-p.ctx.Done():
				logger.Info("stop async producer")
				return
			}
		}
	}(p.asyncProducer)

	go func(producer sarama.AsyncProducer) {
		for {
			select {
			case <-p.ctx.Done():
				logger.Info("stop reporter queue")
				p.isRunning = false
				return
			case m := <-p.messageChan:
				msg := &sarama.ProducerMessage{
					Topic:     m.topic,
					Value:     sarama.ByteEncoder(m.value),
					Timestamp: time.Now(),
				}

				if m.key != "" {
					msg.Key = sarama.StringEncoder(m.key)
				}

				producer.Input() <- msg
				logger.Debugf("sent to kafka, topic: %s, messages_len: %d", msg.Topic, msg.Value.Length())
			}
		}
	}(p.asyncProducer)
}

// Produce 发送消息到队列。仅当需要保证消息顺序时，才使用参数 keys，并且只允许传一个 key
func (p *AsyncProducerClient) Produce(topic string, value []byte, keys ...string) error {
	if !p.isRunning {
		p.RunAsyncProducer()
	}

	msg := &sendMessage{
		topic: topic,
		value: value,
	}

	switch len(keys) {
	case 0:
	case 1:
		msg.key = keys[0]
	default:
		return fmt.Errorf("only need one key")
	}

	p.messageChan <- msg
	return nil
}

func (p *AsyncProducerClient) ProducerErrors() <-chan *sarama.ProducerError {
	return p.asyncError
}

func (p *AsyncProducerClient) CloseProducer() {
	p.isRunning = false
	p.asyncProducer.AsyncClose()
}

func (p *AsyncProducerClient) IsRunning() bool {
	return p.isRunning
}
