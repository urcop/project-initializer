package generator

import (
	"fmt"
	"os"
	"path/filepath"
)

// generateMakefile создает Makefile
func (g *Generator) generateMakefile(config *ProjectConfig) error {
	content := fmt.Sprintf(`# Makefile для %s

# Переменные
APP_NAME=%s
BINARY_NAME=main
DOCKER_IMAGE=%s
VERSION=1.0.0

# Go переменные
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Цвета для вывода
RED=\\033[0;31m
GREEN=\\033[0;32m
YELLOW=\\033[0;33m
BLUE=\\033[0;34m
NC=\\033[0m # No Color

.PHONY: help build run test clean deps tidy lint docker-build docker-run docker-stop docker-clean swagger dev install

# Помощь
help: ## Показать справку
	@echo "$(BLUE)Makefile для $(APP_NAME)$(NC)"
	@echo ""
	@echo "$(YELLOW)Доступные команды:$(NC)"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(GREEN)%%-15s$(NC) %%s\\n", $$1, $$2}' $(MAKEFILE_LIST)

# Установка зависимостей
install: ## Установить зависимости
	@echo "$(BLUE)Установка зависимостей...$(NC)"
	$(GOMOD) download
	$(GOMOD) tidy

# Обновление зависимостей
deps: ## Обновить зависимости
	@echo "$(BLUE)Обновление зависимостей...$(NC)"
	$(GOGET) -u ./...
	$(GOMOD) tidy

# Проверка зависимостей
tidy: ## Очистить и проверить зависимости
	@echo "$(BLUE)Проверка зависимостей...$(NC)"
	$(GOMOD) tidy
	$(GOMOD) verify

# Сборка приложения
build: ## Собрать приложение
	@echo "$(BLUE)Сборка приложения...$(NC)"
	$(GOBUILD) -ldflags="-w -s -X main.version=$(VERSION)" -o $(BINARY_NAME) cmd/main.go
	@echo "$(GREEN)Сборка завершена: $(BINARY_NAME)$(NC)"

# Сборка для разных платформ
build-linux: ## Собрать для Linux
	@echo "$(BLUE)Сборка для Linux...$(NC)"
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags="-w -s" -o $(BINARY_NAME)-linux cmd/main.go

build-windows: ## Собрать для Windows
	@echo "$(BLUE)Сборка для Windows...$(NC)"
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -ldflags="-w -s" -o $(BINARY_NAME)-windows.exe cmd/main.go

build-mac: ## Собрать для macOS
	@echo "$(BLUE)Сборка для macOS...$(NC)"
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -ldflags="-w -s" -o $(BINARY_NAME)-darwin cmd/main.go

build-all: build-linux build-windows build-mac ## Собрать для всех платформ

# Запуск приложения
run: ## Запустить приложение
	@echo "$(BLUE)Запуск приложения...$(NC)"
	$(GOBUILD) -o $(BINARY_NAME) cmd/main.go && ./$(BINARY_NAME)

# Запуск в режиме разработки
dev: ## Запуск в режиме разработки (с автоперезагрузкой)
	@echo "$(BLUE)Запуск в режиме разработки...$(NC)"
	@if command -v air > /dev/null 2>&1; then \\
		air; \\
	else \\
		echo "$(YELLOW)air не установлен. Установите его: go install github.com/cosmtrek/air@latest$(NC)"; \\
		echo "$(YELLOW)Запуск обычным способом...$(NC)"; \\
		$(MAKE) run; \\
	fi

# Тестирование
test: ## Запустить тесты
	@echo "$(BLUE)Запуск тестов...$(NC)"
	$(GOTEST) -v ./...

test-coverage: ## Запустить тесты с покрытием
	@echo "$(BLUE)Запуск тестов с покрытием...$(NC)"
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Отчет о покрытии создан: coverage.html$(NC)"

test-race: ## Запустить тесты с проверкой гонок
	@echo "$(BLUE)Запуск тестов с проверкой гонок...$(NC)"
	$(GOTEST) -race -v ./...

# Линтинг
lint: ## Запустить линтеры
	@echo "$(BLUE)Запуск линтеров...$(NC)"
	@if command -v golangci-lint > /dev/null 2>&1; then \\
		golangci-lint run; \\
	else \\
		echo "$(YELLOW)golangci-lint не установлен. Установите его: https://golangci-lint.run/usage/install/$(NC)"; \\
		echo "$(YELLOW)Используем go vet...$(NC)"; \\
		$(GOCMD) vet ./...; \\
	fi

# Форматирование кода
fmt: ## Форматировать код
	@echo "$(BLUE)Форматирование кода...$(NC)"
	$(GOCMD) fmt ./...

# Swagger документация
swagger: ## Генерировать Swagger документацию
	@echo "$(BLUE)Генерация Swagger документации...$(NC)"
	@if command -v swag > /dev/null 2>&1; then \\
		swag init -g cmd/main.go -o ./docs; \\
		echo "$(GREEN)Swagger документация создана в ./docs$(NC)"; \\
	else \\
		echo "$(YELLOW)swag не установлен. Установите его: go install github.com/swaggo/swag/cmd/swag@latest$(NC)"; \\
	fi

# Docker команды
docker-build: ## Собрать Docker образ
	@echo "$(BLUE)Сборка Docker образа...$(NC)"
	docker build -t $(DOCKER_IMAGE):$(VERSION) -t $(DOCKER_IMAGE):latest .
	@echo "$(GREEN)Docker образ собран: $(DOCKER_IMAGE):$(VERSION)$(NC)"

docker-run: ## Запустить контейнер
	@echo "$(BLUE)Запуск Docker контейнера...$(NC)"
	docker run -d --name $(APP_NAME) -p 8080:8080`, config.Name, config.Name, config.Name)

	// Добавляем gRPC порт если включен
	if config.EnableGRPC {
		content += ` -p 9090:9090`
	}

	content += ` $(DOCKER_IMAGE):latest
	@echo "$(GREEN)Контейнер запущен: $(APP_NAME)$(NC)"

docker-stop: ## Остановить контейнер
	@echo "$(BLUE)Остановка контейнера...$(NC)"
	docker stop $(APP_NAME) || true
	docker rm $(APP_NAME) || true

docker-logs: ## Показать логи контейнера
	docker logs -f $(APP_NAME)

docker-shell: ## Войти в контейнер
	docker exec -it $(APP_NAME) /bin/sh

docker-compose-up: ## Запустить через docker-compose
	@echo "$(BLUE)Запуск через docker-compose...$(NC)"
	docker-compose up -d
	@echo "$(GREEN)Сервисы запущены$(NC)"

docker-compose-down: ## Остановить docker-compose
	@echo "$(BLUE)Остановка docker-compose...$(NC)"
	docker-compose down

docker-compose-logs: ## Показать логи docker-compose
	docker-compose logs -f

docker-clean: ## Очистить Docker ресурсы
	@echo "$(BLUE)Очистка Docker ресурсов...$(NC)"
	docker system prune -f
	docker volume prune -f

# Очистка
clean: ## Очистить собранные файлы
	@echo "$(BLUE)Очистка...$(NC)"
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME)-*
	rm -f coverage.out coverage.html
	@echo "$(GREEN)Очистка завершена$(NC)"

# Проверка безопасности
security: ## Проверка безопасности
	@echo "$(BLUE)Проверка безопасности...$(NC)"
	@if command -v gosec > /dev/null 2>&1; then \\
		gosec ./...; \\
	else \\
		echo "$(YELLOW)gosec не установлен. Установите его: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest$(NC)"; \\
	fi

# Полная проверка
check: lint test security ## Запустить все проверки

# Подготовка к продакшену
release: clean build test lint ## Подготовить релиз
	@echo "$(GREEN)Релиз готов!$(NC)"

# Показать версию
version: ## Показать версию
	@echo "$(BLUE)$(APP_NAME) версия: $(VERSION)$(NC)"

# Статус
status: ## Показать статус приложения
	@echo "$(BLUE)Статус $(APP_NAME):$(NC)"
	@if pgrep -f "$(BINARY_NAME)" > /dev/null; then \\
		echo "$(GREEN)✓ Приложение запущено$(NC)"; \\
	else \\
		echo "$(RED)✗ Приложение остановлено$(NC)"; \\
	fi

# По умолчанию показываем справку
.DEFAULT_GOAL := help
`

	makefilePath := filepath.Join(g.projectPath, "Makefile")
	return os.WriteFile(makefilePath, []byte(content), 0644)
}
