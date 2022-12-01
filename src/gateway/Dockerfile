FROM golang:1.19 AS builder

ENV GO111MODULE=on

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o app ./cmd/gateway

EXPOSE 8080

ENTRYPOINT ["./app"]