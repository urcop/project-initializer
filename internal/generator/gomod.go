package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// generateGoMod создает go.mod файл для проекта
func (g *Generator) generateGoMod(config *ProjectConfig) error {
	dependencies := []string{
		"gopkg.in/yaml.v3 v3.0.1",
		"github.com/swaggo/swag v1.16.2",
	}

	// Добавляем зависимости в зависимости от фреймворка
	switch strings.ToLower(config.Framework) {
	case "gin":
		dependencies = append(dependencies,
			"github.com/gin-gonic/gin v1.9.1",
			"github.com/swaggo/gin-swagger v1.6.0",
			"github.com/swaggo/files v1.0.1",
		)
	case "fiber":
		dependencies = append(dependencies,
			"github.com/gofiber/fiber/v2 v2.52.0",
			"github.com/gofiber/swagger v1.0.0",
		)
	case "echo":
		dependencies = append(dependencies,
			"github.com/labstack/echo/v4 v4.11.4",
			"github.com/swaggo/echo-swagger v1.4.1",
		)
	}

	// Добавляем зависимости для БД
	switch strings.ToLower(config.Database) {
	case "postgresql":
		dependencies = append(dependencies,
			"github.com/lib/pq v1.10.9",
			"gorm.io/gorm v1.25.5",
			"gorm.io/driver/postgres v1.5.4",
		)
	case "mysql":
		dependencies = append(dependencies,
			"github.com/go-sql-driver/mysql v1.7.1",
			"gorm.io/gorm v1.25.5",
			"gorm.io/driver/mysql v1.5.2",
		)
	case "mongodb":
		dependencies = append(dependencies,
			"go.mongodb.org/mongo-driver v1.13.1",
		)
	case "in-memory":
		dependencies = append(dependencies,
			"gorm.io/gorm v1.25.5",
			"gorm.io/driver/sqlite v1.5.4",
		)
	}

	// Добавляем gRPC зависимости если включен
	if config.EnableGRPC {
		dependencies = append(dependencies,
			"google.golang.org/grpc v1.60.1",
			"google.golang.org/protobuf v1.31.0",
		)
	}

	// Общие зависимости
	dependencies = append(dependencies,
		"github.com/sirupsen/logrus v1.9.3",
		"github.com/joho/godotenv v1.4.0",
	)

	// Формируем содержимое go.mod файла
	content := fmt.Sprintf(`module %s

go 1.21

require (
`, config.ModuleName)

	for _, dep := range dependencies {
		content += fmt.Sprintf("\t%s\n", dep)
	}

	content += ")\n"

	// Записываем файл
	goModPath := filepath.Join(g.projectPath, "go.mod")
	return os.WriteFile(goModPath, []byte(content), 0644)
}
