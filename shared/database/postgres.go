package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lavish-gambhir/dashbeam/internal/config"
	"github.com/lavish-gambhir/dashbeam/pkg/apperr"
)

func New(ctx context.Context, conf *config.AppConfig) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, conf.Database.Address())
	if err != nil {
		return nil, apperr.Wrap(err, apperr.DBConnectionFailed, "pgxpool.New")
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, apperr.Wrap(err, apperr.DBConnectionFailed, "db.Ping")
	}
	// createExtension(ctx, pool) -- if any
	return pool, nil
}
