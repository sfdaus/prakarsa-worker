package worker

import (
	"context"
	"encoding/json"
	"prakarsa-app/repository"
)

type Handler interface {
	Name() string
	Concurrency() int
	Handle(ctx context.Context, r repository.Row) error
}

func ParseHeaders(raw json.RawMessage) map[string]string {
	m := map[string]string{}
	_ = json.Unmarshal(raw, &m)
	return m
}
