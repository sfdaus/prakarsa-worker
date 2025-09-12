package worker

import (
	"context"
	"fmt"
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

		items, err := r.Repo.ClaimPending(ctx, r.Batch)
		if err != nil {
			log.Printf("[worker] claim error: %v", err)
			time.Sleep(500 * time.Millisecond)
			continue
		}
		if len(items) == 0 {
			time.Sleep(r.Idle)
			continue
		}
		log.Printf("[worker] claimed %d notifications", len(items))

		for _, it := range items {
			<-sem
			go func(row repository.Row) {
				defer func() {
					if rec := recover(); rec != nil {
						log.Printf("[worker] panic id=%s: %v", row.ID, rec)
						_ = r.Repo.MarkRetry(ctx, row.ID, fmt.Sprint(rec))
					}
					sem <- struct{}{}
				}()

				if err := r.H.Handle(ctx, row); err == nil {
					if e := r.Repo.MarkSent(ctx, row.ID); e != nil {
						log.Printf("[worker] mark sent err id=%s: %v", row.ID, e)
					}
				} else {
					log.Printf("[worker] send err id=%s: %v", row.ID, err)
					if e := r.Repo.MarkRetry(ctx, row.ID, err.Error()); e != nil {
						log.Printf("[worker] mark retry err id=%s: %v", row.ID, e)
					}
				}
			}(it)
		}
	}
}
