/*
@Date: 2021/1/12 下午2:37
@Author: yvan.zhang
@File : dbDemo
@Desc:
*/

package models

import (
	"time"

	"gorm.io/gorm"
)

type Demo struct {
	ID          int64          `json:"id" gorm:"column:id;primaryKey"`                         // 自增主键
	UserName    string         `json:"user_name" gorm:"column:user_name"`                      // 用户名
	CreatedTime time.Time      `json:"created_time" gorm:"column:created_time;autoCreateTime"` // 创建时间
	UpdatedTime time.Time      `json:"updated_time" gorm:"column:updated_time;autoUpdateTime"` // 更新时间
	DeletedTime gorm.DeletedAt `json:"-" gorm:"column:deleted_time"`
}

func (d *Demo) ColumnUserName() string {
	return "user_name"
}
