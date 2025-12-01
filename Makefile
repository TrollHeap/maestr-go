.PHONY: run migrate

build:
	@echo "Running the application..."
	templ generate
	go run cmd/app/main.go

run:
	@echo "Running the application..."
	templ generate
	go run cmd/app/main.go

migrate:
	@echo "Running the application..."
	go run cmd/migrate/main.go
