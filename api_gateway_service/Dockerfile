FROM golang:1.23.6-alpine3.21 as builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /service ./cmd/main.go

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /service .
#COPY --from=builder /app/config ./config

RUN wget https://github.com/fullstorydev/grpcurl/releases/download/v1.9.1/grpcurl_1.9.1_linux_x86_64.tar.gz 
RUN tar -xzf grpcurl_1.9.1_linux_x86_64.tar.gz 
RUN mv grpcurl /usr/local/bin/ 
RUN  rm grpcurl_1.9.1_linux_x86_64.tar.gz
EXPOSE 8085
CMD [ "./service" ]