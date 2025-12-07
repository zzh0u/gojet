package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"gojet/api/v1api"
	"gojet/config"
	"gojet/dao"
	"gojet/models"
	"gojet/router"
	"gojet/service"
	"gojet/util/response"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func serve() {
	service, err := NewService()
	if err != nil {
		slog.Error("❌ 创建服务失败", "错误", err)
		os.Exit(1)
	}

	if err := service.Start(); err != nil {
		slog.Error("❌ 启动服务失败", "错误", err)
		os.Exit(1)
	}
}

// Service 应用服务结构体 - 保存所有服务组件
type Service struct {
	Config     *config.Config
	DB         *gorm.DB
	Logger     *slog.Logger
	UserAPI    *v1api.UserAPI
	HTTPServer *http.Server
}

func NewService() (*Service, error) {
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		return nil, fmt.Errorf("加载配置失败: %w", err)
	}

	var logLevel slog.Level
	switch cfg.Logging.Level {
	case "debug":
		logLevel = slog.LevelDebug
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))
	slog.SetDefault(logger)

	gin.SetMode(cfg.App.Mode)

	// 初始化数据库连接
	db, err := gorm.Open(postgres.Open(cfg.Database.GetDSN()), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	// 自动迁移数据库表结构
	if err := db.AutoMigrate(&models.User{}); err != nil {
		return nil, fmt.Errorf("数据库迁移失败: %w", err)
	}

	// 初始化数据访问层和业务层
	userRepo := dao.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userAPI := v1api.NewUserAPI(userService)

	// 初始化示例数据
	logger.Info("🚀 正在初始化应用示例数据")
	if err := userService.CreateInitialData(); err != nil {
		return nil, fmt.Errorf("初始化示例数据失败: %w", err)
	}

	// 创建 Gin 路由实例
	r := gin.New()

	// 添加中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(loggingMiddleware(logger))

	// 健康检查接口 - 检查数据库连接状态
	r.GET("/health", func(c *gin.Context) {
		sqlDB, err := db.DB()
		if err != nil {
			response.Error(c, 503, "数据库连接失败")
			return
		}

		// 测试数据库连通性
		if err := sqlDB.Ping(); err != nil {
			response.Error(c, 503, "数据库 Ping 失败")
			return
		}

		response.Success(c, "", gin.H{
			"status":    "healthy",
			"timestamp": time.Now().Format(time.RFC3339),
			"version":   cfg.App.Version,
		})
	})

	// 设置应用的所有路由
	router.SetupRoutes(r, userAPI)

	// 创建 HTTP 服务器
	httpServer := &http.Server{
		Addr:    ":" + strconv.Itoa(cfg.App.Port),
		Handler: r,
	}

	return &Service{
		Config:     cfg,
		DB:         db,
		Logger:     logger,
		UserAPI:    userAPI,
		HTTPServer: httpServer,
	}, nil
}

func (s *Service) Start() error {
	s.Logger.Info("🚀 服务器启动中", "端口", s.Config.App.Port)
	s.Logger.Info("💚 健康检查可用", "地址", fmt.Sprintf("http://localhost:%d/health", s.Config.App.Port))

	return s.HTTPServer.ListenAndServe()
}

// Stop 关闭数据库连接
func (s *Service) Stop() error {
	s.Logger.Info("🛑 服务器正在关闭...")

	sqlDB, err := s.DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

// loggingMiddleware 请求日志中间件 - 记录 HTTP 请求详情
func loggingMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		// 记录请求详情
		duration := time.Since(start)
		logger.Info("HTTP Request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"duration", duration.String(),
			"user_agent", c.Request.UserAgent(),
			"ip", c.ClientIP(),
		)
	}
}
