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
}

type OutboxRepo struct{ DB *sql.DB }

func NewOutboxRepo(db *sql.DB) *OutboxRepo { return &OutboxRepo{DB: db} }

func (r *OutboxRepo) PullPending(ctx context.Context, limit int) ([]Row, error) {
	tx, err := r.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	rows, err := tx.QueryContext(ctx, `
    SELECT id, user_id, type, reference_type, reference_id, title, message, headers_json, attempt_count
    FROM notification_outbox
    WHERE status='pending' AND next_attempt_at <= NOW()
    ORDER BY priority ASC, created_at ASC
    LIMIT $1
    FOR UPDATE SKIP LOCKED
  `, limit)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	res := make([]Row, 0, limit)
	for rows.Next() {
		var r Row
		if err = rows.Scan(&r.ID, &r.UserID, &r.Type, &r.RefType, &r.RefID, &r.Title, &r.Message, &r.HeadersJSON, &r.AttemptCount); err != nil {
			return nil, err
		}
		res = append(res, r)
	}
	return res, rows.Err()
}

func (r *OutboxRepo) MarkSent(ctx context.Context, id string) error {
	_, err := r.DB.ExecContext(ctx, `
    UPDATE notification_outbox
    SET status='sent', sent_at=(EXTRACT(EPOCH FROM NOW()))::bigint, error_message=NULL
    WHERE id=$1 AND status IN ('pending','processing')
  `, id)
	return err
}

func (r *OutboxRepo) MarkRetry(ctx context.Context, id string, errMsg string) error {
	_, err := r.DB.ExecContext(ctx, `
    UPDATE notification_outbox
    SET attempt_count = attempt_count + 1,
        next_attempt_at = NOW() + make_interval(mins => LEAST(POWER(2, attempt_count+1), 60)),
        status='pending',
        error_message = left($2, 500)
    WHERE id=$1
  `, id, errMsg)
	return err
}
