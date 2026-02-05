# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

GoJet 是一个基于 Gin 框架的 Go Web 应用模板，采用清晰的分层架构（API/Service/DAO/Models）并集成 PostgreSQL 数据库。代码库使用中文进行注释和 API 响应。

## 开发命令

### 构建与运行

```bash
# 构建 Linux 可执行文件（不使用 Docker）
make build

# 本地运行（需要 PostgreSQL，不使用 Docker）
go run main.go

# 使用开发配置运行
APP_MODE=debug LOG_LEVEL=debug LOG_OUTPUT=stdout go run main.go

# 使用生产配置运行
APP_MODE=release LOG_LEVEL=info LOG_OUTPUT=both go run main.go
```

### Docker 命令

```bash
# 生产环境（使用 docker-compose.yml）
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

目前项目中不存在测试文件。添加测试时，请遵循 Go 约定创建 `*_test.go` 文件。

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./api/v1api

# 运行测试并显示详细输出
go test -v ./...

# 运行测试并生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 架构

### 分层结构

- **models/** - 数据模型定义，包含 GORM 标签和验证标签
- **dao/** - 数据库操作，继承 `dao.BaseRepository` 模式
- **service/** - 业务逻辑实现，通过全局函数注册
- **api/v1api/** - HTTP 处理器，包含参数验证和统一响应格式化
- **router/** - 路由定义，包含 JWT 中间件和白名单配置
- **config/** - YAML 配置文件，支持环境变量覆盖
- **util/** - 响应处理、错误工具和 JWT 中间件

### 关键文件

- `main.go` - 应用入口点，调用 `server()` 函数
- `service.go` - **依赖注入容器**，初始化所有服务组件
- `router/router.go` - 路由设置和中间件配置
- `config/config.yaml` - 默认配置文件
- `config/config.go` - 配置结构定义和加载逻辑
- `util/response/response.go` - 统一响应处理
- `util/apperror/error.go` - 业务错误定义
- `util/jwt/jwt.go` - JWT 中间件和工具

### Service 启动流程

**`service.go` 中的 `newService()` 函数执行顺序**：
1. **加载配置** - 从 `config/config.yaml` 加载，环境变量覆盖
2. **初始化日志** - 根据配置创建 JSON 格式日志处理器
3. **设置 Gin 模式** - `debug` 或 `release` 模式
4. **连接数据库** - PostgreSQL 连接，自动迁移表结构
5. **初始化 DAO 层** - 创建数据访问对象
6. **注册 Service 层** - 通过 `service.InitService()` 和 `service.InitAuth()` 注册
7. **创建初始数据** - 调用 `service.CreateInitialData()`
8. **配置 Gin 路由** - 添加中间件，设置 JWT 白名单
9. **创建 HTTP 服务器** - 绑定端口，启动服务

**依赖注入模式**：
- `Service` 结构体包含所有核心组件：Config, DB, Logger, HTTPServer
- 通过 `newService()` 工厂函数创建和初始化所有依赖
- 数据库连接和配置通过 Gin 上下文传递：`c.Set("db", sqlDB)`, `c.Set("config", cfg)`
- Service 层通过全局函数注册：`service.InitService(userRepo)`

### 添加新功能

1. **定义数据模型** - 在 `models/` 目录创建 Go 结构体，包含 GORM 标签和验证标签
2. **创建数据访问层** - 在 `dao/` 目录实现数据库操作，继承 `dao.BaseRepository` 模式
3. **实现业务逻辑** - 在 `service/` 目录编写业务逻辑，通过 `service.InitService()` 注册
4. **添加 API 端点** - 在 `api/v1api/` 目录创建 HTTP 处理器，使用 `util/response/` 返回统一格式
5. **配置路由** - 在 `router/router.go` 中添加路由定义，支持 JWT 中间件和白名单
6. **初始化组件** - 在 `service.go` 的 `newService()` 函数中初始化新组件

**完整示例流程**：
```go
// 1. models/user_profile.go
type UserProfile struct {
    UserID uint   `gorm:"primaryKey"`
    Bio    string `gorm:"type:text"`
}

// 2. dao/user_profile_repository.go
type UserProfileRepository struct {
    *BaseRepository[models.UserProfile]
}

// 3. service/user_profile_service.go
func CreateUserProfile(userID uint, bio string) error {
    // 业务逻辑
}

// 4. api/v1api/user_profile.go
func UpdateProfile(c *gin.Context) {
    // 参数验证
    // 调用 service.CreateUserProfile()
    // 返回 util.response.Success() 或 .Error()
}

// 5. router/router.go
users.PUT("/:id/profile", v1api.UpdateProfile)

// 6. service.go 的 newService() 函数
// 自动通过 AutoMigrate 创建表
// 通过 service.InitService() 注册服务
```

### 数据库

- 使用 GORM v1.31.1 和 PostgreSQL
- 启动时自动迁移数据库表结构
- 通过环境变量配置连接（DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, DB_SSLMODE）

### API 模式

- **RESTful 端点** - 所有 API 位于 `/v1/` 路径下
- **请求/响应格式** - JSON 格式，统一响应结构：`{"code": 200, "message": "成功", "data": {...}}`
- **错误消息** - 中文错误消息，通过 `util/apperror/` 定义业务错误码
- **健康检查** - `/v1/health` 端点返回应用状态和数据库连接状态
- **认证中间件** - JWT token 验证，白名单路由可跳过验证
- **请求日志** - 自动记录所有 HTTP 请求的详细信息

**错误处理流程**：
1. DAO 层返回原始错误或 `gorm.ErrRecordNotFound`
2. Service 层包装业务错误：`apperror.New("用户不存在", 1001)`
3. API 层通过 `response.Error(c, err)` 返回统一格式
4. 中间件捕获 panic 并返回 500 错误

**JWT 认证系统**：
- 密钥配置在 `config.yaml` 的 `jwt.secret`
- Token 过期时间可配置（默认 24 小时）
- 白名单路由：`/v1/login`, `/v1/register`, `/v1/health`
- Token 存储在请求头：`Authorization: Bearer <token>`
- 用户信息通过 `c.Get("user")` 在上下文中获取

## 日志系统

项目使用 Go 标准库 `log/slog` 的结构化 JSON 日志。

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
go run main.go

# 或使用开发配置
APP_MODE=debug LOG_LEVEL=debug LOG_OUTPUT=stdout go run main.go
```

### 代码质量检查
```bash
# 安装并运行代码检查工具
make install-lint
make lint

# 格式化代码
make install-goimports
make goimports

# 提交前检查
git add .
make lint
go test ./...
```

### 添加新功能的标准流程
1. **定义数据模型** (`models/`) - 包含 GORM 标签和验证标签
2. **创建数据访问层** (`dao/`) - 继承 `dao.BaseRepository` 模式
3. **实现业务逻辑** (`service/`) - 通过 `service.InitService()` 注册
4. **添加 API 端点** (`api/v1api/`) - 使用 `util/response/` 返回统一格式
5. **配置路由** (`router/router.go`) - 支持 JWT 中间件和白名单
6. **运行自动迁移** - 重启应用时自动创建表结构

## Docker 环境配置

项目为开发和生产环境提供独立的 Docker Compose 配置。

### 生产环境

**文件**：`docker-compose.yml`

- **应用模式**：`release`
- **日志级别**：`info`
- **日志输出**：`both`（控制台 + 文件）
- **特性**：性能优化，日志持久化到 `./logs/app.log`
- **使用**：`make up-build`

### 环境文件

- **`.env`**：开发环境变量文件（docker-compose 自动加载）
- **`.env.prod`**：生产环境变量文件（创建用于生产环境覆盖）
- 这些文件不通过 git 跟踪（添加到 `.gitignore`）

## 关键架构决策

### Service.go 架构（待评估）
根据 TODO.md，当前的 `service.go` 实现需要评估：
- **位置**：Service 结构体在 main 包中实现
- **职责**：负责所有组件的初始化和依赖注入
- **待解决问题**：
  - 是否应该将 Service 结构体移到独立的包中？
  - 依赖注入方式是否过于复杂？
  - 是否应该使用更轻量级的初始化方式？

### 依赖注入模式
- **当前实现**：通过 `service.go` 中的 `Service` 结构体管理所有依赖
- **依赖传递**：数据库连接和配置通过 Gin 上下文传递：`c.Set("db", sqlDB)`, `c.Set("config", cfg)`
- **Service 注册**：通过全局函数 `service.InitService()` 和 `service.InitAuth()` 注册

### 错误处理策略
1. **DAO 层**：返回原始错误或 `gorm.ErrRecordNotFound`
2. **Service 层**：包装业务错误：`apperror.New("用户不存在", 1001)`
3. **API 层**：通过 `response.Error(c, err)` 返回统一格式
4. **中间件**：捕获 panic 并返回 500 错误

## 重要说明

1. **日志系统** - 使用 Go 标准库 `log/slog` 的结构化 JSON 日志。支持多输出（stdout/file/both），自动记录 HTTP 请求详情。已移除重复的 gin.Logger() 日志。

2. **架构状态** - 项目采用清晰的分层架构（API/Service/DAO/Models）。Service 层通过全局函数注册，依赖注入通过 `service.go` 中的 `Service` 结构体管理。

3. **代码规范** - 代码库使用中文注释和 API 错误消息。添加新代码时保持这一约定。

4. **JWT 认证** - **已实现**。JWT 认证系统完整集成，包含 token 生成、验证、白名单路由和用户上下文传递。

5. **优先级功能**（根据 TODO.md）：
   - ✅ **高优先级已完成**：JWT 用户认证
   - ⚠️ **待评估**：`service.go` 架构合理性评估
   - ⚠️ **待实现**：所有层的单元测试
   - ⚠️ **待添加**：Husky 配置
   - ⚠️ **待配置**：数据库 Schema 权限

6. **错误处理** - 使用 `util/apperror/` 中的自定义错误处理系统，包含业务错误码和中文错误消息。通过 `util/response/` 返回统一格式。

7. **测试覆盖** - 目前项目中没有测试文件。添加测试时遵循 Go 测试约定，创建 `*_test.go` 文件。

8. **数据库迁移** - 使用 GORM 自动迁移，在启动时自动创建/更新表结构。模型定义在 `models/` 目录。

9. **配置管理** - 支持 YAML 配置文件 + 环境变量覆盖。生产环境建议使用环境变量设置敏感信息（数据库密码、JWT 密钥等）。

10. **Docker 部署** - 提供完整的 Docker Compose 配置，支持开发和生产环境。日志通过卷挂载持久化。