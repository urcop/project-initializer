package generator

import (
	"fmt"
	"os"
	"path/filepath"
)

// ProjectConfig представляет конфигурацию проекта
type ProjectConfig struct {
	Name       string
	ModuleName string
	Framework  string
	Database   string
	EnableGRPC bool
	Path       string
}

// Generator отвечает за генерацию структуры проекта
type Generator struct {
	projectPath string
}

// New создает новый экземпляр генератора
func New(projectPath string) *Generator {
	return &Generator{
		projectPath: projectPath,
	}
}

// Generate генерирует весь проект на основе конфигурации
func (g *Generator) Generate(config *ProjectConfig) error {
	// Создаем базовую структуру директорий
	if err := g.createDirectoryStructure(); err != nil {
		return fmt.Errorf("ошибка создания структуры директорий: %w", err)
	}

	// Создаем go.mod
	if err := g.generateGoMod(config); err != nil {
		return fmt.Errorf("ошибка создания go.mod: %w", err)
	}

	// Создаем конфигурационные файлы
	if err := g.generateConfig(config); err != nil {
		return fmt.Errorf("ошибка создания конфига: %w", err)
	}

	// Создаем основные файлы приложения
	if err := g.generateMainFiles(config); err != nil {
		return fmt.Errorf("ошибка создания основных файлов: %w", err)
	}

	// Создаем handlers
	if err := g.generateHandlers(config); err != nil {
		return fmt.Errorf("ошибка создания handlers: %w", err)
	}

	// Создаем слой БД
	if err := g.generateDatabaseLayer(config); err != nil {
		return fmt.Errorf("ошибка создания слоя БД: %w", err)
	}

	// Создаем gRPC если нужно
	if config.EnableGRPC {
		if err := g.generateGRPC(config); err != nil {
			return fmt.Errorf("ошибка создания gRPC: %w", err)
		}
	}

	// Создаем Docker файлы
	if err := g.generateDockerFiles(config); err != nil {
		return fmt.Errorf("ошибка создания Docker файлов: %w", err)
	}

	// Создаем Makefile
	if err := g.generateMakefile(config); err != nil {
		return fmt.Errorf("ошибка создания Makefile: %w", err)
	}

	return nil
}

// createDirectoryStructure создает базовую структуру директорий
func (g *Generator) createDirectoryStructure() error {
	dirs := []string{
		"cmd",
		"internal/config",
		"internal/handlers",
		"internal/services",
		"internal/repository",
		"internal/models",
		"internal/middleware",
		"pkg/context",
		"pkg/logger",
		"pkg/database",
		"docs",
		"scripts",
		"deployments",
		"api/swagger",
	}

	for _, dir := range dirs {
		fullPath := filepath.Join(g.projectPath, dir)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return fmt.Errorf("ошибка создания директории %s: %w", fullPath, err)
		}
	}

	return nil
}
