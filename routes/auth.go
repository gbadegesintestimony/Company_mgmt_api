package routes

import (
	"company_mgmt_api/config"
	"company_mgmt_api/handlers"
	"company_mgmt_api/middlewares"
	"company_mgmt_api/repositories"
	"database/sql"
	"time"

	"company_mgmt_api/services"

	"github.com/go-chi/chi/v5"
)

func AuthRoutes(
	r chi.Router,
	cfg *config.Config,
	sessionRepo *repositories.SessionRepository,
	userRepo *repositories.UserRepository,
	companyRepo *repositories.CompanyRepository,
	db *sql.DB,
) {
	emailSvc := services.NewEmailService(cfg.ResendAPIKey, cfg.EmailFrom)
	otpRepo := repositories.NewOTPRepository(db)
	otpService := services.NewOTPService(otpRepo, emailSvc)
	// sessionRepo := repositories.NewSessionRepository(db)

	// Assuming session service needs no DB for now
	handler := &handlers.AuthHandler{
		Cfg:         cfg,
		Email:       emailSvc,
		Sessions:    sessionRepo,
		UserRepo:    userRepo,
		CompanyRepo: companyRepo,
		OTPService:  otpService,
	}

	r.Post("/auth/admin/register", handler.Register)

	r.Group(func(r chi.Router) {
		// 10 attempts/minute/IP — login has no other brute-force protection.
		r.Use(middlewares.LoginRateLimiter(10, time.Minute))
		r.Post("/auth/admin/login", handler.AdminLogin)
		r.Post("/auth/employee/login", handler.EmployeeLogin)
	})

	r.Post("/auth/refresh", handler.Refresh)
	r.Post("/auth/logout", handler.Logout)
}
