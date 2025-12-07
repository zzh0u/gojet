# gojet

[![Go Version](https://img.shields.io/badge/go-1.25.5-blue.svg)](https://golang.org/)
[![Gin Web Framework](https://img.shields.io/badge/gin-1.11.0-blue.svg)](https://github.com/gin-gonic/gin)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

gojet 是一个基于 Gin 框架的 Go Web 开发模板项目，包含简单的用户管理 RESTful API。集成了完整的 API 架构、数据库操作、健康检查、日志记录和 Docker 支持，适合学习和快速构建 Web 应用。

## 功能特性

- ✅ **完整的分层架构** - API、Service、DAO、Models 层分离
- ✅ **RESTful API** - 符合 REST 规范的接口设计
- ✅ **数据库支持** - GORM + PostgreSQL，自动迁移
- ✅ **配置管理** - YAML + 环境变量双重配置
- ✅ **结构化日志** - JSON 格式日志，支持日志级别
- ✅ **健康检查** - HTTP 健康检查端点，包含数据库状态
- ✅ **请求追踪** - 自动记录 HTTP 请求日志
- ✅ **Docker 支持** - 完整的 Docker 和 Docker Compose 配置
- ✅ **代码质量工具** - Makefile 集成 golangci-lint 静态检查
- ✅ **API 文档支持** - 支持 Swagger 文档生成
- ✅ **统一响应处理** - 标准化的 API 响应格式和错误消息常量

## 技术栈

- [Gin](https://github.com/gin-gonic/gin) v1.11.0 - HTTP Web 框架
- [GORM](https://gorm.io) v1.31.1 - ORM 框架
- [PostgreSQL](https://www.postgresql.org/) - 关系型数据库
- [Validator](https://github.com/go-playground/validator) v10 - 参数验证
- [log/slog](https://pkg.go.dev/log/slog) - Go 标准库结构化日志
- [Docker](https://www.docker.com/) - 容器化部署

## 快速开始

### 前置要求

- Go 1.25.5 或更高版本
- PostgreSQL 15+
- Make

### 方式一：使用 Docker（推荐）

```bash
# 克隆项目
git clone <repository-url>
cd gojet

# 启动服务
make docker-up-build    # 构建并启动服务
make docker-up          # 启动服务
make docker-down        # 停止服务
make docker-logs        # 查看实时日志
make docker-ps          # 查看服务状态
make docker-restart     # 重启服务
make docker-clean       # 清理容器和数据卷

# 服务将在 http://localhost:8080 运行
# 数据库将在 localhost:5432 运行
```

### 方式二：本地运行

```bash
# 克隆项目
git clone <repository-url>
cd gojet

# 安装依赖
go mod download

# 运行应用
go run main.go

# 或使用 make
make build && ./main
```

## 项目结构

```text
├── api/                  # HTTP API 处理层
├── service/              # 业务逻辑服务层
├── dao/                  # 数据访问对象层
├── models/               # 数据模型定义
├── config/               # 配置文件
├── util/                 # 统一响应处理工具
├── main.go               # 应用入口
├── service.go            # 服务启动逻辑
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
make docker-up          # 启动 Docker Compose 服务
make docker-up-build    # 构建并启动 Docker Compose 服务
make docker-down        # 停止 Docker Compose 服务
make docker-logs        # 查看 Docker Compose 实时日志
make docker-ps          # 查看 Docker Compose 服务状态
make docker-restart     # 重启 Docker Compose 服务
make docker-clean       # 清理 Docker Compose 容器和数据卷
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
4. **添加 API** - 在 `api/` 目录创建 HTTP 处理函数
5. **配置路由** - 在 `router/router.go` 中添加路由
6. **注册组件** - 在 `service.go` 的 `NewService()` 中初始化

## Docker 部署

### 构建 Docker 镜像

```bash
# 使用 docker build
docker build -t gojet:latest .

# 查看镜像
docker images | grep gojet
```

### 运行容器

```bash
# 运行单个容器（需要外部数据库）
docker run -d \
  --name gojet \
  -p 8080:8080 \
  -e DB_HOST=your-db-host \
  -e DB_USER=your-db-user \
  -e DB_PASSWORD=your-db-password \
  -e DB_NAME=gojet \
  gojet:latest

# 查看日志
docker logs -f gojet
```

### 使用 Docker Compose

```bash
# 启动服务（包含 PostgreSQL）
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down

# 停止并删除数据卷
docker-compose down -v
```

### Docker Compose 配置

`docker-compose.yml` 包含两个服务：

- **gojet**: Go Web 应用
- **postgres**: PostgreSQL 数据库

数据库数据会持久化到 `postgres_data` 卷中。

## 许可证

MIT License
Copyright (c) 2025 gojet

## 相关项目

- [Gin Examples](https://github.com/gin-gonic/examples) - Gin 框架示例
- [GORM Guides](https://gorm.io/docs/) - GORM 使用指南
- [Go Project Layout](https://github.com/golang-standards/project-layout) - Go 项目布局标准

## 致谢

- [Gin](https://github.com/gin-gonic/gin) - 高性能 HTTP Web 框架
- [GORM](https://gorm.io) - 优秀的 ORM 框架
- [PostgreSQL](https://www.postgresql.org/) - 强大的开源数据库

如有问题或建议，请提交 Issue。

**Happy Coding!**
