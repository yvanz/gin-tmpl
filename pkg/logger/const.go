package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type ZapConfEncoding string

const (
	ZapEncodeJSON    ZapConfEncoding = "json"
	ZapEncodeConsole ZapConfEncoding = "console"
)

func (code ZapConfEncoding) String() string {
	return string(code)
}

func (code ZapConfEncoding) IsValid() bool {
	return code == ZapEncodeJSON || code == ZapEncodeConsole
}

type LogLevel string

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

type logRotationConfig struct {
	*lumberjack.Logger
}

func (logRotationConfig) Sync() error { return nil }
