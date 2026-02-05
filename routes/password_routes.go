package routes

import (
	"database/sql"

	"github.com/go-chi/chi/v5"

	"company_mgmt_api/config"
	"company_mgmt_api/handlers"
	"company_mgmt_api/repositories"
	"company_mgmt_api/services"
)

func RegisterPasswordRoutes(r chi.Router, db *sql.DB, cfg *config.Config) {
	otpRepo := repositories.NewOTPRepository(db)
	emailService := services.NewEmailService(cfg.ResendAPIKey, cfg.EmailFrom)
	otpService := services.NewOTPService(otpRepo, emailService)

	handler := &handlers.PasswordHandler{
		OTP: otpService,
	}

	r.Route("/password", func(r chi.Router) {
		r.Post("/forgot", handler.Forgot)
		r.Post("/reset", handler.Reset)
	})
}
