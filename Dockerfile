# Build Stage
FROM golang:1.23.1-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to leverage Docker layer caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of your application's source code
COPY . .

# Build the Go application
RUN go build -o d2-loot-backend .

# Run Stage
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/d2-loot-backend .

# Copy the JSON file into the final image (if needed at runtime)
COPY --from=builder /app/cmd/generate_constants/weapons_and_perks.json ./cmd/generate_constants/weapons_and_perks.json

# Expose the port your application listens on (adjust if necessary)
EXPOSE 8080

# Command to run the application
CMD ["./d2-loot-backend"]