package routes

import (
	"company_mgmt_api/config"
	"company_mgmt_api/handlers"
	"company_mgmt_api/repositories"

	"company_mgmt_api/services"

	"github.com/go-chi/chi/v5"
)

func AuthRoutes(
	r chi.Router,
	cfg *config.Config,
	sessionRepo *repositories.SessionRepository,
) {
	// Authentication routes would be defined here
	emailSvc := services.NewEmailService(cfg.ResendAPIKey, cfg.EmailFrom)
	// sessionRepo := repositories.NewSessionRepository(db)

	// Assuming session service needs no DB for now
	handler := &handlers.AuthHandler{
		Cfg:      cfg,
		Email:    emailSvc,
		Sessions: sessionRepo,
	}

	r.Post("/auth/admin/register", handler.Register)
	r.Post("/auth/login", handler.Login)
}
