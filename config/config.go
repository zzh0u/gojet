package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/goccy/go-yaml"
)

// Config 应用配置结构体 - 包含所有配置项
type Config struct {
	App      AppConfig      `yaml:"app"`      // 应用配置
	Database DatabaseConfig `yaml:"database"` // 数据库配置
	Logging  LoggingConfig  `yaml:"logging"`  // 日志配置
}

// AppConfig 应用配置 - 定义应用的基本信息
type AppConfig struct {
	Name    string `yaml:"name"`    // 应用名称
	Version string `yaml:"version"` // 应用版本
	Port    int    `yaml:"port"`    // 服务端口
	Mode    string `yaml:"mode"`    // 运行模式 (debug/release/test)
}

// DatabaseConfig 数据库配置 - PostgreSQL 连接参数
type DatabaseConfig struct {
	Host     string `yaml:"host"`     // 数据库主机地址
	Port     int    `yaml:"port"`     // 数据库端口
	User     string `yaml:"user"`     // 数据库用户名
	Password string `yaml:"password"` // 数据库密码
	DBName   string `yaml:"dbname"`   // 数据库名称
	SSLMode  string `yaml:"sslmode"`  // SSL 连接模式
}

// LoggingConfig 日志配置 - 定义日志行为
type LoggingConfig struct {
	Level  string `yaml:"level"`  // 日志级别 (debug/info/warn/error)
	Format string `yaml:"format"` // 日志格式 (text/json)
	Output string `yaml:"output"` // 日志输出位置 (stdout/file)
}

// LoadConfig 加载配置 - 从 YAML 文件和环境变量读取配置
func LoadConfig(configPath string) (*Config, error) {
	config := &Config{}

	// 从 YAML 文件加载配置
	if configPath != "" {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("读取配置文件失败: %w", err)
		}

		if err := yaml.Unmarshal(data, config); err != nil {
			return nil, fmt.Errorf("解析配置文件失败: %w", err)
		}
	}

	// 使用环境变量覆盖配置文件中的设置
	config.overrideWithEnv()

	return config, nil
}

// overrideWithEnv 使用环境变量覆盖配置 - 优先级：环境变量 > 配置文件
func (c *Config) overrideWithEnv() {
	if val := os.Getenv("APP_NAME"); val != "" {
		c.App.Name = val
	}
	if val := os.Getenv("APP_PORT"); val != "" {
		if port, err := strconv.Atoi(val); err == nil {
			c.App.Port = port
		}
	}
	if val := os.Getenv("APP_MODE"); val != "" {
		c.App.Mode = val
	}

	// 数据库配置
	if val := os.Getenv("DB_HOST"); val != "" {
		c.Database.Host = val
	}
	if val := os.Getenv("DB_PORT"); val != "" {
		if port, err := strconv.Atoi(val); err == nil {
			c.Database.Port = port
		}
	}
	if val := os.Getenv("DB_USER"); val != "" {
		c.Database.User = val
	}
	if val := os.Getenv("DB_PASSWORD"); val != "" {
		c.Database.Password = val
	}
	if val := os.Getenv("DB_NAME"); val != "" {
		c.Database.DBName = val
	}
	if val := os.Getenv("DB_SSLMODE"); val != "" {
		c.Database.SSLMode = val
	}

	// 日志配置
	if val := os.Getenv("LOG_LEVEL"); val != "" {
		c.Logging.Level = val
	}
	if val := os.Getenv("LOG_FORMAT"); val != "" {
		c.Logging.Format = val
	}
}

// GetDSN 获取数据库连接字符串 - 构建 PostgreSQL DSN 连接串
func (db *DatabaseConfig) GetDSN() string {
	// 按照 PostgreSQL 的 DSN 格式拼接连接参数
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s ",
		db.Host, db.User, db.Password, db.DBName, db.Port, db.SSLMode)
}
