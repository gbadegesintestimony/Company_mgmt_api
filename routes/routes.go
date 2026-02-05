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

	r.Route("/v1", func(r chi.Router) {

		r.Get("/health", handlers.Health)

		AuthRoutes(r, cfg, sessionRepo)

		RegisterVerificationRoutes(r, db, cfg)
		RegisterPasswordRoutes(r, db, cfg)

		r.Group(func(r chi.Router) {
			r.Use(middlewares.AuthMiddleware(cfg.JWTAccessSecret))

			EmployeeRoutes(
				r,
				employeeRepo,
				auditRepo,
				sessionRepo,
			)
		})
	})

	return r
}
