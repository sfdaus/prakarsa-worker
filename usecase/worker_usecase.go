package usecase

import (
	"context"
	"prakarsa-app/transport/response"
	"time"

	"prakarsa-app/domain"
	"prakarsa-app/repository/redis"
	"prakarsa-app/transport/request"
)

type ReferenceUsecase struct {
	referenceRepo domain.WorkerRepository
	redisRepo     redis.RedisRepository
	ctxTimeout    time.Duration
}

// NewReferenceUsecase will create new an notificationUsecase object representation of ThreadUsecase interface
func NewReferenceUsecase(referenceRepo domain.WorkerRepository, redisRepo redis.RedisRepository, ctxTimeout time.Duration) *ReferenceUsecase {
	return &ReferenceUsecase{
		referenceRepo: referenceRepo,
		redisRepo:     redisRepo,
		ctxTimeout:    ctxTimeout,
	}
}

func (u *ReferenceUsecase) GetReportReason(c context.Context, request *request.GetReportReasonReferenceReq) (res []response.GetReportReasonRes, meta response.MetaRes, err error) {
	ctx, cancel := context.WithTimeout(c, u.ctxTimeout)
	defer cancel()

	res, meta, err = u.referenceRepo.GetReportReason(ctx, request)
	return
}
