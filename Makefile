.PHONY: build clean install test fmt vet lint release

BINARY_NAME=sannti
VERSION=v0.1.0
BUILD_DIR=build
PLATFORMS=darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet

# Build flags
LDFLAGS=-ldflags "-s -w -X github.com/sannticloud/sannti-cli/cmd.Version=$(VERSION)"

# Default target
all: clean build

# Build for current platform
build:
	@echo "Building $(BINARY_NAME) $(VERSION)..."
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Build for all platforms
release: clean
	@echo "Building releases for version $(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	@for platform in $(PLATFORMS); do \
		os=$${platform%/*}; \
		arch=$${platform#*/}; \
		output=$(BUILD_DIR)/$(BINARY_NAME)-$(VERSION)-$${os}-$${arch}; \
		if [ "$$os" = "windows" ]; then \
			output=$${output}.exe; \
		fi; \
		echo "Building for $$os/$$arch..."; \
		GOOS=$$os GOARCH=$$arch $(GOBUILD) $(LDFLAGS) -o $$output .; \
		if [ $$? -ne 0 ]; then \
			echo "Failed to build for $$os/$$arch"; \
			exit 1; \
		fi; \
	done
	@echo "All releases built successfully!"

# Install to local system
install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "Installed successfully!"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	$(GOCLEAN)
	@echo "Clean complete!"

# Run tests
test:
	$(GOTEST) -v ./...

# Format code
fmt:
	$(GOFMT) ./...

# Run go vet
vet:
	$(GOVET) ./...

# Run linter (requires golangci-lint)
lint:
	@which golangci-lint > /dev/null || (echo "golangci-lint not found. Install from https://golangci-lint.run/usage/install/" && exit 1)
	golangci-lint run

# Download dependencies
deps:
	$(GOGET) -v -d ./...
	$(GOCMD) mod tidy

# Help
help:
	@echo "Sannti CLI Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make build      - Build for current platform"
	@echo "  make release    - Build for all platforms"
	@echo "  make install    - Install to /usr/local/bin"
	@echo "  make clean      - Remove build artifacts"
	@echo "  make test       - Run tests"
	@echo "  make fmt        - Format code"
	@echo "  make vet        - Run go vet"
	@echo "  make lint       - Run linter"
	@echo "  make deps       - Download dependencies"
