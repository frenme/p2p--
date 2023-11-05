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

demo: build
	@echo "Starting demo with two users..."
	@echo "Open two terminals and run:"
	@echo "  Terminal 1: ./$(BUILD_DIR)/$(BINARY_NAME) Alice"
	@echo "  Terminal 2: ./$(BUILD_DIR)/$(BINARY_NAME) Bob"

help:
	@echo "Available targets:"
	@echo "  build    - Build the binary"
	@echo "  test     - Run tests"
	@echo "  clean    - Clean build artifacts"
	@echo "  run      - Build and run with default user"
	@echo "  install  - Install to /usr/local/bin"
	@echo "  deps     - Download dependencies"
	@echo "  lint     - Run linter and formatter"
	@echo "  demo     - Show demo instructions"
	@echo "  help     - Show this help"