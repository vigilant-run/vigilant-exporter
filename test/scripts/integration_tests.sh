#!/bin/bash

set -e

echo "Starting integration tests..."

cd "$(dirname "$0")/../integration"

echo "Stopping and removing existing containers..."
docker compose down -v --remove-orphans

echo "Starting services..."
docker compose up --build -d server

echo "Waiting for health check..."
timeout 30 bash -c 'until curl -f http://localhost:8000/api/health >/dev/null 2>&1; do sleep 2; done'

echo "Running integration tests..."
docker compose up runner

if [ $? -eq 0 ]; then
    echo "Integration tests passed!"
else
    echo "Integration tests failed!"
    exit 1
fi

echo "Cleaning up..."
docker compose down -v

echo "Integration tests completed!"