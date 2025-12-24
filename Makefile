.PHONY: build lint install-lint goimports install-goimports install-swag swag up up-build down logs clean restart

BINARY_NAME := main
GO := go
GOBIN := $(shell go env GOPATH)/bin
LINT := $(GOBIN)/golangci-lint
LINT_VERSION := 2.7.1
SWAG_VERSION := v1.16.6
GOIMPORTS_VERSION := v0.40.0
SWAG := $(GOBIN)/swag
DOCKER_COMPOSE := docker-compose

# 构建命令
build:
	@echo "编译 Linux 可执行文件..."
	@mkdir -p bin
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -ldflags '-extldflags "-static" -s -w' -o bin/$(BINARY_NAME)

# 代码质量工具
lint:
	@which $(LINT) > /dev/null || (echo "golangci-lint 未安装，运行 'make install-lint'" && exit 1)
	$(LINT) run

install-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN) v$(LINT_VERSION)

goimports:
	@which $(GOBIN)/goimports > /dev/null || (echo "goimports 未安装，运行 'make install-goimports'" && exit 1)
	$(GOBIN)/goimports -w .

install-goimports:
	$(GO) install golang.org/x/tools/cmd/goimports@$(GOIMPORTS_VERSION)

# Swagger 文档工具
swag:
	@which $(SWAG) > /dev/null || (echo "swag 未安装，运行 'make install-swag'" && exit 1)
	$(SWAG) init

install-swag:
	$(GO) install github.com/swaggo/swag/cmd/swag@$(SWAG_VERSION)

# Docker 命令
up:
	$(DOCKER_COMPOSE) up -d

up-build:
	$(DOCKER_COMPOSE) up --build -d

down:
	$(DOCKER_COMPOSE) down

logs:
	$(DOCKER_COMPOSE) logs -f

clean:
	$(DOCKER_COMPOSE) down -v --remove-orphans

restart: down clean up-build