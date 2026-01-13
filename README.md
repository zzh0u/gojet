# gojet

[![Go Version](https://img.shields.io/badge/go-1.25.5-blue.svg)](https://golang.org/)
[![Gin Web Framework](https://img.shields.io/badge/gin-1.11.0-blue.svg)](https://github.com/gin-gonic/gin)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

gojet 是一个基于 Gin 框架的 Go Web 开发模板项目，包含简单的用户管理 RESTful API。集成了完整的 API 架构、数据库操作、健康检查、日志记录和 Docker 支持，适合学习和快速构建 Web 应用。

## 功能特性

- **完整的分层架构** - API、Service、DAO、Models 层分离
- **RESTFul API** - 符合 REST 规范的接口设计
- **数据库支持** - GORM + PostgreSQL，自动迁移
- **配置管理** - YAML + 环境变量双重配置
- **结构化日志** - JSON 格式日志，支持日志级别
- **健康检查** - HTTP 健康检查端点，包含数据库状态
- **请求追踪** - 自动记录 HTTP 请求日志
- **JWT 身份认证** - 基于 Token 的认证和授权，支持白名单路由
- **Docker 支持** - 完整的 Docker 和 Docker Compose 配置
- **代码质量工具** - Makefile 集成 golangci-lint 静态检查
- **API 文档支持** - 支持 Swagger 文档生成（需安装 swag 工具并运行 make swag）
- **统一响应处理** - 标准化的 API 响应格式和错误消息常量

## 技术栈

- [Gin](https://github.com/gin-gonic/gin) v1.11.0 - HTTP Web 框架
- [GORM](https://gorm.io) v1.31.1 - ORM 框架
- [PostgreSQL](https://www.postgresql.org/) - 关系型数据库
- [Validator](https://github.com/go-playground/validator) v10 - 参数验证
- [log/slog](https://pkg.go.dev/log/slog) - Go 标准库结构化日志
- [golang-jwt/jwt](https://github.com/golang-jwt/jwt) - JWT 令牌生成和验证（v5）
- [Docker](https://www.docker.com/) - 容器化部署

## 项目结构

```text
├── api/
│   └── v1api/            # HTTP API 处理层
├── service/              # 业务逻辑服务层
├── dao/                  # 数据访问对象层
├── models/               # 数据模型定义
├── config/               # 配置文件
├── router/               # 路由配置
├── util/                 # 工具类
│   ├── apperror/         # 业务错误定义
│   ├── jwt/              # JWT 工具（令牌生成、验证、中间件）
│   └── response/         # 统一响应处理
├── main.go               # 应用入口
├── service.go            # 服务启动和依赖注入逻辑
├── go.mod                # Go 模块定义
├── go.sum                # 依赖版本锁定
├── Dockerfile            # Docker 镜像构建
├── docker-compose.yml    # Docker 服务编排
├── Makefile              # 构建自动化脚本
└── .gitignore            # Git 忽略规则
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

**注意**：使用 `make goimports` 和 `make swag` 前，需要先安装相应工具。

### 代码规范

- 遵循 [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- 使用 golangci-lint 进行代码质量检查
- 提交代码前运行 `make lint` 确保无错误

### 添加新功能

参考以下步骤添加新模块：

1. **定义 Model** - 在 `models/` 目录创建数据模型
2. **创建 DAO** - 在 `dao/` 目录实现数据库操作
3. **实现 Service** - 在 `service/` 目录编写业务逻辑
4. **添加 API** - 在 `api/v1api/` 目录创建 HTTP 处理函数
5. **配置路由** - 在 `router/router.go` 中添加路由
6. **初始化组件** - 服务启动时自动通过 GORM 迁移创建表结构，service 层通过全局函数注册

## 许可证

MIT License
Copyright (c) 2025 gojet

## 相关项目

- [Gin Examples](https://github.com/gin-gonic/examples) - Gin 框架示例
- [GORM Guides](https://gorm.io/docs/) - GORM 使用指南
- [PostgreSQL](https://www.postgresql.org/) - 强大的开源数据库

**Happy Coding!**
