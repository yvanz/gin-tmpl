/*
@Date: 2021/1/12 下午2:23
@Author: yvanz
@File : config
@Desc:
*/

package config

import (
	"encoding/json"

	"github.com/yvanz/gin-tmpl/pkg/apiserver"
	"github.com/yvanz/gin-tmpl/pkg/logger"
)

type Config struct {
	apiserver.APIConfig `yaml:"base"`
}

func (c *Config) String() string {
	configData, err := json.Marshal(c)
	if err != nil {
		logger.Error(err.Error())
	}

	return string(configData)
}

var G = &Config{}
