# Makefile

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=basicLoader
BINARY_UNIX=$(BINARY_NAME)_unix

# All target
all: test build

# Test target
test:
	$(GOTEST) -v ./...

# Build target
build:
	$(GOBUILD) -o $(BINARY_NAME) ./cmd/headless

# Clean target
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

# Run target
run:
	./$(BINARY_NAME)

# Dependency management
deps:
	$(GOGET) -u ./...

.PHONY: all test build clean run deps
