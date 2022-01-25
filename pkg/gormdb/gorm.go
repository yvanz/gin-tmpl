/*
@Date: 2021/10/27 17:49
@Author: yvan.zhang
@File : gorm
*/

package gormdb

import (
	"context"
	"database/sql"

	"github.com/yvanz/gin-tmpl/pkg/gadget"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

var (
	// 仅用于单例模式下
	_default *DB
)

func GetDB() *DB {
	return _default
}

type DB struct {
	db       *gorm.DB
	writeSQL *sql.DB
	ctx      context.Context
}

func (d *DB) Master(ctx context.Context) *gorm.DB {
	spanCtx, err := gadget.ExtractTraceSpan(ctx)
	if err != nil {
		return d.db
	}

	return d.db.WithContext(spanCtx)
}

func (d *DB) Migration(dst ...interface{}) error {
	return d.db.Clauses(dbresolver.Write).AutoMigrate(dst...)
}

func (d *DB) Close() (err error) {
	if d.db != nil {
		err = d.writeSQL.Close()
	}

	return
}
