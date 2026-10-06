package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"gojet/internal/auth"
	"gojet/internal/config"
	"gojet/internal/health"
	"gojet/internal/infra/logger"
	"gojet/internal/infra/middleware"
	"gojet/internal/infra/postgres"
	"gojet/internal/infra/redis"
	"gojet/internal/user"

	"github.com/gin-gonic/gin"
	goredis "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// app 保存进程级依赖。
type app struct {
	config     *config.Config
	db         *gorm.DB
	redis      *redis.Client
	httpServer *http.Server
}

func newApp(cfg *config.Config) (*app, error) {
	if err := logger.Init(cfg.Logging.Level, cfg.Logging.Output, cfg.Logging.FilePath); err != nil {
		return nil, err
	}

	gin.SetMode(cfg.App.Mode)

	db, err := postgres.Open(cfg.Database.GetDSN())
	if err != nil {
		return nil, err
	}

	if err = db.AutoMigrate(&user.User{}); err != nil {
		return nil, fmt.Errorf("数据库迁移失败: %w", err)
	}

	redisClient, err := redis.NewClient(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		slog.Warn("Redis 连接失败，将跳过 Redis 相关功能", "错误", err)
		redisClient = nil
	} else {
		slog.Info("Redis 连接成功")
	}

	userRepo := user.NewUserRepository(db, cacheRaw(redisClient))
	if redisClient != nil {
		slog.Info("用户缓存已启用")
	}

	userSvc := user.New(userRepo)
	slog.Info("正在初始化应用示例数据")
	if err = userSvc.CreateInitialData(); err != nil {
		return nil, fmt.Errorf("初始化示例数据失败: %w", err)
	}

	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(middleware.Logging(slog.Default()))
	engine.Use(middleware.JWT(cfg.JWT.Secret, "login", "health"))

	health.New(dbConn(db), cfg.App.Version).RegisterRoutes(engine)
	user.NewHandler(userSvc).RegisterRoutes(engine)
	auth.NewHandler(auth.New(userRepo, cfg)).RegisterRoutes(engine)

	return &app{
		config: cfg,
		db:     db,
		redis:  redisClient,
		httpServer: &http.Server{
			Addr:    ":" + strconv.Itoa(cfg.App.Port),
			Handler: engine,
		},
	}, nil
}

func (a *app) start() error {
	slog.Info("服务器启动中", "端口", a.config.App.Port)
	return a.httpServer.ListenAndServe()
}

func (a *app) stop() error {
	slog.Info("服务器正在关闭...")

	conn, err := a.db.DB()
	if err != nil {
		return err
	}
	return conn.Close()
}

func cacheRaw(client *redis.Client) *goredis.Client {
	if client == nil {
		return nil
	}
	return client.Raw()
}

func dbConn(db *gorm.DB) *sql.DB {
	conn, err := db.DB()
	if err != nil {
		slog.Error("获取数据库连接失败", "错误", err)
		return nil
	}
	return conn
}
