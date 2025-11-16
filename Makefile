# --- Параметры по умолчанию ---
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_USER ?= postgres
DB_PASSWORD ?= password
DB_NAME ?= pr
SSL_MODE ?= disable

DB_URL := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(SSL_MODE)

# --- Пути в проекте ---
MIGRATIONS_DIR := migrations
APP_DIR := cmd/app

# --- Сборка и запуск приложения ---
run:
	cd $(APP_DIR) && go run main.go

build:
	cd $(APP_DIR) && go build -o /bin/app

# --- Миграции ---
migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down 1

migrate-drop:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" drop -f

migrate-force:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" force $(version)

migrate-version:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" version

migrate-create:
	@read -p "Migration name: " name; \
	migrate create -seq -ext sql -dir $(MIGRATIONS_DIR) $$name

# --- Генерация моков ---
generate:
	go generate ./...

# --- Тесты с покрытием ---
test:
	go test ./... -coverprofile=coverage.out
	@echo "\nAverage coverage:"
	@go tool cover -func=coverage.out \
		| grep -v 'mock_' \
		| grep -v 'api' \
		| grep -v 'main' \
		| grep -v '^total:' \
		| awk '{ s+=$$3; n++ } END { if (n > 0) printf("%.1f%%\n", s/n) }'