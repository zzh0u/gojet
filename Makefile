BINARY_NAME := main
DOCKER_IMAGE := gojet
GO := go
GOBIN := "$(shell go env GOPATH)/bin"
LINT := $(GOBIN)/golangci-lint
LINT_VERSION := 2.7.1
SWAG := $(GOBIN)/swag

build:
	@echo "正在编译 Linux 下可执行文件"
	@mkdir -p bin
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	$(GO) build -ldflags '-extldflags "-static" -s -w' -o bin/$(BINARY_NAME)
	@echo "编译完成"

lint: check-golangci-lint-version
	$(LINT) -c ./.golangci.yaml run --output.text.path=stdout --output.text.print-issued-lines=true --output.text.print-linter-name=true

install-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN) v$(LINT_VERSION)

check-golangci-lint-version:
	@which $(LINT) > /dev/null || (echo "golangci-lint 未安装，运行 'make install-lint' 安装" && exit 1)

GOIMPORTS := $(GOBIN)/goimports

install-goimports:
	@command -v $(GOIMPORTS) >/dev/null || $(GO) install golang.org/x/tools/cmd/goimports@latest

goimports: 
	$(GOIMPORTS) -w .

install-swag:
	@command -v $(SWAG) >/dev/null || $(GO) install github.com/swaggo/swag/cmd/swag@latest

swag:
	$(SWAG) init

# Docker Compose commands
docker-up:
	@echo "启动 Docker Compose 服务"
	docker-compose up -d

docker-up-build:
	@echo "构建并启动 Docker Compose 服务"
	docker-compose up --build -d

docker-down:
	@echo "停止 Docker Compose 服务"
	docker-compose down

docker-logs:
	@echo "查看 Docker Compose 日志"
	docker-compose logs -f

docker-clean:
	@echo "清理 Docker Compose 容器和网络"
	docker-compose down -v --remove-orphans

docker-restart:
	@echo "删除容器与镜像，尝试从注册表拉取镜像，若失败则本地构建并启动"
	# 停止并移除容器、卷与孤儿容器
	docker-compose down
	# 删除本地镜像（若存在）
	docker image rm $(DOCKER_IMAGE):latest || true
	# 使用本地镜像（如不存在则本地构建）并启动
	docker-compose up --build -d

.PHONY: build install-lint check-golangci-lint-version lint goimports install-swag swag docker-up docker-up-build docker-down docker-logs docker-ps docker-clean docker-restart