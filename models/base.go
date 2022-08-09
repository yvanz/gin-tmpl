package models

import (
	"time"

	"gorm.io/gorm"
)

type Meta struct { //nolint:govet
	ID          int64          `json:"Id" gorm:"column:id;primaryKey"`
	CreatedTime time.Time      `json:"CreatedTime" gorm:"column:created_time;autoCreateTime"`
	UpdatedTime time.Time      `json:"UpdatedTime" gorm:"column:updated_time;autoUpdateTime"` // 更新时间
	DeletedTime gorm.DeletedAt `json:"-" gorm:"column:deleted_time"`
}
