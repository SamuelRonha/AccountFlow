package main

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	migratepostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	deliveryhttp "AccountFlow/internal/delivery/http"
	"AccountFlow/internal/infra/postgres"
	"AccountFlow/internal/usecase"
)

func main() {
	db, err := postgres.NewConnection()
	if err != nil {
		log.Fatalf("connecting to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	driver, err := migratepostgres.WithInstance(db, &migratepostgres.Config{})
	if err != nil {
		log.Fatalf("creating migrate driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		log.Fatalf("creating migrate instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("running migrations: %v", err)
	}

	log.Println("migrations applied successfully")

	// Repositories
	accountRepo := postgres.NewAccountRepository(db)
	txRepo := postgres.NewTransactionRepository(db)
	opTypeRepo := postgres.NewOperationTypeRepository(db)

	// Use cases
	accountUC := usecase.NewAccountUseCase(accountRepo)
	transactionUC := usecase.NewTransactionUseCase(txRepo, accountRepo, opTypeRepo)

	// Handlers
	accountHandler := deliveryhttp.NewAccountHandler(accountUC)
	transactionHandler := deliveryhttp.NewTransactionHandler(transactionUC)

	// Router
	router := deliveryhttp.NewRouter(accountHandler, transactionHandler)

	port := getEnv("APP_PORT", "8072")
	log.Printf("server starting on port %s", port)
	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("starting server: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
