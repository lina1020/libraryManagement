# Makefile for Library Management System

# 变量定义
APP_NAME := library-management
DOCKER_IMAGE := $(APP_NAME):latest
DOCKER_COMPOSE_FILE := docker-compose.yml
DOCKER_COMPOSE_DEV_FILE := docker-compose.dev.yml

# Go 相关变量
GO_VERSION := 1.25
GOOS := linux
GOARCH := amd64
CGO_ENABLED := 0

# Docker 相关变量
DOCKER_REGISTRY :=
DOCKER_TAG := latest

# 默认目标
.PHONY: all
all: build

# 构建 Go 应用
.PHONY: build
build:
	@echo "Building $(APP_NAME)..."
	@CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags="-w -s" -o bin/$(APP_NAME) main.go
	@echo "Build completed: bin/$(APP_NAME)"

# 运行单元测试
.PHONY: test
test:
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@echo "Tests completed"

# 生成测试覆盖率报告
.PHONY: coverage
coverage: test
	@echo "Generating coverage report..."
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# 构建 Docker 镜像
.PHONY: docker-build
docker-build:
	@echo "Building Docker image: $(DOCKER_IMAGE)"
	@docker build -t $(DOCKER_IMAGE) .
	@echo "Docker image built successfully: $(DOCKER_IMAGE)"

# 推送 Docker 镜像到仓库
.PHONY: docker-push
docker-push: docker-build
ifneq ($(DOCKER_REGISTRY),)
	@echo "Tagging and pushing to registry..."
	@docker tag $(DOCKER_IMAGE) $(DOCKER_REGISTRY)/$(DOCKER_IMAGE)
	@docker push $(DOCKER_REGISTRY)/$(DOCKER_IMAGE)
	@echo "Image pushed to registry"
else
	@echo "DOCKER_REGISTRY not set, skipping push"
endif

# 使用 Docker Compose 运行整个项目
.PHONY: run
run:
	@echo "Starting all services with Docker Compose..."
	@docker-compose -f $(DOCKER_COMPOSE_FILE) up --build -d
	@echo "Services started. Check status with: make status"

# 开发环境运行
.PHONY: run-dev
run-dev:
	@echo "Starting development environment..."
	@docker-compose -f $(DOCKER_COMPOSE_DEV_FILE) up --build
	@echo "Development environment started"

# 停止所有服务
.PHONY: stop
stop:
	@echo "Stopping all services..."
	@docker-compose -f $(DOCKER_COMPOSE_FILE) down
	@echo "All services stopped"

# 查看服务状态
.PHONY: status
status:
	@echo "Service status:"
	@docker-compose -f $(DOCKER_COMPOSE_FILE) ps

# 查看服务日志
.PHONY: logs
logs:
	@docker-compose -f $(DOCKER_COMPOSE_FILE) logs -f

# 查看特定服务日志
.PHONY: logs-app
logs-app:
	@docker-compose -f $(DOCKER_COMPOSE_FILE) logs -f app

.PHONY: logs-mysql
logs-mysql:
	@docker-compose -f $(DOCKER_COMPOSE_FILE) logs -f mysql

.PHONY: logs-es
logs-es:
	@docker-compose -f $(DOCKER_COMPOSE_FILE) logs -f elasticsearch

# 清理 Docker 资源
.PHONY: clean
clean:
	@echo "Cleaning up Docker resources..."
	@docker-compose -f $(DOCKER_COMPOSE_FILE) down -v --remove-orphans
	@docker system prune -f
	@echo "Cleanup completed"

# 深度清理（包括镜像）
.PHONY: clean-all
clean-all: clean
	@echo "Removing Docker images..."
	@docker rmi $(DOCKER_IMAGE) 2>/dev/null || true
	@docker image prune -a -f
	@echo "All Docker resources cleaned"

# 初始化数据库
.PHONY: init-db
init-db:
	@echo "Initializing database..."
	# 创建数据库
	@docker-compose -f $(DOCKER_COMPOSE_FILE) exec -T mysql mysql -uroot -p10209066 -e "CREATE DATABASE IF NOT EXISTS library CHARACTER SET utf8mb4;"
	# 导入 books.sql
	@echo "Importing books.sql..."
	@docker-compose -f $(DOCKER_COMPOSE_FILE) exec -i -T mysql mysql -uroot -p10209066 library < books.sql
	# 导入 users.sql
	@echo "Importing users.sql..."
	@docker-compose -f $(DOCKER_COMPOSE_FILE) exec -i -T mysql mysql -uroot -p10209066 library < users.sql
	@echo "Database initialized"

# 重启服务
.PHONY: restart
restart: stop run

# 本地开发运行
.PHONY: run-local
run-local:
	@echo "Starting local development server..."
	@go run main.go

# 安装依赖
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies installed"

# 代码格式化
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Code formatted"

# 代码静态检查
.PHONY: lint
lint:
	@echo "Running linter..."
	@golangci-lint run ./...
	@echo "Linting completed"

# 生成 API 文档
.PHONY: docs
docs:
	@echo "Generating API documentation..."
	@swag init -g main.go -o ./docs
	@echo "API documentation generated in ./docs"

# 数据库迁移
.PHONY: migrate
migrate:
	@echo "Running database migrations..."
	@docker-compose -f $(DOCKER_COMPOSE_FILE) exec app ./library-management -migrate
	@echo "Database migration completed"

# 备份数据库
.PHONY: backup-db
backup-db:
	@echo "Creating database backup..."
	@mkdir -p backups
	@docker-compose -f $(DOCKER_COMPOSE_FILE) exec mysql mysqldump -uroot -p10209066 --single-transaction library > backups/library_backup_$(shell date +%Y%m%d_%H%M%S).sql
	@echo "Database backup created in backups/ directory"

# 恢复数据库
.PHONY: restore-db
restore-db:
	@echo "Please specify backup file with: make restore-db BACKUP_FILE=path/to/backup.sql"
ifdef BACKUP_FILE
	@docker-compose -f $(DOCKER_COMPOSE_FILE) exec -T mysql mysql -uroot -p10209066 library < $(BACKUP_FILE)
	@echo "Database restored from $(BACKUP_FILE)"
endif

# 进入应用容器
.PHONY: shell-app
shell-app:
	@docker-compose -f $(DOCKER_COMPOSE_FILE) exec app sh

# 进入数据库容器
.PHONY: shell-db
shell-db:
	@docker-compose -f $(DOCKER_COMPOSE_FILE) exec mysql mysql -uroot -p10209066 library

# ES 相关操作
.PHONY: es-status
es-status:
	@echo "Checking Elasticsearch status..."
	@curl -s "http://localhost:9200/_cluster/health?pretty"

# 重建 ES 索引
.PHONY: reindex
reindex:
	@echo "Reindexing Elasticsearch..."
	@curl -X POST "http://localhost:8080/admin/es/index/reindex" -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
	@echo "Reindexing completed"

# 显示帮助信息
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  build         - Build the Go application"
	@echo "  test          - Run unit tests"
	@echo "  coverage      - Generate test coverage report"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-push   - Push Docker image to registry"
	@echo "  run           - Start all services with Docker Compose"
	@echo "  run-dev       - Start development environment"
	@echo "  run-local     - Run application locally"
	@echo "  stop          - Stop all services"
	@echo "  restart       - Restart all services"
	@echo "  status        - Show service status"
	@echo "  logs          - Show all service logs"
	@echo "  logs-app      - Show application logs"
	@echo "  logs-mysql    - Show MySQL logs"
	@echo "  logs-es       - Show Elasticsearch logs"
	@echo "  clean         - Clean Docker resources"
	@echo "  clean-all     - Clean all Docker resources including images"
	@echo "  init-db       - Initialize database with schema"
	@echo "  backup-db     - Create database backup"
	@echo "  restore-db    - Restore database from backup"
	@echo "  shell-app     - Enter application container"
	@echo "  shell-db      - Enter database container"
	@echo "  es-status     - Check Elasticsearch status"
	@echo "  reindex       - Rebuild Elasticsearch indices"
	@echo "  deps          - Install Go dependencies"
	@echo "  fmt           - Format Go code"
	@echo "  lint          - Run code linter"
	@echo "  docs          - Generate API documentation"
	@echo "  help          - Show this help message"