package project_templates

// Шаблоны для Docker и Docker Compose

const dockerComposeTemplate = `version: '3'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    {{- if .HasDependency "grpc"}}
      - "9090:9090"
    {{- end}}
    environment:
      - SERVER_PORT=8080
    {{- if .HasDependency "postgres"}}
      - POSTGRES_HOST=postgres
      - POSTGRES_PORT=5432
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=postgres
      - POSTGRES_DB={{.GetProjectName}}
    {{- end}}
    {{- if .HasDependency "redis"}}
      - REDIS_HOST=redis
      - REDIS_PORT=6379
    {{- end}}
    {{- if .HasDependency "kafka"}}
      - KAFKA_BROKER=kafka:9092
    {{- end}}
    depends_on:
    {{- if .HasDependency "postgres"}}
      - postgres
    {{- end}}
    {{- if .HasDependency "redis"}}
      - redis
    {{- end}}
    {{- if .HasDependency "kafka"}}
      - kafka
    {{- end}}
    restart: unless-stopped
    networks:
      - app-network

{{- if .HasDependency "postgres"}}
  postgres:
    image: postgres:15-alpine
    ports:
      - "5432:5432"
    environment:
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=postgres
      - POSTGRES_DB={{.GetProjectName}}
    volumes:
      - postgres-data:/var/lib/postgresql/data
    restart: unless-stopped
    networks:
      - app-network
{{- end}}

{{- if .HasDependency "redis"}}
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis-data:/data
    restart: unless-stopped
    networks:
      - app-network
{{- end}}

{{- if .HasDependency "kafka"}}
  zookeeper:
    image: confluentinc/cp-zookeeper:7.3.0
    ports:
      - "2181:2181"
    environment:
      - ZOOKEEPER_CLIENT_PORT=2181
    restart: unless-stopped
    networks:
      - app-network

  kafka:
    image: confluentinc/cp-kafka:7.3.0
    ports:
      - "9092:9092"
    environment:
      - KAFKA_ZOOKEEPER_CONNECT=zookeeper:2181
      - KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://kafka:9092
      - KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR=1
    depends_on:
      - zookeeper
    restart: unless-stopped
    networks:
      - app-network
{{- end}}

networks:
  app-network:
    driver: bridge

volumes:
{{- if .HasDependency "postgres"}}
  postgres-data:
{{- end}}
{{- if .HasDependency "redis"}}
  redis-data:
{{- end}}
`
