# Use a valid Go version (e.g., 1.21)
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o ad-service ./cmd/ad-service

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/ad-service .

EXPOSE 8080

CMD ["./ad-service"]