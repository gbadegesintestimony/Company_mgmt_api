package database

import (
	"company_mgmt_api/config"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func Connect(cfg *config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBSSLMODE)

	var db *sql.DB
	var err error

	// Retry logic — Postgres may not be ready immediately
	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", dsn)
		if err == nil {
			if err = db.Ping(); err == nil {
				break
			}
		}
		log.Printf("Database not ready, retrying in 2s... (%d/10): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	log.Println("Successfully connected to PostgreSQL database")

	// Run migrations
	if err := runMigrations(db, cfg.DBName); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %v", err)
	}

	return db, nil
}

func runMigrations(db *sql.DB, dbName string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://database/migrations",
		dbName,
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %v", err)
	}

	log.Println("Migrations applied successfully")
	return nil
}
