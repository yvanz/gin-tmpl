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

type GetListCrud interface {
	GetList(q BasicQuery, model, list interface{}) (total int64, err error)
}

type GetByIDCrud interface {
	GetByID(model interface{}, id int64) error
}

type GetByConCrud interface {
	GetOneByCon(con, model interface{}, args ...interface{}) error
}

type FindByConCrud interface {
	FindByCon(con, model interface{}, args ...interface{}) error
}

type CreateCrud interface {
	Create(model interface{}) error
}

type UpdateCrud interface {
	UpdateWithMap(model interface{}, u map[string]interface{}) error
}

type DeleteCrud interface {
	Delete(model interface{}, hardDelete bool) error
}
type BasicCrud interface {
	GetListCrud
	GetByIDCrud
	GetByConCrud
	FindByConCrud
	CreateCrud
	UpdateCrud
	DeleteCrud
}
