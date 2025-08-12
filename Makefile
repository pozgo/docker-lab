# LAB Build Configuration
APP_NAME=lab
SOURCE_DIR=./app
OUTPUT_DIR=.

# Build variables
GOOS_LINUX=linux
GOOS_DARWIN=darwin
GOARCH_AMD64=amd64
GOARCH_ARM64=arm64

.PHONY: all build build-linux-amd64 build-linux-arm64 build-darwin-amd64 build-darwin-arm64 clean deps test

# Default target - build for current platform
all: deps build

# Build for current platform
build:
	@echo "🔨 Building LAB for current platform..."
	cd $(SOURCE_DIR) && go build -o ../$(APP_NAME) .
	@echo "✅ Build complete: $(APP_NAME)"

# Build for Linux AMD64
build-linux-amd64:
	@echo "🔨 Building LAB for Linux AMD64..."
	cd $(SOURCE_DIR) && GOOS=$(GOOS_LINUX) GOARCH=$(GOARCH_AMD64) go build -o ../$(APP_NAME)-linux-amd64 .
	@echo "✅ Build complete: $(APP_NAME)-linux-amd64"

# Build for Linux ARM64
build-linux-arm64:
	@echo "🔨 Building LAB for Linux ARM64..."
	cd $(SOURCE_DIR) && GOOS=$(GOOS_LINUX) GOARCH=$(GOARCH_ARM64) go build -o ../$(APP_NAME)-linux-arm64 .
	@echo "✅ Build complete: $(APP_NAME)-linux-arm64"

# Build for macOS AMD64
build-darwin-amd64:
	@echo "🔨 Building LAB for macOS AMD64..."
	cd $(SOURCE_DIR) && GOOS=$(GOOS_DARWIN) GOARCH=$(GOARCH_AMD64) go build -o ../$(APP_NAME)-darwin-amd64 .
	@echo "✅ Build complete: $(APP_NAME)-darwin-amd64"

# Build for macOS ARM64 (M1/M2 Macs)
build-darwin-arm64:
	@echo "🔨 Building LAB for macOS ARM64 (M1/M2)..."
	cd $(SOURCE_DIR) && GOOS=$(GOOS_DARWIN) GOARCH=$(GOARCH_ARM64) go build -o ../$(APP_NAME)-darwin-arm64 .
	@echo "✅ Build complete: $(APP_NAME)-darwin-arm64"

# Build all platforms as required by Phase 2
build-all: build-linux-amd64 build-linux-arm64 build-darwin-amd64 build-darwin-arm64
	@echo "🎉 All binaries built successfully!"
	@echo "📦 Available binaries:"
	@ls -la $(APP_NAME)-* 2>/dev/null || true

# Install dependencies
deps:
	@echo "📥 Installing dependencies..."
	cd $(SOURCE_DIR) && go mod download
	cd $(SOURCE_DIR) && go mod tidy

# Test the application
test:
	@echo "🧪 Running tests..."
	cd $(SOURCE_DIR) && go test -v .

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -f $(APP_NAME) $(APP_NAME)-linux-amd64 $(APP_NAME)-linux-arm64 $(APP_NAME)-darwin-amd64 $(APP_NAME)-darwin-arm64
	@echo "✅ Cleanup complete"

# Show build help
help:
	@echo "LAB Build System"
	@echo "==================="
	@echo ""
	@echo "Available targets:"
	@echo "  build              - Build for current platform"
	@echo "  build-all          - Build for all supported platforms"
	@echo "  build-linux-amd64  - Build for Linux AMD64"
	@echo "  build-linux-arm64  - Build for Linux ARM64"
	@echo "  build-darwin-amd64 - Build for macOS AMD64"
	@echo "  build-darwin-arm64 - Build for macOS ARM64 (M1/M2)"
	@echo "  deps              - Install Go dependencies"
	@echo "  test              - Run tests"
	@echo "  clean             - Clean build artifacts"
	@echo "  help              - Show this help"