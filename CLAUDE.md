# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

GoJet 是一个基于 Gin 框架的 Go Web 应用模板，按 Go 模块的常见组织方式放置代码：命令在 `cmd/api`，不对外暴露的包在 `internal`。代码库使用中文进行注释和 API 响应。

## 开发命令

### 构建与运行

```bash
# 构建 Linux 可执行文件（不使用 Docker）
make build

# 本地运行（需要 PostgreSQL，不使用 Docker）
go run ./cmd/api

# 使用开发配置运行
APP_MODE=debug LOG_LEVEL=debug LOG_OUTPUT=stdout go run ./cmd/api

# 使用生产配置运行
APP_MODE=release LOG_LEVEL=info LOG_OUTPUT=both go run ./cmd/api
```

### Docker 命令

```bash
# 生产环境（使用仓库根目录的 docker-compose.yml，镜像定义在 deploy/Dockerfile）
make up-build      # 构建并启动服务
make up            # 启动服务
make down          # 停止服务
make logs          # 查看日志
make clean         # 清理容器和数据卷
make restart       # 重启服务
```

### 代码质量

```bash
# 安装并运行代码检查工具
make install-lint
make lint

# 格式化导入语句
make install-goimports
make goimports

# 生成 Swagger 文档
make install-swag
make swag
```

### 测试

目前项目中没有 `*_test.go`。`tests/` 只放本地联调命令，不参与编译。添加测试时，请遵循 Go 约定创建 `*_test.go` 文件。

```bash
# 运行所有测试
go test ./...

# 运行测试并显示详细输出
go test -v ./...

# 运行测试并生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 架构

### 目录

```text
cmd/api/main.go          # 加载配置并启动
cmd/api/app.go           # 组装数据库、缓存、路由
configs/config.yaml
deploy/Dockerfile
internal/config/
internal/user/           # 用户：模型、仓储、业务、HTTP
internal/auth/           # 登录
internal/health/
internal/infra/          # 支撑组件，目录本身不是包
  middleware/            # JWT、请求日志
  postgres/
  redis/
  logger/
  jwt/
  httputil/
  apperror/
tests/
```

`internal` 是 Go 工具链强制的：其他模块不能 import 其中的包。这个仓库是应用，没有需要对外发布的库，所以没有 `pkg`。组装写在 `cmd/api` 的 `main` 包里。

同一领域的类型和处理放在一个包中。用户相关代码都在 `internal/user`，登录相关代码都在 `internal/auth`。

### 启动流程

`cmd/api/app.go` 的 `newApp`：

1. `main` 从 `configs/config.yaml` 加载配置，环境变量覆盖
2. 初始化 JSON 日志
3. 设置 Gin 模式
4. 连接 PostgreSQL，并迁移 `user` 表
5. 连接 Redis；失败时继续启动，用户缓存关闭
6. 构造用户仓储和服务，并写入初始示例数据
7. 注册 Recovery、请求日志和 JWT 中间件
8. 调用 health、user、auth 的 `RegisterRoutes`，然后启动 HTTP 服务

### 添加新功能

1. 在 `internal/<name>/` 建一个包，把模型和 HTTP 处理放在一起
2. 在 `cmd/api/app.go` 里构造并调用 `RegisterRoutes`
3. 新表在 `newApp` 里 `AutoMigrate`
4. 响应走 `internal/infra/httputil`，业务错误走 `internal/infra/apperror`

### 数据库

- 使用 GORM v1.31.1 和 PostgreSQL
- 启动时自动迁移数据库表结构
- 通过环境变量配置连接（DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, DB_SSLMODE）

### API 模式

- **RESTful 端点** - 所有 API 位于 `/v1/` 路径下
- **请求/响应格式** - JSON 格式，统一响应结构：`{"code": 200, "message": "成功", "data": {...}}`
- **错误消息** - 中文错误消息，通过 `internal/infra/apperror/` 定义
- **健康检查** - `/v1/health` 端点返回应用状态和数据库连接状态
- **认证中间件** - JWT Access Token 验证，白名单路由可跳过验证（认证相关路由位于 `/v1/auth/` 路径下）
- **请求日志** - 自动记录所有 HTTP 请求的详细信息

**错误处理流程**：

1. repository 返回 `*apperror.Error`（记录不存在时为 404）
2. service 继续包装业务错误：`apperror.New(code, message)` 或 `apperror.Wrap`
3. handler 通过 `httputil.HandleError` 返回统一格式
4. Recovery 中间件捕获 panic 并返回 500 错误

**JWT 认证系统**：

- 双层 Token 架构 - Access Token（短期） + Refresh Token（长期）
- 密钥配置在 `configs/config.yaml` 的 `jwt.secret`
- Access Token 过期时间可配置（默认 24 小时），Refresh Token 过期时间可配置（默认 7 天）
- 白名单路由：路径最后一段为 `login` 或 `health`（即 `/v1/auth/login`、`/v1/health`）
- 用户创建只通过受保护接口 `POST /v1/user` 完成，需要令牌
- Token 存储在请求头：`Authorization: Bearer <access_token>`
- 用户信息通过 `c.Get("userid")` 和 `c.Get("username")` 在上下文中获取
- 登录接口返回双层 Token：Access Token 用于 API 调用，Refresh Token 用于长期保持登录状态

## 日志系统

项目使用 Go 标准库 `log/slog` 的结构化 JSON 日志，封装在 `internal/infra/logger`。

### 日志配置选项

- `LOG_LEVEL` - 日志级别 (debug/info/warn/error)
- `LOG_OUTPUT` - 输出目标 (stdout/file/both)
- `LOG_FILE_PATH` - 日志文件路径（当使用 file/both 输出时）

### 不同环境的日志行为

**开发模式** (`APP_MODE=debug`, config.yaml 中的默认值)：

- Debug 级别日志
- 仅输出到 stdout（便于实时查看）
- Gin debug 模式启用（显示路由信息）

**生产模式** (`APP_MODE=release`, Docker 中使用)：

- Info 级别日志（减少日志量）
- 输出到控制台和文件（./logs/app.log）
- Gin release 模式（性能优化，无调试信息）
- 日志通过 Docker 卷挂载持久化

### 日志特性

- **结构化 JSON 格式**：所有日志均为 JSON 格式，便于日志聚合系统解析
- **HTTP 请求日志**：自动记录所有 HTTP 请求的详细信息（方法、路径、状态码、耗时、客户端 IP）
- **错误上下文**：所有错误日志包含相关上下文（用户 ID、操作等）
- **文件持久化**：日志可写入文件，自动创建目录
- **多输出支持**：支持 stdout 仅输出、文件仅输出、或两者同时输出

### 日志文件管理

- 日志文件自动创建在 `LOG_FILE_PATH` 指定的目录
- 示例：`./logs/app.log` 在 `logs/` 目录创建日志
- Docker Compose 挂载 `./logs` 卷以实现容器重启后的日志持久化
- 日志文件自动追加（不覆盖）
- 将 `logs/` 添加到 `.gitignore` 以防止提交日志文件

## 开发工作流

### 本地开发

```bash
# 1. 安装依赖
go mod download

# 2. 启动本地 PostgreSQL（如果未运行）
# 使用 Docker 启动数据库
docker run -d --name gojet-postgres \
  -e POSTGRES_USER=gojet \
  -e POSTGRES_PASSWORD=gojet123 \
  -e POSTGRES_DB=gojet \
  -p 5432:5432 \
  postgres:15

# 3. 运行应用
go run ./cmd/api

# 或使用开发配置
APP_MODE=debug LOG_LEVEL=debug LOG_OUTPUT=stdout go run ./cmd/api
```

### 代码质量检查

```bash
make install-lint
make lint
make install-goimports
make goimports
go test ./...
```

## Docker 环境配置

镜像定义在 `deploy/Dockerfile`，编排文件是仓库根目录的 `docker-compose.yml`。

### 生产环境

- **应用模式**：`release`
- **日志级别**：`info`
- **日志输出**：`both`（控制台 + 文件）
- **配置挂载**：`./configs` → `/root/configs`
- **使用**：`make up-build`

### 环境文件

- **`.env`**：开发环境变量文件（docker-compose 自动加载）
- **`.env.prod`**：生产环境变量文件（创建用于生产环境覆盖）
- 这些文件不通过 git 跟踪（添加到 `.gitignore`）

## 重要说明

1. **日志系统** - 使用 `internal/infra/logger` 封装的 `log/slog` JSON 日志。支持 stdout/file/both，并由 `internal/infra/middleware` 记录 HTTP 请求。

2. **目录约定** - 入口在 `cmd/api`。不对外的包都在 `internal`。新增领域时加一个 `internal` 包，并在 `cmd/api/app.go` 里挂上路由。

3. **代码规范** - 代码库使用中文注释和 API 错误消息。添加新代码时保持这一约定。

4. **JWT 认证** - 已实现。双层 Token（Access Token + Refresh Token），登录接口返回两份令牌。

5. **错误处理** - 使用 `internal/infra/apperror` 与 `internal/infra/httputil`。

6. **测试覆盖** - 目前没有 `*_test.go`。`tests/` 是本地联调 shell 脚本，说明见 `tests/README.md`。

7. **配置管理** - 默认读取 `configs/config.yaml`，环境变量可覆盖。生产环境用环境变量设置数据库密码和 JWT 密钥。

8. **Docker 部署** - `deploy/Dockerfile` 构建 `./cmd/api`，并把 `configs/` 打进镜像。日志通过卷挂载持久化。
