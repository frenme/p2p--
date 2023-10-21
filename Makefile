.PHONY: build test clean run install deps lint

BINARY_NAME=p2p-chat
BUILD_DIR=bin

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) .

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