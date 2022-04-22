/*
@Date: 2021/10/27 17:49
@Author: yvanz
@File : log
*/

package gormdb

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/yvanz/gin-tmpl/pkg/logger"
	logg "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type DBLog struct {
	logg.Config
}

func (d *DBLog) LogMode(level logg.LogLevel) logg.Interface {
	newLogger := *d
	newLogger.LogLevel = level
	return &newLogger
}

func (d DBLog) Info(ctx context.Context, msg string, data ...interface{}) {
	if d.LogLevel < logg.Info {
		return
	}

	logger.InfofWithTrace(ctx, msg, data)
}

func (d DBLog) Warn(ctx context.Context, msg string, data ...interface{}) {
	if d.LogLevel < logg.Warn {
		return
	}

	logger.WarnfWithTrace(ctx, msg, data)
}

func (d DBLog) Error(ctx context.Context, msg string, data ...interface{}) {
	if d.LogLevel < logg.Error {
		return
	}

	logger.ErrorfWithTrace(ctx, msg, data)
}

func (d DBLog) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if d.LogLevel <= logg.Silent {
		return
	}

	fullFilePath := utils.FileWithLineNum()
	filePathList := strings.Split(fullFilePath, "/")
	filePath := fmt.Sprintf("%s/%s", filePathList[len(filePathList)-2], filePathList[len(filePathList)-1])

	elapsed := time.Since(begin)
	switch {
	case err != nil && d.LogLevel >= logg.Error && (!errors.Is(err, logg.ErrRecordNotFound) || !d.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			logger.ErrorfWithTrace(ctx, "call by %s get an error: %s, cost %f - sql is %s", filePath, err.Error(), float64(elapsed.Nanoseconds())/1e6, sql)
		} else {
			logger.ErrorfWithTrace(ctx, "call by %s get an error: %s, cost %f, affected %d records with sql %s", filePath, err.Error(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > d.SlowThreshold && d.SlowThreshold != 0 && d.LogLevel >= logg.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", d.SlowThreshold)
		if rows == -1 {
			logger.WarnfWithTrace(ctx, "call by %s show %s, cost %f - sql is %s", filePath, slowLog, float64(elapsed.Nanoseconds())/1e6, sql)
		} else {
			logger.WarnfWithTrace(ctx, "call by %s show %s, cost %f, affected %d records with sql %s", filePath, slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case d.LogLevel == logg.Info:
		sql, rows := fc()
		if rows == -1 {
			logger.InfofWithTrace(ctx, "call by %s, cost %f - sql is %s", filePath, float64(elapsed.Nanoseconds())/1e6, sql)
		} else {
			logger.InfofWithTrace(ctx, "call by %s, cost %f, affected %d records with sql %s", filePath, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}

func initLogger(level string) logg.Interface {
	var logLevel logg.LogLevel
	switch level {
	case "silent":
		logLevel = logg.Silent
	case "info":
		logLevel = logg.Info
	case "warn", "waring":
		logLevel = logg.Warn
	case "error":
		logLevel = logg.Error
	default:
		logLevel = logg.Silent
	}

	config := logg.Config{
		SlowThreshold:             time.Second,
		LogLevel:                  logLevel,
		IgnoreRecordNotFoundError: true,
		Colorful:                  false,
	}

	return &DBLog{
		Config: config,
	}
}
