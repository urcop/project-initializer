package generator

import (
	"fmt"
	"os"
	"path/filepath"
)

// generateGRPC создает gRPC сервер файлы
func (g *Generator) generateGRPC(config *ProjectConfig) error {
	// Создаем proto файл
	if err := g.generateProtoFile(config); err != nil {
		return err
	}

	// Создаем gRPC сервер
	if err := g.generateGRPCServer(config); err != nil {
		return err
	}

	// Создаем gRPC клиент (для примера)
	if err := g.generateGRPCClient(config); err != nil {
		return err
	}

	// Создаем Makefile для protobuf
	if err := g.generateProtoMakefile(config); err != nil {
		return err
	}

	return nil
}

// generateProtoFile создает .proto файл
func (g *Generator) generateProtoFile(config *ProjectConfig) error {
	content := fmt.Sprintf(`syntax = "proto3";

package %s;

option go_package = "%s/internal/grpc/pb";

// Сервис для %s
service %sService {
  // Health check
  rpc HealthCheck(HealthCheckRequest) returns (HealthCheckResponse);
  
  // Ping
  rpc Ping(PingRequest) returns (PingResponse);
  
  // Пример CRUD операций
  rpc CreateUser(CreateUserRequest) returns (CreateUserResponse);
  rpc GetUser(GetUserRequest) returns (GetUserResponse);
  rpc UpdateUser(UpdateUserRequest) returns (UpdateUserResponse);
  rpc DeleteUser(DeleteUserRequest) returns (DeleteUserResponse);
  rpc ListUsers(ListUsersRequest) returns (ListUsersResponse);
}

// Health Check
message HealthCheckRequest {}

message HealthCheckResponse {
  string status = 1;
  string service = 2;
  string version = 3;
  int64 timestamp = 4;
}

// Ping
message PingRequest {}

message PingResponse {
  string message = 1;
  string service = 2;
  string version = 3;
}

// User messages
message User {
  int64 id = 1;
  string email = 2;
  string name = 3;
  int64 created_at = 4;
  int64 updated_at = 5;
}

message CreateUserRequest {
  string email = 1;
  string name = 2;
}

message CreateUserResponse {
  User user = 1;
}

message GetUserRequest {
  int64 id = 1;
}

message GetUserResponse {
  User user = 1;
}

message UpdateUserRequest {
  int64 id = 1;
  string email = 2;
  string name = 3;
}

message UpdateUserResponse {
  User user = 1;
}

message DeleteUserRequest {
  int64 id = 1;
}

message DeleteUserResponse {
  bool success = 1;
}

message ListUsersRequest {
  int32 offset = 1;
  int32 limit = 2;
}

message ListUsersResponse {
  repeated User users = 1;
  int32 total = 2;
}
`, config.Name, config.ModuleName, config.Name, config.Name)

	// Создаем директорию для proto файлов
	protoDir := filepath.Join(g.projectPath, "api/proto")
	if err := os.MkdirAll(protoDir, 0755); err != nil {
		return err
	}

	protoPath := filepath.Join(protoDir, fmt.Sprintf("%s.proto", config.Name))
	return os.WriteFile(protoPath, []byte(content), 0644)
}

// generateGRPCServer создает gRPC сервер
func (g *Generator) generateGRPCServer(config *ProjectConfig) error {
	content := fmt.Sprintf(`package grpc

import (
	"context"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"%s/internal/config"
	"%s/internal/grpc/pb"
	"%s/pkg/logger"
)

// Server представляет gRPC сервер
type Server struct {
	cfg      *config.Config
	logger   logger.Logger
	grpcSrv  *grpc.Server
	listener net.Listener
	pb.Unimplemented%sServiceServer
}

// New создает новый gRPC сервер
func New(cfg *config.Config, logger logger.Logger) *Server {
	return &Server{
		cfg:    cfg,
		logger: logger,
	}
}

// Start запускает gRPC сервер
func (s *Server) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%%d", s.cfg.GRPC.Port))
	if err != nil {
		return fmt.Errorf("ошибка создания listener: %%w", err)
	}

	s.listener = lis

	// Создаем gRPC сервер
	s.grpcSrv = grpc.NewServer(
		grpc.ConnectionTimeout(time.Duration(s.cfg.GRPC.MaxConnectionAge)*time.Second),
	)

	// Регистрируем сервис
	pb.Register%sServiceServer(s.grpcSrv, s)

	// Включаем reflection для grpcurl
	reflection.Register(s.grpcSrv)

	s.logger.Info("gRPC сервер запущен", "port", s.cfg.GRPC.Port)

	return s.grpcSrv.Serve(lis)
}

// Stop останавливает gRPC сервер
func (s *Server) Stop() {
	if s.grpcSrv != nil {
		s.grpcSrv.GracefulStop()
	}
}

// HealthCheck реализует health check
func (s *Server) HealthCheck(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	s.logger.Debug("gRPC HealthCheck вызван")

	return &pb.HealthCheckResponse{
		Status:    "ok",
		Service:   s.cfg.App.Name,
		Version:   s.cfg.App.Version,
		Timestamp: time.Now().Unix(),
	}, nil
}

// Ping реализует ping
func (s *Server) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	s.logger.Debug("gRPC Ping вызван")

	return &pb.PingResponse{
		Message: "pong",
		Service: s.cfg.App.Name,
		Version: s.cfg.App.Version,
	}, nil
}

// CreateUser создает пользователя
func (s *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	s.logger.Debug("gRPC CreateUser вызван", "email", req.Email)

	// TODO: Реализовать создание пользователя
	user := &pb.User{
		Id:        1,
		Email:     req.Email,
		Name:      req.Name,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	return &pb.CreateUserResponse{
		User: user,
	}, nil
}

// GetUser получает пользователя
func (s *Server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	s.logger.Debug("gRPC GetUser вызван", "id", req.Id)

	// TODO: Реализовать получение пользователя
	user := &pb.User{
		Id:        req.Id,
		Email:     "user@example.com",
		Name:      "Test User",
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	return &pb.GetUserResponse{
		User: user,
	}, nil
}

// UpdateUser обновляет пользователя
func (s *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	s.logger.Debug("gRPC UpdateUser вызван", "id", req.Id)

	// TODO: Реализовать обновление пользователя
	user := &pb.User{
		Id:        req.Id,
		Email:     req.Email,
		Name:      req.Name,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	return &pb.UpdateUserResponse{
		User: user,
	}, nil
}

// DeleteUser удаляет пользователя
func (s *Server) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	s.logger.Debug("gRPC DeleteUser вызван", "id", req.Id)

	// TODO: Реализовать удаление пользователя

	return &pb.DeleteUserResponse{
		Success: true,
	}, nil
}

// ListUsers возвращает список пользователей
func (s *Server) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	s.logger.Debug("gRPC ListUsers вызван", "offset", req.Offset, "limit", req.Limit)

	// TODO: Реализовать получение списка пользователей
	users := []*pb.User{
		{
			Id:        1,
			Email:     "user1@example.com",
			Name:      "User 1",
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		},
		{
			Id:        2,
			Email:     "user2@example.com",
			Name:      "User 2",
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		},
	}

	return &pb.ListUsersResponse{
		Users: users,
		Total: int32(len(users)),
	}, nil
}
`, config.ModuleName, config.ModuleName, config.ModuleName, config.Name, config.Name)

	// Создаем директорию для gRPC
	grpcDir := filepath.Join(g.projectPath, "internal/grpc")
	if err := os.MkdirAll(grpcDir, 0755); err != nil {
		return err
	}

	serverPath := filepath.Join(grpcDir, "server.go")
	return os.WriteFile(serverPath, []byte(content), 0644)
}

// generateGRPCClient создает gRPC клиент (для примера)
func (g *Generator) generateGRPCClient(config *ProjectConfig) error {
	content := fmt.Sprintf(`package grpc

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"%s/internal/grpc/pb"
	"%s/pkg/logger"
)

// Client представляет gRPC клиент
type Client struct {
	conn   *grpc.ClientConn
	client pb.%sServiceClient
	logger logger.Logger
}

// NewClient создает новый gRPC клиент
func NewClient(address string, logger logger.Logger) (*Client, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к gRPC серверу: %%w", err)
	}

	client := pb.New%sServiceClient(conn)

	return &Client{
		conn:   conn,
		client: client,
		logger: logger,
	}, nil
}

// Close закрывает соединение
func (c *Client) Close() error {
	return c.conn.Close()
}

// HealthCheck выполняет health check
func (c *Client) HealthCheck(ctx context.Context) (*pb.HealthCheckResponse, error) {
	return c.client.HealthCheck(ctx, &pb.HealthCheckRequest{})
}

// Ping выполняет ping
func (c *Client) Ping(ctx context.Context) (*pb.PingResponse, error) {
	return c.client.Ping(ctx, &pb.PingRequest{})
}

// CreateUser создает пользователя
func (c *Client) CreateUser(ctx context.Context, email, name string) (*pb.CreateUserResponse, error) {
	return c.client.CreateUser(ctx, &pb.CreateUserRequest{
		Email: email,
		Name:  name,
	})
}

// GetUser получает пользователя
func (c *Client) GetUser(ctx context.Context, id int64) (*pb.GetUserResponse, error) {
	return c.client.GetUser(ctx, &pb.GetUserRequest{
		Id: id,
	})
}

// UpdateUser обновляет пользователя
func (c *Client) UpdateUser(ctx context.Context, id int64, email, name string) (*pb.UpdateUserResponse, error) {
	return c.client.UpdateUser(ctx, &pb.UpdateUserRequest{
		Id:    id,
		Email: email,
		Name:  name,
	})
}

// DeleteUser удаляет пользователя
func (c *Client) DeleteUser(ctx context.Context, id int64) (*pb.DeleteUserResponse, error) {
	return c.client.DeleteUser(ctx, &pb.DeleteUserRequest{
		Id: id,
	})
}

// ListUsers возвращает список пользователей
func (c *Client) ListUsers(ctx context.Context, offset, limit int32) (*pb.ListUsersResponse, error) {
	return c.client.ListUsers(ctx, &pb.ListUsersRequest{
		Offset: offset,
		Limit:  limit,
	})
}
`, config.ModuleName, config.ModuleName, config.Name, config.Name)

	clientPath := filepath.Join(g.projectPath, "internal/grpc/client.go")
	return os.WriteFile(clientPath, []byte(content), 0644)
}

// generateProtoMakefile создает Makefile для protobuf
func (g *Generator) generateProtoMakefile(config *ProjectConfig) error {
	content := `# Protobuf Makefile

# Переменные
PROTO_DIR=api/proto
GRPC_DIR=internal/grpc/pb

.PHONY: proto-gen proto-clean proto-install

# Генерация Go кода из proto файлов
proto-gen: ## Генерировать Go код из proto файлов
	@echo "Генерация Go кода из proto файлов..."
	@mkdir -p $(GRPC_DIR)
	protoc \
		--go_out=$(GRPC_DIR) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(GRPC_DIR) \
		--go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/*.proto
	@echo "Генерация завершена"

# Установка необходимых инструментов
proto-install: ## Установить protoc и плагины
	@echo "Установка protoc плагинов..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "Плагины установлены"

# Очистка сгенерированных файлов
proto-clean: ## Очистить сгенерированные proto файлы
	@echo "Очистка сгенерированных файлов..."
	rm -rf $(GRPC_DIR)/*.pb.go
	@echo "Очистка завершена"

# Помощь
proto-help: ## Показать справку по proto командам
	@echo "Доступные proto команды:"
	@echo "  proto-install  - Установить protoc плагины"
	@echo "  proto-gen      - Генерировать Go код из proto файлов"
	@echo "  proto-clean    - Очистить сгенерированные файлы"
	@echo ""
	@echo "Требования:"
	@echo "  - protoc должен быть установлен (https://grpc.io/docs/protoc-installation/)"
	@echo "  - Выполните 'make proto-install' для установки Go плагинов"
`

	protoMakefilePath := filepath.Join(g.projectPath, "scripts/proto.mk")
	scriptsDir := filepath.Join(g.projectPath, "scripts")
	if err := os.MkdirAll(scriptsDir, 0755); err != nil {
		return err
	}

	return os.WriteFile(protoMakefilePath, []byte(content), 0644)
}
