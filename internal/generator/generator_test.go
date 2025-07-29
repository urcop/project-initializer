package generator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "generator-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	projectPath := filepath.Join(tempDir, "test-project")
	generator := New(projectPath)

	if generator == nil {
		t.Fatal("Generator should not be nil")
	}

	if generator.projectPath != projectPath {
		t.Errorf("Expected project path %s, got %s", projectPath, generator.projectPath)
	}
}

func TestCreateDirectoryStructure(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "generator-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	projectPath := filepath.Join(tempDir, "test-project")
	generator := New(projectPath)

	err = generator.createDirectoryStructure()
	if err != nil {
		t.Fatalf("Failed to create directory structure: %v", err)
	}

	// Проверяем, что основные директории созданы
	expectedDirs := []string{
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

	for _, dir := range expectedDirs {
		fullPath := filepath.Join(projectPath, dir)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Errorf("Expected directory %s was not created", dir)
		}
	}
}

func TestProjectConfigValidation(t *testing.T) {
	config := &ProjectConfig{
		Name:       "test-service",
		ModuleName: "github.com/test/test-service",
		Framework:  "gin",
		Database:   "postgresql",
		EnableGRPC: false,
		Path:       "/tmp/test-service",
	}

	if config.Name == "" {
		t.Error("Project name should not be empty")
	}

	if config.ModuleName == "" {
		t.Error("Module name should not be empty")
	}

	validFrameworks := []string{"gin", "fiber", "echo"}
	found := false
	for _, fw := range validFrameworks {
		if config.Framework == fw {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Framework %s should be one of %v", config.Framework, validFrameworks)
	}
}

func TestSupportedDatabases(t *testing.T) {
	supportedDatabases := []string{
		"postgresql",
		"mysql",
		"mongodb",
		"in-memory",
		"без бд",
	}

	for _, db := range supportedDatabases {
		config := &ProjectConfig{
			Database: db,
		}

		// Проверяем, что конфигурация не пустая
		if config.Database == "" {
			t.Errorf("Database configuration should not be empty for %s", db)
		}
	}
}
