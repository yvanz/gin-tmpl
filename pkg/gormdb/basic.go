/*
@Date: 2021/12/7 14:30
@Author: yvanz
@File : basic
*/

package gormdb

type BasicQuery struct {
	Fields  []string `json:"Fields"`  // 指定返回字段
	Keyword string   `json:"Keyword"` // 关键词(全局模糊搜索)
	Order   string   `json:"Order"`   // 排序，支持desc和asc
	Query   string   `json:"Query"`   // 自定义查询语句；使用RSQL语法
	IDList  []int64  `json:"IdList"`  // id数组
	Limit   int      `json:"Limit"`   // 分页条数
	Offset  int      `json:"Offset"`  // 分页偏移量
}

type GetListCrud[T any] interface {
	GetList(q BasicQuery, model *T, list *[]T) (total int64, err error)
}

type GetByIDCrud[T any] interface {
	GetByID(model *T, id int64) error
}

type GetOneByConCrud[T any] interface {
	GetOneByCon(con any, model *T, args ...any) error
}

type FindByConCrud[T any] interface {
	FindByCon(con any, model *T, args ...any) error
}

type CreateCrud[T any] interface {
	Create(model *T) error
}

type UpdateCrud[T any] interface {
	UpdateWithMap(model *T, u map[string]any) error
}

type DeleteCrud[T any] interface {
	Delete(model *T, hardDelete bool) error
}

type BasicCrud[T any] interface {
	GetListCrud[T]
	GetByIDCrud[T]
	GetOneByConCrud[T]
	FindByConCrud[T]
	CreateCrud[T]
	UpdateCrud[T]
	DeleteCrud[T]
}
