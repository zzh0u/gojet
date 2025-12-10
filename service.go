package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"gojet/config"
	"gojet/dao"
	"gojet/models"
	"gojet/router"
	"gojet/service"
	"gojet/util/jwt"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func server() {
	newService, err := NewService()
	if err != nil {
		slog.Error("创建服务失败", "错误", err)
		os.Exit(1)
	}

	if err := newService.Start(); err != nil {
		slog.Error("启动服务失败", "错误", err)
		os.Exit(1)
	}
}

// Service 应用服务结构体 - 保存所有服务组件
type Service struct {
	Config     *config.Config
	DB         *gorm.DB
	Logger     *slog.Logger
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

	// 根据配置创建日志处理器（统一使用JSON格式）
	var handler slog.Handler
	switch cfg.Logging.Output {
	case "file":
		// 只输出到文件
		fileHandler, err := createFileHandler(cfg.Logging.FilePath, logLevel)
		if err != nil {
			return nil, fmt.Errorf("创建文件日志处理器失败: %w", err)
		}
		handler = fileHandler
	case "both":
		// 同时输出到控制台和文件
		handler = slog.NewJSONHandler(io.MultiWriter(os.Stdout, fileWriter(cfg.Logging.FilePath)), &slog.HandlerOptions{
			Level: logLevel,
		})
	case "stdout":
		fallthrough
	default:
		// 默认输出到标准输出
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: logLevel,
		})
	}

	logger := slog.New(handler)
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
	service.InitService(userRepo)
	service.InitAuth(cfg)

	// 初始化示例数据
	logger.Info("正在初始化应用示例数据")
	if err := service.CreateInitialData(); err != nil {
		return nil, fmt.Errorf("初始化示例数据失败: %w", err)
	}

	// 创建 Gin 路由实例
	r := gin.New()

	// 配置 JWT 白名单路由（不需要 token 的公开接口）
	jwt.SkipRouter["login"] = true
	jwt.SkipRouter["register"] = true
	jwt.SkipRouter["health"] = true

	// 添加中间件
	r.Use(gin.Recovery())
	r.Use(loggingMiddleware(logger))

	// 设置 JWT secret、数据库连接和配置到 gin 上下文
	r.Use(func(c *gin.Context) {
		c.Set("jwt-secret", cfg.JWT.Secret)
		sqlDB, err := db.DB()
		if err == nil {
			c.Set("db", sqlDB)
		}
		c.Set("config", cfg)
		c.Next()
	})
	r.Use(jwt.Token)

	// 设置应用的所有路由
	router.SetupRoutes(r)

	// 创建 HTTP 服务器
	httpServer := &http.Server{
		Addr:    ":" + strconv.Itoa(cfg.App.Port),
		Handler: r,
	}

	return &Service{
		Config:     cfg,
		DB:         db,
		Logger:     logger,
		HTTPServer: httpServer,
	}, nil
}

func (s *Service) Start() error {
	s.Logger.Info("服务器启动中", "端口", s.Config.App.Port)
	return s.HTTPServer.ListenAndServe()
}

// Stop 关闭数据库连接
func (s *Service) Stop() error {
	s.Logger.Info("服务器正在关闭...")

	sqlDB, err := s.DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

// fileWriter 打开或创建日志文件
func fileWriter(filePath string) io.Writer {
	// 创建日志目录（如果不存在）
	dirPath := filepath.Dir(filePath)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		panic(fmt.Sprintf("创建日志目录失败: %v", err))
	}

	// 打开日志文件（追加模式）
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(fmt.Sprintf("打开日志文件失败: %v", err))
	}
	return file
}

// createFileHandler 创建文件日志处理器
func createFileHandler(filePath string, level slog.Level) (slog.Handler, error) {
	writer := fileWriter(filePath)
	handler := slog.NewJSONHandler(writer, &slog.HandlerOptions{
		Level: level,
	})
	return handler, nil
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
