/*
@Date: 2021/10/27 17:49
@Author: yvanz
@File : config
*/

package gormdb

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/yvanz/gin-tmpl/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
	gormopentracing "gorm.io/plugin/opentracing"
)

type DBConfig struct {
	WriteDBHost     string   `yaml:"write_db_host" env:"MySQLWriteHost" env-description:"mysql master host" json:"write_db_host,omitempty"`
	WriteDBPort     uint16   `yaml:"write_db_port" env:"MySQLWritePort" env-description:"mysql master port" json:"write_db_port,omitempty"`
	WriteDBUser     string   `yaml:"write_db_user" env:"MySQLWriteUser" env-description:"mysql master user" json:"write_db_user,omitempty"`
	WriteDBPassword string   `yaml:"write_db_password" env:"MySQLWritePassword" env-description:"mysql master password" json:"write_db_password,omitempty"`
	WriteDB         string   `yaml:"write_db" env:"MySQLWriteDB" env-description:"mysql master database" json:"write_db,omitempty"`
	ReadDBHostList  []string `yaml:"read_db_host_list" env:"MySQLReadHostList" env-description:"mysql slave host list" json:"read_db_host_list,omitempty"`
	ReadDBPort      uint16   `yaml:"read_db_port" env:"MySQLReadPort" env-description:"mysql slave port" json:"read_db_port,omitempty"`
	ReadDBUser      string   `yaml:"read_db_user" env:"MySQLReadUser" env-description:"mysql slave user" json:"read_db_user,omitempty"`
	ReadDBPassword  string   `yaml:"read_db_password" env:"MySQLReadPassword" env-description:"mysql slave password" json:"read_db_password,omitempty"`
	ReadDB          string   `yaml:"read_db" env:"MySQLReadDB" env-description:"mysql slave database" json:"read_db,omitempty"`
	Prefix          string   `yaml:"table_prefix" json:"prefix,omitempty"`
	MaxIdleConns    int      `yaml:"max_idle_conns" json:"max_idle_conns,omitempty"`
	MaxOpenConns    int      `yaml:"max_open_conns" json:"max_open_conns,omitempty"`
	Logging         bool     `yaml:"logging" json:"logging,omitempty"`
	LogLevel        string   `yaml:"log_level" env:"MySQLLogLevel" env-description:"log level of mysql log: silent/info/warn/error" json:"log_level,omitempty"`
	RawColumn       bool     `yaml:"-" json:"raw_column,omitempty"`
}

func (c *DBConfig) initConfig() (conf *gorm.Config, err error) {
	if c.WriteDBHost == "" {
		return conf, fmt.Errorf("mysql master not found")
	}

	if c.MaxIdleConns == 0 {
		c.MaxIdleConns = 5
	}
	if c.MaxOpenConns == 0 {
		c.MaxOpenConns = 20
	}

	var namingStrategy schema.Namer
	if c.RawColumn {
		namingStrategy = MyNamingStrategy{ns: schema.NamingStrategy{TablePrefix: c.Prefix, SingularTable: true}}
	} else {
		namingStrategy = schema.NamingStrategy{TablePrefix: c.Prefix, SingularTable: true}
	}

	conf = &gorm.Config{
		NamingStrategy: namingStrategy,
	}
	if c.Logging {
		conf.Logger = initLogger(c.LogLevel)
	}

	return
}

func (c DBConfig) BuildMySQLClient(ctx context.Context) (*DB, error) {
	logger.Debug("build mysql client")

	var master *gorm.DB
	var sqlDBMaster *sql.DB

	if _default != nil {
		return _default, nil
	}

	_default = &DB{ctx: ctx}

	gormConfig, err := c.initConfig()
	if err != nil {
		return nil, err
	}

	master, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       createDSN(c.WriteDBUser, c.WriteDBPassword, c.WriteDBHost, c.WriteDB, c.WriteDBPort),
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), gormConfig)
	if err != nil {
		return nil, err
	}

	err = master.Use(gormopentracing.New())
	if err != nil {
		return nil, err
	}

	slaves := make([]gorm.Dialector, 0)
	for _, host := range c.ReadDBHostList {
		if host == "" {
			continue
		}

		dsn := createDSN(c.ReadDBUser, c.ReadDBPassword, host, c.ReadDB, c.ReadDBPort)
		slaves = append(slaves, mysql.Open(dsn))
	}

	if len(slaves) > 0 {
		if master == nil {
			return nil, fmt.Errorf("mysql master init failed")
		}

		err = master.Use(dbresolver.Register(dbresolver.Config{Replicas: slaves, Policy: dbresolver.RandomPolicy{}}).
			SetConnMaxIdleTime(time.Hour).SetConnMaxLifetime(24 * time.Hour).SetMaxIdleConns(c.MaxIdleConns).SetMaxOpenConns(c.MaxOpenConns))
		if err != nil {
			return nil, err
		}
	} else {
		err = master.Use(dbresolver.Register(dbresolver.Config{}).
			SetConnMaxIdleTime(time.Hour).SetConnMaxLifetime(24 * time.Hour).SetMaxIdleConns(c.MaxIdleConns).SetMaxOpenConns(c.MaxOpenConns))
		if err != nil {
			return nil, err
		}
	}

	sqlDBMaster, err = master.DB()
	if err != nil {
		return nil, err
	}

	_default.db = master
	_default.writeSQL = sqlDBMaster

	return _default, nil
}

func createDSN(user, password, host, database string, port uint16) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=5s",
		user,
		password,
		host,
		port,
		database,
	)
}
