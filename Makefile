BINARY_NAME := main
DOCKER_IMAGE := gojet
GO := go
GOBIN := "$(shell go env GOPATH)/bin"
LINT := $(GOBIN)/golangci-lint
LINT_VERSION := 2.7.1
SWAG := $(GOBIN)/swag

build:
	@echo "正在编译Linux下可执行文件"
	@mkdir -p bin
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	$(GO) build -ldflags '-extldflags "-static" -s -w' -o bin/$(BINARY_NAME)
	@echo "编译完成"

lint: check-golangci-lint-version
	$(LINT) run 

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

.PHONY: build install-lint check-golangci-lint-version lint goimports install-swag swag