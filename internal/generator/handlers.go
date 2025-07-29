package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// generateHandlers создает HTTP handlers
func (g *Generator) generateHandlers(config *ProjectConfig) error {
	// Создаем базовый handler
	if err := g.generateBaseHandler(config); err != nil {
		return err
	}

	// Создаем health handler
	if err := g.generateHealthHandler(config); err != nil {
		return err
	}

	// Создаем middleware
	if err := g.generateMiddleware(config); err != nil {
		return err
	}

	return nil
}

// generateBaseHandler создает основной handler файл
func (g *Generator) generateBaseHandler(config *ProjectConfig) error {
	var content string

	switch strings.ToLower(config.Framework) {
	case "gin":
		content = g.generateGinHandler(config)
	case "fiber":
		content = g.generateFiberHandler(config)
	case "echo":
		content = g.generateEchoHandler(config)
	default:
		content = g.generateGinHandler(config) // По умолчанию Gin
	}

	handlerPath := filepath.Join(g.projectPath, "internal/handlers/handler.go")
	return os.WriteFile(handlerPath, []byte(content), 0644)
}

// generateGinHandler генерирует handler для Gin
func (g *Generator) generateGinHandler(config *ProjectConfig) string {
	content := fmt.Sprintf(`package handlers

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"%s/internal/config"
	"%s/pkg/logger"`, config.ModuleName, config.ModuleName)

	if !strings.Contains(strings.ToLower(config.Database), "без") {
		content += fmt.Sprintf(`
	"%s/pkg/database"`, config.ModuleName)
	}

	content += `
)

// Handler представляет HTTP handler
type Handler struct {
	cfg    *config.Config
	logger logger.Logger`

	if !strings.Contains(strings.ToLower(config.Database), "без") {
		content += `
	db     database.Database`
	}

	content += `
}

// New создает новый handler
func New(cfg *config.Config, logger logger.Logger`

	if !strings.Contains(strings.ToLower(config.Database), "без") {
		content += `, db database.Database`
	}

	content += `) *Handler {
	return &Handler{
		cfg:    cfg,
		logger: logger,`

	if !strings.Contains(strings.ToLower(config.Database), "без") {
		content += `
		db:     db,`
	}

	content += `
	}
}

// SetupRoutes настраивает маршруты
func (h *Handler) SetupRoutes() *gin.Engine {
	if h.cfg.App.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	
	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Health check
	router.GET("/health", h.HealthCheck)

	// API группа
	api := router.Group("/api/v1")
	{
		// Здесь будут API маршруты
		api.GET("/ping", h.Ping)
	}

	// Swagger
	if h.cfg.Swagger.Enabled {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	return router
}

// Ping простой ping endpoint
// @Summary Ping
// @Description Проверка работоспособности API
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/v1/ping [get]
func (h *Handler) Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
		"service": h.cfg.App.Name,
		"version": h.cfg.App.Version,
	})
}
`

	return content
}

// generateFiberHandler генерирует handler для Fiber
func (g *Generator) generateFiberHandler(config *ProjectConfig) string {
	content := fmt.Sprintf(`package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"

	"%s/internal/config"
	applogger "%s/pkg/logger"`, config.ModuleName, config.ModuleName)

	if !strings.Contains(strings.ToLower(config.Database), "без") {
		content += fmt.Sprintf(`
	"%s/pkg/database"`, config.ModuleName)
	}

	content += `
)

// Handler представляет HTTP handler
type Handler struct {
	cfg    *config.Config
	logger applogger.Logger`

	if !strings.Contains(strings.ToLower(config.Database), "без") {
		content += `
	db     database.Database`
	}

	content += `
}

// New создает новый handler
func New(cfg *config.Config, logger applogger.Logger`

	if !strings.Contains(strings.ToLower(config.Database), "без") {
		content += `, db database.Database`
	}

	content += `) *Handler {
	return &Handler{
		cfg:    cfg,
		logger: logger,`

	if !strings.Contains(strings.ToLower(config.Database), "без") {
		content += `
		db:     db,`
	}

	content += `
	}
}

// SetupRoutes настраивает маршруты
func (h *Handler) SetupRoutes() *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: h.cfg.App.Name,
	})

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())

	// Health check
	app.Get("/health", h.HealthCheck)

	// API группа
	api := app.Group("/api/v1")
	{
		// Здесь будут API маршруты
		api.Get("/ping", h.Ping)
	}

	// Swagger
	if h.cfg.Swagger.Enabled {
		app.Get("/swagger/*", swagger.HandlerDefault)
	}

	return app
}

// Ping простой ping endpoint
// @Summary Ping
// @Description Проверка работоспособности API
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/v1/ping [get]
func (h *Handler) Ping(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "pong",
		"service": h.cfg.App.Name,
		"version": h.cfg.App.Version,
	})
}
`

	return content
}

// generateEchoHandler генерирует handler для Echo
func (g *Generator) generateEchoHandler(config *ProjectConfig) string {
	content := fmt.Sprintf(`package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	"%s/internal/config"
	"%s/pkg/logger"`, config.ModuleName, config.ModuleName)

	if !strings.Contains(strings.ToLower(config.Database), "без") {
		content += fmt.Sprintf(`
	"%s/pkg/database"`, config.ModuleName)
	}

	content += `
)

// Handler представляет HTTP handler
type Handler struct {
	cfg    *config.Config
	logger logger.Logger`

	if !strings.Contains(strings.ToLower(config.Database), "без") {
		content += `
	db     database.Database`
	}

	content += `
}

// New создает новый handler
func New(cfg *config.Config, logger logger.Logger`

	if !strings.Contains(strings.ToLower(config.Database), "без") {
		content += `, db database.Database`
	}

	content += `) *Handler {
	return &Handler{
		cfg:    cfg,
		logger: logger,`

	if !strings.Contains(strings.ToLower(config.Database), "без") {
		content += `
		db:     db,`
	}

	content += `
	}
}

// SetupRoutes настраивает маршруты
func (h *Handler) SetupRoutes() *echo.Echo {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Health check
	e.GET("/health", h.HealthCheck)

	// API группа
	api := e.Group("/api/v1")
	{
		// Здесь будут API маршруты
		api.GET("/ping", h.Ping)
	}

	// Swagger
	if h.cfg.Swagger.Enabled {
		e.GET("/swagger/*", echoSwagger.WrapHandler)
	}

	return e
}

// Ping простой ping endpoint
// @Summary Ping
// @Description Проверка работоспособности API
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/v1/ping [get]
func (h *Handler) Ping(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "pong",
		"service": h.cfg.App.Name,
		"version": h.cfg.App.Version,
	})
}
`

	return content
}

// generateHealthHandler создает health handler
func (g *Generator) generateHealthHandler(config *ProjectConfig) error {
	var content string

	switch strings.ToLower(config.Framework) {
	case "gin":
		content = g.generateGinHealthHandler(config)
	case "fiber":
		content = g.generateFiberHealthHandler(config)
	case "echo":
		content = g.generateEchoHealthHandler(config)
	default:
		content = g.generateGinHealthHandler(config)
	}

	healthPath := filepath.Join(g.projectPath, "internal/handlers/health.go")
	return os.WriteFile(healthPath, []byte(content), 0644)
}

// generateGinHealthHandler генерирует health handler для Gin
func (g *Generator) generateGinHealthHandler(config *ProjectConfig) string {
	content := `package handlers

import (
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthResponse представляет ответ health check
type HealthResponse struct {
	Status    string            ` + "`" + `json:"status"` + "`" + `
	Service   string            ` + "`" + `json:"service"` + "`" + `
	Version   string            ` + "`" + `json:"version"` + "`" + `
	Timestamp time.Time         ` + "`" + `json:"timestamp"` + "`" + `
	Uptime    string            ` + "`" + `json:"uptime"` + "`" + `
	System    SystemInfo        ` + "`" + `json:"system"` + "`" + `
	Database  *DatabaseStatus   ` + "`" + `json:"database,omitempty"` + "`" + `
}

// SystemInfo информация о системе
type SystemInfo struct {
	GoVersion  string ` + "`" + `json:"go_version"` + "`" + `
	NumCPU     int    ` + "`" + `json:"num_cpu"` + "`" + `
	Goroutines int    ` + "`" + `json:"goroutines"` + "`" + `
}

// DatabaseStatus статус базы данных
type DatabaseStatus struct {
	Connected bool   ` + "`" + `json:"connected"` + "`" + `
	Error     string ` + "`" + `json:"error,omitempty"` + "`" + `
}

var startTime = time.Now()

// HealthCheck проверка работоспособности сервиса
// @Summary Health Check
// @Description Проверка работоспособности и статуса сервиса
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func (h *Handler) HealthCheck(c *gin.Context) {
	response := HealthResponse{
		Status:    "ok",
		Service:   h.cfg.App.Name,
		Version:   h.cfg.App.Version,
		Timestamp: time.Now(),
		Uptime:    time.Since(startTime).String(),
		System: SystemInfo{
			GoVersion:  runtime.Version(),
			NumCPU:     runtime.NumCPU(),
			Goroutines: runtime.NumGoroutine(),
		},
	}`

	if !strings.Contains(strings.ToLower(config.Database), "без") {
		content += `

	// Проверяем статус БД
	dbStatus := &DatabaseStatus{Connected: true}
	if err := h.db.Ping(); err != nil {
		dbStatus.Connected = false
		dbStatus.Error = err.Error()
		response.Status = "degraded"
	}
	response.Database = dbStatus`
	}

	content += `

	status := http.StatusOK
	if response.Status != "ok" {
		status = http.StatusServiceUnavailable
	}

	c.JSON(status, response)
}
`

	return content
}

// generateFiberHealthHandler генерирует health handler для Fiber
func (g *Generator) generateFiberHealthHandler(config *ProjectConfig) string {
	content := `package handlers

import (
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
)

// HealthResponse представляет ответ health check
type HealthResponse struct {
	Status    string            ` + "`" + `json:"status"` + "`" + `
	Service   string            ` + "`" + `json:"service"` + "`" + `
	Version   string            ` + "`" + `json:"version"` + "`" + `
	Timestamp time.Time         ` + "`" + `json:"timestamp"` + "`" + `
	Uptime    string            ` + "`" + `json:"uptime"` + "`" + `
	System    SystemInfo        ` + "`" + `json:"system"` + "`" + `
	Database  *DatabaseStatus   ` + "`" + `json:"database,omitempty"` + "`" + `
}

// SystemInfo информация о системе
type SystemInfo struct {
	GoVersion  string ` + "`" + `json:"go_version"` + "`" + `
	NumCPU     int    ` + "`" + `json:"num_cpu"` + "`" + `
	Goroutines int    ` + "`" + `json:"goroutines"` + "`" + `
}

// DatabaseStatus статус базы данных
type DatabaseStatus struct {
	Connected bool   ` + "`" + `json:"connected"` + "`" + `
	Error     string ` + "`" + `json:"error,omitempty"` + "`" + `
}

var startTime = time.Now()

// HealthCheck проверка работоспособности сервиса
// @Summary Health Check
// @Description Проверка работоспособности и статуса сервиса
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func (h *Handler) HealthCheck(c *fiber.Ctx) error {
	response := HealthResponse{
		Status:    "ok",
		Service:   h.cfg.App.Name,
		Version:   h.cfg.App.Version,
		Timestamp: time.Now(),
		Uptime:    time.Since(startTime).String(),
		System: SystemInfo{
			GoVersion:  runtime.Version(),
			NumCPU:     runtime.NumCPU(),
			Goroutines: runtime.NumGoroutine(),
		},
	}`

	if !strings.Contains(strings.ToLower(config.Database), "без") {
		content += `

	// Проверяем статус БД
	dbStatus := &DatabaseStatus{Connected: true}
	if err := h.db.Ping(); err != nil {
		dbStatus.Connected = false
		dbStatus.Error = err.Error()
		response.Status = "degraded"
	}
	response.Database = dbStatus`
	}

	content += `

	status := fiber.StatusOK
	if response.Status != "ok" {
		status = fiber.StatusServiceUnavailable
	}

	return c.Status(status).JSON(response)
}
`

	return content
}

// generateEchoHealthHandler генерирует health handler для Echo
func (g *Generator) generateEchoHealthHandler(config *ProjectConfig) string {
	content := `package handlers

import (
	"net/http"
	"runtime"
	"time"

	"github.com/labstack/echo/v4"
)

// HealthResponse представляет ответ health check
type HealthResponse struct {
	Status    string            ` + "`" + `json:"status"` + "`" + `
	Service   string            ` + "`" + `json:"service"` + "`" + `
	Version   string            ` + "`" + `json:"version"` + "`" + `
	Timestamp time.Time         ` + "`" + `json:"timestamp"` + "`" + `
	Uptime    string            ` + "`" + `json:"uptime"` + "`" + `
	System    SystemInfo        ` + "`" + `json:"system"` + "`" + `
	Database  *DatabaseStatus   ` + "`" + `json:"database,omitempty"` + "`" + `
}

// SystemInfo информация о системе
type SystemInfo struct {
	GoVersion  string ` + "`" + `json:"go_version"` + "`" + `
	NumCPU     int    ` + "`" + `json:"num_cpu"` + "`" + `
	Goroutines int    ` + "`" + `json:"goroutines"` + "`" + `
}

// DatabaseStatus статус базы данных
type DatabaseStatus struct {
	Connected bool   ` + "`" + `json:"connected"` + "`" + `
	Error     string ` + "`" + `json:"error,omitempty"` + "`" + `
}

var startTime = time.Now()

// HealthCheck проверка работоспособности сервиса
// @Summary Health Check
// @Description Проверка работоспособности и статуса сервиса
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func (h *Handler) HealthCheck(c echo.Context) error {
	response := HealthResponse{
		Status:    "ok",
		Service:   h.cfg.App.Name,
		Version:   h.cfg.App.Version,
		Timestamp: time.Now(),
		Uptime:    time.Since(startTime).String(),
		System: SystemInfo{
			GoVersion:  runtime.Version(),
			NumCPU:     runtime.NumCPU(),
			Goroutines: runtime.NumGoroutine(),
		},
	}`

	if !strings.Contains(strings.ToLower(config.Database), "без") {
		content += `

	// Проверяем статус БД
	dbStatus := &DatabaseStatus{Connected: true}
	if err := h.db.Ping(); err != nil {
		dbStatus.Connected = false
		dbStatus.Error = err.Error()
		response.Status = "degraded"
	}
	response.Database = dbStatus`
	}

	content += `

	status := http.StatusOK
	if response.Status != "ok" {
		status = http.StatusServiceUnavailable
	}

	return c.JSON(status, response)
}
`

	return content
}

// generateMiddleware создает middleware файлы
func (g *Generator) generateMiddleware(config *ProjectConfig) error {
	content := fmt.Sprintf(`package middleware

import (
	"time"

	"%s/pkg/logger"
)

// LoggerMiddleware middleware для логирования
func LoggerMiddleware(logger logger.Logger) interface{} {
	// TODO: Реализовать middleware для выбранного фреймворка
	return nil
}

// AuthMiddleware middleware для аутентификации
func AuthMiddleware(logger logger.Logger) interface{} {
	// TODO: Реализовать аутентификацию
	return nil
}

// CORSMiddleware middleware для CORS
func CORSMiddleware() interface{} {
	// TODO: Реализовать CORS
	return nil
}

// RateLimitMiddleware middleware для ограничения запросов
func RateLimitMiddleware(requests int, window time.Duration) interface{} {
	// TODO: Реализовать rate limiting
	return nil
}
`, config.ModuleName)

	middlewarePath := filepath.Join(g.projectPath, "internal/middleware/middleware.go")
	return os.WriteFile(middlewarePath, []byte(content), 0644)
}
