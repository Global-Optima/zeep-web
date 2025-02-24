# ====================
# Build Stage
# ====================
FROM golang:1.23-alpine AS builder

# Set environment variables for Go
ENV CGO_ENABLED=0 

# Set working directory
WORKDIR /app

# Build arguments for passing env variables at build time
ARG DB_HOST
ARG DB_PORT
ARG DB_USER
ARG DB_PASSWORD
ARG DB_NAME
ARG JWT_SECRET_KEY
ARG REDIS_HOST
ARG REDIS_PORT
ARG REDIS_PASSWORD
ARG REDIS_DB

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
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/internal/localization ./internal/localization
COPY --from=builder /app/pkg/utils/censor/dictionaries.json ./pkg/utils/censor/dictionaries.json

# Copy the .env file from the compose build arguments
ENV DB_HOST=${DB_HOST}
ENV DB_PORT=${DB_PORT}
ENV DB_USER=${DB_USER}
ENV DB_PASSWORD=${DB_PASSWORD}
ENV DB_NAME=${DB_NAME}
ENV JWT_SECRET_KEY=${JWT_SECRET_KEY}
ENV REDIS_HOST=${REDIS_HOST}
ENV REDIS_PORT=${REDIS_PORT}
ENV REDIS_PASSWORD=${REDIS_PASSWORD}
ENV REDIS_DB=${REDIS_DB}

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./server"]
