# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app


RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache postgresql-client

COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations
COPY config/config.yaml ./config/config.yaml

EXPOSE 8080

CMD ["./main"]