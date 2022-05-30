/*
@Date: 2021/11/17 10:50
@Author: yvanz
@File : cli
*/

package rediscache

import (
	"github.com/go-redis/redis/v8"
)

var _rdb *redis.Client

func GetCli() *redis.Client {
	return _rdb
}
