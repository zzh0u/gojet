// Package postgres PostgreSQL 客户端。可复用，不依赖具体业务模块。
package postgres

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Open 使用 DSN 连接 PostgreSQL。
func Open(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}
	return db, nil
}
