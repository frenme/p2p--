.PHONY: build test clean run install deps lint

BINARY_NAME=p2p-chat
BUILD_DIR=bin
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GOVERSION := $(shell go version | awk '{print $$3}')

LDFLAGS := -ldflags "-X p2p-chat/internal/version.GitCommit=$(COMMIT) -X p2p-chat/internal/version.BuildDate=$(DATE) -X p2p-chat/internal/version.GoVersion=$(GOVERSION)"

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .

test:
	@echo "Running tests..."
	@go test ./...

clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@rm -f $(BINARY_NAME)
	@rm -f chat_history_*.json

run: build
	@echo "Running $(BINARY_NAME)..."
	@./$(BUILD_DIR)/$(BINARY_NAME)

install: build
	@echo "Installing $(BINARY_NAME)..."
	@cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/

deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

lint:
	@echo "Running linter..."
	@go vet ./...
	@go fmt ./...