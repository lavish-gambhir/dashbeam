package repository

import (
	"context"

	"github.com/lavish-gambhir/dashbeam/shared/models"
)

// ClickHouse defines the interface for ClickHouse operations
type ClickHouse interface {
	// Event storage
	InsertEvents(ctx context.Context, records []models.AnalyticsRecord) error

	// Metrics operations
	UpsertUserActivityMetrics(ctx context.Context, metrics []models.UserActivityMetric) error
	UpsertQuizMetrics(ctx context.Context, metrics []models.QuizMetric) error
	UpsertSchoolMetrics(ctx context.Context, metrics []models.SchoolMetric) error

	// Query operations for aggregations
	GetUserActivityByDateRange(ctx context.Context, userID string, startDate, endDate string) ([]models.UserActivityMetric, error)
	GetQuizMetricsByDateRange(ctx context.Context, quizID string, startDate, endDate string) ([]models.QuizMetric, error)
	GetSchoolMetricsByDateRange(ctx context.Context, schoolID string, startDate, endDate string) ([]models.SchoolMetric, error)
}
