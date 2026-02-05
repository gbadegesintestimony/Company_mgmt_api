package routes

import (
	"company_mgmt_api/handlers"
	"company_mgmt_api/middlewares"
	"company_mgmt_api/repositories"
	"company_mgmt_api/services"

	"github.com/go-chi/chi/v5"
)

func EmployeeRoutes(
	r chi.Router,
	empRepo *repositories.EmployeeRepository,
	auditRepo *repositories.AuditRepository,
	sessionRepo *repositories.SessionRepository,
) {
	service := services.NewEmployeeService(empRepo, auditRepo)

	admin := &handlers.AdminEmployeeHandler{Service: service}
	me := &handlers.MeHandler{
		Employees: empRepo,
		Sessions:  sessionRepo,
	}

	r.Route("/employees", func(r chi.Router) {
		r.Use(middlewares.RequireAdminRole)

		r.Post("/", admin.Create)
		r.Get("/", admin.List)
		r.Patch("/{id}/deactivate", admin.Deactivate)
	})

	r.Route("/me", func(r chi.Router) {
		r.Put("/profile", me.Profile)
		r.Put("/password", me.ChangePassword)
	})
}
