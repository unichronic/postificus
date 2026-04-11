#!/bin/bash

# Ensure Docker infra is up
echo "Checking Docker infrastructure..."
docker compose up -d db redis rabbitmq

# Stop Docker app services if running
echo "Stopping Docker app services..."
docker compose stop api worker

# Trap Ctrl+C to kill all background processes
trap "kill 0" EXIT

echo "Starting Backend API..."
go run cmd/api/main.go &

echo "Starting Worker..."
go run cmd/worker/main.go &

echo "Services started. Press Ctrl+C to stop."
wait
