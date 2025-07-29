package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// generateDatabaseLayer создает слой базы данных
func (g *Generator) generateDatabaseLayer(config *ProjectConfig) error {
	// Если БД не нужна, пропускаем
	if strings.Contains(strings.ToLower(config.Database), "без") {
		return nil
	}

	// Создаем интерфейс БД
	if err := g.generateDatabaseInterface(config); err != nil {
		return err
	}

	// Создаем реализацию БД
	if err := g.generateDatabaseImplementation(config); err != nil {
		return err
	}

	// Создаем репозитории
	if err := g.generateRepositories(config); err != nil {
		return err
	}

	// Создаем модели
	if err := g.generateModels(config); err != nil {
		return err
	}

	return nil
}

// generateDatabaseInterface создает интерфейс для работы с БД
func (g *Generator) generateDatabaseInterface(config *ProjectConfig) error {
	content := fmt.Sprintf(`package database

import (
	"context"
	"time"

	appcontext "%s/pkg/context"
)

// Database интерфейс для работы с базой данных
type Database interface {
	// Подключение и отключение
	Connect() error
	Close() error
	Ping() error

	// Транзакции
	BeginTx(ctx context.Context) (Tx, error)

	// Миграции
	Migrate() error

	// Статистика
	Stats() Stats
}

// Tx интерфейс для транзакций
type Tx interface {
	Commit() error
	Rollback() error
	Context() context.Context
}

// Stats статистика подключений к БД
type Stats struct {
	OpenConnections int
	InUseConnections int
	IdleConnections int
}

// Repository базовый интерфейс для репозиториев
type Repository interface {
	SetContext(ctx *appcontext.AppContext)
	GetContext() *appcontext.AppContext
}

// BaseRepository базовая реализация репозитория
type BaseRepository struct {
	ctx *appcontext.AppContext
	db  Database
}

// NewBaseRepository создает новый базовый репозиторий
func NewBaseRepository(db Database) *BaseRepository {
	return &BaseRepository{
		db: db,
	}
}

// SetContext устанавливает контекст
func (r *BaseRepository) SetContext(ctx *appcontext.AppContext) {
	r.ctx = ctx
}

// GetContext возвращает контекст
func (r *BaseRepository) GetContext() *appcontext.AppContext {
	return r.ctx
}

// DB возвращает подключение к БД
func (r *BaseRepository) DB() Database {
	return r.db
}

// Logger возвращает логгер из контекста
func (r *BaseRepository) Logger() interface{} {
	if r.ctx != nil {
		return r.ctx.Logger()
	}
	return nil
}
`, config.ModuleName)

	dbInterfacePath := filepath.Join(g.projectPath, "pkg/database/interface.go")
	return os.WriteFile(dbInterfacePath, []byte(content), 0644)
}

// generateDatabaseImplementation создает реализацию БД
func (g *Generator) generateDatabaseImplementation(config *ProjectConfig) error {
	var content string

	switch strings.ToLower(config.Database) {
	case "postgresql":
		content = g.generatePostgreSQLImplementation(config)
	case "mysql":
		content = g.generateMySQLImplementation(config)
	case "mongodb":
		content = g.generateMongoDBImplementation(config)
	case "in-memory":
		content = g.generateSQLiteImplementation(config)
	default:
		content = g.generatePostgreSQLImplementation(config) // По умолчанию PostgreSQL
	}

	dbImplPath := filepath.Join(g.projectPath, "pkg/database/database.go")
	return os.WriteFile(dbImplPath, []byte(content), 0644)
}

// generatePostgreSQLImplementation генерирует реализацию для PostgreSQL
func (g *Generator) generatePostgreSQLImplementation(config *ProjectConfig) string {
	return fmt.Sprintf(`package database

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"%s/internal/config"
)

// PostgreSQLDatabase реализация для PostgreSQL
type PostgreSQLDatabase struct {
	db     *gorm.DB
	config *config.Config
}

// PostgreSQLTx реализация транзакции для PostgreSQL
type PostgreSQLTx struct {
	tx  *gorm.DB
	ctx context.Context
}

// New создает новое подключение к PostgreSQL
func New(cfg *config.Config) (Database, error) {
	dsn := cfg.GetDSN()
	
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	if !cfg.App.Debug {
		gormConfig.Logger = logger.Default.LogMode(logger.Silent)
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к PostgreSQL: %%w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("ошибка получения sql.DB: %%w", err)
	}

	// Настройка пула соединений
	sqlDB.SetMaxOpenConns(cfg.Database.MaxConnections)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConnections)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return &PostgreSQLDatabase{
		db:     db,
		config: cfg,
	}, nil
}

// Connect подключается к БД
func (p *PostgreSQLDatabase) Connect() error {
	sqlDB, err := p.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// Close закрывает подключение
func (p *PostgreSQLDatabase) Close() error {
	sqlDB, err := p.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Ping проверяет подключение
func (p *PostgreSQLDatabase) Ping() error {
	sqlDB, err := p.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// BeginTx начинает транзакцию
func (p *PostgreSQLDatabase) BeginTx(ctx context.Context) (Tx, error) {
	tx := p.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &PostgreSQLTx{
		tx:  tx,
		ctx: ctx,
	}, nil
}

// Migrate выполняет миграции
func (p *PostgreSQLDatabase) Migrate() error {
	// TODO: Добавить модели для миграции
	// return p.db.AutoMigrate(&User{}, &Product{})
	return nil
}

// Stats возвращает статистику
func (p *PostgreSQLDatabase) Stats() Stats {
	sqlDB, err := p.db.DB()
	if err != nil {
		return Stats{}
	}

	stats := sqlDB.Stats()
	return Stats{
		OpenConnections:  stats.OpenConnections,
		InUseConnections: stats.InUse,
		IdleConnections:  stats.Idle,
	}
}

// DB возвращает GORM DB
func (p *PostgreSQLDatabase) DB() *gorm.DB {
	return p.db
}

// Commit подтверждает транзакцию
func (tx *PostgreSQLTx) Commit() error {
	return tx.tx.Commit().Error
}

// Rollback откатывает транзакцию
func (tx *PostgreSQLTx) Rollback() error {
	return tx.tx.Rollback().Error
}

// Context возвращает контекст транзакции
func (tx *PostgreSQLTx) Context() context.Context {
	return tx.ctx
}
`, config.ModuleName)
}

// generateMySQLImplementation генерирует реализацию для MySQL
func (g *Generator) generateMySQLImplementation(config *ProjectConfig) string {
	return fmt.Sprintf(`package database

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"%s/internal/config"
)

// MySQLDatabase реализация для MySQL
type MySQLDatabase struct {
	db     *gorm.DB
	config *config.Config
}

// MySQLTx реализация транзакции для MySQL
type MySQLTx struct {
	tx  *gorm.DB
	ctx context.Context
}

// New создает новое подключение к MySQL
func New(cfg *config.Config) (Database, error) {
	dsn := cfg.GetDSN()
	
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	if !cfg.App.Debug {
		gormConfig.Logger = logger.Default.LogMode(logger.Silent)
	}

	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к MySQL: %%w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("ошибка получения sql.DB: %%w", err)
	}

	// Настройка пула соединений
	sqlDB.SetMaxOpenConns(cfg.Database.MaxConnections)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConnections)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return &MySQLDatabase{
		db:     db,
		config: cfg,
	}, nil
}

// Connect подключается к БД
func (m *MySQLDatabase) Connect() error {
	sqlDB, err := m.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// Close закрывает подключение
func (m *MySQLDatabase) Close() error {
	sqlDB, err := m.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Ping проверяет подключение
func (m *MySQLDatabase) Ping() error {
	sqlDB, err := m.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// BeginTx начинает транзакцию
func (m *MySQLDatabase) BeginTx(ctx context.Context) (Tx, error) {
	tx := m.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &MySQLTx{
		tx:  tx,
		ctx: ctx,
	}, nil
}

// Migrate выполняет миграции
func (m *MySQLDatabase) Migrate() error {
	// TODO: Добавить модели для миграции
	return nil
}

// Stats возвращает статистику
func (m *MySQLDatabase) Stats() Stats {
	sqlDB, err := m.db.DB()
	if err != nil {
		return Stats{}
	}

	stats := sqlDB.Stats()
	return Stats{
		OpenConnections:  stats.OpenConnections,
		InUseConnections: stats.InUse,
		IdleConnections:  stats.Idle,
	}
}

// DB возвращает GORM DB
func (m *MySQLDatabase) DB() *gorm.DB {
	return m.db
}

// Commit подтверждает транзакцию
func (tx *MySQLTx) Commit() error {
	return tx.tx.Commit().Error
}

// Rollback откатывает транзакцию
func (tx *MySQLTx) Rollback() error {
	return tx.tx.Rollback().Error
}

// Context возвращает контекст транзакции
func (tx *MySQLTx) Context() context.Context {
	return tx.ctx
}
`, config.ModuleName)
}

// generateMongoDBImplementation генерирует реализацию для MongoDB
func (g *Generator) generateMongoDBImplementation(config *ProjectConfig) string {
	return fmt.Sprintf(`package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"%s/internal/config"
)

// MongoDatabase реализация для MongoDB
type MongoDatabase struct {
	client   *mongo.Client
	database *mongo.Database
	config   *config.Config
}

// MongoTx реализация транзакции для MongoDB
type MongoTx struct {
	session mongo.Session
	ctx     context.Context
}

// New создает новое подключение к MongoDB
func New(cfg *config.Config) (Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Database.Timeout)*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Database.URI))
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к MongoDB: %%w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("ошибка ping MongoDB: %%w", err)
	}

	database := client.Database(cfg.Database.Name)

	return &MongoDatabase{
		client:   client,
		database: database,
		config:   cfg,
	}, nil
}

// Connect подключается к БД
func (m *MongoDatabase) Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return m.client.Ping(ctx, nil)
}

// Close закрывает подключение
func (m *MongoDatabase) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return m.client.Disconnect(ctx)
}

// Ping проверяет подключение
func (m *MongoDatabase) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return m.client.Ping(ctx, nil)
}

// BeginTx начинает транзакцию (сессию)
func (m *MongoDatabase) BeginTx(ctx context.Context) (Tx, error) {
	session, err := m.client.StartSession()
	if err != nil {
		return nil, err
	}

	if err := session.StartTransaction(); err != nil {
		session.EndSession(ctx)
		return nil, err
	}

	return &MongoTx{
		session: session,
		ctx:     ctx,
	}, nil
}

// Migrate выполняет миграции (создание индексов)
func (m *MongoDatabase) Migrate() error {
	// TODO: Создать индексы
	return nil
}

// Stats возвращает статистику
func (m *MongoDatabase) Stats() Stats {
	// MongoDB не предоставляет такую статистику напрямую
	return Stats{}
}

// Database возвращает MongoDB Database
func (m *MongoDatabase) Database() *mongo.Database {
	return m.database
}

// Client возвращает MongoDB Client
func (m *MongoDatabase) Client() *mongo.Client {
	return m.client
}

// Commit подтверждает транзакцию
func (tx *MongoTx) Commit() error {
	defer tx.session.EndSession(tx.ctx)
	return tx.session.CommitTransaction(tx.ctx)
}

// Rollback откатывает транзакцию
func (tx *MongoTx) Rollback() error {
	defer tx.session.EndSession(tx.ctx)
	return tx.session.AbortTransaction(tx.ctx)
}

// Context возвращает контекст транзакции
func (tx *MongoTx) Context() context.Context {
	return tx.ctx
}
`, config.ModuleName)
}

// generateSQLiteImplementation генерирует реализацию для SQLite (in-memory)
func (g *Generator) generateSQLiteImplementation(config *ProjectConfig) string {
	return fmt.Sprintf(`package database

import (
	"context"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"%s/internal/config"
)

// SQLiteDatabase реализация для SQLite
type SQLiteDatabase struct {
	db     *gorm.DB
	config *config.Config
}

// SQLiteTx реализация транзакции для SQLite
type SQLiteTx struct {
	tx  *gorm.DB
	ctx context.Context
}

// New создает новое подключение к SQLite
func New(cfg *config.Config) (Database, error) {
	dsn := cfg.GetDSN()
	
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	if !cfg.App.Debug {
		gormConfig.Logger = logger.Default.LogMode(logger.Silent)
	}

	db, err := gorm.Open(sqlite.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к SQLite: %%w", err)
	}

	return &SQLiteDatabase{
		db:     db,
		config: cfg,
	}, nil
}

// Connect подключается к БД
func (s *SQLiteDatabase) Connect() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// Close закрывает подключение
func (s *SQLiteDatabase) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Ping проверяет подключение
func (s *SQLiteDatabase) Ping() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// BeginTx начинает транзакцию
func (s *SQLiteDatabase) BeginTx(ctx context.Context) (Tx, error) {
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &SQLiteTx{
		tx:  tx,
		ctx: ctx,
	}, nil
}

// Migrate выполняет миграции
func (s *SQLiteDatabase) Migrate() error {
	// TODO: Добавить модели для миграции
	return nil
}

// Stats возвращает статистику
func (s *SQLiteDatabase) Stats() Stats {
	sqlDB, err := s.db.DB()
	if err != nil {
		return Stats{}
	}

	stats := sqlDB.Stats()
	return Stats{
		OpenConnections:  stats.OpenConnections,
		InUseConnections: stats.InUse,
		IdleConnections:  stats.Idle,
	}
}

// DB возвращает GORM DB
func (s *SQLiteDatabase) DB() *gorm.DB {
	return s.db
}

// Commit подтверждает транзакцию
func (tx *SQLiteTx) Commit() error {
	return tx.tx.Commit().Error
}

// Rollback откатывает транзакцию
func (tx *SQLiteTx) Rollback() error {
	return tx.tx.Rollback().Error
}

// Context возвращает контекст транзакции
func (tx *SQLiteTx) Context() context.Context {
	return tx.ctx
}
`, config.ModuleName)
}

// generateRepositories создает примеры репозиториев
func (g *Generator) generateRepositories(config *ProjectConfig) error {
	content := fmt.Sprintf(`package repository

import (
	"%s/internal/models"
	"%s/pkg/database"
)

// UserRepository интерфейс для работы с пользователями
type UserRepository interface {
	database.Repository
	GetByID(id int64) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id int64) error
	List(offset, limit int) ([]*models.User, error)
}

// UserRepositoryImpl реализация репозитория пользователей
type UserRepositoryImpl struct {
	*database.BaseRepository
}

// NewUserRepository создает новый репозиторий пользователей
func NewUserRepository(db database.Database) UserRepository {
	return &UserRepositoryImpl{
		BaseRepository: database.NewBaseRepository(db),
	}
}

// GetByID получает пользователя по ID
func (r *UserRepositoryImpl) GetByID(id int64) (*models.User, error) {
	// TODO: Реализовать получение пользователя по ID
	return nil, nil
}

// GetByEmail получает пользователя по email
func (r *UserRepositoryImpl) GetByEmail(email string) (*models.User, error) {
	// TODO: Реализовать получение пользователя по email
	return nil, nil
}

// Create создает нового пользователя
func (r *UserRepositoryImpl) Create(user *models.User) error {
	// TODO: Реализовать создание пользователя
	return nil
}

// Update обновляет пользователя
func (r *UserRepositoryImpl) Update(user *models.User) error {
	// TODO: Реализовать обновление пользователя
	return nil
}

// Delete удаляет пользователя
func (r *UserRepositoryImpl) Delete(id int64) error {
	// TODO: Реализовать удаление пользователя
	return nil
}

// List возвращает список пользователей
func (r *UserRepositoryImpl) List(offset, limit int) ([]*models.User, error) {
	// TODO: Реализовать получение списка пользователей
	return nil, nil
}
`, config.ModuleName, config.ModuleName)

	repoPath := filepath.Join(g.projectPath, "internal/repository/user.go")
	return os.WriteFile(repoPath, []byte(content), 0644)
}

// generateModels создает примеры моделей
func (g *Generator) generateModels(config *ProjectConfig) error {
	content := `package models

import (
	"time"
)

// User модель пользователя
type User struct {
	ID        int64     ` + "`" + `json:"id" gorm:"primaryKey;autoIncrement"` + "`" + `
	Email     string    ` + "`" + `json:"email" gorm:"uniqueIndex;not null"` + "`" + `
	Name      string    ` + "`" + `json:"name" gorm:"not null"` + "`" + `
	CreatedAt time.Time ` + "`" + `json:"created_at" gorm:"autoCreateTime"` + "`" + `
	UpdatedAt time.Time ` + "`" + `json:"updated_at" gorm:"autoUpdateTime"` + "`" + `
}

// TableName возвращает имя таблицы
func (User) TableName() string {
	return "users"
}

// Product модель продукта (пример)
type Product struct {
	ID          int64     ` + "`" + `json:"id" gorm:"primaryKey;autoIncrement"` + "`" + `
	Name        string    ` + "`" + `json:"name" gorm:"not null"` + "`" + `
	Description string    ` + "`" + `json:"description"` + "`" + `
	Price       float64   ` + "`" + `json:"price" gorm:"not null"` + "`" + `
	CreatedAt   time.Time ` + "`" + `json:"created_at" gorm:"autoCreateTime"` + "`" + `
	UpdatedAt   time.Time ` + "`" + `json:"updated_at" gorm:"autoUpdateTime"` + "`" + `
}

// TableName возвращает имя таблицы
func (Product) TableName() string {
	return "products"
}
`

	modelsPath := filepath.Join(g.projectPath, "internal/models/models.go")
	return os.WriteFile(modelsPath, []byte(content), 0644)
}
