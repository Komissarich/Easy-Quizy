APP_NAME=auth-service
BINARY=$(APP_NAME)
DB_URL=postgres://postgres:05042007PULlup!@auth-db:5432/postgres?sslmode=disable

all: generate build

# Генерация gRPC кода
generate:
	protoc \
		--proto_path=./api/proto/auth \
		--go_out=./pkg/api/v1 \
		--go-grpc_out=./pkg/api/v1 \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		./api/proto/auth/auth_service.proto

# Миграции в Docker-окружении
migrations-up:
	docker-compose up -d postgres

	docker-compose run --rm auth-service migrate -path /app/migrations -database '$(DB_URL)' up

migrations-down:
	docker-compose run --rm auth-service migrate -path /app/migrations -database '$(DB_URL)' down

# Сборка приложения
build:
	go build -o $(APP_NAME) cmd/main.go
	docker-compose build

# Запуск приложения
run: build migrations-up
	docker-compose up auth-service

# Тестирование
test-integration:
	docker-compose run --rm auth-service go test -v -tags=integration ./tests/integration/auth...

test:
	docker-compose run --rm auth-service go test -v -short ./...

# Очистка
clean:
	docker-compose down -v
	rm -f $(BINARY)

.PHONY: all generate migrations-up migrations-down build run test-integration test clean