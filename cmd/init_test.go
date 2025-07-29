package cmd

import (
	"os"
	"testing"
)

func TestInitCommand(t *testing.T) {
	// Создаем временную директорию для тестов
	tempDir, err := os.MkdirTemp("", "project-initializer-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Меняем рабочую директорию на временную
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working dir: %v", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temp dir: %v", err)
	}

	// Тест создания проекта (будет интерактивный, поэтому пропускаем)
	t.Skip("Skipping interactive test - requires manual input")
}

func TestProjectStructure(t *testing.T) {
	// Проверяем, что существуют необходимые файлы генератора
	requiredFiles := []string{
		"../internal/generator/generator.go",
		"../internal/generator/gomod.go",
		"../internal/generator/config.go",
		"../internal/generator/main.go",
		"../internal/generator/handlers.go",
		"../internal/generator/database.go",
		"../internal/generator/docker.go",
		"../internal/generator/makefile.go",
		"../internal/generator/grpc.go",
	}

	for _, file := range requiredFiles {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			t.Errorf("Required generator file does not exist: %s", file)
		}
	}
}

func TestVersionCommand(t *testing.T) {
	if Version == "" {
		t.Error("Version should not be empty")
	}

	if Version != "1.0.1" {
		t.Errorf("Expected version 1.0.1, got %s", Version)
	}
}
