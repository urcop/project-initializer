package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/urcop/project-initializer/internal/generator"
)

var (
	moduleName string
	framework  string
	database   string
	enableGRPC bool
)

var initCmd = &cobra.Command{
	Use:   "init [project-name]",
	Short: "Инициализировать новый микросервис",
	Long:  `Создает новую структуру микросервиса с выбранными опциями`,
	Args:  cobra.MaximumNArgs(1),
	RunE:  runInit,
}

func init() {
	initCmd.Flags().StringVar(&moduleName, "module", "", "Go module name (например: github.com/yourorg/project)")
	initCmd.Flags().StringVar(&framework, "framework", "", "Веб-фреймворк (gin, fiber, echo)")
	initCmd.Flags().StringVar(&database, "database", "", "База данных (postgresql, mysql, mongodb, in-memory, none)")
	initCmd.Flags().BoolVar(&enableGRPC, "grpc", false, "Включить gRPC сервер")
}

func runInit(cmd *cobra.Command, args []string) error {
	config := &generator.ProjectConfig{}

	// Получаем имя проекта
	if len(args) > 0 {
		config.Name = args[0]
	} else {
		prompt := &survey.Input{
			Message: "Введите имя проекта:",
			Default: "my-service",
		}
		if err := survey.AskOne(prompt, &config.Name); err != nil {
			return fmt.Errorf("ошибка ввода имени проекта: %w", err)
		}
	}

	// Module name
	if moduleName != "" {
		config.ModuleName = moduleName
	} else {
		prompt := &survey.Input{
			Message: "Введите module name для go mod init:",
			Default: fmt.Sprintf("github.com/yourorg/%s", config.Name),
		}
		if err := survey.AskOne(prompt, &config.ModuleName); err != nil {
			return fmt.Errorf("ошибка ввода module name: %w", err)
		}
	}

	// Выбор фреймворка
	if framework != "" {
		config.Framework = framework
	} else {
		frameworkPrompt := &survey.Select{
			Message: "Выберите веб-фреймворк:",
			Options: []string{"Gin", "Fiber", "Echo"},
			Default: "Gin",
		}
		if err := survey.AskOne(frameworkPrompt, &config.Framework); err != nil {
			return fmt.Errorf("ошибка выбора фреймворка: %w", err)
		}
	}

	// Выбор БД
	if database != "" {
		config.Database = database
	} else {
		dbPrompt := &survey.Select{
			Message: "Выберите базу данных:",
			Options: []string{"PostgreSQL", "MySQL", "MongoDB", "In-Memory", "Без БД"},
			Default: "PostgreSQL",
		}
		if err := survey.AskOne(dbPrompt, &config.Database); err != nil {
			return fmt.Errorf("ошибка выбора БД: %w", err)
		}
	}

	// gRPC
	if cmd.Flags().Changed("grpc") {
		config.EnableGRPC = enableGRPC
	} else {
		grpcPrompt := &survey.Confirm{
			Message: "Включить gRPC сервер?",
			Default: false,
		}
		if err := survey.AskOne(grpcPrompt, &config.EnableGRPC); err != nil {
			return fmt.Errorf("ошибка выбора gRPC: %w", err)
		}
	}

	// Путь для создания проекта
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("ошибка получения текущей директории: %w", err)
	}
	config.Path = filepath.Join(currentDir, config.Name)

	// Генерируем проект
	fmt.Printf("\n🚀 Создание проекта '%s' в %s\n", config.Name, config.Path)
	fmt.Printf("📦 Framework: %s\n", config.Framework)
	fmt.Printf("🗄️  Database: %s\n", config.Database)
	fmt.Printf("🌐 gRPC: %t\n", config.EnableGRPC)

	generator := generator.New(config.Path)
	if err := generator.Generate(config); err != nil {
		return fmt.Errorf("ошибка генерации проекта: %w", err)
	}

	fmt.Printf("\n✅ Проект успешно создан!\n")
	fmt.Printf("📁 Директория: %s\n", config.Path)
	fmt.Printf("\nДля начала работы:\n")
	fmt.Printf("  cd %s\n", config.Name)
	fmt.Printf("  go mod tidy\n")
	fmt.Printf("  make run\n")

	return nil
}
