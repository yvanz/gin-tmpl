package models

import (
	"time"

	"gorm.io/gorm"
)

type Meta struct { //nolint:govet
	ID          int64          `json:"id" gorm:"column:id;primaryKey"`
	CreatedTime time.Time      `json:"created_time" gorm:"column:created_time;autoCreateTime"`
	UpdatedTime time.Time      `json:"updated_time" gorm:"column:updated_time;autoUpdateTime"` // 更新时间
	DeletedTime gorm.DeletedAt `json:"-" gorm:"column:deleted_time"`
}
