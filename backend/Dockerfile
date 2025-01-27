# backend/Dockerfile

# Build Stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o server cmd/main.go

# Runtime Stage
FROM alpine:latest

WORKDIR /app

# Copy the compiled binary, .env file, and migrations directory
COPY --from=builder /app/server .
COPY --from=builder /app/.env .env
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/internal/localization ./internal/localization

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./server"]
