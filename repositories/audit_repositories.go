package repositories

import (
	"company_mgmt_api/models"
	"context"
	"database/sql"
	"encoding/json"
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
	metadata map[string]string,
) error {
	metaJSON, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	_, err = r.DB.ExecContext(ctx,
		`INSERT INTO audit_logs (actor_user_id, action, target_id, metadata, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		`,
		actorID,
		action,
		targetID,
		metaJSON,
	)
	return err
}

func (r *AuditRepository) ListByCompany(
	ctx context.Context,
	companyID string,
) ([]*models.AuditLog, error) {
	rows, err := r.DB.QueryContext(ctx,
		`SELECT al.id, al.actor_user_id, al.action, al.target_id, COALESCE(al.metadata, '{}'), al.created_at
		FROM audit_logs al
		JOIN users u ON al.actor_user_id = u.id
		WHERE u.company_id = $1
		ORDER BY al.created_at DESC
		LIMIT 100
		`,
		companyID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*models.AuditLog
	for rows.Next() {
		log := &models.AuditLog{}
		if err := rows.Scan(
			&log.ID,
			&log.ActorID,
			&log.Action,
			&log.TargetID,
			&log.Metadata,
			&log.CreatedAt,
		); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, rows.Err()
}
