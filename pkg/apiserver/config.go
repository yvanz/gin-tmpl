/*
@Date: 2021/11/10 11:27
@Author: yvanz
@File : config
*/

package apiserver

import (
	"encoding/json"
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/spf13/cobra"
	"github.com/yvanz/gin-tmpl/pkg/gormdb"
	"github.com/yvanz/gin-tmpl/pkg/kafka"
	"github.com/yvanz/gin-tmpl/pkg/logger"
	"github.com/yvanz/gin-tmpl/pkg/tracer"
)

type APIConfig struct {
	App    AppConfig       `yaml:"app"`
	Log    logger.Config   `yaml:"log"`
	MySQL  gormdb.DBConfig `yaml:"mysql"`
	Kafka  kafka.Config    `yaml:"kafka"`
	Tracer tracer.Config   `yaml:"tracer"`
}

type AppConfig struct {
	ServiceName string `yaml:"service_name" env-default:"gin-project" env-description:"the name of the service"`
	HostIP      string `yaml:"local_ip" env:"HostIP" env-default:"0.0.0.0" env-description:"listening on which IP"`
	APIPort     int    `yaml:"api_port" env:"APIPort" env-default:"8000" env-description:"listening on which port"`
	AdminPort   int    `yaml:"admin_port" env:"AdminPort" env-default:"8001" env-description:"listening on which port of admin service"`
	RunMode     string `yaml:"run_mode" env:"RunMode" env-description:"run mode of the service"`
	CertFile    string `yaml:"cert_file" env:"CertFile" env-description:"cert file if server need to use tls"`
	KeyFile     string `yaml:"key_file" env:"KeyFile" env-description:"key file if server need to use tls"`
}

func (c *APIConfig) buildLogger() *logger.DemoLog {
	if c.Log.LogName == "" {
		c.Log.LogName = c.App.ServiceName
	}

	return logger.ConfigureLogger(&logger.Options{Config: c.Log})
}

func (c *APIConfig) String() string {
	configData, err := json.Marshal(c)
	if err != nil {
		fmt.Println(err)
	}

	return string(configData)
}

func NewConfigEnvCommand(c interface{}) *cobra.Command {
	return &cobra.Command{
		Use:   "env",
		Short: "Prints environment variables.",
		Run: func(*cobra.Command, []string) {
			help, _ := cleanenv.GetDescription(c, nil)
			fmt.Println(help)
		},
	}
}
