package clickhouse

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/lavish-gambhir/dashbeam/pkg/apperr"
	"github.com/lavish-gambhir/dashbeam/shared/config"
)

type DB struct {
	conn   clickhouse.Conn
	logger *slog.Logger
}

func New(cfg config.AnalyticsConfig, logger *slog.Logger) (*DB, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{cfg.ClickHouseURL},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		Auth: clickhouse.Auth{
			Database: "dashbeam_analytics",
			Username: "default",
			Password: "password",
		},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, apperr.Wrap(err, apperr.Internal, "failed to connect to ClickHouse")
	}

	db := &DB{
		conn:   conn,
		logger: logger.With("component", "clickhouse"),
	}

	// Create tables on initialization
	if err := db.createTables(context.Background()); err != nil {
		return nil, apperr.Wrap(err, apperr.Internal, "failed to create ClickHouse tables")
	}

	return db, nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) Ping(ctx context.Context) error {
	return db.conn.Ping(ctx)
}

func (db *DB) Conn() clickhouse.Conn {
	return db.conn
}

func (db *DB) Query(ctx context.Context, query string, args ...interface{}) (driver.Rows, error) {
	return db.conn.Query(ctx, query, args...)
}

func (db *DB) Exec(ctx context.Context, query string, args ...interface{}) error {
	return db.conn.Exec(ctx, query, args...)
}

func (db *DB) PrepareBatch(ctx context.Context, query string) (driver.Batch, error) {
	return db.conn.PrepareBatch(ctx, query)
}

func (db *DB) createTables(ctx context.Context) error {
	queries := []string{
		// Events table
		`CREATE TABLE IF NOT EXISTS events (
			event_id String,
			event_type String,
			category String,
			action String,
			user_id UUID,
			school_id UUID,
			session_id Nullable(UUID),
			quiz_id Nullable(UUID),
			question_id Nullable(UUID),
			value Nullable(Float64),
			metadata String,
			timestamp DateTime64(3),
			processed_at DateTime64(3)
		) ENGINE = MergeTree()
		PARTITION BY toYYYYMM(timestamp)
		ORDER BY (school_id, user_id, timestamp)`,

		// User activity metrics table
		`CREATE TABLE IF NOT EXISTS user_activity_metrics (
			user_id UUID,
			school_id UUID,
			date Date,
			login_count UInt32,
			session_time_minutes UInt32,
			quiz_count UInt32,
			updated_at DateTime64(3)
		) ENGINE = ReplacingMergeTree(updated_at)
		PARTITION BY toYYYYMM(date)
		ORDER BY (school_id, user_id, date)`,

		// Quiz metrics table
		`CREATE TABLE IF NOT EXISTS quiz_metrics (
			quiz_id UUID,
			school_id UUID,
			date Date,
			participant_count UInt32,
			completion_rate Float64,
			average_score Float64,
			average_response_time_ms UInt32,
			updated_at DateTime64(3)
		) ENGINE = ReplacingMergeTree(updated_at)
		PARTITION BY toYYYYMM(date)
		ORDER BY (school_id, quiz_id, date)`,

		// School metrics table
		`CREATE TABLE IF NOT EXISTS school_metrics (
			school_id UUID,
			date Date,
			active_users UInt32,
			total_quizzes UInt32,
			total_events UInt32,
			updated_at DateTime64(3)
		) ENGINE = ReplacingMergeTree(updated_at)
		PARTITION BY toYYYYMM(date)
		ORDER BY (school_id, date)`,
	}

	for _, query := range queries {
		if err := db.conn.Exec(ctx, query); err != nil {
			return fmt.Errorf("failed to execute query: %w", err)
		}
	}

	db.logger.Info("ClickHouse tables created successfully")
	return nil
}
