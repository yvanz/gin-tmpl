// Copyright (c) 2016-2019 Uber Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package logger

// This package wraps logger functionality that is being used
// in kraken providing seamless migration tooling if needed
// and hides out some initialization details

import (
	"fmt"
	"os"
	"path"
	"reflect"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	DefaultLog    *DemoLog
	defaultConfig *Options
)

type DemoLog struct {
	sync.Mutex
	*zap.SugaredLogger
	config      *zap.Config
	logDir      string
	logBaseName string
	logFullPath string
	createTime  time.Time
	isRotate    bool
	rate        time.Duration
}

type Options struct {
	zapConfig   zap.Config `yaml:"-"`
	Level       LogLevel   `yaml:"level" json:"level"`
	Development bool       `yaml:"development" json:"development"`
	LogPath     string     `yaml:"log_path" json:"log_path"`
	LogName     string     `yaml:"log_name" json:"log_name"`
	Rotate      RotateRole `yaml:"rotate" json:"rotate"`
}

type RotateRole string
type LogLevel string

func (r RotateRole) Parse() time.Duration {
	switch string(r) {
	case `minute`:
		return 60 * time.Second
	case `hour`:
		return time.Hour
	case `day`:
		return 24 * time.Hour
	case `week`:
		return 7 * 24 * time.Hour
	default:
		return time.Hour
	}
}

func (l LogLevel) parse() zapcore.Level {
	switch string(l) {
	case `debug`:
		return zapcore.DebugLevel
	case `info`:
		return zapcore.InfoLevel
	case `warn`:
		return zapcore.WarnLevel
	case `error`:
		return zapcore.ErrorLevel
	case `dpanic`:
		return zapcore.DPanicLevel
	case `panic`:
		return zapcore.PanicLevel
	case `fatal`:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func (l LogLevel) Level() zap.AtomicLevel {
	return zap.NewAtomicLevelAt(l.parse())
}

// configure a default logger
func init() {
	defaultConfig = &Options{}
	defaultConfig.zapConfig = zap.NewProductionConfig()
	defaultConfig.zapConfig.Encoding = "console"
	defaultConfig.zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	defaultConfig.zapConfig.DisableStacktrace = true

	DefaultLog = new(DemoLog)
	DefaultLog.config = &defaultConfig.zapConfig
	DefaultLog.isRotate = false
	ConfigureLogger(defaultConfig)
}

func ConfigureLogger(logOptions *Options) *DemoLog {
	if !reflect.DeepEqual(logOptions, Options{}) {
		DefaultLog.config.Level = logOptions.Level.Level()
		DefaultLog.config.Development = logOptions.Development
		if logOptions.LogPath != "" {
			DefaultLog.isRotate = true
			pwd, _ := os.Getwd()
			logPath := path.Join(pwd, logOptions.LogPath)
			DefaultLog.logDir = logPath
			if string(logPath[len(logPath)-1]) == "/" {
				DefaultLog.isRotate = false
			} else {
				if err := os.MkdirAll(DefaultLog.logDir, 0755); err != nil {
					panic(err)
				}

				if logOptions.LogName == "" {
					DefaultLog.logBaseName = "app.log"
				} else {
					DefaultLog.logBaseName = logOptions.LogName + ".log"
				}

				name := fmt.Sprintf("%s.%s", DefaultLog.logBaseName, time.Now().Format("200601021504"))
				DefaultLog.logFullPath = path.Join(DefaultLog.logDir, name)
				DefaultLog.createTime = time.Now()
				DefaultLog.rate = logOptions.Rotate.Parse()
				// DefaultLog.config.OutputPaths = append(DefaultLog.config.OutputPaths, DefaultLog.logFullPath)
				DefaultLog.config.OutputPaths = []string{DefaultLog.logFullPath}
				DefaultLog.config.ErrorOutputPaths = []string{DefaultLog.logFullPath}
			}
		}
	}
	logger, err := DefaultLog.config.Build()
	if err != nil {
		panic(err)
	}

	// Skip this wrapper in a call stack.
	logger = logger.WithOptions(zap.AddCallerSkip(1))
	DefaultLog.SugaredLogger = logger.Sugar()
	// DefaultLog.config = &zapConfig

	return DefaultLog
}

// todo 每次写入日志之前进行检查并切割
func (d *DemoLog) check() {
	d.Lock()
	defer d.Unlock()
	if !d.isRotate {
		return
	}

	if d.logBaseName == "" {
		return
	}

	oldCfg := *d.config
	oldLogName := d.logFullPath
	oldCreateTime := d.createTime

	if time.Since(d.createTime) >= d.rate {
		for i := range d.config.OutputPaths {
			if d.config.OutputPaths[i] == d.logFullPath {
				d.logFullPath = path.Join(d.logDir, fmt.Sprintf("%s.%s", d.logBaseName, time.Now().Format("200601021504")))
				d.config.OutputPaths[i] = d.logFullPath
				d.createTime = time.Now()
			}
		}

		_ = d.Sync()
		newLogger, err := d.config.Build()
		if err != nil {
			d.logFullPath = oldLogName
			d.config = &oldCfg
			d.createTime = oldCreateTime
		} else {
			newLogger = newLogger.WithOptions(zap.AddCallerSkip(1))
			d.SugaredLogger = newLogger.Sugar()
		}
	}
}

func Default() *zap.SugaredLogger {
	DefaultLog.check()
	return DefaultLog.SugaredLogger
}

// Debug uses fmt.Sprint to construct and log a message.
func Debug(args ...interface{}) {
	Default().Debug(args...)
}

// Info uses fmt.Sprint to construct and log a message.
func Info(args ...interface{}) {
	Default().Info(args...)
}

// Warn uses fmt.Sprint to construct and log a message.
func Warn(args ...interface{}) {
	Default().Warn(args...)
}

// Error uses fmt.Sprint to construct and log a message.
func Error(args ...interface{}) {
	Default().Error(args...)
}

// Panic uses fmt.Sprint to construct and log a message, then panics.
func Panic(args ...interface{}) {
	Default().Panic(args...)
}

// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
func Fatal(args ...interface{}) {
	Default().Fatal(args...)
}

// Debugf uses fmt.Sprintf to log a templated message.
func Debugf(template string, args ...interface{}) {
	Default().Debugf(template, args...)
}

// Infof uses fmt.Sprintf to log a templated message.
func Infof(template string, args ...interface{}) {
	Default().Infof(template, args...)
}

// Warnf uses fmt.Sprintf to log a templated message.
func Warnf(template string, args ...interface{}) {
	Default().Warnf(template, args...)
}

// Errorf uses fmt.Sprintf to log a templated message.
func Errorf(template string, args ...interface{}) {
	Default().Errorf(template, args...)
}

// Panicf uses fmt.Sprintf to log a templated message, then panics.
func Panicf(template string, args ...interface{}) {
	Default().Panicf(template, args...)
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func Fatalf(template string, args ...interface{}) {
	Default().Fatalf(template, args...)
}

// Debugw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
//
// When debug-level logging is disabled, this is much faster than
//  s.With(keysAndValues).Debug(msg)
func Debugw(msg string, keysAndValues ...interface{}) {
	Default().Debugw(msg, keysAndValues...)
}

// Infow logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Infow(msg string, keysAndValues ...interface{}) {
	Default().Infow(msg, keysAndValues...)
}

// Warnw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Warnw(msg string, keysAndValues ...interface{}) {
	Default().Warnw(msg, keysAndValues...)
}

// Errorw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Errorw(msg string, keysAndValues ...interface{}) {
	Default().Errorw(msg, keysAndValues...)
}

// Panicw logs a message with some additional context, then panics. The
// variadic key-value pairs are treated as they are in With.
func Panicw(msg string, keysAndValues ...interface{}) {
	Default().Panicw(msg, keysAndValues...)
}

// Fatalw logs a message with some additional context, then calls os.Exit. The
// variadic key-value pairs are treated as they are in With.
func Fatalw(msg string, keysAndValues ...interface{}) {
	Default().Fatalw(msg, keysAndValues...)
}

// With adds a variadic number of fields to the logging context.
// It accepts a mix of strongly-typed zapcore.Field objects and loosely-typed key-value pairs.
func With(args ...interface{}) *zap.SugaredLogger {
	return Default().With(args...)
}

type Logger interface {
	// Error logs a message at error priority
	Error(msg string)

	// Infof logs a message at info priority
	Infof(msg string, args ...interface{})

	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}

func (d *DemoLog) Error(msg string) {
	Default().Error(msg)
}

func (d *DemoLog) Infof(template string, args ...interface{}) {
	Default().Infof(template, args...)
}

func (d *DemoLog) Print(v ...interface{}) {
	Info(v)
}

func (d *DemoLog) Printf(format string, v ...interface{}) {
	bs := []byte(format)
	length := len(bs)
	format = string(bs[:length-1])
	Info(strings.TrimSpace(fmt.Sprintf(format, v)))
}

func (d *DemoLog) Println(v ...interface{}) {
	Info(v)
}

func (d *DemoLog) GetLogDir() string {
	return d.logDir
}

func SetLevel(level LogLevel) {
	DefaultLog.config.Level.SetLevel(level.parse())
	logger, err := DefaultLog.config.Build()
	if err != nil {
		panic(err)
	}
	logger = logger.WithOptions(zap.AddCallerSkip(1))
	DefaultLog.SugaredLogger = logger.Sugar()
}
