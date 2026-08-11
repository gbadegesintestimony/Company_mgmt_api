package routes

import (
	"company_mgmt_api/handlers"
	"company_mgmt_api/repositories"
	"company_mgmt_api/services"
	"database/sql"

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

	r.Post("/", admin.Create)
	r.Get("/", admin.List)
	r.Get("/{userId}", admin.Get)
	r.Patch("/{userId}", admin.Update)
	r.Patch("/{userId}/deactivate", admin.Deactivate)
	r.Post("/{userId}/reactivate", admin.Reactivate)
	r.Delete("/{userId}", admin.Delete)
}

func MeRoutes(
	r chi.Router,
	empRepo *repositories.EmployeeRepository,
	sessionRepo *repositories.SessionRepository,
	db *sql.DB,

) {
	service := services.NewEmployeeService(empRepo, nil)

	me := &handlers.MeHandler{
		Employees: empRepo,
		Sessions:  sessionRepo,
		Services:  service,
		UserRepo:  repositories.NewUserRepository(db),
	}
	r.Get("/", me.Get)
	r.Patch("/profile", me.Profile)
	r.Patch("/password", me.ChangePassword)

}
