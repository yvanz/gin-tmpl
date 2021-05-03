/*
@Date: 2021/1/12 下午2:23
@Author: yvan.zhang
@File : config
@Desc:
*/

package config

import (
	"gin-tmpl/pkg/kafka"
	"gin-tmpl/pkg/logger"
	"gin-tmpl/pkg/xormmysql"

	"github.com/gin-gonic/gin"
)

const (
	ServerName = `gin-demo`
	Port       = 80
	Version    = `v0.0.1`
)

type Config struct {
	App   AppConfig          `yaml:"app"`
	MySQL xormmysql.DBConfig `yaml:"mysql"`
	Log   *logger.Options    `yaml:"log"`
	Kafka kafka.Config       `yaml:"kafka"`
}

type AppConfig struct {
	ServiceName string `yaml:"service_name"`
	LocalIP     string `yaml:"local_ip"`
	APIPort     int64  `yaml:"api_port"`
	RunMode     string `yaml:"run_mode"`
	KafkaTopic  string `yaml:"kafka_topic"`
	KafkaGroup  string `yaml:"kafka_group"`
	Version     string
}

func (c *AppConfig) Init() {
	if c.ServiceName == "" {
		c.ServiceName = ServerName
	}

	if c.LocalIP == "" {
		c.LocalIP = "0.0.0.0"
	}

	if c.APIPort == 0 {
		c.APIPort = Port
	}

	if c.RunMode != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	if c.KafkaGroup == "" {
		c.KafkaGroup = "gin-demo"
	}
	if c.KafkaTopic == "" {
		c.KafkaTopic = "gin-demo"
	}

	c.Version = Version
}

func (g *Config) Parse() error {
	g.App.Init()
	g.Log.LogName = g.App.ServiceName

	return nil
}

var DefaultConfig = new(Config)
