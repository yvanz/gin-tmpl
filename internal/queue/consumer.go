/*
@Date: 2021/1/14 上午10:28
@Author: yvan.zhang
@File : consumer
@Desc:
*/

package queue

import (
	"encoding/json"
	"fmt"
	"gin-tmpl/internal/config"
	"gin-tmpl/models"
	"gin-tmpl/pkg/kafka"
	"gin-tmpl/pkg/logger"
	"strings"

	"github.com/Shopify/sarama"
)

func RunConsume(consumer *kafka.ConsumerClient) (err error) {
	topicList := strings.Split(config.DefaultConfig.App.KafkaTopic, ",")

	hand := kafka.NewConsumerGroupHandler(handler)
	err = consumer.RunConsumer(config.DefaultConfig.App.KafkaGroup, topicList, hand)
	if err != nil {
		logger.Errorf("Failed to consume: %s", err.Error())
		return err
	}

	if consumer.IsRunning() {
		logger.Infof("Kafka consumer is ready")
	} else {
		return fmt.Errorf("queue consumer is not running")
	}

	return nil
}

type DemoMessages struct {
	UserName string `json:"user_name"`
}

func handler(message *sarama.ConsumerMessage) {
	var tmp DemoMessages

	err := json.Unmarshal(message.Value, &tmp)
	if err != nil {
		logger.Errorf("Unmarshal %s of message offset %d with partition %d failed: %s", string(message.Value), message.Offset, message.Partition, err.Error())
	} else {
		err = consumerPurchase(tmp)
		if err != nil {
			logger.Errorf("create data failed, message offset is %d with partition %d: %s", message.Offset, message.Partition, err.Error())
		}
	}
}

func consumerPurchase(data DemoMessages) error {
	tmp := models.Demo{
		UserName: data.UserName,
	}

	return tmp.Add()
}
