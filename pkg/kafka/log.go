/*
@Date: 2021/11/9 16:10
@Author: yvanz
@File : log
*/

package kafka

import (
	"strings"

	"github.com/yvanz/gin-tmpl/pkg/logger"
)

const LogDebug = "debug"

type kafkaLog struct {
	Level string
}

func newKafkaLog(level string) *kafkaLog {
	return &kafkaLog{
		Level: level,
	}
}

func (d *kafkaLog) Print(v ...interface{}) {
	if d.Level == LogDebug {
		logger.Default().Debug(v...)
	} else {
		logger.Default().Info(v...)
	}
}

func (d *kafkaLog) Printf(format string, v ...interface{}) {
	if d.Level == LogDebug {
		logger.Default().Debugf(strings.TrimSpace(format), v...)
	} else {
		logger.Default().Infof(strings.TrimSpace(format), v...)
	}
}

func (d *kafkaLog) Println(v ...interface{}) {
	if d.Level == LogDebug {
		logger.Default().Debug(v...)
	} else {
		logger.Default().Info(v...)
	}
}
