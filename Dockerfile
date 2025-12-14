# Stage 1: Build
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install dependencies (cached)
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build binaries (Static linking for Alpine)
RUN CGO_ENABLED=0 GOOS=linux go build -o postificus-api ./cmd/api/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o postificus-worker ./cmd/worker/main.go

# Stage 2: Run (Tiny Image)
FROM alpine:latest

WORKDIR /root/
ENV APP_ENV=production

# Install Chromium/Chrome dependencies for Rod (CRITICAL for Alpine)
# Rod needs these libraries to run the browser even if Chrome is downloaded separately
RUN apk add --no-cache \
    chromium \
    nss \
    freetype \
    harfbuzz \
    ca-certificates \
    ttf-freefont

# Copy binaries from builder
COPY --from=builder /app/postificus-api .
COPY --from=builder /app/postificus-worker .


# Expose API port
EXPOSE 8080

# Default command (can be overridden in compose)
CMD ["./postificus-api"]
