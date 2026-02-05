package repositories

import (
	"company_mgmt_api/models"
	"context"
	"database/sql"
)

type UserRepository struct {
	// Methods for user data access would be defined here
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// Additional user-related methods would go here

func (r *UserRepository) FindByEmailAndCompany(
	ctx context.Context,
	email, companyID string,
) (*models.User, error) {
	row := r.DB.QueryRowContext(ctx,
		`SELECT id, company_id, email, password_hash, first_name, last_name, is_active, role, created_at, updated_at
		FROM users
		WHERE email = $1 AND company_id = $2 AND deleted_at IS NULL
		`,
		email,
		companyID,
	)

	var user models.User
	err := row.Scan(
		&user.ID,
		&user.CompanyID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
