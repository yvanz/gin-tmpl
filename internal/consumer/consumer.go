/*
@Date: 2021/1/14 上午10:28
@Author: yvanz
@File : consumer
@Desc:
*/

package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/yvanz/gin-tmpl/models"
	"github.com/yvanz/gin-tmpl/pkg/gormdb"
	"github.com/yvanz/gin-tmpl/pkg/kafka"
	"github.com/yvanz/gin-tmpl/pkg/logger"
)

func RunConsume(consumer *kafka.ConsumerClient) (err error) {
	hand := kafka.NewConsumerGroup(handler)
	err = consumer.RunConsumer("test-group", []string{"test"}, hand)
	if err != nil {
		logger.Errorf("failed to consume: %s", err.Error())
		return err
	}

	if !consumer.IsRunning() {
		return fmt.Errorf("queue consumer is not running")
	}

	logger.Infof("Kafka consumer is ready")
	return nil
}

type DemoMessages struct {
	UserName string `json:"user_name"`
}

func handler(message *sarama.ConsumerMessage) {
	switch message.Topic {
	case "test":
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
	default:
		logger.Errorf("unknown topic [%s] with message %+v", message.Topic, message.Value)
	}
}

func consumerPurchase(data DemoMessages) error {
	db := gormdb.GetDB().Master(context.TODO())
	crud := gormdb.NewCRUD(db)

	tmp := &models.Demo{
		UserName: data.UserName,
	}

	return crud.Create(tmp)
}
