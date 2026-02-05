package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv  string
	AppPort string

	DBHost    string
	DBPort    string
	DBUser    string
	DBPass    string
	DBName    string
	DBSSLMODE string

	JWTAccessSecret  string
	JWTRefreshSecret string

	EmailFrom    string
	ResendAPIKey string
}

func LoadConfig() *Config {
	// Implementation for loading configuration settings
	_ = godotenv.Load()

	cfg := &Config{
		AppEnv:  os.Getenv("APP_ENV"),
		AppPort: os.Getenv("APP_PORT"),

		DBHost:           os.Getenv("DB_HOST"),
		DBPort:           os.Getenv("DB_PORT"),
		DBUser:           os.Getenv("DB_USER"),
		DBPass:           os.Getenv("DB_PASSWORD"),
		DBName:           os.Getenv("DB_NAME"),
		DBSSLMODE:        os.Getenv("DB_SSLMODE"),
		JWTAccessSecret:  os.Getenv("JWT_ACCESS_SECRET"),
		JWTRefreshSecret: os.Getenv("JWT_REFRESH_SECRET"),

		EmailFrom:    os.Getenv("EMAIL_FROM"),
		ResendAPIKey: os.Getenv("RESEND_API_KEY"),
	}
	validate(cfg)
	return cfg
}

func validate(cfg *Config) {
	// Implementation for validating configuration
	if cfg.AppEnv == "" {
		log.Fatal("APP_ENV is not set")
	}
	if cfg.AppPort == "" {
		log.Fatal("APP_PORT is not set")
	}

	if cfg.DBHost == "" || cfg.DBPort == "" || cfg.DBUser == "" || cfg.DBPass == "" || cfg.DBName == "" || cfg.DBSSLMODE == "" {
		log.Fatal("Database configuration is incomplete")
	}
	if cfg.JWTAccessSecret == "" || cfg.JWTRefreshSecret == "" {
		log.Fatal("JWT secrets must be provided")
	}
	if cfg.EmailFrom == "" || cfg.ResendAPIKey == "" {
		log.Fatal("Email configuration is incomplete")
	}
	if os.Getenv("APP_ENV") == "production" {
		cfg.DBSSLMODE = "require"
	}

}
