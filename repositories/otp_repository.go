package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"company_mgmt_api/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type OTPRepository struct {
	// Define methods for OTP management
	DB *sql.DB
}

func NewOTPRepository(db *sql.DB) *OTPRepository {
	return &OTPRepository{DB: db}
}

func (r *OTPRepository) Create(
	ctx context.Context,
	o *models.OTP,
) error {
	// Ensure ID is set
	id := o.ID
	if id == "" {
		id = uuid.New().String()
	}

	_, err := r.DB.ExecContext(ctx,
		`INSERT INTO otps (id, user_id, purpose, code_hash, expires_at) 
		VALUES ($1, $2, $3, $4, $5)
		`,
		id,
		o.UserID,
		o.Purpose,
		o.CodeHash,
		o.ExpiresAt,
	)
	return err
}

func (r *OTPRepository) CountRecent(
	ctx context.Context,
	userID string,
	since time.Duration,
) (int, error) {
	// Implementation for counting recent OTPs
	var count int
	minutes := int(since.Minutes())
	query := fmt.Sprintf(`
        SELECT COUNT(*) FROM otps 
        WHERE user_id = $1 AND created_at >= NOW() - INTERVAL '%d minutes'
    `, minutes)
	err := r.DB.QueryRowContext(ctx, query,
		userID,
	).Scan(&count)
	return count, err
}

func (r *OTPRepository) GetActive(
	ctx context.Context,
	userID, purpose string,
) (*models.OTP, error) {
	var id string
	var hash string
	var expiresAt time.Time

	err := r.DB.QueryRowContext(ctx,
		`SELECT id, code_hash, expires_at FROM otps 
		WHERE user_id = $1 
		AND purpose = $2 
		AND consumed_at IS NULL
		ORDER BY created_at DESC
		LIMIT 1
	`, userID, purpose).Scan(&id, &hash, &expiresAt)
	if err != nil {
		return nil, err
	}
	return &models.OTP{ID: id, UserID: userID, CodeHash: hash, Purpose: purpose, ExpiresAt: expiresAt}, nil
}

func (r *OTPRepository) Consume(
	ctx context.Context,
	userID, purpose string,
) error {
	// Implementation for consuming an OTP
	_, err := r.DB.ExecContext(ctx, `
		UPDATE otps 
		SET consumed_at = NOW()
		WHERE user_id = $1 AND purpose = $2 AND consumed_at IS NULL
	`, userID, purpose)
	return err
}

func (r *OTPRepository) VerifyAndConsumeTx(
	ctx context.Context,
	tx *sql.Tx,
	userID, purpose, rawCode string,
) error {

	var id string
	var dbHash string
	var expiresAt time.Time

	err := tx.QueryRowContext(ctx,
		`SELECT id, code_hash, expires_at 
		FROM otps 
		WHERE user_id = $1
		AND purpose = $2
		AND consumed_at IS NULL
		ORDER BY created_at DESC
		LIMIT 1
		FOR UPDATE
	`, userID, purpose).Scan(&id, &dbHash, &expiresAt)

	if err != nil {
		return err
	}

	if time.Now().After(expiresAt) {
		return errors.New("OTP expired")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbHash), []byte(rawCode)); err != nil {
		return errors.New("invalid OTP")
	}

	_, err = tx.ExecContext(ctx, `
		UPDATE otps SET consumed_at = NOW() WHERE id = $1
	`, id)

	return err
}
