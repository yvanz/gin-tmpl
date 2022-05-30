/*
@Date: 2021/12/17 14:21
@Author: yvanz
@File : repo
*/

package rediscache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisCrud struct {
	Ctx context.Context
	Rdb *redis.Client
}

func NewCRUD(ctx context.Context, cli *redis.Client) BasicCrud {
	return &RedisCrud{Ctx: ctx, Rdb: cli}
}

func (c *RedisCrud) Get(key string) (val string, err error) {
	if c.Rdb == nil {
		return "", fmt.Errorf("redis client is not initialized yet")
	}

	val, err = c.Rdb.Get(c.Ctx, key).Result()

	return
}

func (c *RedisCrud) Set(key string, value interface{}, timeOut time.Duration) (err error) {
	if c.Rdb == nil {
		return fmt.Errorf("redis client is not initialized yet")
	}

	return c.Rdb.Set(c.Ctx, key, value, timeOut).Err()
}
