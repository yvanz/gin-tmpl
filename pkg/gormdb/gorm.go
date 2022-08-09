/*
@Date: 2021/10/27 17:49
@Author: yvanz
@File : gorm
*/

package gormdb

import (
	"context"
	"database/sql"
	"errors"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/yvanz/gin-tmpl/pkg/gadget"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

var (
	// 仅用于单例模式下
	_default  *DB
	ErrClient = errors.New("mysql client is not initialized yet")
)

func GetDB() *DB {
	if _default == nil {
		return &DB{}
	}

	return _default
}

// Cli is a shortcut
func Cli(ctx context.Context) *gorm.DB {
	return GetDB().Master(ctx)
}

type DB struct {
	db       *gorm.DB
	writeSQL *sql.DB
	mock     sqlmock.Sqlmock
	ctx      context.Context
}

// Master check *gorm.DB if is nil
func (d *DB) Master(ctx context.Context) *gorm.DB {
	if d == nil {
		return nil
	}

	spanCtx, err := gadget.ExtractTraceSpan(ctx)
	if err != nil {
		return d.db
	}

	return d.db.WithContext(spanCtx)
}

func (d *DB) Migration(dst ...interface{}) error {
	if d == nil {
		return ErrClient
	}

	return d.db.Clauses(dbresolver.Write).AutoMigrate(dst...)
}

func (d *DB) Close() (err error) {
	if d == nil {
		return nil
	}

	if d.db != nil {
		err = d.writeSQL.Close()
	}

	return
}
