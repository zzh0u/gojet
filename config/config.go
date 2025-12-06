package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/goccy/go-yaml"
)

type Config struct {
	App      AppConfig      `yaml:"app"`
	Database DatabaseConfig `yaml:"database"`
	Logging  LoggingConfig  `yaml:"logging"`
}

type AppConfig struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	Port    int    `yaml:"port"`
	Mode    string `yaml:"mode"`
}

type DatabaseConfig struct {
	Host            string        `yaml:"host"`
	Port            int           `yaml:"port"`
	User            string        `yaml:"user"`
	Password        string        `yaml:"password"`
	DBName          string        `yaml:"dbname"`
	SSLMode         string        `yaml:"sslmode"`
	Timezone        string        `yaml:"timezone"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
}

type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
	Output string `yaml:"output"`
}

// LoadConfig loads configuration from file and environment variables
func LoadConfig(configPath string) (*Config, error) {
	config := &Config{}

	// Load from YAML file
	if configPath != "" {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}

		if err := yaml.Unmarshal(data, config); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %w", err)
		}
	}

	// Override with environment variables
	config.overrideWithEnv()

	return config, nil
}

// overrideWithEnv overrides config values with environment variables
func (c *Config) overrideWithEnv() {
	// App config
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

	// Database config
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

	// Logging config
	if val := os.Getenv("LOG_LEVEL"); val != "" {
		c.Logging.Level = val
	}
	if val := os.Getenv("LOG_FORMAT"); val != "" {
		c.Logging.Format = val
	}
}

// GetDSN returns the database connection string
func (db *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		db.Host, db.User, db.Password, db.DBName, db.Port, db.SSLMode, db.Timezone)
}
