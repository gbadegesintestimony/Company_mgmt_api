package routes

import (
	"database/sql"

	"github.com/go-chi/chi/v5"

	"company_mgmt_api/config"
	"company_mgmt_api/handlers"
	"company_mgmt_api/middlewares"
	"company_mgmt_api/repositories"
	"company_mgmt_api/services"
)

func RegisterVerificationRoutes(r chi.Router, db *sql.DB, cfg *config.Config) {
	// 1. Build repository
	otpRepo := repositories.NewOTPRepository(db)

	// 2. Build email service (already created earlier)
	emailService := services.NewEmailService(cfg.ResendAPIKey, cfg.EmailFrom)
	// emailService := services.NewEmailService(cfg.ResendAPIKey, cfg.EmailFrom, cfg.DevMode) // for dev mode
	// 3. Build OTP service
	otpService := services.NewOTPService(otpRepo, emailService)

	// 4. Build handler
	handler := &handlers.VerificationHandler{
		OTP:         otpService,
		UserRepo:    repositories.NewUserRepository(db),
		CompanyRepo: repositories.NewCompanyRepository(db),
		Sessions:    repositories.NewSessionRepository(db),
		Cfg:         cfg,
	}

	r.Post("/verification/email/confirm", handler.Confirm)

	r.Group(func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware(cfg.JWTAccessSecret))

		r.Post("/verification/email/request", handler.Request)
	})
}
