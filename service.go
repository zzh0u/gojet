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

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func serve() {
	service, err := NewService()
	if err != nil {
		slog.Error("❌ Failed to create service", "error", err)
		os.Exit(1)
	}
	if err := service.Start(); err != nil {
		slog.Error("❌ Failed to start service", "error", err)
		os.Exit(1)
	}
}

// Service represents the application service
type Service struct {
	Config     *config.Config
	DB         *gorm.DB
	Logger     *slog.Logger
	UserAPI    *v1api.UserAPI
	HTTPServer *http.Server
}

// NewService creates a new service instance
func NewService() (*Service, error) {
	// Load configuration
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	// Initialize structured logger
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

	// Set Gin mode
	gin.SetMode(cfg.App.Mode)

	// Initialize database
	db, err := gorm.Open(postgres.Open(cfg.Database.GetDSN()), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto migrate the schema
	if err := db.AutoMigrate(&models.User{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	// Initialize repository and services
	userRepo := dao.NewUserRepository(db)
	userService := service.NewUserService(userRepo, logger)
	userAPI := v1api.NewUserAPI(userService)

	// Initialize with sample data
	logger.Info("🚀 Initializing application with sample data")
	if err := userService.CreateInitialData(); err != nil {
		return nil, fmt.Errorf("failed to initialize sample data: %w", err)
	}

	// Create Gin router
	r := gin.New()

	// Add middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(loggingMiddleware(logger))

	// Health check endpoint with database check
	r.GET("/health", func(c *gin.Context) {
		// Check database connection
		sqlDB, err := db.DB()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":    "unhealthy",
				"error":     "database connection failed",
				"timestamp": time.Now().Format(time.RFC3339),
			})
			return
		}

		if err := sqlDB.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":    "unhealthy",
				"error":     "database ping failed",
				"timestamp": time.Now().Format(time.RFC3339),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().Format(time.RFC3339),
			"version":   cfg.App.Version,
		})
	})

	// Setup all application routes
	router.SetupRoutes(r, userAPI)

	// Create HTTP server
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

// Start starts the service
func (s *Service) Start() error {
	s.Logger.Info("🚀 Server starting", "port", s.Config.App.Port)
	s.Logger.Info("💚 Health check available", "url", fmt.Sprintf("http://localhost:%d/health", s.Config.App.Port))

	return s.HTTPServer.ListenAndServe()
}

// Stop gracefully stops the service
func (s *Service) Stop() error {
	s.Logger.Info("🛑 Server shutting down...")

	// Close database connection
	sqlDB, err := s.DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

// loggingMiddleware adds structured logging to requests
func loggingMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Log request details
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
