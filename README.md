# Golang Initializr

A modern project generator for Go applications based on clean architecture principles.

## Overview

Golang Initializr is a web application that helps you quickly bootstrap Go projects with a clean architecture structure and common dependencies. Similar to Spring Initializr for Java, this tool generates a complete project structure with all the necessary files and configurations.

## Features

- **Clean Architecture**: All generated projects follow clean architecture principles
- **Uber FX**: Dependency injection using Uber FX
- **Zap Logger**: Structured logging with Uber's Zap
- **Default Components**: Includes default usecase and repository with two endpoints (create and get user)
- **Unified Configuration**: Creates a single config that contains configs for all dependencies
- **Environment Variables**: Generates `.env.example` file for configuration
- **Build Tools**: Includes Makefile with commands for building, testing, and more
- **Mock Generation**: Adds mockery directives for interfaces
- **GitHub Actions**: Includes GitHub Pages workflow configuration

## Available Dependencies

- **Databases**:

  - PostgreSQL (with goqu)
  - Redis

- **Messaging**:

  - Kafka

- **API**:

  - HTTP (Echo framework)
  - gRPC

- **Tools**:
  - Docker

## Getting Started

### Prerequisites

- Go 1.24 or later
- [templ](https://github.com/a-h/templ) for template generation

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/malinatrash/golang-initializr.git
   cd golang-initializr
   ```

2. Install dependencies:

   ```bash
   go mod download
   ```

3. Generate templ templates:

   ```bash
   make templ
   ```

4. Run the application:

   ```bash
   make run
   ```

5. Open your browser and navigate to `http://localhost:8081`

### Development

For development with hot reload:

```bash
make dev
```

To watch and automatically recompile templ templates:

```bash
make templ-watch
```

## Project Structure

```
├── internal/
│   ├── app/         # Application entry point
│   ├── bootstrap/   # Application bootstrap components
│   ├── config/      # Configuration
│   ├── domain/      # Domain models and interfaces
│   ├── usecase/     # Business logic
│   ├── repository/  # Data access layer
│   └── delivery/    # API layer (HTTP, gRPC)
├── .env.example     # Example environment variables
├── go.mod           # Go modules file
├── Makefile         # Build and development commands
└── Dockerfile       # Docker configuration
```

## Building and Deployment

### Build the Application

```bash
make build
```

### Run Tests

```bash
make test
```

### Deploy to GitHub Pages

```bash
make deploy-gh-pages
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- [Uber FX](https://github.com/uber-go/fx) for dependency injection
- [Zap](https://github.com/uber-go/zap) for logging
- [Echo](https://echo.labstack.com/) for HTTP server
- [templ](https://github.com/a-h/templ) for HTML templates
