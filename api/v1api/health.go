package v1api

import (
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	"gojet/config"
	"gojet/util/response"

	"github.com/gin-gonic/gin"
)

type HealthStatus struct {
	Status    string   `json:"status"`
	Timestamp string   `json:"timestamp"`
	Version   string   `json:"version"`
	Database  DBStatus `json:"database"`
}

type DBStatus struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

func HealthCheck(c *gin.Context) {

	db, exists := c.Get("db")
	if !exists {
		slog.Error("数据库连接未配置在 gin context 中")
		response.Error(c, http.StatusServiceUnavailable, "数据库连接未初始化")
		return
	}

	sqlDB, ok := db.(*sql.DB)
	if !ok {
		slog.Error("gin context 中的数据库连接类型错误")
		response.Error(c, http.StatusServiceUnavailable, "数据库连接类型错误")
		return
	}

	// 测试数据库连通性
	if err := sqlDB.Ping(); err != nil {
		slog.Error("数据库 Ping 失败", "error", err)
		response.Error(c, http.StatusServiceUnavailable, "数据库连接失败")
		return
	}

	// 从 gin context 获取配置
	cfg, exists := c.Get("config")
	if !exists {
		slog.Error("配置未设置")
		response.Error(c, http.StatusInternalServerError, "配置未初始化")
		return
	}

	appConfig, ok := cfg.(*config.Config)
	if !ok {
		slog.Error("gin context 中的配置类型错误")
		response.Error(c, http.StatusInternalServerError, "配置类型错误")
		return
	}

	health := HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now().Format(time.RFC3339),
		Version:   appConfig.App.Version,
		Database: DBStatus{
			Status: "healthy",
		},
	}

	response.Success(c, "", health)
}
