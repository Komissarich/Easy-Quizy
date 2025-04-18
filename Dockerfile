FROM golang:1.24.2-alpine AS builder
RUN apk add --no-cache protoc git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN protoc -I . --go_out=. --go-grpc_out=. --grpc-gateway_out=. ./quiz.proto
RUN go build -o quiz_app cmd/quiz/main.go

FROM alpine:latest
RUN apk add --no-cache postgresql-client
WORKDIR /app
COPY --from=builder /app/quiz_app .
COPY db/migrations/init_up.sql /app/db/migrations/
COPY init-db.sh /app/
COPY config/config.yaml /app/config/
RUN chmod +x /app/init-db.sh
EXPOSE 50051 8080
CMD ["sh", "-c", "/app/init-db.sh && /app/quiz_app"]