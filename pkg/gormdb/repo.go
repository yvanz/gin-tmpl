/*
@Date: 2021/11/9 16:48
@Author: yvan.zhang
@File : repo
*/

package gormdb

import (
	"fmt"
	"strings"

	"github.com/yvanz/gin-tmpl/pkg/gadget"
	"github.com/yvanz/gin-tmpl/pkg/rsql"
	"gorm.io/gorm"
)

type CRUDImpl struct {
	Conn *gorm.DB
}

func NewCRUD(conn *gorm.DB) BasicCrud {
	return &CRUDImpl{Conn: conn}
}

// model and list must be a pointer
func (c *CRUDImpl) GetList(q BasicQuery, model, list interface{}) (total int64, err error) {
	db := c.Conn.Model(model)

	// 指定字段
	if len(q.Fields) > 0 {
		db.Select(strings.Join(q.Fields, ", "))
	}

	// 基于id查询
	if len(q.IDList) > 0 {
		db.Where("`Id` IN ?", q.IDList)
	}

	// 全局模糊
	if q.Keyword != "" {
		fields := gadget.FieldsFromModel(model, db, true).GetStringField()
		db.Scopes(KeywordGenerator(fields, q.Keyword))
	}

	parseColumnFunc := func(s string) string { return c.Conn.NamingStrategy.ColumnName("", s) }

	// 自定义查询条件
	if q.Query != "" {
		// 把传递过来的Query字段通过gorm的字段命名策略转义成数据库字段
		preParser, err := rsql.NewPreParser(rsql.MysqlPre(parseColumnFunc))
		if err != nil {
			return total, err
		}

		preStmt, values, err := preParser.ProcessPre(q.Query)
		if err != nil {
			return total, err
		}

		db.Where(preStmt, values...)
	}

	// 排序
	if q.Order != "" {
		orderKey := strings.Split(q.Order, " ")
		switch len(orderKey) {
		case 1:
			columnName := parseColumnFunc(orderKey[0])
			db.Order(columnName)
		case 2:
			columnName := parseColumnFunc(orderKey[0])
			order := strings.ToLower(orderKey[1])
			if order != "desc" && order != "asc" {
				order = "asc"
			}

			db.Order(fmt.Sprintf("%s %s", columnName, order))
		}
	}

	// 计数
	db = db.Count(&total)

	// 分页
	if q.Limit > 0 && q.Offset >= 0 {
		db.Limit(q.Limit).Offset(q.Offset)
	}

	err = db.Find(list).Error

	return total, err
}

// model must be a pointer
func (c *CRUDImpl) GetByID(model interface{}, id int64) error {
	q := c.Conn.First(model, id)
	if q.Error != nil {
		return q.Error
	}

	return nil
}

// conditions could be pointer of a model struct, map or string
// model must be a pointer
func (c *CRUDImpl) GetOneByCon(con, model interface{}, args ...interface{}) error {
	var q *gorm.DB

	if len(args) > 0 {
		q = c.Conn.Where(con, args...).First(model)
	} else {
		q = c.Conn.Where(con).First(model)
	}

	if q.Error != nil {
		return q.Error
	}

	return nil
}

// conditions could be pointer of a model struct, map or string
// model must be a pointer
func (c *CRUDImpl) FindByCon(con, model interface{}, args ...interface{}) error {
	var q *gorm.DB

	if len(args) > 0 {
		q = c.Conn.Where(con, args...).Find(model)
	} else {
		q = c.Conn.Where(con).Find(model)
	}

	if q.Error != nil {
		return q.Error
	}

	return nil
}

// model must be a pointer
func (c *CRUDImpl) Create(model interface{}) error {
	err := c.Conn.Create(model).Error
	return err
}

// model must be a pointer
func (c *CRUDImpl) UpdateWithMap(model interface{}, u map[string]interface{}) error {
	return c.Conn.Model(model).Updates(u).Error
}

// model must be a pointer
func (c *CRUDImpl) Delete(m interface{}, hardDelete bool) error {
	tx := c.Conn
	if hardDelete {
		tx = tx.Unscoped()
	}

	return tx.Delete(m).Error
}

func KeywordGenerator(columnList []string, keyword string) func(db *gorm.DB) *gorm.DB {
	var values []interface{}
	stmt := "1 AND ("

	length := len(columnList) - 1
	for i := range columnList {
		if columnList[i] == "id" {
			continue
		}

		stmt += fmt.Sprintf("`%s` LIKE ? ", columnList[i])
		values = append(values, fmt.Sprintf("%%%s%%", keyword))
		if i != length {
			stmt += "OR "
		}
	}

	stmt += ") AND 1"
	fs := func(db *gorm.DB) *gorm.DB {
		return db.Where(stmt, values...)
	}

	return fs
}