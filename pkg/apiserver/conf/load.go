/*
@Date: 2021/11/11 14:13
@Author: yvan.zhang
@File : load
*/

package conf

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

func LoadConfig(configFile string, c interface{}) error {
	err := cleanenv.ReadConfig(configFile, c)
	if err != nil {
		return fmt.Errorf("read config file %s failed: %s", configFile, err.Error())
	}

	return nil
}
