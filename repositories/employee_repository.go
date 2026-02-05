package repositories

import (
	"company_mgmt_api/models"
	"context"
	"database/sql"
	"strconv"
)

type EmployeeRepository struct {
	DB *sql.DB
}

func NewEmployeeRepository(db *sql.DB) *EmployeeRepository {
	return &EmployeeRepository{DB: db}
}

func (r *EmployeeRepository) Create(
	ctx context.Context,
	id, companyID, email, passwordHash, firstname, lastname, role string,
) error {
	_, err := r.DB.ExecContext(ctx,
		`INSERT INTO users (id, company_id, email, password_hash, first_name, last_name, role, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`,
		id,
		companyID,
		email,
		passwordHash,
		firstname,
		lastname,
		role,
	)
	return err
}

func (r *EmployeeRepository) FindByID(
	ctx context.Context,
	id, companyID string,
) (*models.Employee, error) {
	row := r.DB.QueryRowContext(ctx,
		`SELECT id, company_id, email, first_name, last_name, is_active, role, created_at, updated_at
		FROM users
		WHERE id = $1 AND company_id = $2 AND deleted_at IS NULL
		`,
		id,
		companyID,
	)

	var emp models.Employee
	if err := row.Scan(
		&emp.ID,
		&emp.CompanyID,
		&emp.Email,
		&emp.FirstName,
		&emp.LastName,
		&emp.IsActive,
		&emp.Role,
		&emp.CreatedAt,
		&emp.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &emp, nil
}

func (r *EmployeeRepository) ListWithCount(
	ctx context.Context,
	companyID string,
	search string,
	isActive *bool,
	limit, offset int,
) ([]*models.Employee, int, error) {

	where := `WHERE company_id = $1 AND deleted_at IS NULL
	And (email ILIKE '%' || $2 || '%' 
	OR first_name ILIKE '%' || $2 || '%' 
	OR last_name ILIKE '%' || $2 || '%')
	`

	if isActive != nil {
		where += " AND is_active = $3 "

	}

	var total int
	countQuery := `SELECT COUNT(*) FROM users ` + where

	args := []any{companyID, search}
	argPos := 3

	if isActive != nil {
		where += " AND is_active = $" + strconv.Itoa(argPos)
		args = append(args, *isActive)
		argPos++
	}

	err := r.DB.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}
	query := `SELECT id, company_id, email, first_name, last_name, role, is_active, created_at, updated_at
	FROM users 
	` + where + ` ORDER BY created_at DESC 
	LIMIT $3 OFFSET $4
	`

	args = []any{companyID, search, limit, offset}
	if isActive != nil {
		args = append(args, *isActive)
	}
	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var out []*models.Employee
	for rows.Next() {
		emp := &models.Employee{}
		{
			if err := rows.Scan(
				&emp.ID,
				&emp.CompanyID,
				&emp.Email,
				&emp.FirstName,
				&emp.LastName,
				&emp.Role,
				&emp.IsActive,
			); err != nil {
				return nil, 0, err
			}
			out = append(out, emp)
		}

	}
	return out, total, nil
}

func (r *EmployeeRepository) UpdateProfile(
	ctx context.Context,
	id, companyID, firstname, lastname string,
) error {
	_, err := r.DB.ExecContext(ctx,
		`UPDATE users
		SET first_name = $1, last_name = $2, updated_at = NOW()
		WHERE id = $3 AND company_id = $4 AND deleted_at IS NULL
		`,
		firstname,
		lastname,
		id,
		companyID,
	)
	return err
}

func (r *EmployeeRepository) SetActive(
	ctx context.Context,
	id, companyID string,
	active bool,
) error {
	_, err := r.DB.ExecContext(ctx,
		`UPDATE users
		SET is_active = $1, updated_at = NOW()
		WHERE id = $2 AND company_id = $3 AND deleted_at IS NULL
		`,
		active,
		id,
		companyID,
	)
	return err
}

func (r *EmployeeRepository) SoftDelete(
	ctx context.Context,
	id, companyID string,
) error {
	_, err := r.DB.ExecContext(ctx,
		`UPDATE users SET deleted_at = NOW()
		WHERE id = $1 AND company_id = $2 AND deleted_at IS NULL
		`,
		id,
		companyID,
	)
	return err
}
