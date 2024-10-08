# Go-related variables
GO := go
GOBIN := $(shell $(GO) env GOBIN)
GOPATH := $(shell $(GO) env GOPATH)
GOFILES := $(shell find . -name '*.go' -not -path "./vendor/*")
APP_NAME := busgopher

# Set the default action when `make` is run without arguments
.DEFAULT_GOAL := build

# Compile and build the binary
build:
	@echo "Building the binary..."
	$(GO) build -o $(APP_NAME) .

# Run the Go application
run: build
	@echo "Running the application..."
	./$(APP_NAME)

# Run tests
test:
	@echo "Running tests..."
	$(GO) test ./...

# Run tests with coverage
test-cover:
	@echo "Running tests with coverage..."
	$(GO) test ./... -cover
