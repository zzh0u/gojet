package health

import (
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	"gojet/internal/infra/httputil"

	"github.com/gin-gonic/gin"
)

// Handler 提供应用健康检查。
type Handler struct {
	db      *sql.DB
	version string
}

// New 创建健康检查处理器。
func New(db *sql.DB, version string) *Handler {
	return &Handler{db: db, version: version}
}

// RegisterRoutes 注册健康检查路由。
func (h *Handler) RegisterRoutes(router gin.IRouter) {
	router.GET("/v1/health", h.Check)
}

type healthStatus struct {
	Status    string   `json:"status"`
	Timestamp string   `json:"timestamp"`
	Version   string   `json:"version"`
	Database  dbStatus `json:"database"`
}

type dbStatus struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// Check 返回应用状态和数据库连通性。
func (h *Handler) Check(c *gin.Context) {
	if h.db == nil {
		slog.Error("数据库连接未配置")
		httputil.Error(c, http.StatusServiceUnavailable, "数据库连接未初始化")
		return
	}

	if err := h.db.Ping(); err != nil {
		slog.Error("数据库 Ping 失败", "error", err)
		httputil.Error(c, http.StatusServiceUnavailable, "数据库连接失败")
		return
	}

	httputil.Success(c, "", healthStatus{
		Status:    "healthy",
		Timestamp: time.Now().Format(time.RFC3339),
		Version:   h.version,
		Database: dbStatus{
			Status: "healthy",
		},
	})
}
