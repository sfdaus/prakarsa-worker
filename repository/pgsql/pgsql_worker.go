package pgsql

import (
	"context"
	"database/sql"
	"fmt"
	"prakarsa-app/transport/request"
	"prakarsa-app/transport/response"
	"strings"
)

type pgsqlReferenceRepository struct {
	db *sql.DB
}

// NewPgsqlReferenceRepository will create new an todoRepository object representation of ReferenceRepository interface
func NewPgsqlReferenceRepository(db *sql.DB) *pgsqlReferenceRepository {
	return &pgsqlReferenceRepository{
		db: db,
	}
}

func (r *pgsqlReferenceRepository) GetReportReason(ctx context.Context, request *request.GetReportReasonReferenceReq) (res []response.GetReportReasonRes, meta response.MetaRes, err error) {
	// 1. Build WHERE clauses
	wheres := []string{}
	args := []interface{}{}
	idx := 1

	if request.Label != "" {
		wheres = append(wheres, fmt.Sprintf("label ILIKE $%d", idx))
		args = append(args, "%"+request.Label+"%")
		idx++
	}

	if request.IsActive != nil {
		wheres = append(wheres, fmt.Sprintf("is_active = $%d", idx))
		args = append(args, request.IsActive)
		idx++
	}

	whereSQL := ""
	if len(wheres) > 0 {
		whereSQL = "WHERE " + strings.Join(wheres, " AND ")
	}

	// --- 2. Hitung totalCount dulu (tanpa LIMIT/OFFSET) ---
	countQuery := fmt.Sprintf(
		"SELECT COUNT(*) FROM report_reasons %s",
		whereSQL,
	)

	if err = r.db.QueryRowContext(ctx, countQuery, args...).Scan(&meta.TotalData); err != nil {
		return nil, meta, err
	}

	// 2. Calculate LIMIT & OFFSET
	perPage := request.PerPage
	if perPage <= 0 {
		perPage = 10
	}
	page := request.Page
	if page <= 0 {
		page = 1
	}

	// total pages = ceil(total / perPage)
	meta.Page = page
	meta.PerPage = perPage
	meta.TotalPages = (meta.TotalData + perPage - 1) / perPage

	offset := (page - 1) * perPage

	// add LIMIT & OFFSET to args
	args = append(args, perPage, offset)
	limitPos, offsetPos := idx, idx+1

	// 3. Final query
	query := fmt.Sprintf(`
        SELECT
            id, label, description, is_active, created_at
        FROM report_reasons
        %s
        ORDER BY created_at DESC
        LIMIT $%d OFFSET $%d
    `, whereSQL, limitPos, offsetPos)

	// 4. Execute
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, meta, err
	}
	defer rows.Close()

	// 5. Scan results
	for rows.Next() {
		var item response.GetReportReasonRes

		if err := rows.Scan(
			&item.ID,
			&item.Label,
			&item.Description,
			&item.IsActive,
			&item.CreatedAt,
		); err != nil {
			return nil, meta, err
		}

		res = append(res, item)
	}
	if errRow := rows.Err(); errRow != nil {
		return nil, meta, errRow
	}

	return
}
