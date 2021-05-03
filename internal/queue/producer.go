/*
@Date: 2021/1/14 上午10:28
@Author: yvan.zhang
@File : producer
@Desc:
*/

package queue

import (
	"encoding/json"
	"gin-tmpl/pkg/kafka"
)

var kafkaProducer *kafka.AsyncProducer

func SendMessage(topic string, msg interface{}) error {
	producer := *kafkaProducer
	js, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return producer.Produce(topic, js)
}

func NewProducer(async *kafka.AsyncProducer) {
	kafkaProducer = async
}
