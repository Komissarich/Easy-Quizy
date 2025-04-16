FROM golang:1.23.5

WORKDIR /app

COPY . .
RUN go mod download

COPY *.go ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /stat_service ./cmd/main.go

EXPOSE 8080
EXPOSE 50051
# ENTRYPOINT [ "/cmd/bin" ]

# Run
CMD ["/stat_service"]