# Golang Initializr Makefile

.PHONY: all build run clean test templ templ-watch dev

# Go related variables
GO=go
GOBUILD=$(GO) build
GOCLEAN=$(GO) clean
GOTEST=$(GO) test
GOGET=$(GO) get
GOFMT=$(GO) fmt
BINARY_NAME=golang-initializr
BUILD_DIR=./build

# Templ related variables
TEMPL=templ
TEMPL_DIR=./templates

# Default target
all: templ build

# Build the application
build:
	@echo "Building..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -v

# Run the application
run: templ
	@echo "Running..."
	$(GO) run main.go

# Clean build files
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

# Run tests
test:
	@echo "Testing..."
	$(GOTEST) -v ./...

# Format code
fmt:
	@echo "Formatting..."
	$(GOFMT) ./...

# Compile templ templates
templ:
	@echo "Compiling templ templates..."
	$(TEMPL) generate $(TEMPL_DIR)

# Watch templ templates for changes and recompile
templ-watch:
	@echo "Watching templ templates..."
	$(TEMPL) generate --watch $(TEMPL_DIR)

# Development mode - run with hot reload
dev: templ
	@echo "Starting development server..."
	air

# Install dependencies
deps:
	@echo "Installing dependencies..."
	$(GOGET) -u github.com/a-h/templ/cmd/templ
	$(GOGET) -u github.com/cosmtrek/air

# Deploy to GitHub Pages
deploy-gh-pages:
	@echo "Deploying to GitHub Pages..."
	git subtree push --prefix build origin gh-pages

# Help command
help:
	@echo "Available commands:"
	@echo "  make all          - Compile templates and build the application"
	@echo "  make build        - Build the application"
	@echo "  make run          - Run the application"
	@echo "  make clean        - Clean build files"
	@echo "  make test         - Run tests"
	@echo "  make fmt          - Format code"
	@echo "  make templ        - Compile templ templates"
	@echo "  make templ-watch  - Watch and compile templ templates on changes"
	@echo "  make dev          - Run with hot reload (requires air)"
	@echo "  make deps         - Install development dependencies"
	@echo "  make deploy-gh-pages - Deploy to GitHub Pages"
