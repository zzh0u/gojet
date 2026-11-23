package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Config 应用配置结构体 - 包含所有配置项
type Config struct {
	App      AppConfig      `mapstructure:"app"`      // 应用配置
	Database DatabaseConfig `mapstructure:"database"` // 数据库配置
	Logging  LoggingConfig  `mapstructure:"logging"`  // 日志配置
	JWT      JWTConfig      `mapstructure:"jwt"`      // JWT 配置
	Redis    RedisConfig    `mapstructure:"redis"`    // Redis 配置
}

// AppConfig 应用配置 - 定义应用的基本信息
type AppConfig struct {
	Name    string `mapstructure:"name"`    // 应用名称
	Version string `mapstructure:"version"` // 应用版本
	Port    int    `mapstructure:"port"`    // 服务端口
	Mode    string `mapstructure:"mode"`    // 运行模式 (debug/release/test)
}

// DatabaseConfig 数据库配置 - PostgreSQL 连接参数
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`     // 数据库主机地址
	Port     int    `mapstructure:"port"`     // 数据库端口
	User     string `mapstructure:"user"`     // 数据库用户名
	Password string `mapstructure:"password"` // 数据库密码
	DBName   string `mapstructure:"dbname"`   // 数据库名称
	SSLMode  string `mapstructure:"sslmode"`  // SSL 连接模式
}

// LoggingConfig 日志配置 - 定义日志行为
type LoggingConfig struct {
	Level    string `mapstructure:"level"`     // 日志级别 (debug/info/warn/error)
	Output   string `mapstructure:"output"`    // 日志输出位置 (stdout/file/both)
	FilePath string `mapstructure:"file_path"` // 日志文件路径
}

// JWTConfig JWT 配置 - 定义 JWT token 相关参数
type JWTConfig struct {
	Secret            string `mapstructure:"secret"`              // JWT 签名密钥
	AccessExpireHours int    `mapstructure:"access_expire_hours"` // Access Token 过期时间（小时）
	RefreshExpireDays int    `mapstructure:"refresh_expire_days"` // Refresh Token 过期时间（天）
}

// RedisConfig Redis 配置 - 定义 Redis 连接参数
type RedisConfig struct {
	Host     string `mapstructure:"host"`     // Redis 主机地址
	Port     int    `mapstructure:"port"`     // Redis 端口
	Password string `mapstructure:"password"` // Redis 密码（无密码留空）
	DB       int    `mapstructure:"db"`       // Redis 数据库编号
}

// envBindings 配置键与现有环境变量名的对应关系。
var envBindings = []struct {
	key string
	env string
}{
	{"app.name", "APP_NAME"},
	{"app.port", "APP_PORT"},
	{"app.mode", "APP_MODE"},
	{"database.host", "DB_HOST"},
	{"database.port", "DB_PORT"},
	{"database.user", "DB_USER"},
	{"database.password", "DB_PASSWORD"},
	{"database.dbname", "DB_NAME"},
	{"database.sslmode", "DB_SSLMODE"},
	{"logging.level", "LOG_LEVEL"},
	{"logging.output", "LOG_OUTPUT"},
	{"logging.file_path", "LOG_FILE_PATH"},
	{"jwt.secret", "JWT_SECRET"},
	{"jwt.access_expire_hours", "JWT_ACCESS_EXPIRE_HOURS"},
	{"jwt.refresh_expire_days", "JWT_REFRESH_EXPIRE_DAYS"},
	{"redis.host", "REDIS_HOST"},
	{"redis.port", "REDIS_PORT"},
	{"redis.password", "REDIS_PASSWORD"},
	{"redis.db", "REDIS_DB"},
}

// LoadConfig 加载配置。Viper 先读 YAML，再按 envBindings 用环境变量覆盖。
func LoadConfig(configPath string) (*Config, error) {
	v := viper.New()
	v.SetConfigType("yaml")
	// 空环境变量视为未设置，避免 REDIS_PASSWORD= 清掉文件里的密码。
	v.AllowEmptyEnv(false)

	for _, binding := range envBindings {
		if err := v.BindEnv(binding.key, binding.env); err != nil {
			return nil, fmt.Errorf("绑定环境变量 %s 失败: %w", binding.env, err)
		}
	}

	if configPath != "" {
		if _, err := os.Stat(configPath); err != nil {
			return nil, fmt.Errorf("读取配置文件失败: %w", err)
		}
		v.SetConfigFile(configPath)
		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("解析配置文件失败: %w", err)
		}
	}

	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
	}
	return cfg, nil
}

// GetDSN 获取数据库连接字符串 - 构建 PostgreSQL DSN 连接串
func (db *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s ",
		db.Host, db.User, db.Password, db.DBName, db.Port, db.SSLMode)
}
