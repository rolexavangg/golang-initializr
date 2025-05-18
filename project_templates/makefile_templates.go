package project_templates

const makefileTemplate = `# Makefile for {{.GetProjectName}}

.PHONY: all build run test clean lint mock proto docker docker-compose

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
GOGET=$(GOCMD) get
GOLINT=golangci-lint

# Binary name
BINARY_NAME={{.GetProjectName}}

all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v

run:
	$(GORUN) main.go

test:
	$(GOTEST) -v ./...

clean:
	rm -f $(BINARY_NAME)
	rm -f coverage.out

lint:
	$(GOLINT) run ./...

tidy:
	$(GOMOD) tidy

update:
	$(GOMOD) tidy
	$(GOGET) -u ./...

{{- if .HasDependency "grpc"}}
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		api/proto/*.proto
{{- end}}

mock:
	mockery --all --keeptree --dir=internal/domain --output=internal/mocks

docker:
	docker build -t $(BINARY_NAME) .

docker-compose:
	docker-compose up -d
`
