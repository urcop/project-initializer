package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// generateConfig создает конфигурационные файлы
func (g *Generator) generateConfig(config *ProjectConfig) error {
	// Создаем config.yaml
	if err := g.generateConfigYAML(config); err != nil {
		return err
	}

	// Создаем config.go
	if err := g.generateConfigGo(config); err != nil {
		return err
	}

	return nil
}

// generateConfigYAML создает файл config.yaml
func (g *Generator) generateConfigYAML(config *ProjectConfig) error {
	content := fmt.Sprintf(`# Конфигурация для %s
app:
  name: "%s"
  version: "1.0.0"
  debug: true
  port: 8080

`, config.Name, config.Name)

	// Добавляем конфигурацию БД если нужно
	if !strings.Contains(strings.ToLower(config.Database), "без") {
		switch strings.ToLower(config.Database) {
		case "postgresql":
			content += `database:
  type: "postgres"
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "password"
  name: "` + config.Name + `"
  ssl_mode: "disable"
  max_connections: 100
  max_idle_connections: 10

`
		case "mysql":
			content += `database:
  type: "mysql"
  host: "localhost"
  port: 3306
  user: "root"
  password: "password"
  name: "` + config.Name + `"
  charset: "utf8mb4"
  max_connections: 100
  max_idle_connections: 10

`
		case "mongodb":
			content += `database:
  type: "mongodb"
  uri: "mongodb://localhost:27017"
  name: "` + config.Name + `"
  timeout: 30

`
		case "in-memory":
			content += `database:
  type: "sqlite"
  path: ":memory:"
  max_connections: 1

`
		}
	}

	// Добавляем gRPC конфигурацию если включен
	if config.EnableGRPC {
		content += `grpc:
  enabled: true
  port: 9090
  max_connection_age: 30
  max_connection_idle: 30

`
	}

	content += `logger:
  level: "debug"
  format: "json"
  output: "stdout"

swagger:
  enabled: true
  title: "` + config.Name + ` API"
  description: "API документация для ` + config.Name + `"
  version: "1.0.0"
  host: "localhost:8080"
  base_path: "/api/v1"
`

	configPath := filepath.Join(g.projectPath, "config.yaml")
	return os.WriteFile(configPath, []byte(content), 0644)
}

// generateConfigGo создает файл internal/config/config.go
func (g *Generator) generateConfigGo(config *ProjectConfig) error {
	content := fmt.Sprintf(`package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config представляет конфигурацию приложения
type Config struct {
	App      AppConfig      %sconfig:"app" yaml:"app"%s
	Database DatabaseConfig %sconfig:"database" yaml:"database"%s
	Logger   LoggerConfig   %sconfig:"logger" yaml:"logger"%s
	Swagger  SwaggerConfig  %sconfig:"swagger" yaml:"swagger"%s`, "`", "`", "`", "`", "`", "`", "`", "`")

	if config.EnableGRPC {
		content += fmt.Sprintf(`
	GRPC     GRPCConfig     %sconfig:"grpc" yaml:"grpc"%s`, "`", "`")
	}

	content += `
}

// AppConfig конфигурация приложения
type AppConfig struct {
	Name    string ` + "`" + `config:"name" yaml:"name"` + "`" + `
	Version string ` + "`" + `config:"version" yaml:"version"` + "`" + `
	Debug   bool   ` + "`" + `config:"debug" yaml:"debug"` + "`" + `
	Port    int    ` + "`" + `config:"port" yaml:"port"` + "`" + `
}

// DatabaseConfig конфигурация базы данных
type DatabaseConfig struct {
	Type               string ` + "`" + `config:"type" yaml:"type"` + "`" + `
	Host               string ` + "`" + `config:"host" yaml:"host"` + "`" + `
	Port               int    ` + "`" + `config:"port" yaml:"port"` + "`" + `
	User               string ` + "`" + `config:"user" yaml:"user"` + "`" + `
	Password           string ` + "`" + `config:"password" yaml:"password"` + "`" + `
	Name               string ` + "`" + `config:"name" yaml:"name"` + "`" + `
	SSLMode            string ` + "`" + `config:"ssl_mode" yaml:"ssl_mode"` + "`" + `
	URI                string ` + "`" + `config:"uri" yaml:"uri"` + "`" + `
	Path               string ` + "`" + `config:"path" yaml:"path"` + "`" + `
	Charset            string ` + "`" + `config:"charset" yaml:"charset"` + "`" + `
	MaxConnections     int    ` + "`" + `config:"max_connections" yaml:"max_connections"` + "`" + `
	MaxIdleConnections int    ` + "`" + `config:"max_idle_connections" yaml:"max_idle_connections"` + "`" + `
	Timeout            int    ` + "`" + `config:"timeout" yaml:"timeout"` + "`" + `
}

// LoggerConfig конфигурация логгера
type LoggerConfig struct {
	Level  string ` + "`" + `config:"level" yaml:"level"` + "`" + `
	Format string ` + "`" + `config:"format" yaml:"format"` + "`" + `
	Output string ` + "`" + `config:"output" yaml:"output"` + "`" + `
}

// SwaggerConfig конфигурация Swagger
type SwaggerConfig struct {
	Enabled     bool   ` + "`" + `config:"enabled" yaml:"enabled"` + "`" + `
	Title       string ` + "`" + `config:"title" yaml:"title"` + "`" + `
	Description string ` + "`" + `config:"description" yaml:"description"` + "`" + `
	Version     string ` + "`" + `config:"version" yaml:"version"` + "`" + `
	Host        string ` + "`" + `config:"host" yaml:"host"` + "`" + `
	BasePath    string ` + "`" + `config:"base_path" yaml:"base_path"` + "`" + `
}`

	if config.EnableGRPC {
		content += `

// GRPCConfig конфигурация gRPC сервера
type GRPCConfig struct {
	Enabled           bool ` + "`" + `config:"enabled" yaml:"enabled"` + "`" + `
	Port              int  ` + "`" + `config:"port" yaml:"port"` + "`" + `
	MaxConnectionAge  int  ` + "`" + `config:"max_connection_age" yaml:"max_connection_age"` + "`" + `
	MaxConnectionIdle int  ` + "`" + `config:"max_connection_idle" yaml:"max_connection_idle"` + "`" + `
}`
	}

	content += `

// Load загружает конфигурацию из файла
func Load(configPath string) (*Config, error) {
	config := &Config{}

	// Читаем файл конфигурации
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения файла конфигурации: %w", err)
	}

	// Парсим YAML
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("ошибка парсинга конфигурации: %w", err)
	}

	// Валидируем конфигурацию
	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("ошибка валидации конфигурации: %w", err)
	}

	return config, nil
}

// validate проверяет корректность конфигурации
func (c *Config) validate() error {
	if c.App.Name == "" {
		return fmt.Errorf("имя приложения не может быть пустым")
	}

	if c.App.Port <= 0 || c.App.Port > 65535 {
		return fmt.Errorf("порт приложения должен быть в диапазоне 1-65535")
	}

	return nil
}

// GetDSN возвращает строку подключения к БД
func (c *Config) GetDSN() string {
	switch c.Database.Type {
	case "postgres":
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			c.Database.Host, c.Database.Port, c.Database.User, c.Database.Password, c.Database.Name, c.Database.SSLMode)
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
			c.Database.User, c.Database.Password, c.Database.Host, c.Database.Port, c.Database.Name, c.Database.Charset)
	case "mongodb":
		return c.Database.URI
	case "sqlite":
		return c.Database.Path
	default:
		return ""
	}
}
`

	configPath := filepath.Join(g.projectPath, "internal/config/config.go")
	return os.WriteFile(configPath, []byte(content), 0644)
}
