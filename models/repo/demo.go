package repo

import "github.com/yvanz/gin-tmpl/pkg/gormdb"

type DemoRepo interface {
	gormdb.GetListCrud
	gormdb.GetByIDCrud
	gormdb.CreateCrud
	gormdb.UpdateCrud
	Deletes([]int64) (err error)
}
