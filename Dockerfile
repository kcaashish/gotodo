# syntax=docker/dockerfile:1

FROM golang:1.18-alpine AS prepare

WORKDIR /app

COPY vendor .

FROM prepare AS build

COPY . .

RUN CGO_ENABLED=0 go build -o /gotodo ./cmd/gotodo/main.go

EXPOSE 8000

CMD ["/gotodo"]
