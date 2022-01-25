/*
@Date: 2021/11/9 16:23
@Author: yvan.zhang
@File : consumer
*/

package kafka

import (
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/yvanz/gin-tmpl/pkg/logger"
)

type Consumer interface {
	RunConsumer(group string, topic []string, handler sarama.ConsumerGroupHandler) error // 运行消费者线程
	Close()                                                                              // 关闭线程
	IsRunning() bool                                                                     // 运行状态
}

type ConsumerClient struct {
	kafkaOptions *CliCfg
	isRunning    bool
	group        map[string]sarama.ConsumerGroup
}

func (cc *ConsumerClient) RunConsumer(group string, topic []string, groupHandler sarama.ConsumerGroupHandler) error {
	if _, ok := cc.group[group]; ok {
		return fmt.Errorf("consumer of group %s exists already", group)
	}

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
			err := consumerGroupClient.Consume(cc.kafkaOptions.ctx, topic, groupHandler)
			if err != nil {
				logger.Errorf("error from consumer group: %s", err.Error())
			}
		}
	}()

	return nil
}

func (cc *ConsumerClient) Close() {
	cc.isRunning = false

	for _, g := range cc.group {
		_ = g.Close()
	}
}

func (cc *ConsumerClient) IsRunning() bool {
	return cc.isRunning
}

type ConsumerGroup struct {
	handler func(*sarama.ConsumerMessage)
}

func NewConsumerGroup(handler func(*sarama.ConsumerMessage)) sarama.ConsumerGroupHandler {
	return &ConsumerGroup{handler: handler}
}

func (ConsumerGroup) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (ConsumerGroup) Cleanup(sarama.ConsumerGroupSession) error { return nil }
func (h ConsumerGroup) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		logger.Debugf("find message of topic: %q, partition: %d, offset: %d", msg.Topic, msg.Partition, msg.Offset)
		h.handler(msg)
		sess.MarkMessage(msg, "")
	}

	return nil
}
