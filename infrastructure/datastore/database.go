package datastore

import (
	"context"
	"database/sql"
	"net/url"
	"time"

	_ "github.com/lib/pq"
)

func MustOpen(databaseURL string) (db *sql.DB, err error) {
	parseDBUrl, _ := url.Parse(databaseURL)
	db, err = sql.Open(parseDBUrl.Scheme, databaseURL)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}
	return
}
