FROM golang:1.20-alpine

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o ad-service ./cmd/ad-service

EXPOSE 8080

CMD ["./ad-service"]