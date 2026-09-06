package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool interface {
	Ping(context.Context) error
	Close()
}

func Open(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, err
	}
	config.MaxConns = 8
	config.MinConns = 0
	return pgxpool.NewWithConfig(ctx, config)
}
