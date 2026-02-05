package database

import (
	"company_mgmt_api/config"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func Connect(cfg *config.Config) (*sql.DB, error) {
	// Implementation for connecting to PostgreSQL database

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBSSLMODE)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	//  Verify the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	log.Println("Successfully connected to PostgreSQL database")
	return db, nil
}
