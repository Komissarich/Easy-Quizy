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

# сборка приложения
build:
	go build -o $(APP_NAME) cmd/main.go
	docker build -t divineheart/auth-service .
	
# запуск приложения
run: build
	./$(BINARY)

.PHONY: all generate build run
