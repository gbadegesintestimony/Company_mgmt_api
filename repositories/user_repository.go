package repositories

import (
	"company_mgmt_api/models"
	"context"
	"database/sql"
	"time"
)

type UserRepository struct {
	// Methods for user data access would be defined here
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) CreateAdmin(
	ctx context.Context,
	id string,
	email string,
	passwordHash string,
	companyID string,
) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		INSERT INTO users (id, company_id, email, password_hash, role, is_active, created_at)
		VALUES ($1, $2, $3, $4, 'admin', false, $5)
	`, id, companyID, email, passwordHash, time.Now())
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO profiles (user_id, created_at)
		VALUES ($1, $2)
	`, id, time.Now())
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *UserRepository) FindByEmailAndCompany(
	ctx context.Context,
	email, companyID string,
) (*models.User, error) {
	row := r.DB.QueryRowContext(ctx,
		`SELECT 
		u.id, u.company_id, u.email, u.password_hash, 
		p.first_name, p.last_name, 
		u.is_active, u.role, u.created_at, u.updated_at, u.email_verified_at
		FROM users u
		LEFT JOIN profiles p ON p.user_id = u.id
		WHERE u.email = $1 AND u.company_id = $2 AND u.deleted_at IS NULL
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
		&user.FirstName,
		&user.LastName,
		&user.IsActive,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.EmailVerifiedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindIDByEmail(
	ctx context.Context,
	email string,
) (string, error) {
	row := r.DB.QueryRowContext(ctx,
		`SELECT id FROM users WHERE email = $1 AND deleted_at IS NULL`,
		email,
	)

	var userID string
	if err := row.Scan(&userID); err != nil {
		return "", err
	}
	return userID, nil
}

func (r *UserRepository) IsVerified(ctx context.Context, userID string) (bool, error) {
	var emailVerifiedAt *time.Time
	var isActive bool

	err := r.DB.QueryRowContext(ctx, `
		SELECT email_verified_at, is_active FROM users WHERE id = $1
	`, userID).Scan(&emailVerifiedAt, &isActive)
	if err != nil {
		return false, err
	}

	return emailVerifiedAt != nil && isActive, nil
}

func (r *UserRepository) DeleteByID(ctx context.Context, userID string) error {
	_, err := r.DB.ExecContext(ctx, `
		DELETE FROM users WHERE id = $1
	`, userID)
	return err
}

func (r *UserRepository) SetEmailVerified(ctx context.Context, userID string) error {
	_, err := r.DB.ExecContext(ctx, `
		UPDATE users SET email_verified_at = NOW(), is_active = true, updated_at = NOW()
		WHERE id = $1
	`, userID)
	return err
}

func (r *UserRepository) FindByID(ctx context.Context, userID string) (*models.User, error) {
	row := r.DB.QueryRowContext(ctx,
		`SELECT 
		u.id, u.company_id, u.email, u.password_hash,
		p.first_name, p.last_name, p.phone, p.job_title, p.department,
		u.is_active, u.role, u.created_at, u.updated_at, u.email_verified_at
		FROM users u
		LEFT JOIN profiles p ON p.user_id = u.id
		WHERE u.id = $1 AND u.deleted_at IS NULL
		`,
		userID,
	)

	var user models.User
	err := row.Scan(
		&user.ID,
		&user.CompanyID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&user.JobTitle,
		&user.Department,
		&user.IsActive,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.EmailVerifiedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
