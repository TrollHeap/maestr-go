.PHONY: run migrate

run:
	@echo "Running the application..."
	go run cmd/app/main.go

migrate:
	@echo "Running the application..."
	go run cmd/migrate/main.go
