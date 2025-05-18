package project_templates

const grpcServerTemplate = `package grpc

import (
	"context"
	"fmt"
	"net"
	
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	
	"{{.Name}}/internal/config"
)


type GRPCConfig struct {
	Port int
}


func NewGRPCConfig(cfg *config.Config) *GRPCConfig {
	return &GRPCConfig{
		Port: cfg.GetEnvAsInt("GRPC_PORT", 9090),
	}
}


func NewGRPCServer(lc fx.Lifecycle, cfg *GRPCConfig, logger *zap.Logger, userService *UserService) *grpc.Server {
	server := grpc.NewServer()
	
	
	RegisterUserServiceServer(server, userService)
	
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			addr := fmt.Sprintf(":%d", cfg.Port)
			listener, err := net.Listen("tcp", addr)
			if err != nil {
				return err
			}
			
			logger.Info("Starting gRPC server", zap.String("addr", addr))
			
			go func() {
				if err := server.Serve(listener); err != nil {
					logger.Error("Failed to start gRPC server", zap.Error(err))
				}
			}()
			
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping gRPC server")
			server.GracefulStop()
			return nil
		},
	})
	
	return server
}
`

const grpcUserServiceTemplate = `package grpc

import (
	"context"
	
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	
	"{{.Name}}/internal/domain"
)


type UserService struct {
	UnimplementedUserServiceServer
	useCase domain.UserUseCase
	logger  *zap.Logger
}


func NewUserService(useCase domain.UserUseCase, logger *zap.Logger) *UserService {
	return &UserService{
		useCase: useCase,
		logger:  logger,
	}
}


func (s *UserService) CreateUser(ctx context.Context, req *CreateUserRequest) (*UserResponse, error) {
	user := &domain.User{
		Username: req.Username,
		Email:    req.Email,
	}
	
	if err := s.useCase.Create(user); err != nil {
		s.logger.Error("Failed to create user", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to create user")
	}
	
	return &UserResponse{
		Id:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}


func (s *UserService) GetUser(ctx context.Context, req *GetUserRequest) (*UserResponse, error) {
	user, err := s.useCase.GetByID(req.Id)
	if err != nil {
		s.logger.Error("Failed to get user", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to get user")
	}
	
	if user == nil {
		return nil, status.Error(codes.NotFound, "User not found")
	}
	
	return &UserResponse{
		Id:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}


func (s *UserService) ListUsers(ctx context.Context, req *ListUsersRequest) (*ListUsersResponse, error) {
	users, err := s.useCase.List()
	if err != nil {
		s.logger.Error("Failed to list users", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to list users")
	}
	
	response := &ListUsersResponse{
		Users: make([]*UserResponse, 0, len(users)),
	}
	
	for _, user := range users {
		response.Users = append(response.Users, &UserResponse{
			Id:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		})
	}
	
	return response, nil
}


func (s *UserService) UpdateUser(ctx context.Context, req *UpdateUserRequest) (*UserResponse, error) {
	user := &domain.User{
		ID:       req.Id,
		Username: req.Username,
		Email:    req.Email,
	}
	
	if err := s.useCase.Update(user); err != nil {
		s.logger.Error("Failed to update user", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to update user")
	}
	
	return &UserResponse{
		Id:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}


func (s *UserService) DeleteUser(ctx context.Context, req *DeleteUserRequest) (*DeleteUserResponse, error) {
	if err := s.useCase.Delete(req.Id); err != nil {
		s.logger.Error("Failed to delete user", zap.Error(err))
		return nil, status.Error(codes.Internal, "Failed to delete user")
	}
	
	return &DeleteUserResponse{Success: true}, nil
}
`

const userProtoTemplate = `syntax = "proto3";

package user;

option go_package = "internal/delivery/grpc";

service UserService {
  rpc CreateUser(CreateUserRequest) returns (UserResponse) {}
  rpc GetUser(GetUserRequest) returns (UserResponse) {}
  rpc ListUsers(ListUsersRequest) returns (ListUsersResponse) {}
  rpc UpdateUser(UpdateUserRequest) returns (UserResponse) {}
  rpc DeleteUser(DeleteUserRequest) returns (DeleteUserResponse) {}
}

message CreateUserRequest {
  string username = 1;
  string email = 2;
}

message GetUserRequest {
  string id = 1;
}

message UserResponse {
  string id = 1;
  string username = 2;
  string email = 3;
  string created_at = 4;
  string updated_at = 5;
}

message ListUsersRequest {
  
}

message ListUsersResponse {
  repeated UserResponse users = 1;
}

message UpdateUserRequest {
  string id = 1;
  string username = 2;
  string email = 3;
}

message DeleteUserRequest {
  string id = 1;
}

message DeleteUserResponse {
  bool success = 1;
}
`
