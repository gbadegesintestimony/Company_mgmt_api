package routes

import (
	"company_mgmt_api/handlers"
	"company_mgmt_api/middlewares"
	"company_mgmt_api/repositories"

	"github.com/go-chi/chi/v5"
)

func CompanyRoutes(r chi.Router, handler *handlers.CompanyHandler, audit *handlers.AuditHandler, empRepo *repositories.EmployeeRepository,
	auditRepo *repositories.AuditRepository,
	sessionRepo *repositories.SessionRepository) {

	r.Route("/companies/{companyId}", func(r chi.Router) {
		r.Use(middlewares.RequireOwnCompany)

		r.Get("/", handler.Get)
		r.Patch("/", handler.Update)

		r.Get("/audit-logs", audit.List)

		r.Route("/employees", func(r chi.Router) {
			r.Use(middlewares.RequireAdminRole)
			EmployeeRoutes(r, empRepo, auditRepo, sessionRepo)
		})
	})
}
