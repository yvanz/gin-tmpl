/*
@Date: 2021/1/12 下午2:37
@Author: yvan.zhang
@File : dbDemo
@Desc:
*/

package models

import (
	"gin-tmpl/pkg/tools"
	"gin-tmpl/pkg/xormmysql"
	"time"
)

type Demo struct {
	ID          int64     `json:"id" xorm:"'id' pk autoincr"`                 // 自增主键
	UserName    string    `json:"user_name" xorm:"'user_name'"`               // 用户名
	CreatedTime time.Time `json:"created_time" xorm:"'created_time' created"` // 创建时间
	UpdatedTime time.Time `json:"updated_time" xorm:"'updated_time' updated"` // 更新时间
	DeletedTime time.Time `json:"-" xorm:"'deleted_time' deleted"`
}

func (d *Demo) List(whereLike tools.WhereLike, searchCol, keyword string, pageNum, pageSize int, sortBy string, orderBy string) (int64, interface{}, error) {
	table := new(Demo)
	list := make([]Demo, 0)
	return tools.TaleList(whereLike, searchCol, keyword, pageNum, pageSize, sortBy, orderBy, table, &list)
}

func (d *Demo) GetByID(id int64) (*Demo, error) {
	p := new(Demo)
	has, err := xormmysql.My().Slave().Where(`id = ?`, id).Get(p)
	if err != nil {
		return nil, err
	}

	if has {
		return p, nil
	}

	return nil, nil
}

func (d *Demo) Add() error {
	_, err := xormmysql.My().Master().Insert(d)
	return err
}

func (d *Demo) DeleteByIDs(ids []int64) error {
	_, err := xormmysql.My().Master().In(`id`, ids).Delete(new(Demo))
	return err
}

func (d *Demo) Update(up map[string]interface{}) error {
	_, err := xormmysql.My().Master().Table(d).Where(`id = ?`, d.ID).Update(up)
	return err
}
