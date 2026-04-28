# ──────────────────────────────────────────────────────────────
#  AccountFlow — Makefile
# ──────────────────────────────────────────────────────────────

APP_NAME   := accountflow
CMD_PATH   := ./cmd/api
BINARY     := bin/$(APP_NAME)
DOCKER_IMG := accountflow-api

.DEFAULT_GOAL := help

# ── Help ──────────────────────────────────────────────────────
.PHONY: help
help:
	@echo ""
	@echo "  AccountFlow — available targets"
	@echo ""
	@echo "  Development"
	@echo "    make run          Run the app locally (requires local PostgreSQL)"
	@echo "    make build        Compile the binary to ./bin/accountflow"
	@echo "    make clean        Remove compiled binary"
	@echo ""
	@echo "  Docker"
	@echo "    make up           Build images and start all services"
	@echo "    make up-d         Same as up but in background (detached)"
	@echo "    make down         Stop all services"
	@echo "    make reset        Stop, wipe database volume and rebuild"
	@echo "    make logs         Tail API container logs"
	@echo "    make ps           Show running containers"
	@echo ""
	@echo "  Tests"
	@echo "    make test         Run all unit tests"
	@echo "    make test-v       Run all tests (verbose)"
	@echo "    make test-cover   Run tests and open HTML coverage report"
	@echo "    make test-domain  Run domain tests only"
	@echo "    make test-usecase Run use case tests only"
	@echo ""
	@echo "  Code quality"
	@echo "    make fmt          Format all Go files"
	@echo "    make vet          Run go vet"
	@echo "    make lint         Run go vet + fmt check"
	@echo ""

# ── Development ───────────────────────────────────────────────
.PHONY: run
run:
	go run $(CMD_PATH)/main.go

.PHONY: build
build:
	mkdir -p bin
	go build -o $(BINARY) $(CMD_PATH)

.PHONY: clean
clean:
	rm -rf bin/
	rm -f coverage.out

# ── Docker ────────────────────────────────────────────────────
.PHONY: up
up:
	docker compose up --build

.PHONY: up-d
up-d:
	docker compose up --build -d

.PHONY: down
down:
	docker compose down

.PHONY: reset
reset:
	docker compose down -v
	docker compose up --build

.PHONY: logs
logs:
	docker compose logs -f api

.PHONY: ps
ps:
	docker compose ps

# ── Tests ─────────────────────────────────────────────────────
.PHONY: test
test:
	go test ./...

.PHONY: test-v
test-v:
	go test -v ./...

.PHONY: test-cover
test-cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

.PHONY: test-domain
test-domain:
	go test -v ./internal/domain/...

.PHONY: test-usecase
test-usecase:
	go test -v ./internal/usecase/...

# ── Code quality ──────────────────────────────────────────────
.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: lint
lint: vet
	@echo "Checking formatting..."
	@test -z "$$(gofmt -l .)" || (echo "Files not formatted: $$(gofmt -l .)"; exit 1)
	@echo "All good."

# ── Dependencies ──────────────────────────────────────────────
.PHONY: tidy
tidy:
	go mod tidy

