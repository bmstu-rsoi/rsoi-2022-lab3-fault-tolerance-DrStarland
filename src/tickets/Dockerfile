FROM golang:1.19 AS builder

ENV GO111MODULE=on

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o app ./cmd/tickets

EXPOSE 8070
ENTRYPOINT ["./app"]