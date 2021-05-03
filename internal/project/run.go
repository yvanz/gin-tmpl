/*
@Date: 2021/1/12 下午2:23
@Author: yvan.zhang
@File : run
@Desc:
*/

package project

import (
	"context"
	"fmt"
	"gin-tmpl/internal/common"
	"gin-tmpl/internal/config"
	"gin-tmpl/internal/project/server"
	"gin-tmpl/pkg/configutil"
	"gin-tmpl/pkg/logger"

	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

var (
	configFile string
	rootCmd    = &cobra.Command{
		Short: "gin-demo",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "configs/dev.yaml", "configuration file path")
}

func run() (err error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		if err != nil {
			cancel()
		}
	}()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	if configFile != "" {
		if err = configutil.Load(configFile, &config.DefaultConfig); err != nil {
			return fmt.Errorf("加载配置文件失败，%s", err.Error())
		}
	}

	// 解析配置文件
	if err = config.DefaultConfig.Parse(); err != nil {
		return fmt.Errorf("解析配置文件失败,%s, config: %+v", err, config.DefaultConfig)
	}

	zapLog := logger.ConfigureLogger(config.DefaultConfig.Log)
	defer func() {
		_ = zapLog.Sync()
	}()

	logger.Debugf("%+v", config.DefaultConfig)

	db, err := config.DefaultConfig.MySQL.BuildMySQLClient()
	if err != nil {
		logger.Fatalf("Error connect db by db config: %s", err)
	}
	defer db.Close()

	_, err = config.DefaultConfig.Kafka.BuildKafka(ctx)
	if err != nil {
		logger.Fatalf("Error build kafka client by config: %s", err)
	}

	// 初始化 validator 翻译器
	if err := common.InitTrans("zh"); err != nil {
		logger.Fatalf("init trans failed, err: %v\n", err)
		return err
	}

	srv := server.New(ctx, config.DefaultConfig, zapLog)
	if err := srv.Init(); err != nil {
		logger.Fatalf("init srv failed: %s", err.Error())
	}

	srv.Start()

	<-exit
	srv.Stop()
	cancel()
	logger.Infof("srv shutdown")

	return err
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
