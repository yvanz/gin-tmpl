/*
@Date: 2021/1/12 下午2:23
@Author: yvan.zhang
@File : run
@Desc:
*/

package app

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yvanz/gin-tmpl/internal/common"
	"github.com/yvanz/gin-tmpl/internal/config"
	"github.com/yvanz/gin-tmpl/internal/consumer"
	"github.com/yvanz/gin-tmpl/internal/handler"
	"github.com/yvanz/gin-tmpl/internal/producer"
	"github.com/yvanz/gin-tmpl/models"
	"github.com/yvanz/gin-tmpl/pkg/apiserver"
	"github.com/yvanz/gin-tmpl/pkg/apiserver/conf"
	"github.com/yvanz/gin-tmpl/pkg/kafka"
	"github.com/yvanz/gin-tmpl/pkg/logger"
	"github.com/yvanz/gin-tmpl/pkg/version"
)

var (
	configFile string
	rootCmd    = &cobra.Command{
		Short: "gin-demo",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
	versionCommand = version.NewVerCommand("gin-demo")
	envCommand     = apiserver.NewConfigEnvCommand(config.G)
)

func init() {
	rootCmd.AddCommand(versionCommand, envCommand)
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "configs/dev.yaml", "configuration file path")
}

func run() (err error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 解析配置文件
	err = conf.LoadConfig(configFile, config.G)
	if err != nil {
		logger.Fatalf("config file init failed: %s", err.Error())
		return
	}

	// 数据表迁移，新增表时修改 AllTables
	m := apiserver.Migration(models.AllTables)
	server := apiserver.CreateNewServer(ctx, config.G.APIConfig, m)
	defer server.Stop()

	logger.Debugf("%+v", config.G)

	group := server.AddGinGroup("/api")
	tra := server.GetTracer()
	handler.RegisterRouter(tra, group)
	producer.NewProducer(config.G.Kafka)

	c, err := kafka.Default().NewConsumer()
	if err != nil {
		logger.Fatal(err.Error())
		return
	}
	if err = consumer.RunConsume(c); err != nil {
		logger.Fatal(err.Error())
		return
	}

	// 初始化 validator 翻译器
	if err = common.InitTrans("zh"); err != nil {
		logger.Fatalf("init trans failed, err: %s", err.Error())
		return err
	}

	server.Start()
	return err
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
