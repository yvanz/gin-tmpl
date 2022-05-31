/*
@Date: 2021/1/12 下午2:23
@Author: yvanz
@File : run
@Desc:
*/

package app

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	"github.com/yvanz/gin-tmpl/internal/common"
	"github.com/yvanz/gin-tmpl/internal/config"
	"github.com/yvanz/gin-tmpl/internal/handler"
	"github.com/yvanz/gin-tmpl/models"
	"github.com/yvanz/gin-tmpl/pkg/apiserver"
	"github.com/yvanz/gin-tmpl/pkg/apiserver/conf"
	"github.com/yvanz/gin-tmpl/pkg/logger"
	"github.com/yvanz/gin-tmpl/pkg/version"
)

const projectName = "gin-demo"

var (
	configFile string
	rootCmd    = &cobra.Command{
		Short: projectName,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
	versionCommand = version.NewVerCommand(projectName)
	envCommand     = apiserver.NewConfigEnvCommand(config.G)
	initDB         = models.NewCreateDatabaseCommand(&configFile)
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "configs/dev.yaml", "configuration file path")
	rootCmd.AddCommand(versionCommand, envCommand, initDB)
}

func run() (err error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// read config file
	err = conf.LoadConfig(configFile, config.G)
	if err != nil {
		logger.Errorf("config file init failed: %s", err.Error())
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

	// uncomment or delete code below as you need
	// producer.NewProducer(config.G.Kafka)
	// c, err := kafka.Default().NewConsumer()
	// if err != nil {
	// 	logger.Fatal(err.Error())
	// 	return
	// }
	// if err = consumer.RunConsume(c); err != nil {
	// 	logger.Fatal(err.Error())
	// 	return
	// }

	// 初始化 validator 翻译器
	if err = common.InitTrans("zh"); err != nil {
		logger.Errorf("init trans failed, err: %s", err.Error())
		return err
	}

	server.Start()
	return err
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
