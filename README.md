# gojet

[![Go Version](https://img.shields.io/badge/go-1.25.5-blue.svg)](https://golang.org/)
[![Gin Web Framework](https://img.shields.io/badge/gin-1.11.0-blue.svg)](https://github.com/gin-gonic/gin)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![zread](https://img.shields.io/badge/Ask_Zread-_.svg?style=flat&color=00b0aa&labelColor=000000&logo=data%3Aimage%2Fsvg%2Bxml%3Bbase64%2CPHN2ZyB3aWR0aD0iMTYiIGhlaWdodD0iMTYiIHZpZXdCb3g9IjAgMCAxNiAxNiIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPHBhdGggZD0iTTQuOTYxNTYgMS42MDAxSDIuMjQxNTZDMS44ODgxIDEuNjAwMSAxLjYwMTU2IDEuODg2NjQgMS42MDE1NiAyLjI0MDFWNC45NjAxQzEuNjAxNTYgNS4zMTM1NiAxLjg4ODEgNS42MDAxIDIuMjQxNTYgNS42MDAxSDQuOTYxNTZDNS4zMTUwMiA1LjYwMDEgNS42MDE1NiA1LjMxMzU2IDUuNjAxNTYgNC45NjAxVjIuMjQwMUM1LjYwMTU2IDEuODg2NjQgNS4zMTUwMiAxLjYwMDEgNC45NjE1NiAxLjYwMDFaIiBmaWxsPSIjZmZmIi8%2BCjxwYXRoIGQ9Ik00Ljk2MTU2IDEwLjM5OTlIMi4yNDE1NkMxLjg4ODEgMTAuMzk5OSAxLjYwMTU2IDEwLjY4NjQgMS42MDE1NiAxMS4wMzk5VjEzLjc1OTlDMS42MDE1NiAxNC4xMTM0IDEuODg4MSAxNC4zOTk5IDIuMjQxNTYgMTQuMzk5OUg0Ljk2MTU2QzUuMzE1MDIgMTQuMzk5OSA1LjYwMTU2IDE0LjExMzQgNS42MDE1NiAxMy43NTk5VjExLjAzOTlDNS42MDE1NiAxMC42ODY0IDUuMzE1MDIgMTAuMzk5OSA0Ljk2MTU2IDEwLjM5OTlaIiBmaWxsPSIjZmZmIi8%2BCjxwYXRoIGQ9Ik0xMy43NTg0IDEuNjAwMUgxMS4wMzg0QzEwLjY4NSAxLjYwMDEgMTAuMzk4NCAxLjg4NjY0IDEwLjM5ODQgMi4yNDAxVjQuOTYwMUMxMC4zOTg0IDUuMzEzNTYgMTAuNjg1IDUuNjAwMSAxMS4wMzg0IDUuNjAwMUgxMy43NTg0QzE0LjExMTkgNS42MDAxIDE0LjM5ODQgNS4zMTM1NiAxNC4zOTg0IDQuOTYwMVYyLjI0MDFDMTQuMzk4NCAxLjg4NjY0IDE0LjExMTkgMS42MDAxIDEzLjc1ODQgMS42MDAxWiIgZmlsbD0iI2ZmZiIvPgo8cGF0aCBkPSJNNCAxMkwxMiA0TDQgMTJaIiBmaWxsPSIjZmZmIi8%2BCjxwYXRoIGQ9Ik00IDEyTDEyIDQiIHN0cm9rZT0iI2ZmZiIgc3Ryb2tlLXdpZHRoPSIxLjUiIHN0cm9rZS1saW5lY2FwPSJyb3VuZCIvPgo8L3N2Zz4K&logoColor=ffffff)](https://zread.ai/zzh0u/gojet)

gojet 是一个基于 Gin 框架的 Go Web 开发模板项目，解决以下问题：

- **基础架构代码** - 新建 Go Web 项目需要重新配置数据库连接、日志系统、路由、中间件等基础设施
- **标准化项目结构** - 命令放在 `cmd`，不对外的包放在 `internal`
- **环境配置** - 自动处理配置文件、环境变量、Docker 部署等运维相关设置
- **认证系统实现** - JWT 认证、Token 刷新、权限控制等安全功能实现

适合以下场景：
- 学习 Go Web 开发的目录约定和模块化组织
- 快速启动新项目，无需从零配置基础设施
- 作为团队项目模板，统一代码结构和开发规范

## 解决方案

本项目提供开箱即用的基础架构，包含：

- **按领域分包** - 用户、登录各自一个包，由 `cmd/api` 组装后启动
- **RESTFul API** - 符合 REST 规范的接口设计
- **数据库支持** - GORM + PostgreSQL，自动迁移
- **配置管理** - YAML + 环境变量双重配置
- **结构化日志** - JSON 格式日志，支持日志级别
- **健康检查** - HTTP 健康检查端点，包含数据库状态
- **请求追踪** - 自动记录 HTTP 请求日志
- **JWT 身份认证** - 基于双层 Token 的认证系统
- **Docker 支持** - 完整的 Docker 和 Docker Compose 配置
- **代码质量工具** - Makefile 集成 golangci-lint 静态检查
- **API 文档支持** - 支持 Swagger 文档生成（需安装 swag 工具并运行 make swag）
- **统一响应处理** - 标准化的 API 响应格式和错误消息常量

## 项目结构

这是一个应用模块，不是给外部 import 的库。目录按 [Organizing a Go module](https://go.dev/doc/modules/layout) 组织：命令放在 `cmd/api`，其余包放在 `internal`，这样别的模块不能引用它们。

启动链路：`cmd/api/main.go` 加载配置，`cmd/api/app.go` 连接数据库和 Redis、挂上路由并启动。

```text
cmd/api/                # 唯一的 main：加载配置、组装依赖、启动 HTTP
configs/                # YAML 配置；由环境变量覆盖
deploy/Dockerfile       # 镜像构建
internal/               # 本模块私有包
  config/               # 配置加载
  user/                 # 用户模型、存储、业务和 HTTP
  auth/                 # 登录
  health/               # 健康检查
  infra/                # 支撑组件，目录本身不是包
    middleware/         # JWT、请求日志
    postgres/           # PostgreSQL 客户端
    redis/              # Redis 客户端
    logger/             # JSON 日志
    jwt/                # JWT 签发与解析
    httputil/           # 统一 HTTP 响应
    apperror/           # 业务错误
tests/                  # 本地联调 shell 脚本；不参与编译
```

## 开发指南

### 使用 Makefile

项目提供了丰富的 Makefile 命令，简化开发流程：

```bash
# 代码质量工具
make install-lint        # 安装 golangci-lint 代码检查工具
make lint               # 运行代码静态检查
make install-goimports  # 安装 goimports 工具（格式化 Go 导入语句）
make goimports          # 格式化代码导入语句
make install-swag       # 安装 swag 工具（生成 Swagger 文档）
make swag               # 生成 Swagger 文档

# 构建和运行
make build              # 编译 Linux 可执行文件

# Docker Compose 命令
make up                 # 启动 Docker Compose 服务
make up-build           # 构建并启动 Docker Compose 服务
make down               # 停止 Docker Compose 服务
make logs               # 查看 Docker Compose 实时日志
make restart            # 重启 Docker Compose 服务
make clean              # 清理 Docker Compose 容器和数据卷
# 查看服务状态可使用: docker-compose ps 或 docker ps
```

> 使用 `make goimports` 和 `make swag` 前，需要先安装相应工具。

### 代码规范

- 遵循 [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- 使用 golangci-lint 进行代码质量检查
- 提交代码前运行 `make lint` 确保无错误

### 添加新功能

1. 在 `internal/<name>/` 增加一个包，模型和 HTTP 处理放在同一个包里
2. 在 `cmd/api/app.go` 里构造它，并调用 `RegisterRoutes`
3. 如有新表，在 `newApp` 里对该模型执行 `AutoMigrate`
4. HTTP 响应使用 `internal/infra/httputil`，业务错误使用 `internal/infra/apperror`

本地接口检查见 [tests/README.md](tests/README.md)。

## 许可证

MIT License
Copyright (c) 2025 gojet

## 相关项目

- [Gin Examples](https://github.com/gin-gonic/examples) - Gin 框架示例
- [GORM Guides](https://gorm.io/docs/) - GORM 使用指南
- [PostgreSQL](https://www.postgresql.org/) - 强大的开源数据库

**Happy Coding!**
