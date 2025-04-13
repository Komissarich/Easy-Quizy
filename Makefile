# BINARY_NAME = app
# PROTO_FILE = order.proto
# GENERATED_DIR = ./pkg/api/v1
# GOOGLEAPIS_PATH = ./api/google/api
# GOPATH = $(shell go env GOPATH)
# MAIN_GO_PATH = ./cmd/app/main.go
# FULL_DIR = ./pkg/api

# .PHONY: all generate build run clean

# all: build

# generate:
# 	@if not exist "$(GENERATED_DIR)" mkdir "$(GENERATED_DIR)"
# 	protoc -I. \
# 	-I"$(GOOGLEAPIS_PATH)" \
# 	--plugin=protoc-gen-go="$(GOPATH)/bin/protoc-gen-go.exe" \
# 	--plugin=protoc-gen-go-grpc="$(GOPATH)/bin/protoc-gen-go-grpc.exe" \
# 	--plugin=protoc-gen-grpc-gateway="$(GOPATH)/bin/protoc-gen-grpc-gateway.exe" \
# 	--go_out="$(GENERATED_DIR)" \
# 	--go_opt=paths=source_relative \
# 	--go-grpc_out="$(GENERATED_DIR)" \
# 	--go-grpc_opt=paths=source_relative \
# 	--grpc-gateway_out=logtostderr=true,paths=source_relative:"$(GENERATED_DIR)" \
# 	"$(PROTO_FILE)"

# build: generate
# 	go build -o "$(BINARY_NAME)" "$(MAIN_GO_PATH)"

# run: build
# 	.\$(BINARY_NAME)

# clean:
# 	del /f "$(BINARY_NAME)" 2> nul
# 	if exist "$(FULL_DIR)" rmdir /s /q "$(FULL_DIR)"


# Переменные
GO=go#$(shell go env GOPATH)#C:\Program Files\Go\bin\go.exe
PROTOC=protoc
MODULE_NAME=grpc_service
PATH=cmd
PROTO_PATH=api/proto
APP_NAME=main

# Цель по умолчанию
all: proto build run

# Проверка и установка зависимостей
deps:
	@echo "Checking dependencies..."
	$(GO) mod tidy
	$(GO) mod download

# Генерация .proto файлов
proto:
	@echo "Generating .proto files..."
	protoc --go_out=. --go-grpc_out=. --grpc-gateway_out=. api/proto/statistics/session.proto api/proto/quizzes/quizzes.proto api/proto/users/users.proto

# Сборка приложения
build:
	@echo "Building application..."
	$(GO) build -o $(PATH)\\bin\\$(APP_NAME).exe $(PATH)\\app\\$(APP_NAME).go

# Запуск приложения
run:
	@echo "Running application..."
	$(PATH)\\bin\\$(APP_NAME).exe

# Очистка
clean:
	@echo "Cleaning up..."
	del /q $(PATH)\\bin\\$(APP_NAME).exe

.PHONY: all deps proto build run clean