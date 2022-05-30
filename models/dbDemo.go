/*
@Date: 2021/1/12 下午2:37
@Author: yvanz
@File : dbDemo
@Desc:
*/

package models

type Demo struct {
	Meta
	UserName string `json:"user_name" gorm:"column:user_name"` // 用户名
}

func (d *Demo) ColumnUserName() string {
	return "user_name"
}
