package routes

import (
	"company_mgmt_api/config"
	"company_mgmt_api/handlers"
	"company_mgmt_api/middlewares"
	"company_mgmt_api/repositories"
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func RegisterRoutes(cfg *config.Config, db *sql.DB) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middlewares.RequestID)

	sessionRepo := repositories.NewSessionRepository(db)
	employeeRepo := repositories.NewEmployeeRepository(db)
	auditRepo := repositories.NewAuditRepository(db)
	companyRepo := repositories.NewCompanyRepository(db)
	userRepo := repositories.NewUserRepository(db)

	companyHandler := &handlers.CompanyHandler{Service: nil}
	auditHandler := &handlers.AuditHandler{Repo: auditRepo}

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", handlers.Health)

		// Public routes
		AuthRoutes(r, cfg, sessionRepo, userRepo, companyRepo, db)
		RegisterVerificationRoutes(r, db, cfg)
		RegisterPasswordRoutes(r, db, cfg)

		// Authenticated routes
		r.Group(func(r chi.Router) {
			r.Use(middlewares.AuthMiddleware(cfg.JWTAccessSecret))

			// /v1/me
			r.Route("/me", func(r chi.Router) {
				MeRoutes(r, employeeRepo, sessionRepo, db)
			})

			// /v1/companies/{companyId}/...
			CompanyRoutes(r, companyHandler, auditHandler, employeeRepo, auditRepo, sessionRepo)
		})
	})

	return r
}
