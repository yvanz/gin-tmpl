/*
@Date: 2021/12/7 14:30
@Author: yvan.zhang
@File : basic
*/

package gormdb

type BasicQuery struct {
	IDList  []int64  `json:"IdList"`  // id数组
	Fields  []string `json:"Fields"`  // 指定返回字段
	Keyword string   `json:"Keyword"` // 关键词(全局模糊搜索)
	Order   string   `json:"Order"`   // 排序，支持desc和asc
	Limit   int      `json:"Limit"`   // 分页条数
	Offset  int      `json:"Offset"`  // 分页偏移量
	Query   string   `json:"Query"`   // 自定义查询语句；使用RSQL语法，具体见: https://cmdb-web.ucloudadmin.com/docs/#api-appendix-query-syntax
}

type BasicCrud interface {
	GetList(q BasicQuery, model, list interface{}) (total int64, err error)
	GetByID(model interface{}, id int64) error
	GetOneByCon(con, model interface{}, args ...interface{}) error
	FindByCon(con, model interface{}, args ...interface{}) error
	Create(model interface{}) error
	UpdateWithMap(model interface{}, u map[string]interface{}) error
	Delete(model interface{}, hardDelete bool) error
}
