/*
@Date: 2021/11/10 11:27
@Author: yvanz
@File : config
*/

package apiserver

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/spf13/cobra"
	"github.com/yvanz/gin-tmpl/pkg/gormdb"
	"github.com/yvanz/gin-tmpl/pkg/kafka"
	"github.com/yvanz/gin-tmpl/pkg/logger"
	"github.com/yvanz/gin-tmpl/pkg/rediscache"
	"github.com/yvanz/gin-tmpl/pkg/tracer"
)

const (
	RunModeDebug      = "debug"
	RunModeTest       = "test"
	RunModeDev        = "dev"
	RunModeProd       = "prod"
	RunModeProduction = "production"
	RunModeRelease    = "release"
)

type APIConfig struct {
	App    AppConfig         `yaml:"app" json:"app,omitempty"`
	Log    logger.Config     `yaml:"log" json:"log,omitempty"`
	MySQL  gormdb.DBConfig   `yaml:"mysql" json:"mysql,omitempty"`
	Redis  rediscache.Config `yaml:"redis" json:"redis,omitempty"`
	Kafka  kafka.Config      `yaml:"kafka" json:"kafka,omitempty"`
	Tracer tracer.Config     `yaml:"tracer" json:"tracer,omitempty"`
}

type AppConfig struct {
	ServiceName string `yaml:"service_name" env-default:"gin-project" env-description:"the name of the service" json:"service_name,omitempty"`
	HostIP      string `yaml:"local_ip" env:"HostIP" env-default:"0.0.0.0" env-description:"listening on which IP" json:"host_ip,omitempty"`
	APIPort     int    `yaml:"api_port" env:"APIPort" env-default:"8000" env-description:"listening on which port" json:"api_port,omitempty"`
	AdminPort   int    `yaml:"admin_port" env:"AdminPort" env-default:"8001" env-description:"listening on which port of admin service" json:"admin_port,omitempty"`
	RunMode     string `yaml:"run_mode" env:"RunMode" env-description:"run mode of the service" json:"run_mode,omitempty"`
	CertFile    string `yaml:"cert_file" env:"CertFile" env-description:"cert file if server need to use tls" json:"cert_file,omitempty"`
	KeyFile     string `yaml:"key_file" env:"KeyFile" env-description:"key file if server need to use tls" json:"key_file,omitempty"`
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

func (c *APIConfig) initService(ctx context.Context, opts *serverOptions) (err error) {
	if c.MySQL.WriteDBHost != "" {
		c.MySQL.RawColumn = opts.tableColumnWithRaw
		db, e := c.MySQL.BuildMySQLClient(ctx)
		if e != nil {
			err = e
			return
		}

		if len(opts.migrationList) > 0 {
			err = db.Migration(opts.migrationList...)
			if err != nil {
				return
			}
		}
	}

	if c.Redis.Addr != "" {
		err = c.Redis.NewRedisCli(ctx)
		if err != nil {
			return
		}
	}

	if c.Kafka.Addr != "" {
		_, err = c.Kafka.BuildKafka(ctx)
		if err != nil {
			return
		}
	}

	return err
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
