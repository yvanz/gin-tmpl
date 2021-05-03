/*
@Date: 2021/1/15 下午6:57
@Author: yvan.zhang
@File : config
@Desc:
*/

package xormmysql

type (
	DBConfig struct {
		WriteDB      MySQLConfig   `yaml:"write_db"`
		ReadDB       []MySQLConfig `yaml:"read_db"`
		Prefix       string        `yaml:"table_prefix"`
		MaxIdleConns int           `yaml:"max_idle_conns"`
		MaxOpenConns int           `yaml:"max_open_conns"`
		Logging      bool          `yaml:"logging"`
		LogLevel     string        `yaml:"log_level"`
		LogDir       string        `yaml:"log_dir"`
	}

	MySQLConfig struct {
		Host     string `yaml:"host"`
		Port     uint16 `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
	}
)
