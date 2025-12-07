package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id"`         // 用户 ID，主键
	Name      string         `json:"name"`       // 用户姓名
	DeletedAt gorm.DeletedAt `json:"-"`          // 软删除时间，GORM 自动管理，不序列化到 JSON
	CreatedAt time.Time      `json:"created_at"` // 创建时间，GORM 自动管理
	UpdatedAt time.Time      `json:"updated_at"` // 更新时间，GORM 自动管理
}
