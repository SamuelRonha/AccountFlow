#!/usr/bin/env pwsh
# run.ps1 — AccountFlow development helper
#
# Usage:
#   .\run.ps1            → start API + PostgreSQL via Docker Compose
#   .\run.ps1 local      → run the API locally (requires PostgreSQL running)
#   .\run.ps1 test       → run all unit tests
#   .\run.ps1 test:v     → run tests with verbose output
#   .\run.ps1 test:cover → run tests and open HTML coverage report
#   .\run.ps1 down       → stop and remove containers
#   .\run.ps1 reset      → stop containers and wipe the database volume

$command = if ($args.Count -gt 0) { $args[0] } else { "docker" }

switch ($command) {

    "docker" {
        Write-Host "▶  Starting AccountFlow (Docker Compose)..." -ForegroundColor Cyan
        docker compose up --build
    }

    "local" {
        Write-Host "▶  Starting AccountFlow locally..." -ForegroundColor Cyan
        $env:DB_HOST     = if ($env:DB_HOST)     { $env:DB_HOST }     else { "localhost" }
        $env:DB_PORT     = if ($env:DB_PORT)     { $env:DB_PORT }     else { "5432" }
        $env:DB_USER     = if ($env:DB_USER)     { $env:DB_USER }     else { "postgres" }
        $env:DB_PASSWORD = if ($env:DB_PASSWORD) { $env:DB_PASSWORD } else { "postgres" }
        $env:DB_NAME     = if ($env:DB_NAME)     { $env:DB_NAME }     else { "accountflow" }
        $env:DB_SSLMODE  = if ($env:DB_SSLMODE)  { $env:DB_SSLMODE }  else { "disable" }
        $env:APP_PORT    = if ($env:APP_PORT)    { $env:APP_PORT }    else { "8072" }
        go run ./cmd/api/main.go
    }

    "test" {
        Write-Host "▶  Running unit tests..." -ForegroundColor Cyan
        go test ./...
    }

    "test:v" {
        Write-Host "▶  Running unit tests (verbose)..." -ForegroundColor Cyan
        go test -v ./...
    }

    "test:cover" {
        Write-Host "▶  Running tests with coverage..." -ForegroundColor Cyan
        go test -coverprofile=coverage.out ./...
        go tool cover -html=coverage.out
    }

    "down" {
        Write-Host "▶  Stopping containers..." -ForegroundColor Yellow
        docker compose down
    }

    "reset" {
        Write-Host "▶  Stopping containers and wiping database..." -ForegroundColor Red
        docker compose down -v
    }

    default {
        Write-Host @"
AccountFlow — run script

Usage:
  .\run.ps1            start API + PostgreSQL via Docker Compose
  .\run.ps1 local      run the API locally (needs PostgreSQL running)
  .\run.ps1 test       run all unit tests
  .\run.ps1 test:v     run tests with verbose output
  .\run.ps1 test:cover run tests and open HTML coverage report
  .\run.ps1 down       stop containers
  .\run.ps1 reset      stop containers and wipe database volume
"@
    }
}

