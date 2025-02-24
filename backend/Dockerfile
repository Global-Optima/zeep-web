# backend/Dockerfile

# ====================
# Build Stage
# ====================
FROM golang:1.23-alpine AS builder

# Set environment variables for Go
ENV CGO_ENABLED=0 

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files first to leverage caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy the source code
COPY . .

# Build the application
RUN go build -ldflags="-s -w" -o server cmd/main.go

# ====================
# Runtime Stage
# ====================
FROM alpine:latest AS runtime

# Install necessary tools (e.g., ca-certificates for SSL)
RUN apk add --no-cache bash ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy the compiled binary and required assets from the builder stage
COPY --from=builder /app/server .
COPY --from=builder /app/.env .env
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/internal/localization ./internal/localization
COPY --from=builder /app/pkg/utils/censor/dictionaries.json ./pkg/utils/censor/dictionaries.json

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./server"]