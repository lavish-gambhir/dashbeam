package repositories

import (
	"context"
	"encoding/json"

	"github.com/lavish-gambhir/dashbeam/pkg/apperr"
	"github.com/lavish-gambhir/dashbeam/shared/database/clickhouse"
	"github.com/lavish-gambhir/dashbeam/shared/models"
)

type ClickHouseRepository struct {
	db *clickhouse.DB
}

func NewClickHouseRepository(db *clickhouse.DB) *ClickHouseRepository {
	return &ClickHouseRepository{
		db: db,
	}
}

func (r *ClickHouseRepository) InsertEvents(ctx context.Context, records []models.AnalyticsRecord) error {
	if len(records) == 0 {
		return nil
	}

	batch, err := r.db.PrepareBatch(ctx, "INSERT INTO events")
	if err != nil {
		return apperr.Wrap(err, apperr.Internal, "failed to prepare batch")
	}

	for _, record := range records {
		metadataJSON, _ := json.Marshal(record.Metadata)

		err := batch.Append(
			record.EventID,
			record.EventType,
			record.Category,
			record.Action,
			record.UserID,
			record.SchoolID,
			record.SessionID,
			record.QuizID,
			record.QuestionID,
			record.Value,
			string(metadataJSON),
			record.Timestamp,
			record.ProcessedAt,
		)
		if err != nil {
			return apperr.Wrap(err, apperr.Internal, "failed to append to batch")
		}
	}

	return batch.Send()
}

func (r *ClickHouseRepository) UpsertUserActivityMetrics(ctx context.Context, metrics []models.UserActivityMetric) error {
	if len(metrics) == 0 {
		return nil
	}

	batch, err := r.db.PrepareBatch(ctx, "INSERT INTO user_activity_metrics")
	if err != nil {
		return apperr.Wrap(err, apperr.Internal, "failed to prepare batch")
	}

	for _, metric := range metrics {
		err := batch.Append(
			metric.UserID,
			metric.SchoolID,
			metric.Date,
			metric.LoginCount,
			metric.SessionTime,
			metric.QuizCount,
			metric.UpdatedAt,
		)
		if err != nil {
			return apperr.Wrap(err, apperr.Internal, "failed to append to batch")
		}
	}

	return batch.Send()
}

func (r *ClickHouseRepository) UpsertQuizMetrics(ctx context.Context, metrics []models.QuizMetric) error {
	if len(metrics) == 0 {
		return nil
	}

	batch, err := r.db.PrepareBatch(ctx, "INSERT INTO quiz_metrics")
	if err != nil {
		return apperr.Wrap(err, apperr.Internal, "failed to prepare batch")
	}

	for _, metric := range metrics {
		err := batch.Append(
			metric.QuizID,
			metric.SchoolID,
			metric.Date,
			metric.ParticipantCount,
			metric.CompletionRate,
			metric.AverageScore,
			metric.AverageResponseTime,
			metric.UpdatedAt,
		)
		if err != nil {
			return apperr.Wrap(err, apperr.Internal, "failed to append to batch")
		}
	}

	return batch.Send()
}

func (r *ClickHouseRepository) UpsertSchoolMetrics(ctx context.Context, metrics []models.SchoolMetric) error {
	if len(metrics) == 0 {
		return nil
	}

	batch, err := r.db.PrepareBatch(ctx, "INSERT INTO school_metrics")
	if err != nil {
		return apperr.Wrap(err, apperr.Internal, "failed to prepare batch")
	}

	for _, metric := range metrics {
		err := batch.Append(
			metric.SchoolID,
			metric.Date,
			metric.ActiveUsers,
			metric.TotalQuizzes,
			metric.TotalEvents,
			metric.UpdatedAt,
		)
		if err != nil {
			return apperr.Wrap(err, apperr.Internal, "failed to append to batch")
		}
	}

	return batch.Send()
}

func (r *ClickHouseRepository) GetUserActivityByDateRange(ctx context.Context, userID string, startDate, endDate string) ([]models.UserActivityMetric, error) {
	query := `
		SELECT user_id, school_id, date, login_count, session_time_minutes, quiz_count, updated_at
		FROM user_activity_metrics
		WHERE user_id = ? AND date >= ? AND date <= ?
		ORDER BY date DESC
	`

	rows, err := r.db.Query(ctx, query, userID, startDate, endDate)
	if err != nil {
		return nil, apperr.Wrap(err, apperr.Internal, "failed to query user activity metrics")
	}
	defer rows.Close()

	var metrics []models.UserActivityMetric
	for rows.Next() {
		var metric models.UserActivityMetric
		err := rows.Scan(
			&metric.UserID,
			&metric.SchoolID,
			&metric.Date,
			&metric.LoginCount,
			&metric.SessionTime,
			&metric.QuizCount,
			&metric.UpdatedAt,
		)
		if err != nil {
			return nil, apperr.Wrap(err, apperr.Internal, "failed to scan user activity metric")
		}
		metrics = append(metrics, metric)
	}

	return metrics, nil
}

func (r *ClickHouseRepository) GetQuizMetricsByDateRange(ctx context.Context, quizID string, startDate, endDate string) ([]models.QuizMetric, error) {
	query := `
		SELECT quiz_id, school_id, date, participant_count, completion_rate, average_score, average_response_time_ms, updated_at
		FROM quiz_metrics
		WHERE quiz_id = ? AND date >= ? AND date <= ?
		ORDER BY date DESC
	`

	rows, err := r.db.Query(ctx, query, quizID, startDate, endDate)
	if err != nil {
		return nil, apperr.Wrap(err, apperr.Internal, "failed to query quiz metrics")
	}
	defer rows.Close()

	var metrics []models.QuizMetric
	for rows.Next() {
		var metric models.QuizMetric
		err := rows.Scan(
			&metric.QuizID,
			&metric.SchoolID,
			&metric.Date,
			&metric.ParticipantCount,
			&metric.CompletionRate,
			&metric.AverageScore,
			&metric.AverageResponseTime,
			&metric.UpdatedAt,
		)
		if err != nil {
			return nil, apperr.Wrap(err, apperr.Internal, "failed to scan quiz metric")
		}
		metrics = append(metrics, metric)
	}

	return metrics, nil
}

func (r *ClickHouseRepository) GetSchoolMetricsByDateRange(ctx context.Context, schoolID string, startDate, endDate string) ([]models.SchoolMetric, error) {
	query := `
		SELECT school_id, date, active_users, total_quizzes, total_events, updated_at
		FROM school_metrics
		WHERE school_id = ? AND date >= ? AND date <= ?
		ORDER BY date DESC
	`

	rows, err := r.db.Query(ctx, query, schoolID, startDate, endDate)
	if err != nil {
		return nil, apperr.Wrap(err, apperr.Internal, "failed to query school metrics")
	}
	defer rows.Close()

	var metrics []models.SchoolMetric
	for rows.Next() {
		var metric models.SchoolMetric
		err := rows.Scan(
			&metric.SchoolID,
			&metric.Date,
			&metric.ActiveUsers,
			&metric.TotalQuizzes,
			&metric.TotalEvents,
			&metric.UpdatedAt,
		)
		if err != nil {
			return nil, apperr.Wrap(err, apperr.Internal, "failed to scan school metric")
		}
		metrics = append(metrics, metric)
	}

	return metrics, nil
}
