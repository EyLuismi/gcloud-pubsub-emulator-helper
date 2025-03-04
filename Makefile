# Makefile

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=basicLoader
BUILD_DIR=build
BINARY_PATH=$(BUILD_DIR)/$(BINARY_NAME)
BINARY_UNIX=$(BINARY_PATH)_unix

# All target
all: test build

# Test target
test:
	$(GOTEST) -v ./...

# Build target
build:
	mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BINARY_PATH) ./cmd/headless

# Clean target
clean:
	$(GOCLEAN)
	rm -f $(BINARY_PATH)
	rm -f $(BINARY_UNIX)
	rm -rf $(BUILD_DIR)

# Run target
run: build
	./$(BINARY_PATH) -config ./full.example.json

# Dependency management
deps:
	$(GOGET) -u ./...

.PHONY: all test build clean run deps
