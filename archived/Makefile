.PHONY: build run test clean fmt lint dev install

# Build
build:
	@mkdir -p bin
	go build -o bin/maestro-server ./cmd/maestro-server
	@echo "✅ Built: bin/maestro-server"

# Run the server
run: build
	@echo "🎯 Starting Maestro Server..."
	./bin/maestro-server -port 8080

# Run with custom port
run-port:
	@read -p "Enter port [8080]: " port; \
	port=$${port:-8080}; \
	./bin/maestro-server -port $$port

# Run with custom data directory
run-custom:
	@read -p "Enter data directory: " datadir; \
	./bin/maestro-server -data-dir "$$datadir"

# Test
test:
	@echo "🧪 Running tests..."
	go test -v ./...

# Test with coverage
test-coverage:
	@echo "🧪 Running tests with coverage..."
	go test -cover ./...

# Test specific package
test-domain:
	go test -v ./internal/domain/...

test-storage:
	go test -v ./internal/storage/...

# Format code
fmt:
	@echo "📝 Formatting code..."
	go fmt ./...

# Lint code
lint:
	@echo "🔍 Linting code..."
	golangci-lint run ./... || echo "Install golangci-lint: https://golangci-lint.run/usage/install/"

# Clean
clean:
	@echo "🧹 Cleaning..."
	rm -rf bin/
	go clean -modcache
	@echo "✅ Done"

# Mod download and tidy
deps:
	@echo "📦 Downloading dependencies..."
	go mod download
	go mod tidy

# Dev mode (requires air for hot reload)
dev:
	@echo "🔄 Dev mode (hot reload)..."
	air || echo "Install air: go install github.com/cosmtrek/air@latest"

# Help
help:
	@echo "Maestro Makefile commands:"
	@echo ""
	@echo "  make build           - Build maestro-server binary"
	@echo "  make run             - Build and run on port 8080"
	@echo "  make run-port        - Run on custom port"
	@echo "  make run-custom      - Run with custom data directory"
	@echo "  make test            - Run all tests"
	@echo "  make test-coverage   - Run tests with coverage"
	@echo "  make test-domain     - Test domain package only"
	@echo "  make test-storage    - Test storage package only"
	@echo "  make fmt             - Format code"
	@echo "  make lint            - Lint code"
	@echo "  make clean           - Clean build artifacts"
	@echo "  make deps            - Download/tidy dependencies"
	@echo "  make dev             - Dev mode with hot reload"
	@echo "  make help            - Show this help"

.DEFAULT_GOAL := help
