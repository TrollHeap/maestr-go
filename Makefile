.PHONY: run build clean dev css test migrate

# === RUN ===
run:
	@echo "Running the application..."
	@templ generate
	@go run cmd/app/main.go

# === BUILD ===
build:
	@echo "Building application..."
	@templ generate
	@./scripts/build-css.sh
	@go build -o bin/maestro cmd/app/main.go
	@echo "✅ Build complete: bin/maestro"

# === DEV MODE ===
dev:
	@echo "🚀 Starting dev mode..."
	@trap 'kill 0' EXIT; \
	./scripts/watch-css.sh & \
	templ generate --watch --proxy="http://localhost:8080" --cmd="go run cmd/app/main.go"

# === CSS ===
css:
	@./scripts/build-css.sh

css-watch:
	@./scripts/watch-css.sh

# === MIGRATION ===
migrate:
	@echo "Running migration..."
	@go run cmd/migrate/main.go

# === TEST ===
test:
	@go test ./...

# === CLEAN ===
clean:
	@rm -rf bin/
	@rm -rf public/css/
	@echo "✅ Clean complete"

# === INSTALL DEPS ===
install:
	@echo "Installing dependencies..."
	@go mod download
	@templ generate
	@./scripts/install-tailwind.sh
	@echo "✅ Installation complete"
