package main

import (
	"log/slog"
	"os"

	"gojet/internal/config"
)

func main() {
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		slog.Error("加载配置失败", "错误", err)
		os.Exit(1)
	}

	application, err := newApp(cfg)
	if err != nil {
		slog.Error("创建服务失败", "错误", err)
		os.Exit(1)
	}

	if err = application.start(); err != nil {
		slog.Error("启动服务失败", "错误", err)
		os.Exit(1)
	}
}
