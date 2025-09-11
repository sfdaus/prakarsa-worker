package request

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

// GetReportReasonReferenceReq request body
type GetReportReasonReferenceReq struct {
	PerPage  int64  `query:"per_page"`
	Page     int64  `query:"page"`
	Label    string `query:"label"`
	IsActive *bool  `query:"is_active"`
}

func (request GetReportReasonReferenceReq) Validate() error {
	return validation.ValidateStruct(
		&request,
	)
}
