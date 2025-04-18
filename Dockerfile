FROM golang:1.24.2-alpine3.21 as builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /stat_service ./cmd/app/main.go

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /stat_service .
COPY --from=builder /app/config ./config
COPY --from=builder /app/db/migrations ./db/migrations

EXPOSE 50051
EXPOSE 8080

CMD [ "./stat_service" ]