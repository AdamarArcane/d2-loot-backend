# Build Stage
FROM golang:1.23.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o d2-loot-backend .

# Run Stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/d2-loot-backend .

EXPOSE 8080

CMD ["./d2-loot-backend"]