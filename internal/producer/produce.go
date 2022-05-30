/*
@Date: 2022/1/25 14:49
@Author: yvanz
@File : produce_task
*/

package producer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/yvanz/gin-tmpl/pkg/kafka"
	"github.com/yvanz/gin-tmpl/pkg/logger"
)

const (
	TaskTopic = "test"
)

var kafkaProducer *kafka.AsyncProducer

func SendMessage(msg interface{}, keys ...string) error {
	if kafkaProducer == nil {
		return fmt.Errorf("kakfa producer is not initialized yet")
	}

	producer := *kafkaProducer
	js, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	logger.Debugf("send message is %v", string(js))

	return producer.Produce(TaskTopic, js, keys...)
}

func NewProducer(conf kafka.Config) {
	var err error

	if kafkaProducer != nil {
		return
	}

	cli := kafka.Default()
	if cli == nil {
		cli, err = conf.BuildKafka(context.TODO())
		if err != nil {
			logger.Fatal(err.Error())
			return
		}
	}

	p, err := cli.NewAsyncProducerClient()
	if err != nil {
		logger.Fatal(err.Error())
		return
	}

	p.RunAsyncProducer()
	go func() {
		for {
			e := <-p.ProducerErrors()
			logger.Errorf("message %v produce failed", e.Msg)
		}
	}()

	kafkaProducer = &p
}
