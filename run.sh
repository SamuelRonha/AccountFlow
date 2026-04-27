#!/usr/bin/env bash
# run.sh — AccountFlow development helper
#
# Usage:
#   ./run.sh            → start API + PostgreSQL via Docker Compose
#   ./run.sh local      → run the API locally (requires PostgreSQL running)
#   ./run.sh test       → run all unit tests
#   ./run.sh test:v     → run tests with verbose output
#   ./run.sh test:cover → run tests and open HTML coverage report
#   ./run.sh down       → stop and remove containers
#   ./run.sh reset      → stop containers and wipe the database volume

set -e

CMD="${1:-docker}"

case "$CMD" in
  docker)
    echo "▶  Starting AccountFlow (Docker Compose)..."
    docker compose up --build
    ;;

  local)
    echo "▶  Starting AccountFlow locally..."
    export DB_HOST="${DB_HOST:-localhost}"
    export DB_PORT="${DB_PORT:-5432}"
    export DB_USER="${DB_USER:-postgres}"
    export DB_PASSWORD="${DB_PASSWORD:-postgres}"
    export DB_NAME="${DB_NAME:-accountflow}"
    export DB_SSLMODE="${DB_SSLMODE:-disable}"
    export APP_PORT="${APP_PORT:-8072}"
    go run ./cmd/api/main.go
    ;;

  test)
    echo "▶  Running unit tests..."
    go test ./...
    ;;

  test:v)
    echo "▶  Running unit tests (verbose)..."
    go test -v ./...
    ;;

  test:cover)
    echo "▶  Running tests with coverage..."
    go test -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out
    ;;

  down)
    echo "▶  Stopping containers..."
    docker compose down
    ;;

  reset)
    echo "▶  Stopping containers and wiping database..."
    docker compose down -v
    ;;

  *)
    echo "AccountFlow — run script"
    echo ""
    echo "Usage:"
    echo "  ./run.sh            start API + PostgreSQL via Docker Compose"
    echo "  ./run.sh local      run the API locally (needs PostgreSQL running)"
    echo "  ./run.sh test       run all unit tests"
    echo "  ./run.sh test:v     run tests with verbose output"
    echo "  ./run.sh test:cover run tests and open HTML coverage report"
    echo "  ./run.sh down       stop containers"
    echo "  ./run.sh reset      stop containers and wipe database volume"
    ;;
esac

