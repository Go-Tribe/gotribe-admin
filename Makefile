# ==============================================================================
# 定义全局 Makefile 变量方便后面引用
PROJECT_NAME := gotribe-admin
COMMON_SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
# 项目根目录
ROOT_DIR := $(abspath $(shell cd $(COMMON_SELF_DIR)/ && pwd -P))
# 构建产物、临时文件存放目录
OUTPUT_DIR := $(ROOT_DIR)/_output
# 前端项目目录
WEB_DIR := $(ROOT_DIR)/web/admin
# 包管理器
PNPM ?= pnpm
# 版本信息
VERSION := $(shell git describe --tags --always --dirty)
VERSION_PACKAGE := gotribe-admin/internal/pkg/common
GO ?= go
DOCKER_COMPOSE ?= docker compose
# Go 编译并发数限制（默认使用 CPU 核心数的一半，可通过环境变量覆盖）
# 如果 nproc 不可用，默认使用 2
GO_BUILD_PARALLEL ?= $(shell nproc 2>/dev/null | awk '{print int($$1/2+1)}' || echo 2)

# ==============================================================================
## 检查代码仓库是否是 dirty（默认dirty）
GIT_TREE_STATE:="dirty"
ifeq (, $(shell git status --porcelain 2>/dev/null))
	GIT_TREE_STATE="clean"
endif
GIT_COMMIT:=$(shell git rev-parse HEAD)

GO_LDFLAGS += \
	-X $(VERSION_PACKAGE).GitVersion=$(VERSION) \
	-X $(VERSION_PACKAGE).GitCommit=$(GIT_COMMIT) \
	-X $(VERSION_PACKAGE).GitTreeState=$(GIT_TREE_STATE) \
	-X $(VERSION_PACKAGE).BuildDate=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

# ==============================================================================
# 定义 Makefile all 伪目标，执行 `make` 时，会默认会执行 all 伪目标
.PHONY: all
all: format build

.PHONY: run
run: format build-web # 格式化并编译前后端，然后运行（开发模式）
	@echo ">>> Starting server..."
	@$(GO) run $(ROOT_DIR)/$(PROJECT_NAME).go

.PHONY: dev
dev: # 直接运行 Go（不重新编译前端，假设 dist 已存在）
	@echo ">>> Starting server..."
	@$(GO) run $(ROOT_DIR)/$(PROJECT_NAME).go

.PHONY: dev-frontend
dev-frontend: # 启动前端开发服务器（热重载，需配合 Go 服务使用）
	@cd $(WEB_DIR) && $(PNPM) run dev

# ==============================================================================
# 定义其他需要的伪目标

.PHONY: build-web
build-web: # 编译前端
	@echo ">>> Building frontend..."
	@if [ ! -d "$(WEB_DIR)/node_modules" ]; then \
		echo ">>> Installing frontend dependencies..."; \
		cd $(WEB_DIR) && $(PNPM) install; \
	fi
	@cd $(WEB_DIR) && $(PNPM) run build
	@echo ">>> Frontend build complete."

.PHONY: build
build: build-web # 编译前后端
	@echo ">>> Building Go backend..."
	@mkdir -p $(OUTPUT_DIR)
	@CGO_ENABLED=0 $(GO) build -v -ldflags "$(GO_LDFLAGS)" -o $(OUTPUT_DIR)/$(PROJECT_NAME) $(ROOT_DIR)/$(PROJECT_NAME).go
	@echo ">>> Build complete: $(OUTPUT_DIR)/$(PROJECT_NAME)"

.PHONY: linux
linux: build-web # 交叉编译 Linux 可执行文件（包含前端）
	@mkdir -p $(OUTPUT_DIR)
	@echo ">>> Cross-compiling Linux version..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -ldflags "$(GO_LDFLAGS)" -o $(OUTPUT_DIR)/$(PROJECT_NAME)-linux $(ROOT_DIR)/$(PROJECT_NAME).go
	@echo ">>> Build complete: $(OUTPUT_DIR)/$(PROJECT_NAME)-linux"

.PHONY: format
format: # 格式化 Go 源码.
	@gofmt -s -w ./

.PHONY: add-copyright
add-copyright: # 添加版权头信息.
	@addlicense -v -f $(ROOT_DIR)/scripts/boilerplate.txt $(ROOT_DIR) --skip-dirs=third_party,vendor,resources,doc,template,tmp,static,.idea,$(OUTPUT_DIR)


.PHONY: tidy
tidy: # 自动添加/移除依赖包.
	@$(GO) mod tidy

.PHONY: clean
clean: # 清理构建产物、临时文件等.
	@-rm -vrf $(OUTPUT_DIR)

.PHONY: test
test: # 运行测试
	@$(GO) test -v ./...

.PHONY: test-all
test-all: # 运行所有测试用例（包括单元测试和基准测试）
	@echo "运行所有单元测试..."
	@$(GO) test -v ./...
	@echo "\n运行基准测试..."
	@$(GO) test -bench=. -v ./...
	@echo "\n测试完成！"

.PHONY: test-coverage
test-coverage: # 运行测试并生成覆盖率报告
	@$(GO) test -v -coverprofile=coverage.out ./...
	@$(GO) tool cover -html=coverage.out -o coverage.html

.PHONY: lint
lint: # 代码检查
	@golangci-lint run

.PHONY: fmt
fmt: # 格式化代码
	@$(GO) fmt ./...

.PHONY: vet
vet: # 静态分析
	@$(GO) vet ./...

.PHONY: mod
mod: # 模块管理
	@$(GO) mod download
	@$(GO) mod verify

.PHONY: migrate
migrate: # 数据库迁移
	@$(GO) run $(ROOT_DIR)/$(PROJECT_NAME).go migrate

.PHONY: seed
seed: # 初始化数据
	@$(GO) run $(ROOT_DIR)/$(PROJECT_NAME).go seed

.PHONY: docker
docker: # 构建 Docker 镜像
	@docker build -t $(PROJECT_NAME):$(VERSION) .

.PHONY: docker-run
docker-run: # 运行 Docker 容器
	@$(DOCKER_COMPOSE) up -d

.PHONY: docker-stop
docker-stop: # 停止 Docker 容器
	@$(DOCKER_COMPOSE) down

.PHONY: docker-clean
docker-clean: # 清理 Docker 资源
	@$(DOCKER_COMPOSE) down -v
	@docker system prune -f

.PHONY: install
install: build # 安装到系统
	@cp $(OUTPUT_DIR)/$(PROJECT_NAME) /usr/local/bin/

.PHONY: uninstall
uninstall: # 从系统卸载
	@rm -f /usr/local/bin/$(PROJECT_NAME)

.PHONY: swagger
swagger: # 生成 Swagger 文档
	@swag init -g ./gotribe-admin.go -o ./docs/swagger

.PHONY: verify
verify: format vet test build # 本地校验常用组合

.PHONY: swagger-clean
swagger-clean: # 清理 Swagger 文档
	@-rm -vrf $(ROOT_DIR)/docs/swagger

.PHONY: help
help: # 显示帮助信息
	@echo "Available targets:"
	@echo "  all                - 构建项目 (默认)"
	@echo "  build              - 编译前后端（先编译前端 dist，再编译 Go 二进制）"
	@echo "  build-web          - 编译前端（输出到 web/admin/dist）"
	@echo "  linux              - 交叉编译 Linux 可执行文件（包含前端）"
	@echo "  run                - 格式化并编译前后端，然后运行（开发模式，推荐）"
	@echo "  dev                - 直接运行 Go（不重新编译前端，需确保 dist 已存在）"
	@echo "  dev-frontend       - 启动前端开发服务器（热重载，需配合 Go 服务使用）"
	@echo "  test               - 运行测试"
	@echo "  test-all           - 运行所有测试用例（包括单元测试和基准测试）"
	@echo "  test-coverage      - 运行测试并生成覆盖率报告"
	@echo "  lint               - 代码检查"
	@echo "  fmt                - 使用 go fmt 格式化代码"
	@echo "  format             - 使用 gofmt -s 格式化代码"
	@echo "  vet                - 静态分析"
	@echo "  mod                - 下载并校验模块"
	@echo "  tidy               - 整理 go.mod/go.sum"
	@echo "  verify             - 执行 format、vet、test、build"
	@echo "  migrate            - 数据库迁移"
	@echo "  seed               - 初始化数据"
	@echo "  docker             - 构建 Docker 镜像"
	@echo "  docker-run         - 运行 Docker 容器"
	@echo "  docker-stop        - 停止 Docker 容器"
	@echo "  docker-clean       - 清理 Docker 资源"
	@echo "  install            - 安装到系统"
	@echo "  uninstall          - 从系统卸载"
	@echo "  clean              - 清理构建产物"
	@echo "  swagger            - 生成 Swagger 文档"
	@echo "  swagger-clean      - 清理 Swagger 文档"
	@echo "  add-copyright      - 添加版权头信息"
	@echo "  help               - 显示帮助信息"
