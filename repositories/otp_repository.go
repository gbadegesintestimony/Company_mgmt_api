package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
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
	userID, code, purpose string,
) error {
	// Implementation for creating an OTP
	_, err := r.DB.ExecContext(ctx,
		`INSERT INTO otps (id, user_id, code, purpose, expires_at) 
		VALUES ($1, $2, $3, $4, $5)
		`,

		uuid.New(),
		userID,
		code,
		purpose,
		time.Now().Add(10*time.Minute),
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
	err := r.DB.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM otps 
		WHERE user_id = $1 AND created_at >= NOW() - INTERVAL $2 * INTERVAL '1 minutes'
	`, userID, int(since.Minutes()),
	).Scan(&count)
	return count, err
}

func (r *OTPRepository) GetActive(
	ctx context.Context,
	userID, purpose string,
) (string, time.Time, error) {
	// Implementation for getting an active OTP
	var hash string
	var expiresAt time.Time

	err := r.DB.QueryRowContext(ctx,
		`SELECT code, expires_at FROM otps 
		WHERE user_id = $1 
		AND purpose = $2 
		AND consumed_at IS NULL
		ORDER BY created_at DESC
		LIMIT 1
	`, userID, purpose).Scan(&hash, &expiresAt)
	return hash, expiresAt, err
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
	userID, purpose, hash string,
) error {

	var expiresAt time.Time
	err := tx.QueryRowContext(ctx,
		`SELECT expires_at 
		FROM otps 
		WHERE user_id = $1
		AND purpose = $2
		AND code_hash = $3
		AND consumed_at IS NULL
		ORDER BY created_at DESC
		LIMIT 1
		FOR UPDATE
	`, userID, purpose, hash).Scan(&expiresAt)

	if err != nil {
		return err
	}

	if time.Now().After(expiresAt) {
		return errors.New("OTP expired")
	}

	_, err = tx.ExecContext(ctx, `
	UPDATE otps
	SET consumed_at = NOW()
	WHERE user_id = $ 1 AND purpose = $2 AND code_hash = $3
	`, userID, purpose, hash)

	return err
}
