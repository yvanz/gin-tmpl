/*
@Date: 2021/10/12 16:26
@Author: yvanz
@File : config
*/

package logger

import (
	"fmt"
	"os"
	"path"

	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	defaultConfig *Options
)

type Config struct {
	Level       LogLevel        `yaml:"level" env:"LogLevel" env-default:"info" env-description:"log level" json:"level,omitempty"`
	Encoding    ZapConfEncoding `yaml:"encoding" env:"LogEncoding" env-default:"console" env-description:"log encoding" json:"encoding,omitempty"`
	LogPath     string          `yaml:"log_path" env:"LogPath" env-description:"which path the log file should be" json:"log_path,omitempty"`
	LogName     string          `yaml:"log_name" env:"LogFileName" env-description:"which file name the log file should be" json:"log_name,omitempty"`
	MaxSize     int             `yaml:"max_size" env:"LogMaxSize" env-description:"max size of rotating" json:"max_size,omitempty"`
	MaxAge      int             `yaml:"max_age" env:"LogMaxAge" env-description:"max age of rotating" json:"max_age,omitempty"`
	LocalTime   bool            `yaml:"localtime" json:"local_time,omitempty"`
	Compress    bool            `yaml:"compress" env:"LogCompress" env-description:"compress old log files or not" json:"compress,omitempty"`
	Development bool            `yaml:"development" json:"development,omitempty"`
	EnableTrace bool            `yaml:"enable_trace" json:"enable_trace,omitempty"`
	DisableStd  bool            `yaml:"disable_std" json:"disable_std,omitempty"`
}

type Options struct {
	zapConfig zap.Config
	Config
}

func (o *Options) CompareOptions() string {
	return fmt.Sprintf("level is %s, encoding is %s, log path is %s,", o.Level, o.Encoding.String(), o.LogPath)
}

func (o *Options) GenLogPath() (rotatePath string) {
	pwd, _ := os.Getwd()
	logDir := path.Join(pwd, o.LogPath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic(err)
	}

	var logName string
	if o.LogName == "" {
		logName = "app.log"
	} else {
		logName = o.LogName + ".log"
	}

	fullPath := path.Join(logDir, logName)
	if fullPath[0:1] == "/" {
		rotatePath = fmt.Sprintf("rotate:/%%2F%s", fullPath[1:])
	} else {
		rotatePath = fmt.Sprintf("rotate:/%s", fullPath)
	}

	return
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
