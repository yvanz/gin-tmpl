/*
@Date: 2021/11/18 15:24
@Author: yvan.zhang
@File : main
*/

package main

import (
	"context"

	"github.com/yvanz/gin-tmpl/pkg/kafka"
	"github.com/yvanz/gin-tmpl/pkg/logger"
)

var exampleConfig = kafka.Config{
	Addr:         "localhost:9092",
	KafkaVersion: "",
	EnableLog:    true,
	LogLevel:     "info",
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cli, err := exampleConfig.BuildKafka(ctx)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info(cli.Address())
	producer, err := cli.NewAsyncProducerClient()
	if err != nil {
		logger.Fatal(err)
	}

	producer.IsRunning()
}
