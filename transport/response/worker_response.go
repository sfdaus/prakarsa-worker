package response

// Get List Response
type GetReportReasonRes struct {
	ID          string `json:"id"`
	Label       string `json:"label"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   int64  `json:"created_at"`
}
