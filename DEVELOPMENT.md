# Development Guide

Everything a developer needs to work on, extend, and maintain AccountFlow.

---

## Table of Contents

- [Prerequisites](#prerequisites)
- [Project Setup](#project-setup)
- [Project Structure](#project-structure)
- [Architecture Rules](#architecture-rules)
- [Amount Sign Contract](#amount-sign-contract)
- [Running the Application](#running-the-application)
- [Running Tests](#running-tests)
- [Writing Tests](#writing-tests)
- [Adding a New Feature](#adding-a-new-feature)
- [Adding a New Migration](#adding-a-new-migration)
- [Environment Variables](#environment-variables)
- [Docker](#docker)
- [Code Style & Conventions](#code-style--conventions)
- [Common Mistakes](#common-mistakes)

---

## Prerequisites

| Tool           | Version              | Download                                  |
|----------------|----------------------|-------------------------------------------|
| Go             | 1.25+                | https://go.dev/dl/                        |
| Docker         | 24+                  | https://www.docker.com/                   |
| Docker Compose | 2+                   | https://docs.docker.com/compose/install/  |
| PostgreSQL     | 16 (optional, local) | https://www.postgresql.org/               |

> **IDE:** GoLand or VS Code with the Go extension are recommended.
> If the IDE fails to resolve imports, run `go mod tidy` and invalidate caches.

---

## Project Setup

```bash
# 1. Clone the repository
git clone https://github.com/your-org/AccountFlow.git
cd AccountFlow

# 2. Install all dependencies
go mod tidy

# 3. Verify the project builds
go build ./...

# 4. Run all tests
go test ./...
```

---

## Project Structure

```
AccountFlow/
|-- cmd/
|   `-- api/
|       `-- main.go                            # Entry point: wires all layers together
|-- internal/
|   |-- domain/                                # Layer 1 - Business rules (zero external deps)
|   |   |-- account.go                         # Account entity
|   |   |-- operation_type.go                  # OperationType entity + sign validation
|   |   |-- transaction.go                     # Transaction entity + amount sign contract
|   |   `-- errors.go                          # Sentinel domain errors
|   |-- usecase/                               # Layer 2 - Application logic
|   |   |-- account_usecase.go                 # CreateAccount, GetByID
|   |   |-- transaction_usecase.go             # CreateTransaction
|   |   `-- mocks/
|   |       `-- repositories.go                # Manual mocks for unit tests
|   |-- repository/                            # Layer 3 - Repository interfaces (contracts)
|   |   |-- account_repository.go
|   |   |-- operation_type_repository.go
|   |   `-- transaction_repository.go
|   |-- delivery/
|   |   `-- http/                              # Layer 4 - HTTP handlers (Gin)
|   |       |-- router.go
|   |       |-- account_handler.go
|   |       |-- transaction_handler.go
|   |       `-- response.go
|   `-- infra/
|       `-- postgres/                          # Layer 5 - Concrete implementations
|           |-- db.go
|           |-- account_repo.go
|           |-- operation_type_repo.go
|           `-- transaction_repo.go
|-- migrations/
|   |-- 000001_create_users.up.sql             # Creates operation_types + seeds 4 rows
|   |-- 000001_create_users.down.sql
|   |-- 000002_create_accounts.up.sql          # Creates accounts table
|   |-- 000002_create_accounts.down.sql
|   |-- 000003_create_transactions.up.sql      # Creates transactions table
|   `-- 000003_create_transactions.down.sql
|-- run.ps1                                    # Windows run helper
|-- run.sh                                     # Linux/Mac run helper
|-- Dockerfile
|-- docker-compose.yml
|-- go.mod
`-- go.sum
```

---

## Architecture Rules

This project follows **Clean Architecture**. The golden rule:

> **Dependencies always point inward.**
> Outer layers may depend on inner layers. Inner layers must NEVER depend on outer layers.

```
delivery/http  ->  usecase  ->  domain
infra/postgres ->  repository interfaces  <-  usecase
```

### Allowed
- `usecase` importing from `domain` and `repository`
- `delivery/http` importing from `usecase`
- `infra/postgres` implementing interfaces from `repository`
- `cmd/api/main.go` importing from any layer (composition root)

### Never allowed
- `domain` importing from `usecase`, `delivery`, or `infra`
- `usecase` importing from `delivery` or `infra`
- `usecase` importing `gin` or any HTTP framework
- `domain` importing any external library (except `github.com/google/uuid`)

---

## Amount Sign Contract

The most important business rule before touching any transaction code.

| Operation Type           | ID | Stored amount | Example |
|--------------------------|----|---------------|---------|
| Normal Purchase          | 1  | negative      | -50.0   |
| Purchase with Installments | 2 | negative     | -23.5   |
| Withdrawal               | 3  | negative      | -18.7   |
| Credit Voucher           | 4  | positive      | +60.0   |

**The HTTP caller is responsible for sending the correct sign.**
`NewTransaction()` validates the sign and returns `ErrInvalidAmount` if it contradicts the operation type.

- Where to look: `internal/domain/transaction.go` -> `NewTransaction()`
- Where to change if the rule changes: **only that one function**

---

## Running the Application

### With the run script (recommended)

```bash
# Windows PowerShell
.\run.ps1              # start API + PostgreSQL via Docker Compose
.\run.ps1 local        # run the API locally (PostgreSQL must be running)
.\run.ps1 test         # run all unit tests
.\run.ps1 test:v       # run tests verbose
.\run.ps1 test:cover   # run tests + open HTML coverage report
.\run.ps1 down         # stop containers
.\run.ps1 reset        # stop containers and wipe the database volume

# Linux / Mac
chmod +x run.sh
./run.sh               # same commands
```

### With Docker directly

```bash
docker compose up --build        # start everything
docker compose up --build -d     # start in background
docker compose logs -f api       # follow API logs
docker compose down              # stop
docker compose down -v           # stop + wipe database
```

### Locally without Docker

```powershell
# PowerShell
$env:DB_HOST="localhost"; $env:DB_PORT="5432"; $env:DB_USER="postgres"
$env:DB_PASSWORD="postgres"; $env:DB_NAME="accountflow"; $env:APP_PORT="8072"
go run ./cmd/api/main.go
```

```bash
# Bash
export DB_HOST=localhost DB_PORT=5432 DB_USER=postgres \
       DB_PASSWORD=postgres DB_NAME=accountflow APP_PORT=8072
go run ./cmd/api/main.go
```

> Migrations run automatically on every startup via `golang-migrate`.

---

## Running Tests

```bash
# All tests
go test ./...

# Verbose — shows every subtest name
go test -v ./...

# Specific package
go test -v ./internal/domain/...
go test -v ./internal/usecase/...

# Single test or subtest
go test -v -run TestCreateAccount ./internal/usecase/...
go test -v -run "TestCreateTransaction/normal_purchase" ./internal/usecase/...

# Coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

> Unit tests do NOT require a running database. All repository calls use mocks.

---

## Writing Tests

### Style: table-driven with t.Run

All tests use `t.Run` subtests so failures read like plain English:

```go
func TestCreateTransaction(t *testing.T) {
    cases := []struct {
        name    string
        amount  float64
        wantErr error
    }{
        {"normal purchase stores -50.0",        -50.0, nil},
        {"positive amount on debit rejected",    50.0, domain.ErrInvalidAmount},
        {"zero amount rejected",                  0.0, domain.ErrInvalidAmount},
    }
    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            // arrange, act, assert
        })
    }
}
```

### Domain tests (internal/domain/)

Test entity behaviour and business rules directly. No mocks needed.

```go
func TestNewAccount(t *testing.T) {
    t.Run("valid document creates account", func(t *testing.T) {
        acc, err := domain.NewAccount("12345678900")
        if err != nil {
            t.Fatalf("unexpected error: %v", err)
        }
        if acc.AccountID == uuid.Nil {
            t.Error("AccountID must not be nil")
        }
    })
}
```

### Use case tests (internal/usecase/)

Use the mocks in `internal/usecase/mocks/`.
Only set the `Fn` fields the test actually needs.
**Unset fields panic with a descriptive message** — this catches missing setup early.

```go
repo := &mocks.MockAccountRepository{
    CreateFn: func(_ context.Context, _ *domain.Account) error { return nil },
    // FindByIDFn intentionally unset — will panic with a clear message if called
}

// Assert call counts after the test:
if repo.CreateCalls != 1 {
    t.Errorf("Create called %d times, want 1", repo.CreateCalls)
}
```

### Mock structure (internal/usecase/mocks/repositories.go)

```go
type MockAccountRepository struct {
    CreateFn   func(ctx context.Context, account *domain.Account) error
    FindByIDFn func ctx context.Context, accountID uuid.UUID) (*domain.Account, error)

    // Incremented automatically on every call
    CreateCalls   int
    FindByIDCalls int
}
```

---

## Adding a New Feature

Follow this order to keep architecture clean:

### 1. Domain — entity or rule

Add entity in `internal/domain/`. Add errors in `internal/domain/errors.go`.

```go
var ErrExampleNotFound = errors.New("example not found")
```

### 2. Repository — define the interface

Create `internal/repository/example_repository.go`:

```go
type ExampleRepository interface {
    Create(ctx context.Context, e *domain.Example) error
    FindByID(ctx context.Context, id uuid.UUID) (*domain.Example, error)
}
```

### 3. Use case — implement business logic

Create `internal/usecase/example_usecase.go`.
Depend only on repository interfaces — never on infra directly.

### 4. Mock — add test double

Add to `internal/usecase/mocks/repositories.go` following the same pattern
(Fn fields + call counters + descriptive panic on nil call).

### 5. Tests — write use case tests

Create `internal/usecase/example_usecase_test.go`.
Use `t.Run` subtests. Cover success + every error path.

### 6. Infra — implement the repository

Create `internal/infra/postgres/example_repo.go` implementing the interface with `*sql.DB`.

### 7. Delivery — add the HTTP handler

Create `internal/delivery/http/example_handler.go` (Gin).
Register the route in `router.go`.
Map new domain errors in `response.go` -> `mapDomainError()`.

### 8. Wire — connect in main.go

Instantiate and inject all new dependencies in `cmd/api/main.go`.

---

## Adding a New Migration

```
migrations/
  000004_create_examples.up.sql
  000004_create_examples.down.sql
```

Naming: `{sequence}_{short_description}.up.sql`

```sql
-- 000004_create_examples.up.sql
CREATE TABLE examples (
    id         UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    name       VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- 000004_create_examples.down.sql
DROP TABLE IF EXISTS examples;
```

Always create both `.up.sql` and `.down.sql`.

Migrations run automatically on startup. To run manually:

```bash
migrate -path ./migrations \
  -database "postgres://postgres:postgres@localhost:5432/accountflow?sslmode=disable" up
```

---

## Environment Variables

| Variable      | Description              | Default       |
|---------------|--------------------------|---------------|
| `DB_HOST`     | PostgreSQL host          | `localhost`   |
| `DB_PORT`     | PostgreSQL port          | `5432`        |
| `DB_USER`     | PostgreSQL user          | `postgres`    |
| `DB_PASSWORD` | PostgreSQL password      | `postgres`    |
| `DB_NAME`     | PostgreSQL database name | `accountflow` |
| `DB_SSLMODE`  | PostgreSQL SSL mode      | `disable`     |
| `APP_PORT`    | Application HTTP port    | `8072`        |

---

## Docker

### Dockerfile

Multi-stage build:

1. **Builder** (`golang:1.25-alpine`) — compiles binary with `CGO_ENABLED=0`
2. **Runtime** (`alpine:3.19`) — ships only the compiled binary, minimal image size

### docker-compose services

| Service    | Description             | Port   |
|------------|-------------------------|--------|
| `postgres` | PostgreSQL 16           | `5432` |
| `api`      | AccountFlow application | `8072` |

The `api` service has a `depends_on` healthcheck and only starts after PostgreSQL is ready.

---

## Code Style & Conventions

- **Format:** run `go fmt ./...` before every commit
- **Vet:** run `go vet ./...` to catch common issues
- **Errors:** never ignore errors silently — always return or log them
- **Domain errors:** define in `internal/domain/errors.go` as sentinel `var` errors
- **HTTP error mapping:** only in `internal/delivery/http/response.go` inside `mapDomainError()`
- **Context:** always pass `context.Context` as the first parameter in repo and use case methods
- **Naming:**
  - Files: `snake_case.go`
  - Types / Functions: `PascalCase`
  - Variables / Parameters: `camelCase`
- **Tests:** file ends with `_test.go`, package is `<pkg>_test` (black-box preferred)
- **Test style:** use `t.Run` subtests with descriptive names; prefer table-driven for multiple cases

---

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| Adding `gin` import to `usecase` or `domain` | Move HTTP logic to `delivery/http` |
| Calling infra directly from usecase | Always depend on the repository interface; inject via constructor |
| Hardcoding DB credentials | Use env vars — `getEnv()` is in `cmd/api/main.go` |
| Forgetting `.down.sql` migration | Always create both `up` and `down` files |
| Sending positive amount for debit op (1,2,3) | Debit ops require negative amounts — see Amount Sign Contract |
| Sending negative amount for credit op (4) | Credit op requires positive amount — see Amount Sign Contract |
| Not writing tests for a new use case | Every `_usecase.go` must have a `_usecase_test.go` |
| IDE not resolving imports | Run `go mod tidy`, then invalidate IDE caches (GoLand: File -> Invalidate Caches) |
