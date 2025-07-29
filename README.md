# Go Microservice Initializer

[![CI](https://github.com/urcop/project-initializer/workflows/CI/badge.svg)](https://github.com/urcop/project-initializer/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/urcop/project-initializer)](https://goreportcard.com/report/github.com/urcop/project-initializer)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Release](https://img.shields.io/github/release/urcop/project-initializer.svg)](https://github.com/urcop/project-initializer/releases)

**CLI инструмент для быстрого создания архитектуры микросервисов на Go**

Аналог `django-admin startproject`, но для Go микросервисов. Автоматически генерирует полную структуру проекта с выбранными технологиями.

## ✨ Возможности

- 🎯 **Интерактивный CLI** - простые prompts для выбора настроек
- 🚀 **Веб-фреймворки**: Gin, Fiber, Echo
- 🗄️ **Базы данных**: PostgreSQL, MySQL, MongoDB, SQLite (in-memory), без БД
- 🌐 **gRPC поддержка** - опциональный gRPC сервер с proto файлами
- 📝 **Swagger документация** - автоматическая генерация API docs
- 🐳 **Docker** - готовые Dockerfile и docker-compose.yml
- 🔧 **Makefile** - команды для разработки и деплоя
- ⚙️ **YAML конфигурация** - гибкая настройка через config.yaml
- 🏗️ **Чистая архитектура** - интерфейсы для БД и внешних сервисов
- 🎯 **Health endpoint** - /health с подробной информацией о сервисе

## 🚀 Установка

### Через go install (рекомендуется)

```bash
# Установка последней версии
go install github.com/urcop/project-initializer@latest

# Или установка конкретной версии
go install github.com/urcop/project-initializer@v1.0.0
```

После установки команда `project-initializer` будет доступна глобально:

```bash
project-initializer init my-service
project-initializer --version
project-initializer --help
```

> **Примечание**: Убедитесь, что `$GOPATH/bin` или `$GOBIN` добавлен в ваш `$PATH`

### Скачивание бинарных файлов

Скачайте готовые бинарные файлы для вашей платформы из [GitHub Releases](https://github.com/urcop/project-initializer/releases):

- **Linux**: `project-initializer_Linux_x86_64.tar.gz`
- **macOS**: `project-initializer_Darwin_x86_64.tar.gz` (Intel) / `project-initializer_Darwin_arm64.tar.gz` (Apple Silicon)
- **Windows**: `project-initializer_Windows_x86_64.zip`

```bash
# Пример для Linux
wget https://github.com/urcop/project-initializer/releases/latest/download/project-initializer_Linux_x86_64.tar.gz
tar -xzf project-initializer_Linux_x86_64.tar.gz
sudo mv project-initializer /usr/local/bin/
```

### Из исходного кода

```bash
git clone https://github.com/urcop/project-initializer.git
cd project-initializer
go build -o project-initializer .
./project-initializer init my-service
```

### Homebrew (планируется)

```bash
# Скоро будет доступно
brew install project-initializer
```

## 📖 Использование

### Создание нового микросервиса

```bash
project-initializer init
```

Или с указанием имени проекта:

```bash
project-initializer init my-service
```

### Интерактивные вопросы

1. **Module name** - для `go mod init` (например: `github.com/myorg/my-service`)
2. **Веб-фреймворк** - Gin, Fiber или Echo
3. **База данных** - PostgreSQL, MySQL, MongoDB, SQLite или без БД
4. **gRPC сервер** - включить или нет

### Запуск созданного проекта

```bash
cd my-service
go mod tidy
make run
```

## 📁 Структура созданного проекта

```
my-service/
├── cmd/
│   └── main.go                    # Точка входа приложения
├── internal/
│   ├── app/
│   │   └── app.go                # Основная логика приложения
│   ├── config/
│   │   └── config.go             # Конфигурация из YAML
│   ├── handlers/
│   │   ├── handler.go            # HTTP обработчики
│   │   └── health.go             # Health check endpoint
│   ├── middleware/
│   │   └── middleware.go         # HTTP middleware
│   ├── models/
│   │   └── models.go             # Модели данных
│   ├── repository/
│   │   └── user.go               # Репозитории для работы с БД
│   ├── services/
│   └── grpc/                     # gRPC сервер (опционально)
│       ├── server.go
│       ├── client.go
│       └── pb/                   # Сгенерированные protobuf файлы
├── pkg/
│   ├── context/
│   │   └── context.go            # Кастомный контекст приложения
│   ├── database/
│   │   ├── interface.go          # Интерфейсы для БД
│   │   └── database.go           # Реализация для выбранной БД
│   └── logger/
│       └── logger.go             # Интерфейс и реализация логгера
├── api/
│   ├── proto/                    # Proto файлы (если включен gRPC)
│   └── swagger/                  # Swagger документация
├── docs/                         # Документация
├── scripts/
│   └── proto.mk                  # Makefile для protobuf
├── deployments/                  # Конфигурации для деплоя
├── config.yaml                   # Основной конфиг файл
├── Dockerfile                    # Multi-stage Docker build
├── docker-compose.yml            # Оркестрация с БД
├── Makefile                      # Команды разработки
├── .dockerignore
└── go.mod
```

## 🔧 Команды Makefile

```bash
# Разработка
make install      # Установить зависимости
make run          # Запустить приложение
make dev          # Запуск с hot reload (air)
make build        # Собрать бинарный файл

# Тестирование
make test         # Запустить тесты
make test-coverage # Тесты с покрытием
make lint         # Линтеры

# Docker
make docker-build    # Собрать Docker образ
make docker-run      # Запустить контейнер
make docker-compose-up # Запустить через docker-compose

# Swagger
make swagger      # Генерировать Swagger документацию

# Protobuf (если включен gRPC)
make proto-gen    # Генерировать Go код из proto файлов
```

## 🌐 API Endpoints

Каждый созданный микросервис включает:

### Health Check
```bash
GET /health
```

Возвращает подробную информацию о состоянии сервиса:

```json
{
  "status": "ok",
  "service": "my-service",
  "version": "1.0.0",
  "timestamp": "2024-01-01T12:00:00Z",
  "uptime": "1h30m45s",
  "system": {
    "go_version": "go1.21",
    "num_cpu": 8,
    "goroutines": 15
  },
  "database": {
    "connected": true
  }
}
```

### Ping
```bash
GET /api/v1/ping
```

```json
{
  "message": "pong",
  "service": "my-service",
  "version": "1.0.0"
}
```

### Swagger UI (если настроен)
```bash
GET /swagger/
```

## 🗄️ Поддерживаемые базы данных

### PostgreSQL
- GORM ORM
- Пул соединений
- Миграции

### MySQL  
- GORM ORM
- Конфигурация charset
- Пул соединений

### MongoDB
- Official Go driver
- Контекстные операции
- Индексы

### SQLite (In-Memory)
- Для тестирования
- GORM поддержка

### Без БД
- Чистый HTTP сервер
- Без зависимостей БД

## 🌐 gRPC поддержка

При включении gRPC автоматически генерируется:

- **Proto файлы** с базовыми сервисами
- **gRPC сервер** с реализацией методов
- **gRPC клиент** для примера
- **Makefile команды** для protobuf
- **Health check и CRUD** операции

### Пример использования gRPC

```bash
# Генерация protobuf файлов
make proto-install  # Один раз
make proto-gen     # При изменении .proto файлов

# Тестирование с grpcurl
grpcurl -plaintext localhost:9090 list
grpcurl -plaintext localhost:9090 myservice.MyServiceService/HealthCheck
```

## ⚙️ Конфигурация

Все настройки в `config.yaml`:

```yaml
app:
  name: "my-service"
  version: "1.0.0"
  debug: true
  port: 8080

database:
  type: "postgres"
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "password"
  name: "my-service"

logger:
  level: "debug"
  format: "json"
  output: "stdout"

swagger:
  enabled: true
  title: "My Service API"
  host: "localhost:8080"
  base_path: "/api/v1"

grpc:
  enabled: true
  port: 9090
```

## 🎯 Философия дизайна

### Интерфейсы везде
- Все взаимодействия с БД через интерфейсы
- Легкое тестирование и мокирование
- Слабая связанность компонентов

### Контекст-ориентированный подход
- Кастомный `AppContext` с логгером
- Трассировка запросов
- Управление таймаутами

### Готовность к продакшену
- Graceful shutdown
- Health checks
- Структурированное логирование
- Docker контейнеризация
- Метрики и мониторинг

## 🛠️ Технологический стек

- **Языки**: Go 1.21+
- **CLI**: Cobra + Survey (интерактивные prompts)
- **HTTP**: Gin / Fiber / Echo
- **gRPC**: google.golang.org/grpc
- **БД**: GORM, MongoDB Driver
- **Логирование**: Logrus
- **Конфигурация**: YAML
- **Документация**: Swagger/OpenAPI
- **Контейнеризация**: Docker, Docker Compose

## 📋 Требования

- **Go 1.21+** - для установки через `go install`
- **Docker** (опционально) - для запуска сгенерированных docker-compose файлов
- **protoc** (опционально) - для gRPC генерации protobuf файлов

### Проверка установки

После установки проверьте что CLI работает:

```bash
project-initializer --version
project-initializer --help
```

Если команда не найдена, убедитесь что `$GOPATH/bin` добавлен в `$PATH`:

```bash
# Для bash/zsh
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
source ~/.bashrc

# Для fish
fish_add_path (go env GOPATH)/bin
```

## 🚀 Пример создания полного микросервиса

```bash
# 1. Создание проекта
project-initializer init user-service

# Выбираем:
# - Module: github.com/mycompany/user-service  
# - Framework: Gin
# - Database: PostgreSQL
# - gRPC: Yes

# 2. Запуск окружения
cd user-service
docker-compose up -d postgres  # Запуск БД

# 3. Установка зависимостей и запуск
go mod tidy
make run

# 4. Тестирование
curl http://localhost:8080/health
curl http://localhost:8080/api/v1/ping
curl http://localhost:8080/swagger/

# 5. gRPC тестирование
grpcurl -plaintext localhost:9090 list
```

## 📚 Roadmap

- [ ] Поддержка Redis
- [ ] Метрики (Prometheus)
- [ ] Авторизация/JWT
- [ ] Миграции БД
- [ ] CI/CD темплейты
- [ ] Kubernetes манифесты
- [ ] Трассировка (Jaeger)
- [ ] GraphQL поддержка

## 🤝 Вклад в проект

1. Fork репозитория
2. Создайте feature branch
3. Добавьте изменения
4. Создайте Pull Request