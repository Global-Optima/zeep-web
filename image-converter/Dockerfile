# Stage 1: Build the binary
FROM golang:1.24.1-alpine AS builder

# Install build tools and the libwebp development package.
RUN apk update && apk add --no-cache \
    build-base \
    libwebp-dev

WORKDIR /app

# Enable CGO explicitly. 
ENV CGO_ENABLED=1

# Copy go.mod and go.sum for dependency caching.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of your source code into the container.
COPY . .

# Build your application. Adjust the path to your main.go if needed.
RUN go build -o image-converter main.go

# Stage 2: Create a minimal runtime image.
FROM alpine:latest

# Install runtime dependencies: libwebp + CA certificates (if HTTPS calls are needed).
RUN apk --no-cache add ca-certificates libwebp

WORKDIR /root

# Copy the compiled binary from the builder stage.
COPY --from=builder /app/image-converter .

EXPOSE 8082

CMD ["./image-converter"]
