/*
@Date: 2021/10/29 11:15
@Author: yvan.zhang
@File : main
*/

package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/yvanz/gin-tmpl/pkg/gormdb"
	"github.com/yvanz/gin-tmpl/pkg/logger"
	"gorm.io/gorm"
)

var testConfig = gormdb.DBConfig{
	WriteDBHost:     "localhost",
	WriteDBPort:     3306,
	WriteDBUser:     "root",
	WriteDBPassword: "root",
	WriteDB:         "gorm",
	ReadDBHostList:  []string{"localhost"},
	ReadDBPort:      3306,
	ReadDBUser:      "root",
	ReadDBPassword:  "root",
	ReadDB:          "gorm2",
	Prefix:          "cmp_bpmn_",
	MaxIdleConns:    10,
	MaxOpenConns:    100,
	Logging:         true,
}

type User struct {
	gorm.Model
	Name string
}

func main() {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())

	_, err := testConfig.BuildMySQLClient(ctx)
	if err != nil {
		logger.Fatal(err)
	}
	err = gormdb.GetDB().Migration(&User{})
	if err != nil {
		logger.Fatal(err)
	}

	var user User
	data := gormdb.GetDB().Master(ctx).First(&user)
	if data.Error != nil {
		logger.Fatal(data.Error)
	}

	logger.Info(user.Name)

	<-exit
	cancel()
}
