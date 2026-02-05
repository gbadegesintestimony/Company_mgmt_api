package main

import (
	"company_mgmt_api/config"
	"company_mgmt_api/database"
	"company_mgmt_api/routes"
	"log"
	"net/http"
)

func main() {
	// Entry point of the application
	cfg := config.LoadConfig()
	log.Printf("Starting application in %s mode on port %s", cfg.AppEnv, cfg.AppPort)

	db, err := database.Connect(cfg)

	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	defer db.Close()

	// Further initialization and server start logic goes here
	handler := routes.RegisterRoutes(cfg, db)

	server := &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: handler,
	}
	log.Println("Server started on port", cfg.AppPort)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on port %s: %v\n", cfg.AppPort, err)
	}
}
