package services

import (
	"company_mgmt_api/repositories"
	"company_mgmt_api/utils"
	"context"
	"database/sql"

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
	actorID string,
) error {

	hash, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	id := uuid.New().String()

	if err := s.Employees.Create(
		ctx,
		id,
		companyID,
		email,
		hash,
		firstName,
		lastName,
		"employee",
	); err != nil {
		return err
	}

	return s.Audit.Log(
		ctx,
		"create_employee",
		actorID,
		id,
		"Created new employee with ID "+id,
	)
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

	if err := s.Audit.Log(ctx, "EMPLOYEE_CREATED", actorID, id, ""); err != nil {
		return err
	}

	return tx.Commit()
}
