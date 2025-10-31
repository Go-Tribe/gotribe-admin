# ==============================================================================
# 定义全局 Makefile 变量方便后面引用
PROJECT_NAME = "gotribe-admin"
COMMON_SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
# 项目根目录
ROOT_DIR := $(abspath $(shell cd $(COMMON_SELF_DIR)/ && pwd -P))
# 构建产物、临时文件存放目录
OUTPUT_DIR:= $(ROOT_DIR)/_output
# 版本信息
VERSION := $(shell git describe --tags --always --dirty)
VERSION_PACKAGE := gotribe-admin/internal/pkg/common

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
all: add-copyright  format build

.PHONY: run
run: tidy  format dev

# ==============================================================================
# 定义其他需要的伪目标

.PHONY: build
build: tidy # 编译源码，依赖 tidy 目标自动添加/移除依赖包.
	@go build -v -ldflags "$(GO_LDFLAGS)" -o $(OUTPUT_DIR)/$(PROJECT_NAME) $(ROOT_DIR)/$(PROJECT_NAME).go

.PHONY: format
format: # 格式化 Go 源码.
	@gofmt -s -w ./

.PHONY: add-copyright
add-copyright: # 添加版权头信息.
	@addlicense -v -f $(ROOT_DIR)/scripts/boilerplate.txt $(ROOT_DIR) --skip-dirs=third_party,vendor,resources,doc,template,tmp,static,.idea,$(OUTPUT_DIR)


.PHONY: tidy
tidy: # 自动添加/移除依赖包.
	@go mod tidy

.PHONY: clean
clean: # 清理构建产物、临时文件等.
	@-rm -vrf $(OUTPUT_DIR)

.PHONY: dev
dev: # 开发运行
	@go run $(ROOT_DIR)/$(PROJECT_NAME).go

.PHONY: test
test: # 运行测试
	@go test -v ./...

.PHONY: test-coverage
test-coverage: # 运行测试并生成覆盖率报告
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

.PHONY: lint
lint: # 代码检查
	@golangci-lint run

.PHONY: fmt
fmt: # 格式化代码
	@go fmt ./...

.PHONY: vet
vet: # 静态分析
	@go vet ./...

.PHONY: mod
mod: # 模块管理
	@go mod download
	@go mod verify

.PHONY: migrate
migrate: # 数据库迁移
	@go run $(ROOT_DIR)/$(PROJECT_NAME).go migrate

.PHONY: seed
seed: # 初始化数据
	@go run $(ROOT_DIR)/$(PROJECT_NAME).go seed

.PHONY: docker
docker: # 构建 Docker 镜像
	@docker build -t $(PROJECT_NAME):$(VERSION) .

.PHONY: docker-run
docker-run: # 运行 Docker 容器
	@docker-compose up -d

.PHONY: docker-stop
docker-stop: # 停止 Docker 容器
	@docker-compose down

.PHONY: docker-clean
docker-clean: # 清理 Docker 资源
	@docker-compose down -v
	@docker system prune -f

.PHONY: install
install: build # 安装到系统
	@cp $(OUTPUT_DIR)/$(PROJECT_NAME) /usr/local/bin/

.PHONY: uninstall
uninstall: # 从系统卸载
	@rm -f /usr/local/bin/$(PROJECT_NAME)

.PHONY: swagger
swagger: # 生成 Swagger 文档
	@swag init -g ./gotribe-admin.go -o ./docs

.PHONY: swagger-clean
swagger-clean: # 清理 Swagger 文档
	@-rm -vrf $(ROOT_DIR)/docs

.PHONY: help
help: # 显示帮助信息
	@echo "Available targets:"
	@echo "  all          - 构建项目 (默认)"
	@echo "  build        - 编译源码"
	@echo "  run          - 开发运行"
	@echo "  dev          - 开发运行"
	@echo "  test         - 运行测试"
	@echo "  test-coverage- 运行测试并生成覆盖率报告"
	@echo "  lint         - 代码检查"
	@echo "  fmt          - 格式化代码"
	@echo "  vet          - 静态分析"
	@echo "  mod          - 模块管理"
	@echo "  migrate      - 数据库迁移"
	@echo "  seed         - 初始化数据"
	@echo "  docker       - 构建 Docker 镜像"
	@echo "  docker-run   - 运行 Docker 容器"
	@echo "  docker-stop  - 停止 Docker 容器"
	@echo "  docker-clean - 清理 Docker 资源"
	@echo "  install      - 安装到系统"
	@echo "  uninstall    - 从系统卸载"
	@echo "  clean        - 清理构建产物"
	@echo "  swagger      - 生成 Swagger 文档"
	@echo "  swagger-clean- 清理 Swagger 文档"
	@echo "  help         - 显示帮助信息"
