package project_templates

// Шаблоны для файлов чистой архитектуры

const appTemplate = `package app

import (
	"go.uber.org/fx"

	"%s/internal/bootstrap"
	"%s/internal/delivery/http"
	"%s/internal/usecase"
	"%s/internal/repository"
	"%s/internal/config"
)

// Module provides dependencies for the application
var Module = fx.Options(
	// Core dependencies
	bootstrap.Module,
	// Provide all usecases
	usecase.Module,
	// Provide all repositories
	repository.Module,
	// Provide HTTP handlers
	http.Module,
	// Register HTTP routes
	fx.Invoke(http.RegisterRoutes),
)
`

const oldConfigTemplate = `package config

import (
	"os"
	"strconv"
)

// Config содержит конфигурацию приложения
type Config struct {
	Server ServerConfig
}

// ServerConfig содержит конфигурацию сервера
type ServerConfig struct {
	Port int
}

// NewConfig создает новый экземпляр конфигурации
func NewConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnvAsInt("SERVER_PORT", 8080),
		},
	}
}

// getEnvAsInt получает переменную окружения как int
func getEnvAsInt(name string, defaultVal int) int {
	valueStr := os.Getenv(name)
	if valueStr == "" {
		return defaultVal
	}
	
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultVal
	}
	
	return value
}
`

const userDomainTemplate = `package domain

import (
	"time"
)

type User struct {
	ID        string    ` + "`json:\"id\"`" + `
	Username  string    ` + "`json:\"username\"`" + `
	Email     string    ` + "`json:\"email\"`" + `
	CreatedAt time.Time ` + "`json:\"created_at\"`" + `
	UpdatedAt time.Time ` + "`json:\"updated_at\"`" + `
}

//go:generate mockery --name=UserRepository --output=../mocks --outpkg=mocks
type UserRepository interface {
	Create(user *User) error
	GetByID(id string) (*User, error)
	List() ([]*User, error)
	Update(user *User) error
	Delete(id string) error
}

//go:generate mockery --name=UserUseCase --output=../mocks --outpkg=mocks
type UserUseCase interface {
	Create(user *User) error
	GetByID(id string) (*User, error)
	List() ([]*User, error)
	Update(user *User) error
	Delete(id string) error
}
`

const userUsecaseTemplate = `package usecase

import (
	"time"
	
	"go.uber.org/fx"
	"go.uber.org/zap"
	
	"{{.Name}}/internal/domain"
)

var Module = fx.Options(
	fx.Provide(NewUserUseCase),
)

type userUseCase struct {
	repo   domain.UserRepository
	logger *zap.Logger
}

func NewUserUseCase(repo domain.UserRepository, logger *zap.Logger) domain.UserUseCase {
	return &userUseCase{
		repo:   repo,
		logger: logger,
	}
}

func (u *userUseCase) Create(user *domain.User) error {
	u.logger.Info("Creating new user", zap.String("username", user.Username))
	
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	
	return u.repo.Create(user)
}

func (u *userUseCase) GetByID(id string) (*domain.User, error) {
	u.logger.Info("Getting user by ID", zap.String("id", id))
	return u.repo.GetByID(id)
}

func (u *userUseCase) List() ([]*domain.User, error) {
	u.logger.Info("Getting list of users")
	return u.repo.List()
}

func (u *userUseCase) Update(user *domain.User) error {
	u.logger.Info("Updating user", zap.String("id", user.ID))
	
	user.UpdatedAt = time.Now()
	
	return u.repo.Update(user)
}

func (u *userUseCase) Delete(id string) error {
	u.logger.Info("Deleting user", zap.String("id", id))
	return u.repo.Delete(id)
}
`

const userRepositoryTemplate = `package repository

import (
	"go.uber.org/fx"
	
	"{{.Name}}/internal/domain"
)

var Module = fx.Options(
	fx.Provide(NewUserRepository),
)

func NewUserRepository() domain.UserRepository {


	return NewInMemoryUserRepository()
}

type InMemoryUserRepository struct {
	users map[string]*domain.User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]*domain.User),
	}
}

func (r *InMemoryUserRepository) Create(user *domain.User) error {
	r.users[user.ID] = user
	return nil
}

func (r *InMemoryUserRepository) GetByID(id string) (*domain.User, error) {
	user, exists := r.users[id]
	if !exists {
		return nil, nil
	}
	return user, nil
}

func (r *InMemoryUserRepository) List() ([]*domain.User, error) {
	users := make([]*domain.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}
	return users, nil
}

func (r *InMemoryUserRepository) Update(user *domain.User) error {
	r.users[user.ID] = user
	return nil
}

func (r *InMemoryUserRepository) Delete(id string) error {
	delete(r.users, id)
	return nil
}
`
