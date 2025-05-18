package project_templates

const readmeTemplate = `# {{.GetProjectName}}

Проект создан с помощью Golang Initializr.

## Особенности

- Чистая архитектура
- Uber FX для внедрения зависимостей
- Zap Logger для логирования
{{- if .HasDependency "postgres"}}
- PostgreSQL для хранения данных
{{- end}}
{{- if .HasDependency "redis"}}
- Redis для кеширования
{{- end}}
{{- if .HasDependency "kafka"}}
- Kafka для обмена сообщениями
{{- end}}
{{- if .HasDependency "http"}}
- HTTP API (Echo framework)
{{- end}}
{{- if .HasDependency "grpc"}}
- gRPC API
{{- end}}

## Запуск

### Локальная разработка

` + "```bash" + `
go run main.go
` + "```" + `

{{- if .HasDependency "docker"}}

### С использованием Docker

` + "```bash" + `
docker-compose up -d
` + "```" + `
{{- end}}

## Структура проекта

Проект следует принципам чистой архитектуры:

- domain - бизнес-сущности
- usecase - бизнес-логика
- repository - слой доступа к данным
- delivery - слой доставки (HTTP, gRPC и т.д.)
`

const goModTemplate = `module {{.Name}}

go 1.24

require (
	go.uber.org/fx v1.20.1
	go.uber.org/zap v1.27.0
{{- if .HasDependency "http"}}
	github.com/labstack/echo/v4 v4.13.3
{{- end}}
{{- if .HasDependency "postgres"}}
	github.com/doug-martin/goqu/v9 v9.19.0
	github.com/jackc/pgx/v5 v5.5.5
{{- end}}
{{- if .HasDependency "redis"}}
	github.com/redis/go-redis/v9 v9.5.1
{{- end}}
{{- if .HasDependency "kafka"}}
	github.com/segmentio/kafka-go v0.4.47
{{- end}}
{{- if .HasDependency "grpc"}}
	google.golang.org/grpc v1.62.1
	google.golang.org/protobuf v1.33.0
{{- end}}
)
`

const mainTemplate = `package main

import (
	"{{.Name}}/internal/bootstrap"
)

func main() {
	
	bootstrap.BuildApp().Run()
}
`

const gitignoreTemplate = `# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib
*.db

# Test binary, built with 'go test -c'
*.test

# Output of the go coverage tool
*.out

# Dependency directories
vendor/

# IDE files
.idea/
.vscode/
*.swp
*.swo

# OS files
.DS_Store
Thumbs.db

# Environment variables
.env
`

const dockerfileTemplate = `FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/app .
{{- if .HasDependency "grpc"}}
COPY --from=builder /app/api ./api
{{- end}}

EXPOSE 8080
{{- if .HasDependency "grpc"}}
EXPOSE 9090
{{- end}}

CMD ["./app"]
`
