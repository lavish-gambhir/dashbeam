package postgres

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/lavish-gambhir/dashbeam/pkg/apperr"
	"github.com/lavish-gambhir/dashbeam/shared/config"
)

func Connect(ctx context.Context, conf *config.AppConfig) (*pgxpool.Pool, error) {
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

type DB struct {
	pool *pgxpool.Pool

	log *slog.Logger
}

func New(p *pgxpool.Pool, logger *slog.Logger) *DB {
	return &DB{
		pool: p,
		log:  logger,
	}
}

func (db DB) TransactionContext(ctx context.Context) (context.Context, error) {
	tx, err := db.Conn(ctx).Begin(ctx)
	if err != nil {
		return nil, err
	}
	return context.WithValue(ctx, txCtx{}, tx), nil
}

func (db DB) Commit(ctx context.Context) error {
	if tx, ok := ctx.Value(txCtx{}).(pgx.Tx); ok && tx != nil {
		return tx.Commit(ctx)
	}
	return errors.New("context has no transaction")
}

func (db DB) Rollback(ctx context.Context) error {
	if tx, ok := ctx.Value(txCtx{}).(pgx.Tx); ok && tx != nil {
		return tx.Rollback(ctx)
	}
	return errors.New("context has no transaction")
}

func (db DB) WithAcquire(ctx context.Context) (context.Context, error) {
	if _, ok := ctx.Value(connCtx{}).(*pgxpool.Conn); ok {
		return nil, errors.New("context already has a connecting acquired")
	}
	res, err := db.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	return context.WithValue(ctx, connCtx{}, res), nil
}

func (db DB) Release(ctx context.Context) {
	if res, ok := ctx.Value(connCtx{}).(*pgxpool.Conn); ok && res != nil {
		res.Release()
	}
}

func (db DB) Conn(ctx context.Context) PGXQuerier {
	if tx, ok := ctx.Value(txCtx{}).(pgx.Tx); ok && tx != nil {
		return tx
	}
	if res, ok := ctx.Value(connCtx{}).(*pgxpool.Conn); ok && res != nil {
		return res
	}
	return db.pool
}

// GetTx retrieves the transaction from the context
func (db DB) Tx(ctx context.Context) (pgx.Tx, error) {
	tx, ok := ctx.Value(txCtx{}).(pgx.Tx)
	if !ok {
		return nil, errors.New("no transaction found in context")
	}
	return tx, nil
}

// CommitOrRollback commits the transaction if no error occurred or rolls back on error
func (db DB) CommitOrRollback(ctx context.Context, err *error) error {
	if *err != nil {
		rollbackErr := db.Rollback(ctx)
		if rollbackErr != nil {
			return rollbackErr
		}
	} else {
		commitErr := db.Commit(ctx)
		if commitErr != nil {
			return commitErr
		}
	}
	return nil
}

type txCtx struct{}
type connCtx struct{}
