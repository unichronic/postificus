# Build Stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final Stage
FROM alpine:latest

# Install Chromium for go-rod
RUN apk add --no-cache chromium

# Set environment variables for go-rod
ENV ROD_BIN=/usr/bin/chromium-browser

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]
