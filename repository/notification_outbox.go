package repository

import (
	"context"
	"database/sql"
	"encoding/json"
)

type Row struct {
	ID           string
	UserID       string
	Type         string
	RefType      string
	RefID        string
	Title        string
	Message      string
	HeadersJSON  json.RawMessage
	AttemptCount int
	Priority     string
	ActionURL    sql.NullString
}

type OutboxRepo struct{ DB *sql.DB }

func NewOutboxRepo(db *sql.DB) *OutboxRepo { return &OutboxRepo{DB: db} }

// ClaimPending: atomically klaim N baris (status: pending -> processing) lalu return datanya
func (r *OutboxRepo) ClaimPending(ctx context.Context, limit int) ([]Row, error) {
	rows, err := r.DB.QueryContext(ctx, `
    WITH pick AS (
      SELECT id
      FROM notification_outbox
      WHERE status = 'pending' AND next_attempt_at <= NOW()
      ORDER BY priority ASC, created_at ASC
      LIMIT $1
      FOR UPDATE SKIP LOCKED
    )
    UPDATE notification_outbox n
    SET status = 'processing'
    FROM pick
    WHERE n.id = pick.id
    RETURNING n.id, n.user_id, n.type, n.reference_type, n.reference_id,
              n.title, n.message, n.headers_json, n.attempt_count, n.priority, n.action_url
  `, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Row
	for rows.Next() {
		var rRow Row
		if err := rows.Scan(&rRow.ID, &rRow.UserID, &rRow.Type, &rRow.RefType, &rRow.RefID,
			&rRow.Title, &rRow.Message, &rRow.HeadersJSON, &rRow.AttemptCount, &rRow.Priority, &rRow.ActionURL); err != nil {
			return nil, err
		}
		out = append(out, rRow)
	}
	return out, rows.Err()
}

func (r *OutboxRepo) MarkSent(ctx context.Context, id string) error {
	_, err := r.DB.ExecContext(ctx, `
    UPDATE notification_outbox
    SET status='sent', sent_at=(EXTRACT(EPOCH FROM NOW()))::bigint, error_message=NULL
    WHERE id=$1 AND status='processing'
  `, id)
	return err
}

func (r *OutboxRepo) MarkRetry(ctx context.Context, id string, errMsg string) error {
	_, err := r.DB.ExecContext(ctx, `
				UPDATE notification_outbox
				SET attempt_count = attempt_count + 1,
					-- backoff: 1m,2m,4m.. capped 60m
					next_attempt_at = NOW() + (LEAST(POWER(2, (attempt_count+1))::int, 60)) * INTERVAL '1 minute',
					status='pending',
					error_message = left($2, 500)
				WHERE id=$1
				`, id, errMsg)
	return err
}
