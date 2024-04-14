FROM golang:1.21.8-alpine3.19 as builder
RUN apk add --no-cache git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd/
COPY internal ./internal/
COPY pkg ./pkg/
COPY ./tests ./tests
