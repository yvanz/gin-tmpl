/*
@Date: 2021/11/17 10:48
@Author: yvanz
@File : config
*/

package rediscache

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/yvanz/gin-tmpl/pkg/logger"
)

type (
	Config struct {
		Addr           string         `yaml:"host_and_port" env:"RedisHostAndPort" end-description:"redis host and port, seems like 127.0.0.1:6379"`
		Username       string         `yaml:"user_name" env:"RedisUsername"`
		Password       string         `yaml:"password" env:"RedisPassword"`
		DB             int            `yaml:"db" env:"RedisDB"`
		PoolSize       int            `yaml:"pool_size" env:"RedisPoolSize"`
		SentinelConfig sentinelConfig `yaml:"sentinel"`
		ServerType     string         `yaml:"server_type" env:"RedisServerType" env-default:"standalone" end-description:"redis type, support standalone/sentinel only"`
	}
	sentinelConfig struct {
		Addrs      []string `yaml:"sentinel_addrs" env:"RedisSentinelAddrs"`
		MasterName string   `yaml:"sentinel_master_name" env:"RedisSentinelMasterName"`
		Username   string   `yaml:"sentinel_username" env:"RedisSentinelUsername"`
		Password   string   `yaml:"sentinel_password" env:"RedisSentinelPassword"`
	}
)

func (c *Config) NewRedisCli(ctx context.Context) error {
	if _rdb != nil {
		return nil
	}

	logger.Debug("build redis cli")
	var rdb *redis.Client
	switch c.ServerType {
	case "standalone":
		clientOpts := &redis.Options{
			Addr:     c.Addr,
			Username: c.Username,
			Password: c.Password,
			DB:       c.DB,
			PoolSize: c.PoolSize,
		}
		rdb = redis.NewClient(clientOpts)
	case "sentinel":
		failoverOptions := &redis.FailoverOptions{
			MasterName:       c.SentinelConfig.MasterName,
			SentinelAddrs:    c.SentinelConfig.Addrs,
			SentinelUsername: c.SentinelConfig.Username,
			SentinelPassword: c.SentinelConfig.Password,
			DB:               c.DB,
			PoolSize:         c.PoolSize,
		}
		rdb = redis.NewFailoverClient(failoverOptions)
	default:
		return fmt.Errorf("unsupported server type: %s", c.ServerType)
	}

	err := rdb.Ping(ctx).Err()
	if err != nil {
		return err
	}

	_rdb = rdb
	return nil
}
