.PHONY: build lint install-lint goimports swag run-dev run-prod up up-build down logs ps clean restart
BINARY_NAME := main
GO := go
GOBIN := "$(shell go env GOPATH)/bin"
LINT := $(GOBIN)/golangci-lint
LINT_VERSION := 2.7.1
SWAG := $(GOBIN)/swag
DOCKER_COMPOSE := docker-compose

build:
	@echo "编译 Linux 可执行文件..."
	@mkdir -p bin
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	$(GO) build -ldflags '-extldflags "-static" -s -w' -o bin/$(BINARY_NAME)

lint:
	@which $(LINT) > /dev/null || (echo "golangci-lint 未安装，运行 'make install-lint'" && exit 1)
	$(LINT) run

install-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN) v$(LINT_VERSION)

goimports: $(GOBIN)/goimports
	$(GOBIN)/goimports -w .

$(GOBIN)/goimports:
	$(GO) install golang.org/x/tools/cmd/goimports@latest

install-swag: $(GOBIN)/swag

$(GOBIN)/swag:
	$(GO) install github.com/swaggo/swag/cmd/swag@latest

swag: install-swag
	$(SWAG) init

# Docker commands
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

restart: down up-build