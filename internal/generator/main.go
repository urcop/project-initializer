package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// generateMainFiles создает основные файлы приложения
func (g *Generator) generateMainFiles(config *ProjectConfig) error {
	// Создаем main.go
	if err := g.generateMainGo(config); err != nil {
		return err
	}

	// Создаем app.go
	if err := g.generateAppGo(config); err != nil {
		return err
	}

	// Создаем контекст
	if err := g.generateContext(config); err != nil {
		return err
	}

	// Создаем логгер
	if err := g.generateLogger(config); err != nil {
		return err
	}

	return nil
}

// generateMainGo создает файл cmd/main.go
func (g *Generator) generateMainGo(config *ProjectConfig) error {
	content := fmt.Sprintf(`package main

import (
	"log"

	"%s/internal/config"
	"%s/internal/app"
)

// @title %s API
// @version 1.0
// @description API документация для %s
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Загружаем конфигурацию
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %%v", err)
	}

	// Создаем и запускаем приложение
	application := app.New(cfg)
	if err := application.Run(); err != nil {
		log.Fatalf("Ошибка запуска приложения: %%v", err)
	}
}
`, config.ModuleName, config.ModuleName, config.Name, config.Name)

	mainPath := filepath.Join(g.projectPath, "cmd/main.go")
	return os.WriteFile(mainPath, []byte(content), 0644)
}

// generateAppGo создает файл internal/app/app.go
func (g *Generator) generateAppGo(config *ProjectConfig) error {
	var content string

	if strings.ToLower(config.Framework) == "fiber" {
		content = g.generateFiberApp(config)
	} else {
		content = g.generateStandardApp(config)
	}

	appPath := filepath.Join(g.projectPath, "internal/app")
	if err := os.MkdirAll(appPath, 0755); err != nil {
		return err
	}

	appGoPath := filepath.Join(appPath, "app.go")
	return os.WriteFile(appGoPath, []byte(content), 0644)
}

// generateFiberApp генерирует app.go для Fiber
func (g *Generator) generateFiberApp(config *ProjectConfig) string {
	content := fmt.Sprintf(`package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"%s/internal/config"
	"%s/internal/handlers"
	"%s/pkg/logger"`, config.ModuleName, config.ModuleName, config.ModuleName)

	if !strings.Contains(strings.ToLower(config.Database), "без") {
		content += fmt.Sprintf(`
	"%s/pkg/database"`, config.ModuleName)
	}

	content += `
)

// App представляет приложение
type App struct {
	cfg    *config.Config
	logger logger.Logger
	app    *fiber.App
}

// New создает новое приложение
func New(cfg *config.Config) *App {
	// Инициализируем логгер
	loggerConfig := logger.LoggerConfig{
		Level:  cfg.Logger.Level,
		Format: cfg.Logger.Format,
		Output: cfg.Logger.Output,
	}
	log := logger.New(loggerConfig)

	return &App{
		cfg:    cfg,
		logger: log,
	}
}

// Run запускает приложение
func (a *App) Run() error {`

	if !strings.Contains(strings.ToLower(config.Database), "без") {
		content += `
	// Инициализируем базу данных
	db, err := database.New(a.cfg)
	if err != nil {
		return fmt.Errorf("ошибка подключения к БД: %w", err)
	}
	defer db.Close()
	
	a.logger.Info("Подключение к базе данных установлено")

	// Создаем Fiber приложение
	handler := handlers.New(a.cfg, a.logger, db)`
	} else {
		content += `
	// Создаем Fiber приложение
	handler := handlers.New(a.cfg, a.logger)`
	}

	content += `
	a.app = handler.SetupRoutes()

	// Запускаем Fiber сервер в горутине
	go func() {
		a.logger.Info("Запуск HTTP сервера", "port", a.cfg.App.Port)
		if err := a.app.Listen(fmt.Sprintf(":%d", a.cfg.App.Port)); err != nil {
			a.logger.Error("Ошибка HTTP сервера", "error", err)
		}
	}()

	// Ожидаем сигналы завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	a.logger.Info("Получен сигнал завершения, останавливаем сервер...")

	// Graceful shutdown
	if err := a.app.Shutdown(); err != nil {
		return fmt.Errorf("ошибка остановки сервера: %w", err)
	}

	a.logger.Info("Сервер успешно остановлен")
	return nil
}
`

	return content
}

// generateStandardApp генерирует app.go для Gin/Echo
func (g *Generator) generateStandardApp(config *ProjectConfig) string {
	content := fmt.Sprintf(`package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"%s/internal/config"
	"%s/internal/handlers"
	"%s/pkg/logger"`, config.ModuleName, config.ModuleName, config.ModuleName)

	if !strings.Contains(strings.ToLower(config.Database), "без") {
		content += fmt.Sprintf(`
	"%s/pkg/database"`, config.ModuleName)
	}

	content += `
)

// App представляет приложение
type App struct {
	cfg    *config.Config
	logger logger.Logger
	server *http.Server
}

// New создает новое приложение
func New(cfg *config.Config) *App {
	// Инициализируем логгер
	loggerConfig := logger.LoggerConfig{
		Level:  cfg.Logger.Level,
		Format: cfg.Logger.Format,
		Output: cfg.Logger.Output,
	}
	log := logger.New(loggerConfig)

	return &App{
		cfg:    cfg,
		logger: log,
	}
}

// Run запускает приложение
func (a *App) Run() error {`

	if !strings.Contains(strings.ToLower(config.Database), "без") {
		content += `
	// Инициализируем базу данных
	db, err := database.New(a.cfg)
	if err != nil {
		return fmt.Errorf("ошибка подключения к БД: %w", err)
	}
	defer db.Close()
	
	a.logger.Info("Подключение к базе данных установлено")

	// Создаем HTTP сервер
	handler := handlers.New(a.cfg, a.logger, db)`
	} else {
		content += `
	// Создаем HTTP сервер
	handler := handlers.New(a.cfg, a.logger)`
	}

	content += `
	
	a.server = &http.Server{
		Addr:         fmt.Sprintf(":%d", a.cfg.App.Port),
		Handler:      handler.SetupRoutes(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Запускаем HTTP сервер в горутине
	go func() {
		a.logger.Info("Запуск HTTP сервера", "port", a.cfg.App.Port)
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Error("Ошибка HTTP сервера", "error", err)
		}
	}()

	// Ожидаем сигналы завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	a.logger.Info("Получен сигнал завершения, останавливаем сервер...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("ошибка остановки сервера: %w", err)
	}

	a.logger.Info("Сервер успешно остановлен")
	return nil
}
`

	return content
}

// generateContext создает файл pkg/context/context.go
func (g *Generator) generateContext(config *ProjectConfig) error {
	content := fmt.Sprintf(`package context

import (
	"context"
	"time"

	"%s/pkg/logger"
)

// AppContext представляет контекст приложения
type AppContext struct {
	ctx    context.Context
	logger logger.Logger
	userID string
	traceID string
}

// New создает новый контекст приложения
func New(ctx context.Context, logger logger.Logger) *AppContext {
	return &AppContext{
		ctx:    ctx,
		logger: logger,
	}
}

// Context возвращает базовый контекст
func (c *AppContext) Context() context.Context {
	return c.ctx
}

// Logger возвращает логгер
func (c *AppContext) Logger() logger.Logger {
	return c.logger
}

// WithUserID устанавливает ID пользователя
func (c *AppContext) WithUserID(userID string) *AppContext {
	newCtx := *c
	newCtx.userID = userID
	return &newCtx
}

// UserID возвращает ID пользователя
func (c *AppContext) UserID() string {
	return c.userID
}

// WithTraceID устанавливает ID трассировки
func (c *AppContext) WithTraceID(traceID string) *AppContext {
	newCtx := *c
	newCtx.traceID = traceID
	return &newCtx
}

// TraceID возвращает ID трассировки
func (c *AppContext) TraceID() string {
	return c.traceID
}

// WithTimeout создает контекст с таймаутом
func (c *AppContext) WithTimeout(timeout time.Duration) (*AppContext, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(c.ctx, timeout)
	newCtx := *c
	newCtx.ctx = ctx
	return &newCtx, cancel
}

// WithDeadline создает контекст с дедлайном
func (c *AppContext) WithDeadline(deadline time.Time) (*AppContext, context.CancelFunc) {
	ctx, cancel := context.WithDeadline(c.ctx, deadline)
	newCtx := *c
	newCtx.ctx = ctx
	return &newCtx, cancel
}

// Done возвращает канал завершения
func (c *AppContext) Done() <-chan struct{} {
	return c.ctx.Done()
}

// Err возвращает ошибку контекста
func (c *AppContext) Err() error {
	return c.ctx.Err()
}
`, config.ModuleName)

	contextPath := filepath.Join(g.projectPath, "pkg/context/context.go")
	return os.WriteFile(contextPath, []byte(content), 0644)
}

// generateLogger создает файл pkg/logger/logger.go
func (g *Generator) generateLogger(config *ProjectConfig) error {
	content := `package logger

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

// Logger интерфейс для логгирования
type Logger interface {
	Debug(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Fatal(msg string, fields ...interface{})
	WithField(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger
}

// LogrusLogger реализация Logger на основе logrus
type LogrusLogger struct {
	entry *logrus.Entry
}

// LoggerConfig конфигурация логгера
type LoggerConfig struct {
	Level  string
	Format string
	Output string
}

// New создает новый логгер
func New(config LoggerConfig) Logger {
	log := logrus.New()

	// Устанавливаем уровень
	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	log.SetLevel(level)

	// Устанавливаем формат
	switch config.Format {
	case "json":
		log.SetFormatter(&logrus.JSONFormatter{})
	default:
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	// Устанавливаем вывод
	switch config.Output {
	case "stdout":
		log.SetOutput(os.Stdout)
	case "stderr":
		log.SetOutput(os.Stderr)
	default:
		log.SetOutput(os.Stdout)
	}

	return &LogrusLogger{
		entry: logrus.NewEntry(log),
	}
}

// Debug логирует отладочное сообщение
func (l *LogrusLogger) Debug(msg string, fields ...interface{}) {
	l.entry.WithFields(l.parseFields(fields...)).Debug(msg)
}

// Info логирует информационное сообщение
func (l *LogrusLogger) Info(msg string, fields ...interface{}) {
	l.entry.WithFields(l.parseFields(fields...)).Info(msg)
}

// Warn логирует предупреждение
func (l *LogrusLogger) Warn(msg string, fields ...interface{}) {
	l.entry.WithFields(l.parseFields(fields...)).Warn(msg)
}

// Error логирует ошибку
func (l *LogrusLogger) Error(msg string, fields ...interface{}) {
	l.entry.WithFields(l.parseFields(fields...)).Error(msg)
}

// Fatal логирует фатальную ошибку
func (l *LogrusLogger) Fatal(msg string, fields ...interface{}) {
	l.entry.WithFields(l.parseFields(fields...)).Fatal(msg)
}

// WithField добавляет поле к логгеру
func (l *LogrusLogger) WithField(key string, value interface{}) Logger {
	return &LogrusLogger{
		entry: l.entry.WithField(key, value),
	}
}

// WithFields добавляет поля к логгеру
func (l *LogrusLogger) WithFields(fields map[string]interface{}) Logger {
	return &LogrusLogger{
		entry: l.entry.WithFields(fields),
	}
}

// parseFields парсит поля из slice
func (l *LogrusLogger) parseFields(fields ...interface{}) logrus.Fields {
	parsed := make(logrus.Fields)
	
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			key := fmt.Sprintf("%v", fields[i])
			parsed[key] = fields[i+1]
		}
	}
	
	return parsed
}
`

	loggerPath := filepath.Join(g.projectPath, "pkg/logger/logger.go")
	return os.WriteFile(loggerPath, []byte(content), 0644)
}
