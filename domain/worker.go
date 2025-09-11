package domain

import (
	"context"
	"prakarsa-app/transport/request"
	"prakarsa-app/transport/response"
)

// // WorkerRepository represent the worker repository contract
type WorkerRepository interface {
	GetReportReason(ctx context.Context, request *request.GetReportReasonReferenceReq) ([]response.GetReportReasonRes, response.MetaRes, error)
}

// WorkerUsecase represent the worker usecase contract
type WorkerUsecase interface {
	GetReportReason(ctx context.Context, request *request.GetReportReasonReferenceReq) ([]response.GetReportReasonRes, response.MetaRes, error)
}
