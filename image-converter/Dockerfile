# Stage 1: Build the binary
FROM golang:1.20-alpine AS builder
WORKDIR /app

# Install build tools and the libwebp development package.
RUN apk update && apk add --no-cache build-base libwebp-dev

# Copy go.mod and go.sum first for caching dependencies.
COPY go.mod go.sum ./
RUN go mod download

# Copy all source code (adjust path if needed).
COPY . .

# Build the image-converter binary.
# Adjust the path to your main package file if it's not in "backend/image-converter/main.go"
RUN CGO_ENABLED=1 GOOS=linux go build -o image-converter ./backend/image-converter/main.go

# Stage 2: Create a minimal runtime image.
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy the binary from the builder stage.
COPY --from=builder /app/image-converter .

# Expose the port your image-converter will run on (e.g., 8001).
EXPOSE 8001

# Run the image-converter binary.
CMD ["./image-converter"]
