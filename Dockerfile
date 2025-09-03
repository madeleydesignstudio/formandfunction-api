# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server

# Release stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 3000

CMD ["/app/server"]
