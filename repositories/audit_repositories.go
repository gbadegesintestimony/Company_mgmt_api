package repositories

import (
	"context"
	"database/sql"
)

type AuditRepository struct {
	DB *sql.DB
}

func NewAuditRepository(db *sql.DB) *AuditRepository {
	return &AuditRepository{DB: db}
}

func (r *AuditRepository) Log(
	ctx context.Context,
	action string,
	actorID string,
	targetID string,
	metadata string,
) error {
	_, err := r.DB.ExecContext(ctx,
		`INSERT INTO audit_logs (action, actor_id, target_id, metadata, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		`,
		action,
		actorID,
		targetID,
		metadata,
	)
	return err
}
