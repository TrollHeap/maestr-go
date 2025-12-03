.PHONY: run build clean dev css css-watch test migrate

# === DEV MODE (Recommandé) ===
dev:
	@echo "🚀 Mode développement..."
	@make css
	@trap 'kill 0' EXIT; \
	./scripts/watch-css.sh & \
	templ generate --watch --proxy="http://localhost:7331" --cmd="go run cmd/app/main.go"

# === CSS BUILD ===
css:
	@echo "🎨 Building CSS..."
	@./bin/tailwindcss -i ./public/css/input.css -o ./public/css/style.css --minify
	@echo "✅ CSS compilé: $(shell pwd)/public/css/style.css"

# === CSS WATCH (debug) ===
css-watch:
	@echo "👀 Watching CSS..."
	@./bin/tailwindcss -i ./public/css/input.css -o ./public/css/style.css --watch

# === RUN (sans watch) ===
run:
	@echo "🚀 Running..."
	@make css
	@templ generate
	@go run cmd/app/main.go

# === BUILD PROD ===
build:
	@echo "🔨 Building for production..."
	@templ generate
	@./bin/tailwindcss -i ./public/css/input.css -o ./public/css/style.css --minify
	@go build -ldflags="-s -w" -o bin/maestro cmd/app/main.go
	@echo "✅ Build complete: bin/maestro"

# === CLEAN ===
clean:
	@rm -rf bin/
	@rm -f public/css/style.css
	@find . -name "*_templ.go" -delete
	@echo "✅ Clean complete"

# === TEST CSS SIZE ===
css-size:
	@echo "📊 CSS Size:"
	@ls -lh public/css/style.css | awk '{print $$5}'
	@echo "📊 Gzip Size:"
	@gzip -c public/css/style.css | wc -c | numfmt --to=iec-i --suffix=B


migrate:
	rm -rf data/maestro.db && rm -rf data/maestro.db-shm && rm -rf data/maestro.db-wal
	go run cmd/migrate/main.go
