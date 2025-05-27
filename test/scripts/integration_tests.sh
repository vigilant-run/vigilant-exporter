#!/bin/bash

set -e

PROJECT_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
mkdir -p "${PROJECT_ROOT}/build"

cd "${PROJECT_ROOT}"
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o "${PROJECT_ROOT}/build/vigilant-exporter" ./cmd/exporter
chmod +x "${PROJECT_ROOT}/build/vigilant-exporter"

cd "${PROJECT_ROOT}/test/integration"
docker compose down -v --remove-orphans
docker compose up --build -d server

timeout 30 bash -c 'until curl -sf http://localhost:8000/api/health >/dev/null 2>&1; do echo "Waiting for server..." && sleep 2; done'

docker compose up --build runner
exit_code=$(docker compose ps -aq runner | xargs docker inspect -f '{{.State.ExitCode}}')

docker compose down -v --remove-orphans

if [ "$exit_code" -eq 0 ]; then
    echo "Integration tests passed"
    exit 0
else
    echo "Integration tests failed"
    exit 1
fi