package services

import (
	"company_mgmt_api/models"
	"company_mgmt_api/repositories"
	"company_mgmt_api/utils"
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type EmployeeService struct {
	// Implementation goes here
	Employees *repositories.EmployeeRepository
	Audit     *repositories.AuditRepository
}

func NewEmployeeService(
	employeeRepo *repositories.EmployeeRepository,
	auditRepo *repositories.AuditRepository,
) *EmployeeService {
	return &EmployeeService{
		Employees: employeeRepo,
		Audit:     auditRepo,
	}
}

func (s *EmployeeService) CreateEmployees(
	ctx context.Context,
	companyID string,
	password string,
	email string,
	firstName string,
	lastName string,
	role string,
	actorID string,
) (*models.Employee, error) {

	hash, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	id := uuid.New().String()
	now := time.Now()

	if err := s.Employees.Create(ctx, id, companyID, email, hash, firstName, lastName, role); err != nil {
		return nil, err
	}

	_ = s.Audit.Log(ctx, "employee_created", actorID, id, map[string]string{
		"email":      email,
		"first_name": firstName,
		"last_name":  lastName,
		"role":       role,
	})
	return &models.Employee{
		ID:        id,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Role:      role,
		CreatedAt: now,
	}, nil
}

func (s *EmployeeService) CreateEmployeeTx(
	ctx context.Context,
	db *sql.DB,
	companyID, email, password, firstName, lastName, actorID string,
) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	hash, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	id := uuid.New().String()

	if err := s.Employees.Create(ctx, id, companyID, email, hash, firstName, lastName, "employee"); err != nil {
		return err
	}

	if err := s.Audit.Log(ctx, "employee_created", actorID, id, map[string]string{
		"email": email,
	}); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *EmployeeService) Get(ctx context.Context, userID, companyID string) (*models.Employee, error) {
	return s.Employees.FindByIDAndCompany(ctx, userID, companyID)
}

func (s *EmployeeService) Update(ctx context.Context, userID, companyID, firstName, lastName, role string) error {
	p := &models.Profile{FirstName: firstName, LastName: lastName}
	return s.Employees.UpdateProfile(ctx, userID, companyID, p)
}

func (s *EmployeeService) Reactivate(ctx context.Context, userID, companyID string) error {
	return s.Employees.SetActive(ctx, userID, companyID, true)
}

func (s *EmployeeService) Delete(ctx context.Context, userID, companyID string) error {
	return s.Employees.SoftDelete(ctx, userID, companyID)
}
