package cmd

import (
	"github.com/spf13/cobra"
)

var Version = "1.0.0"

var rootCmd = &cobra.Command{
	Use:   "project-initializer",
	Short: "Инициализатор архитектуры микросервисов на Go",
	Long: `CLI инструмент для создания структуры микросервисов на Go.
Аналогично django-admin startproject, но для Go микросервисов.

Поддерживает:
- Выбор БД (PostgreSQL, MySQL, MongoDB, in-memory, без БД)
- Выбор фреймворка (Fiber, Gin, Echo)
- Опциональный gRPC сервер
- Автогенерация Swagger, Dockerfile, Makefile
- Конфигурация через YAML
- Интерфейсы для БД и внешних сервисов

Установка:
  go install github.com/urcop/project-initializer@latest

Использование:
  project-initializer init my-service`,
	Version: Version,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(initCmd)
}
