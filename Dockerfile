FROM golang:alpine as build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o app .

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/app .

EXPOSE 8080

CMD ["app"]
