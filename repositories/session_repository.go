package repositories

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	"company_mgmt_api/models"

	"github.com/google/uuid"
)

type SessionRepository struct {
	DB *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{DB: db}
}

/*
CreateSession stores a hashed refresh token.
NEVER store raw refresh tokens.
*/
func (r *SessionRepository) CreateSession(
	ctx context.Context,
	userID string,
	companyID string,
	role string,
	refreshToken string,
	expiry time.Time,
) error {

	hash := hashToken(refreshToken)

	_, err := r.DB.ExecContext(ctx, `
		INSERT INTO sessions 
		(id, user_id, company_id, role, refresh_token_hash, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`,
		uuid.New().String(),
		userID,
		companyID,
		role,
		hash,
		expiry,
	)

	return err
}

/*
GetByRefreshToken validates:
- token exists
- not revoked
- not expired
*/
func (r *SessionRepository) GetByRefreshToken(
	ctx context.Context,
	refreshToken string,
) (*models.Session, error) {

	hash := hashToken(refreshToken)

	row := r.DB.QueryRowContext(ctx, `
		SELECT id, user_id, company_id, role, refresh_token_hash, expires_at, revoked_at
		FROM sessions
		WHERE refresh_token_hash = $1
	`, hash)

	var s models.Session

	err := row.Scan(
		&s.ID,
		&s.UserID,
		&s.CompanyID,
		&s.Role,
		&s.RefreshTokenHash,
		&s.ExpiresAt,
		&s.RevokedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("invalid session")
		}
		return nil, err
	}

	if s.RevokedAt != nil {
		return nil, errors.New("session revoked")
	}

	if time.Now().After(s.ExpiresAt) {
		return nil, errors.New("session expired")
	}

	return &s, nil
}

/*
Logout single device
*/
func (r *SessionRepository) RevokeByRefreshToken(
	ctx context.Context,
	refreshToken string,
) error {

	hash := hashToken(refreshToken)

	_, err := r.DB.ExecContext(ctx, `
		UPDATE sessions
		SET revoked_at = NOW()
		WHERE refresh_token_hash = $1
		AND revoked_at IS NULL
	`, hash)

	return err
}

/*
Logout all devices
*/
func (r *SessionRepository) RevokeAll(
	ctx context.Context,
	userID string,
) error {

	_, err := r.DB.ExecContext(ctx, `
		UPDATE sessions
		SET revoked_at = NOW()
		WHERE user_id = $1
		AND revoked_at IS NULL
	`, userID)

	return err
}

/*
Simple SHA256 hash helper (internal)
*/
func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
