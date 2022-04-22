/*
@Date: 2021/10/12 16:26
@Author: yvanz
@File : config
*/

package logger

import (
	"fmt"

	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	defaultConfig *Options
)

type Config struct {
	Level       LogLevel        `yaml:"level" env:"LogLevel" env-default:"info" env-description:"log level"`
	Encoding    ZapConfEncoding `yaml:"encoding" env:"LogEncoding" env-default:"console" env-description:"log encoding"`
	Development bool            `yaml:"development"`
	EnableTrace bool            `yaml:"enable_trace"`
	LogPath     string          `yaml:"log_path" env:"LogPath" env-description:"which path the log file should be"`
	LogName     string          `yaml:"log_name" env:"LogFileName" env-description:"which file name the log file should be"`
	MaxSize     int             `yaml:"max_size" env:"LogMaxSize" env-description:"max size of rotating"`
	MaxAge      int             `yaml:"max_age" env:"LogMaxAge" env-description:"max age of rotating"`
	LocalTime   bool            `yaml:"localtime"`
	Compress    bool            `yaml:"compress" env:"LogCompress" env-description:"compress old log files or not"`
}

type Options struct {
	Config
	zapConfig zap.Config
}

func (o *Options) CompareOptions() string {
	return fmt.Sprintf("level is %s, encoding is %s, log path is %s,", o.Level, o.Encoding.String(), o.LogPath)
}

func initLumberjackConf(o *Options) *lumberjack.Logger {
	if o.MaxSize == 0 {
		o.MaxSize = 200
		o.LocalTime = true
		o.Compress = true
	}
	if o.MaxAge == 0 {
		o.MaxAge = 28
	}

	return &lumberjack.Logger{
		MaxSize:   o.MaxSize, // megabytes
		MaxAge:    o.MaxAge,  // days
		LocalTime: o.LocalTime,
		Compress:  o.Compress,
	}
}
