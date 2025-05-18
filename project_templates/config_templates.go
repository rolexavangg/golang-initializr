package project_templates

const configMainTemplate = `package config

import (
	"sync"
)

// Config holds all configuration for the application
type Config struct {
	App      AppConfig
	Server   ServerConfig
{{- if .HasDependency "postgres"}}
	Postgres PostgresConfig
{{- end}}
{{- if .HasDependency "redis"}}
	Redis    RedisConfig
{{- end}}
{{- if .HasDependency "kafka"}}
	Kafka    KafkaConfig
{{- end}}
{{- if .HasDependency "grpc"}}
	GRPC     GRPCConfig
{{- end}}
}

var (
	instance *Config
	once     sync.Once
)

// GetConfig returns a singleton instance of Config
func GetConfig() *Config {
	once.Do(func() {
		instance = NewConfig()
	})
	return instance
}

// NewConfig creates a new configuration instance with values from environment variables
func NewConfig() *Config {
	return &Config{
		App:    NewAppConfig(),
		Server: NewServerConfig(),
{{- if .HasDependency "postgres"}}
		Postgres: NewPostgresConfig(),
{{- end}}
{{- if .HasDependency "redis"}}
		Redis: NewRedisConfig(),
{{- end}}
{{- if .HasDependency "kafka"}}
		Kafka: NewKafkaConfig(),
{{- end}}
{{- if .HasDependency "grpc"}}
		GRPC: NewGRPCConfig(),
{{- end}}
	}
}`

const configAppTemplate = `package config

// AppConfig holds application-specific configuration
type AppConfig struct {
	Name        string
	Environment string
	Debug       bool
}


func NewAppConfig() AppConfig {
	return AppConfig{
		Name:        getEnv("APP_NAME", "app"),
		Environment: getEnv("APP_ENV", "development"),
		Debug:       getEnvAsBool("APP_DEBUG", true),
	}
}`

const configServerTemplate = `package config


type ServerConfig struct {
	Host string
	Port int
}


func NewServerConfig() ServerConfig {
	return ServerConfig{
		Host: getEnv("SERVER_HOST", "localhost"),
		Port: getEnvAsInt("SERVER_PORT", 8080),
	}
}`

const configPostgresTemplate = `package config


type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string
}


func NewPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Host:     getEnv("POSTGRES_HOST", "localhost"),
		Port:     getEnvAsInt("POSTGRES_PORT", 5432),
		User:     getEnv("POSTGRES_USER", "postgres"),
		Password: getEnv("POSTGRES_PASSWORD", "postgres"),
		Database: getEnv("POSTGRES_DB", "postgres"),
		SSLMode:  getEnv("POSTGRES_SSLMODE", "disable"),
	}
}`

const configRedisTemplate = `package config


type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}


func NewRedisConfig() RedisConfig {
	return RedisConfig{
		Host:     getEnv("REDIS_HOST", "localhost"),
		Port:     getEnvAsInt("REDIS_PORT", 6379),
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       getEnvAsInt("REDIS_DB", 0),
	}
}`

const configKafkaTemplate = `package config

import "strings"


type KafkaConfig struct {
	Brokers []string
	Topic   string
	GroupID string
}


func NewKafkaConfig() KafkaConfig {
	return KafkaConfig{
		Brokers: strings.Split(getEnv("KAFKA_BROKERS", "localhost:9092"), ","),
		Topic:   getEnv("KAFKA_TOPIC", "default-topic"),
		GroupID: getEnv("KAFKA_GROUP_ID", "default-group"),
	}
}`

const configGRPCTemplate = `package config


type GRPCConfig struct {
	Host string
	Port int
}


func NewGRPCConfig() GRPCConfig {
	return GRPCConfig{
		Host: getEnv("GRPC_HOST", "localhost"),
		Port: getEnvAsInt("GRPC_PORT", 50051),
	}
}`

const configUtilsTemplate = `package config

import (
	"os"
	"strconv"
	"strings"
)


func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}


func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}


func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}


func getEnvAsSlice(key string, defaultValue []string, sep string) []string {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	return strings.Split(valueStr, sep)
}`

const exampleEnvTemplate = `# Application
APP_NAME={{.GetProjectName}}
APP_ENV=development
APP_DEBUG=true

# Server
SERVER_HOST=localhost
SERVER_PORT=8080

{{- if .HasDependency "postgres"}}
# PostgreSQL
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB={{.GetProjectName}}
POSTGRES_SSLMODE=disable
{{- end}}

{{- if .HasDependency "redis"}}
# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
{{- end}}

{{- if .HasDependency "kafka"}}
# Kafka
KAFKA_BROKERS=localhost:9092
KAFKA_TOPIC={{.GetProjectName}}-topic
KAFKA_GROUP_ID={{.GetProjectName}}-consumer
{{- end}}

{{- if .HasDependency "grpc"}}
# gRPC
GRPC_HOST=localhost
GRPC_PORT=50051
{{- end}}
`
