// Package logger 全局 slog 日志封装，与具体业务无关。
package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

// Init 按级别和输出目标初始化默认 JSON 日志。
func Init(level, output, filePath string) error {
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	writer, err := logWriter(strings.ToLower(output), filePath)
	if err != nil {
		return err
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(writer, &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: true,
	})))
	return nil
}

func logWriter(output, filePath string) (io.Writer, error) {
	switch output {
	case "file", "both":
		fileW, err := fileWriter(filePath)
		if err != nil {
			return nil, fmt.Errorf("创建日志文件失败: %w", err)
		}
		if output == "file" {
			return fileW, nil
		}
		return io.MultiWriter(os.Stdout, fileW), nil
	default:
		return os.Stdout, nil
	}
}

func fileWriter(filePath string) (*os.File, error) {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("创建日志目录失败: %w", err)
	}
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("打开日志文件失败: %w", err)
	}
	return f, nil
}
