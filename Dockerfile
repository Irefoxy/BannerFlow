FROM golang:1.21.8-alpine3.19 as builder
RUN apk add --no-cache git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd/
COPY internal ./internal/
COPY pkg ./pkg/
RUN go test -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -v -installsuffix cgo -o bannerflow ./cmd/bannerflow

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/bannerflow .

CMD ["./bannerflow"]
