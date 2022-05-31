/*
@Date: 2021/1/12 下午2:37
@Author: yvanz
@File : init
@Desc:
*/

package models

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yvanz/gin-tmpl/internal/config"
	"github.com/yvanz/gin-tmpl/pkg/apiserver/conf"
)

var AllTables = []interface{}{
	&Demo{},
}

// NewCreateDatabaseCommand is prepared for creating database when init project
func NewCreateDatabaseCommand(configFile *string) *cobra.Command {
	return &cobra.Command{
		Use:   "create_db",
		Short: "create database which project needed",
		RunE: func(*cobra.Command, []string) error {
			err := conf.LoadConfig(*configFile, config.G)
			if err != nil {
				return err
			}

			if config.G.MySQL.WriteDB == "" {
				return fmt.Errorf("no database, please check the config file or command flag")
			}

			ctx := context.Background()
			dbShouldCreate := config.G.MySQL.WriteDB

			// clear db name which config specified, otherwise build client will fail
			config.G.MySQL.WriteDB = ""
			config.G.MySQL.LogLevel = "silent"
			dbCli, err := config.G.MySQL.BuildMySQLClient(ctx)
			if err != nil {
				return err
			}

			defer dbCli.Close()
			db := dbCli.Master(ctx)
			dbSQL := fmt.Sprintf(`CREATE DATABASE IF NOT EXISTS %s DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci`, dbShouldCreate)
			err = db.Exec(dbSQL).Error
			if err != nil {
				return err
			}

			fmt.Printf("well done...\ndatabase %s create successfully or exists already\n", dbShouldCreate)
			return nil
		},
	}
}
