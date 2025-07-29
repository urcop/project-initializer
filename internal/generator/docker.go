package generator

import (
	"fmt"
	"os"
	"path/filepath"
)

// generateDockerFiles создает Docker файлы
func (g *Generator) generateDockerFiles(config *ProjectConfig) error {
	// Создаем Dockerfile
	if err := g.generateDockerfile(config); err != nil {
		return err
	}

	// Создаем docker-compose.yml
	if err := g.generateDockerCompose(config); err != nil {
		return err
	}

	// Создаем .dockerignore
	if err := g.generateDockerIgnore(); err != nil {
		return err
	}

	return nil
}

// generateDockerfile создает Dockerfile
func (g *Generator) generateDockerfile(config *ProjectConfig) error {
	content := fmt.Sprintf(`# Build stage
FROM golang:1.21-alpine AS builder

# Устанавливаем git для go mod download
RUN apk add --no-cache git

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go mod файлы
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main cmd/main.go

# Production stage
FROM alpine:latest

# Устанавливаем ca-certificates для HTTPS запросов
RUN apk --no-cache add ca-certificates tzdata

# Создаем пользователя
RUN addgroup -g 1001 -S %s && \
    adduser -S %s -u 1001 -G %s

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем бинарный файл из build stage
COPY --from=builder /app/main .
COPY --from=builder /app/config.yaml .

# Устанавливаем владельца файлов
RUN chown -R %s:%s /app

# Переключаемся на пользователя
USER %s

# Открываем порт
EXPOSE 8080

# Команда запуска
CMD ["./main"]
`, config.Name, config.Name, config.Name, config.Name, config.Name, config.Name)

	dockerfilePath := filepath.Join(g.projectPath, "Dockerfile")
	return os.WriteFile(dockerfilePath, []byte(content), 0644)
}

// generateDockerCompose создает docker-compose.yml
func (g *Generator) generateDockerCompose(config *ProjectConfig) error {
	content := fmt.Sprintf(`version: '3.8'

services:
  %s:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "8080:8080"`, config.Name)

	// Добавляем gRPC порт если включен
	if config.EnableGRPC {
		content += `
      - "9090:9090"`
	}

	content += `
    environment:
      - APP_ENV=production
    volumes:
      - ./config.yaml:/app/config.yaml:ro
    depends_on:`

	// Добавляем БД если нужно
	switch config.Database {
	case "PostgreSQL":
		content += `
      - postgres
    networks:
      - app-network

  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: ` + config.Name + `
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - app-network`

	case "MySQL":
		content += `
      - mysql
    networks:
      - app-network

  mysql:
    image: mysql:8.0
    environment:
      MYSQL_DATABASE: ` + config.Name + `
      MYSQL_ROOT_PASSWORD: password
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql
    networks:
      - app-network`

	case "MongoDB":
		content += `
      - mongodb
    networks:
      - app-network

  mongodb:
    image: mongo:7
    environment:
      MONGO_INITDB_DATABASE: ` + config.Name + `
    ports:
      - "27017:27017"
    volumes:
      - mongodb_data:/data/db
    networks:
      - app-network`
	default:
		content += `
    networks:
      - app-network`
	}

	content += `

networks:
  app-network:
    driver: bridge
`

	// Добавляем volumes если есть БД
	if config.Database != "In-Memory" && config.Database != "Без БД" {
		content += `
volumes:`
		switch config.Database {
		case "PostgreSQL":
			content += `
  postgres_data:`
		case "MySQL":
			content += `
  mysql_data:`
		case "MongoDB":
			content += `
  mongodb_data:`
		}
	}

	dockerComposePath := filepath.Join(g.projectPath, "docker-compose.yml")
	return os.WriteFile(dockerComposePath, []byte(content), 0644)
}

// generateDockerIgnore создает .dockerignore
func (g *Generator) generateDockerIgnore() error {
	content := `# Git
.git
.gitignore

# Documentation
README.md
docs/

# Docker
Dockerfile
.dockerignore
docker-compose.yml

# Development files
.env
.env.local
.env.development
.env.test

# IDE
.vscode/
.idea/
*.swp
*.swo

# OS
.DS_Store
Thumbs.db

# Logs
*.log
logs/

# Test files
*_test.go
test/
coverage.out

# Build artifacts
main
*.exe
dist/
build/

# Temporary files
tmp/
temp/
`

	dockerIgnorePath := filepath.Join(g.projectPath, ".dockerignore")
	return os.WriteFile(dockerIgnorePath, []byte(content), 0644)
}
