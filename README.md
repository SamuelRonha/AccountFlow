# AccountFlow
A financial account management system built with **Go**, following **Clean Architecture (Hexagonal)** principles.
It supports account creation, financial transactions (purchases, withdrawals, credit vouchers) and transaction history retrieval.
---
## 🏗️ Architecture
This project follows **Clean Architecture**, separating concerns into distinct layers where dependencies always point **inward**.
```
dependency flow:
delivery → usecase → domain ← repository
```
```
AccountFlow/
├── cmd/
│   └── api/
│       └── main.go                            # Application entry point
├── internal/
│   ├── domain/
│   │   ├── account.go                         # Account entity
│   │   ├── operation_type.go                  # OperationType entity + sign logic
│   │   ├── transaction.go                     # Transaction entity
│   │   └── errors.go                          # Domain-level typed errors
│   ├── usecase/
│   │   ├── account_usecase.go                 # Create account, get by ID
│   │   ├── transaction_usecase.go             # Create transaction
│   │   └── mocks/
│   │       └── repositories.go                # Mock implementations for unit tests
│   ├── repository/
│   │   ├── account_repository.go              # AccountRepository interface
│   │   ├── operation_type_repository.go       # OperationTypeRepository interface
│   │   └── transaction_repository.go          # TransactionRepository interface
│   ├── delivery/
│   │   └── http/
│   │       ├── router.go                      # Gin route definitions
│   │       ├── account_handler.go             # Account HTTP handlers
│   │       ├── transaction_handler.go         # Transaction HTTP handlers
│   │       └── response.go                    # Error mapping to HTTP status codes
│   └── infra/
│       └── postgres/
│           ├── db.go                          # PostgreSQL connection
│           ├── account_repo.go                # AccountRepository implementation
│           ├── operation_type_repo.go         # OperationTypeRepository implementation
│           └── transaction_repo.go            # TransactionRepository implementation
├── migrations/
│   ├── 000001_create_operation_types.up.sql             # Creates operation_types + seeds data
│   ├── 000002_create_accounts.up.sql          # Creates accounts table
│   └── 000003_create_transactions.up.sql      # Creates transactions table
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── go.sum
```
### Layer Responsibilities
| Layer            | Responsibility                                                        |
|------------------|-----------------------------------------------------------------------|
| `domain`         | Entities, business rules, domain errors — zero external dependencies |
| `usecase`        | Orchestrates business logic by calling repository interfaces          |
| `repository`     | Interfaces for data access (contracts only)                           |
| `infra/postgres` | PostgreSQL implementations of repository interfaces                   |
| `delivery/http`  | Gin handlers: parse requests, call use cases, return responses        |
---
## 🗃️ Data Model
### operation_types (seeded, read-only)
| operation_type_id | description                |
|-------------------|----------------------------|
| 1                 | Normal Purchase            |
| 2                 | Purchase with Installments |
| 3                 | Withdrawal                 |
| 4                 | Credit Voucher             |
> Operations **1, 2, 3** always produce a **negative** amount.
> Operation **4** always produces a **positive** amount.
> The sign is enforced by the domain — regardless of what the caller sends.
### accounts
| column          | type        | notes       |
|-----------------|-------------|-------------|
| account_id      | UUID        | primary key |
| document_number | VARCHAR(50) | unique      |
| created_at      | TIMESTAMPTZ |             |
### transactions
| column            | type          | notes                           |
|-------------------|---------------|---------------------------------|
| transaction_id    | UUID          | primary key                     |
| account_id        | UUID          | FK → accounts                   |
| operation_type_id | INT           | FK → operation_types            |
| amount            | NUMERIC(15,2) | negative=debit, positive=credit |
| event_date        | TIMESTAMPTZ   |                                 |
---
## 🚀 Features
| Feature               | Description                                           |
|-----------------------|-------------------------------------------------------|
| Account Creation      | Create an account linked to a document number         |
| Account Lookup        | Retrieve an account by its ID                         |
| Transaction Creation  | Post a transaction using an operation type            |
| Auto Sign Enforcement | Amount sign is always derived from operation type     |
---
## 🛠️ Tech Stack
| Technology     | Usage                           |
|----------------|---------------------------------|
| Go 1.25+       | Main language                   |
| Gin            | HTTP framework                  |
| PostgreSQL 16  | Relational database             |
| Docker         | Containerization                |
| Docker Compose | Local environment orchestration |
---
## 🔌 API Endpoints
Base path: `/api/v1`
| Method | Endpoint               | Description              |
|--------|------------------------|--------------------------|
| POST   | `/accounts`            | Create a new account     |
| GET    | `/accounts/:accountId` | Get account by ID        |
| POST   | `/transactions`        | Create a new transaction |
---
## 📋 Request & Response Examples
### Create Account
```http
POST /api/v1/accounts
Content-Type: application/json
{
  "document_number": "12345678900"
}
```
```json
{
  "account_id": "a1b2c3d4-e5f6-...",
  "document_number": "12345678900",
  "created_at": "2026-01-01T10:00:00Z"
}
```
### Get Account
```http
GET /api/v1/accounts/a1b2c3d4-e5f6-...
```
```json
{
  "account_id": "a1b2c3d4-e5f6-...",
  "document_number": "12345678900",
  "created_at": "2026-01-01T10:00:00Z"
}
```
### Create Transaction
```http
POST /api/v1/transactions
Content-Type: application/json
{
  "account_id": "a1b2c3d4-e5f6-...",
  "operation_type_id": 4,
  "amount": 123.45
}
```
```json
{
  "transaction_id": "f7g8h9i0-...",
  "account_id": "a1b2c3d4-e5f6-...",
  "operation_type_id": 4,
  "amount": 123.45,
  "event_date": "2026-01-05T09:34:18Z"
}
```
> **Note:** For debit operations (type 1, 2, 3), the stored `amount` is always **negative**:
> ```json
> { "operation_type_id": 1, "amount": 50.0 }  →  stored as -50.0
> ```
### Error Response
```json
{ "error": "account not found" }
```
Validation error:
```json
{
  "errors": [
    "'document_number' is required"
  ]
}
```
---
## 🔴 HTTP Error Mapping
| Domain Error               | HTTP Status                |
|----------------------------|----------------------------|
| `ErrAccountNotFound`       | `404 Not Found`            |
| `ErrDocumentAlreadyUsed`   | `409 Conflict`             |
| `ErrOperationTypeNotFound` | `422 Unprocessable Entity` |
| `ErrInvalidAmount`         | `400 Bad Request`          |
| `ErrInvalidField`          | `400 Bad Request`          |
| Internal errors            | `500 Internal Server Error`|
---
## ⚙️ Environment Variables
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
## 🐳 Running with Docker
### Prerequisites
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
```bash
# Build and start all services (API + PostgreSQL)
docker-compose up --build
# Run in background
docker-compose up --build -d
# Stop services
docker-compose down
# Stop and remove volumes (resets database)
docker-compose down -v
```
Services started:
- **PostgreSQL** → `localhost:5432`
- **AccountFlow API** → `localhost:8072`
> Database migrations are applied automatically on startup.
---
## 💻 Running Locally
### Prerequisites
- Go 1.25+
- PostgreSQL running locally
```bash
# Clone the repository
git clone https://github.com/your-org/AccountFlow.git
cd AccountFlow
# Install dependencies
go mod tidy
# Set environment variables (PowerShell)
$env:DB_HOST="localhost"; $env:DB_PORT="5432"; $env:DB_USER="postgres"
$env:DB_PASSWORD="postgres"; $env:DB_NAME="accountflow"; $env:APP_PORT="8072"
# Run the application
go run ./cmd/api/main.go
```
---
## 🧪 Tests
Unit tests cover all layers from domain to use cases using manual mocks — **no real database required**.
```bash
# Run all unit tests
go test ./...
# Run with verbose output
go test -v ./...
# Run with coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```
### Test Results (18 tests)
#### Domain (`internal/domain`)
| Test | Description |
|------|-------------|
| `TestNewAccount_Success` | Creates account with valid document |
| `TestNewAccount_EmptyDocument` | Rejects empty document |
| `TestNewAccount_WhitespaceDocument` | Rejects whitespace-only document |
| `TestNewTransaction_DebitSign` | Amount is negative for debit ops (1,2,3) |
| `TestNewTransaction_CreditSign` | Amount is positive for credit op (4) |
| `TestNewTransaction_NegativeInputNormalized` | Sign always derived from op type |
| `TestNewTransaction_ZeroAmount` | Rejects zero amount |
| `TestOperationType_IsDebit` | Correctly identifies debit vs credit types |
#### Use Cases (`internal/usecase`)
| Test | Description |
|------|-------------|
| `TestCreateAccount_Success` | Creates account and returns it |
| `TestCreateAccount_EmptyDocument` | Rejects empty document |
| `TestCreateAccount_DocumentAlreadyUsed` | Returns conflict on duplicate document |
| `TestGetByID_Success` | Returns account by ID |
| `TestGetByID_NotFound` | Returns not found error |
| `TestCreateTransaction_Debit_Success` | Debit transaction stores negative amount |
| `TestCreateTransaction_Credit_Success` | Credit transaction stores positive amount |
| `TestCreateTransaction_AccountNotFound` | Fails when account does not exist |
| `TestCreateTransaction_InvalidOperationType` | Fails for unknown operation type |
| `TestCreateTransaction_ZeroAmount` | Rejects zero amount |
### Mock Strategy
```go
type MockAccountRepository struct {
    CreateFn   func(ctx context.Context, account *domain.Account) error
    FindByIDFn func(ctx context.Context, accountID uuid.UUID) (*domain.Account, error)
}
type MockOperationTypeRepository struct {
    FindByIDFn func(ctx context.Context, id int) (*domain.OperationType, error)
}
```
Each test injects only the functions it needs, keeping tests focused and explicit.
---
## 📐 Design Decisions
- **Clean Architecture**: outer layers depend on inner layers, never the reverse — `domain` has zero external dependencies
- **No User entity**: accounts are identified by `document_number` directly — simplicity over extra abstraction
- **Operation type seeding**: the 4 operation types are immutable reference data seeded in migrations, never mutated at runtime
- **Sign enforcement in domain**: `NewTransaction` always applies the correct sign from the operation type — the HTTP layer cannot bypass this rule
- **Domain errors**: typed sentinel errors in `domain/errors.go`, mapped to HTTP status codes only at the delivery layer
- **Repository pattern**: data access is abstracted behind interfaces for easy mocking and future swaps
- **No framework in domain/usecase**: zero dependency on Gin or any HTTP concern in business logic
- **UUID primary keys**: all entities use UUID v4
- **Auto migrations**: `golang-migrate` runs pending migrations on every startup
---
## 📄 License
This project is licensed under the MIT License.
