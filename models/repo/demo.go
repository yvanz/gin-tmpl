package repo

import "github.com/yvanz/gin-tmpl/pkg/gormdb"

type DemoRepo[T any] interface {
	gormdb.GetListCrud[T]
	gormdb.GetByIDCrud[T]
	gormdb.CreateCrud[T]
	gormdb.UpdateCrud[T]
	Deletes([]int64) (err error)
}