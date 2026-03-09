package repositories

import (
	"company_mgmt_api/models"
	"context"
	"database/sql"
	"strconv"
	"time"
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
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		INSERT INTO users (id, company_id, email, password_hash, role, is_active, created_at)
		VALUES ($1, $2, $3, $4, $5, true, $6)
	`, id, companyID, email, passwordHash, role, time.Now())
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO profiles (user_id, first_name, last_name, created_at)
		VALUES ($1, $2, $3, $4)
	`, id, firstname, lastname, time.Now())
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *EmployeeRepository) FindByIDAndCompany(
	ctx context.Context,
	id, companyID string,
) (*models.Employee, error) {
	row := r.DB.QueryRowContext(ctx,
		`SELECT u.id, u.company_id, u.email, p.first_name, p.last_name,u.password_hash, u.is_active, u.role, u.created_at, u.updated_at
		FROM users u
		LEFT JOIN profiles p ON p.user_id = u.id
		WHERE u.id = $1 AND u.company_id = $2 AND u.deleted_at IS NULL
		`,
		id,
		companyID,
	)

	var emp models.Employee
	var firstName, lastName *string
	if err := row.Scan(
		&emp.ID,
		&emp.CompanyID,
		&emp.Email,
		&firstName,
		&lastName,
		&emp.PasswordHash,
		&emp.IsActive,
		&emp.Role,
		&emp.CreatedAt,
		&emp.UpdatedAt,
	); err != nil {
		return nil, err
	}

	if firstName != nil {
		emp.FirstName = *firstName
	}

	if lastName != nil {
		emp.LastName = *lastName
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

	args := []any{companyID, search}

	where := `
		WHERE u.company_id = $1 
		AND u.deleted_at IS NULL
		AND u.role = 'employee'
		AND (
			u.email ILIKE '%' || $2 || '%'
			OR p.first_name ILIKE '%' || $2 || '%'
			OR p.last_name ILIKE '%' || $2 || '%'
		)
	`

	if isActive != nil {
		where += " AND u.is_active = $3"
		args = append(args, *isActive)
	}

	var total int
	countQuery := `SELECT COUNT(*) FROM users u LEFT JOIN profiles p ON p.user_id = u.id ` + where
	if err := r.DB.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Add limit and offset to args
	limitPos := len(args) + 1
	offsetPos := len(args) + 2
	args = append(args, limit, offset)

	query := `
		SELECT u.id, u.company_id, u.email, p.first_name, p.last_name, u.role, u.is_active, u.created_at, u.updated_at
		FROM users u
		LEFT JOIN profiles p ON p.user_id = u.id
		` + where + ` ORDER BY u.created_at DESC
		LIMIT $` + itoa(limitPos) + ` OFFSET $` + itoa(offsetPos)

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var out []*models.Employee
	for rows.Next() {
		emp := &models.Employee{}
		var firstName, lastName *string
		if err := rows.Scan(
			&emp.ID,
			&emp.CompanyID,
			&emp.Email,
			&firstName,
			&lastName,
			&emp.Role,
			&emp.IsActive,
			&emp.CreatedAt,
			&emp.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}

		if firstName != nil {
			emp.FirstName = *firstName
		}

		if lastName != nil {
			emp.LastName = *lastName
		}
		out = append(out, emp)
	}
	return out, total, nil
}

func (r *EmployeeRepository) UpdateProfile(
	ctx context.Context,
	id, companyID string,
	p *models.Profile,
) error {
	_, err := r.DB.ExecContext(ctx, `
		UPDATE profiles
		SET 
		first_name = CASE WHEN $1 != '' THEN $1 ELSE first_name END,
		last_name = CASE WHEN $2 != '' THEN $2 ELSE last_name END, 
		phone = CASE WHEN $3 != '' THEN $3 ELSE phone END,
		job_title = CASE WHEN $4 != '' THEN $4 ELSE job_title END,
		department = CASE WHEN $5 != '' THEN $5 ELSE department END
		WHERE user_id = $6
		AND EXISTS (
			SELECT 1 FROM users 
			WHERE id = $6 AND company_id = $7 AND deleted_at IS NULL
		)
	`, p.FirstName, p.LastName, p.Phone, p.JobTitle, p.Department, id, companyID)
	return err
}

func (r *EmployeeRepository) SetActive(
	ctx context.Context,
	id, companyID string,
	active bool,
) error {
	_, err := r.DB.ExecContext(ctx, `
		UPDATE users
		SET is_active = $1, updated_at = NOW()
		WHERE id = $2 AND company_id = $3 AND deleted_at IS NULL
	`, active, id, companyID)
	return err
}

func (r *EmployeeRepository) SoftDelete(
	ctx context.Context,
	id, companyID string,
) error {
	_, err := r.DB.ExecContext(ctx, `
		UPDATE users SET deleted_at = NOW()
		WHERE id = $1 AND company_id = $2 AND deleted_at IS NULL
	`, id, companyID)
	return err
}

func itoa(i int) string {
	return strconv.Itoa(i)
}

func (r *EmployeeRepository) UpdatePassword(ctx context.Context, userID, companyID, newHash string) error {
	_, err := r.DB.ExecContext(ctx, `
		UPDATE users SET password_hash = $1, updated_at = NOW()
		WHERE id = $2 AND company_id = $3 AND deleted_at IS NULL
	`, newHash, userID, companyID)
	return err
}
