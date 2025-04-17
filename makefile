APP_NAME=auth-service
BINARY=$(APP_NAME)

all: generate build



# generate gRPC 
generate:
	protoc \
		--proto_path=./api/proto/auth \
		--go_out=./pkg/api/v1 \
		--go-grpc_out=./pkg/api/v1 \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		./api/proto/auth/auth_service.proto

# прогонка миграций
migrations-up:
	migrate -path migrations -database 'postgres://postgres:pass@localhost:5432/postgres?sslmode=disable' up

migrations-down:
	migrate -path migrations -database 'postgres://postgres:pass@localhost:5432/postgres?sslmode=disable' down


# сборка приложения
build:
	go build -o $(APP_NAME) cmd/main.go
	docker build -t auth-service .
	

# запуск приложения
run: build
	./$(BINARY)

# тестирование
# test: 

# интеграционное тестирование
test-integration:
	go test -v -tags=integration ./tests/integration/auth...

# юнит тестирование
test:
	go test -v -short ./...

.PHONY: all generate run-migrations build run test-integration test