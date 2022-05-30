/*
@Date: 2021/12/17 14:17
@Author: yvanz
@File : basic
*/

package rediscache

import (
	"time"
)

type BasicCrud interface {
	Set(key string, value interface{}, timeOut time.Duration) (err error)
	Get(key string) (val string, err error)
}
