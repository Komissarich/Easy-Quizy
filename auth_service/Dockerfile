FROM golang:1.24.2-alpine AS builder

WORKDIR /app
COPY . .

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/auth-service ./cmd/main.go

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/auth-service .
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate
COPY --from=builder /app/migrations ./migrations

EXPOSE 50051
CMD ["./auth-service"]