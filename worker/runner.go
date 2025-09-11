package worker

import (
	"context"
	"log"
	"time"

	"prakarsa-app/repository"
)

type Runner struct {
	Repo  *repository.OutboxRepo
	H     Handler
	Batch int
	Idle  time.Duration
}

func (r *Runner) Start(ctx context.Context) error {
	sem := make(chan struct{}, r.H.Concurrency())
	for i := 0; i < r.H.Concurrency(); i++ {
		sem <- struct{}{}
	}
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		rows, err := r.Repo.PullPending(ctx, r.Batch)
		if err != nil {
			log.Printf("pull error: %v", err)
			time.Sleep(500 * time.Millisecond)
			continue
		}
		if len(rows) == 0 {
			time.Sleep(r.Idle)
			continue
		}

		for _, row := range rows {
			<-sem
			go func(row repository.Row) {
				defer func() { sem <- struct{}{} }()
				if err := r.H.Handle(ctx, row); err == nil {
					_ = r.Repo.MarkSent(ctx, row.ID)
				} else {
					_ = r.Repo.MarkRetry(ctx, row.ID, err.Error())
				}
			}(row)
		}
	}
}
