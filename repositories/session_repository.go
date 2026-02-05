package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type SessionRepository struct {
	// Define methods for session management
	DB *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{DB: db}
}

func (r *SessionRepository) CreateSession(
	ctx context.Context,
	userID string,
	refreshHash string,
) error {
	// Implementation for creating a session
	_, err := r.DB.ExecContext(ctx,
		`INSERT INTO sessions (id, user_id, refresh_token_hash, expires_at) VALUES ($1, $2, $3, $4)`,
		uuid.New(),
		userID,
		refreshHash,
		time.Now().Add(30*24*time.Hour),
	)
	return err
}

func (r *SessionRepository) RevokeAll(
	ctx context.Context,
	userID string,
) error {
	// Implementation for revoking a sessionRepository
	_, err := r.DB.ExecContext(ctx, `
		UPDATE sessions
		SET revoked_at = NOW()
		WHERE user_id = $1 AND revoked_at IS NULL
	`, userID)
	return err
}
