FROM golang:1.23.6-alpine3.21 as builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /service ./cmd/quiz/main.go

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /service .
COPY --from=builder /app/config ./config

EXPOSE 50053
EXPOSE 8080
CMD [ "./service" ]