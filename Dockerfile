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



# -------------------------------------------------------
# Stage 2: Run (Alpine)
# -------------------------------------------------------
FROM alpine:latest

WORKDIR /app

# Install Chromium/Chrome dependencies for Rod (CRITICAL for Alpine)
RUN apk add --no-cache \
    chromium \
    nss \
    freetype \
    harfbuzz \
    ca-certificates \
    ttf-freefont

# Copy binaries
COPY --from=builder /app/postificus-api .
COPY --from=builder /app/postificus-worker .

# Expose API port
EXPOSE 8080

# Command is set by docker-compose (either api or worker)
CMD ["./postificus-api"]
