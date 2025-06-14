package analytics

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/lavish-gambhir/dashbeam/pkg/apperr"
	"github.com/lavish-gambhir/dashbeam/services/analytics/repository"
	"github.com/lavish-gambhir/dashbeam/shared/models"
	"github.com/lavish-gambhir/dashbeam/shared/streaming"
)

type EventProcessor struct {
	clickhouseRepo repository.ClickHouse
	logger         *slog.Logger
	batchSize      uint
}

func NewEventProcessor(
	clickhouseRepo repository.ClickHouse,
	logger *slog.Logger,
	batchSize uint,
) *EventProcessor {
	return &EventProcessor{
		clickhouseRepo: clickhouseRepo,
		logger:         logger.With("component", "event_processor"),
		batchSize:      batchSize,
	}
}

func (ep *EventProcessor) ProcessEvents(ctx context.Context, events []streaming.Event) error {
	if len(events) == 0 {
		return nil
	}

	ep.logger.Info("processing events batch", slog.Int("count", len(events)))

	// Transform events into analytics records
	records := make([]models.AnalyticsRecord, 0, len(events))
	for _, event := range events {
		record, err := ep.transformEvent(event)
		if err != nil {
			ep.logger.Warn("failed to transform event",
				slog.String("eventID", event.ID.String()),
				slog.Any("error", err))
			continue
		}
		records = append(records, record)
	}

	if len(records) == 0 {
		ep.logger.Warn("no valid records to process")
		return nil
	}

	// Store in ClickHouse
	if err := ep.clickhouseRepo.InsertEvents(ctx, records); err != nil {
		return apperr.Wrap(err, apperr.Internal, "failed to insert events into ClickHouse")
	}

	// Update aggregated metrics
	if err := ep.updateMetrics(ctx, records); err != nil {
		ep.logger.Error("failed to update metrics", slog.Any("error", err))
		// Don't fail the entire batch for metrics update failure
	}

	ep.logger.Info("successfully processed events batch",
		slog.Int("processed", len(records)),
		slog.Int("skipped", len(events)-len(records)))

	return nil
}

func (ep *EventProcessor) transformEvent(event streaming.Event) (models.AnalyticsRecord, error) {
	// Base record with common fields
	record := models.AnalyticsRecord{
		EventID:     event.ID,
		EventType:   event.Type.String(),
		UserID:      event.UserID,
		SchoolID:    event.SchoolID,
		Timestamp:   event.Timestamp,
		ProcessedAt: time.Now().UTC(),
		Metadata:    make(map[string]any),
	}

	// Add session ID from metadata if available
	if event.Metadata.SessionID != nil {
		if sessionUUID, err := uuid.Parse(*event.Metadata.SessionID); err == nil {
			record.SessionID = &sessionUUID
		}
	}

	// Transform based on event type
	switch event.Type {
	case streaming.UserLogin:
		return ep.transformUserLoginEvent(record, event)
	case streaming.UserLogout:
		return ep.transformUserLogoutEvent(record, event)
	case streaming.QuizSessionStarted:
		return ep.transformQuizStartEvent(record, event)
	case streaming.QuizSessionCompleted:
		return ep.transformQuizSubmitEvent(record, event)
	case streaming.QuizAnswerSubmitted:
		return ep.transformQuizAnswerEvent(record, event)
	case streaming.QuizQuestionShown:
		return ep.transformQuizQuestionEvent(record, event)
	case streaming.AppInteraction:
		return ep.transformAppInteractionEvent(record, event)
	case streaming.AppNavigation:
		return ep.transformAppNavigationEvent(record, event)
	case streaming.SystemStartup, streaming.SystemShutdown:
		return ep.transformSystemEvent(record, event)
	default:
		return record, apperr.Newf(apperr.BadRequest, "unsupported event type: %s", event.Type)
	}
}

func (ep *EventProcessor) transformUserLoginEvent(base models.AnalyticsRecord, event streaming.Event) (models.AnalyticsRecord, error) {
	base.Category = "user_activity"
	base.Action = "login"

	if payload, ok := event.Payload.(streaming.UserLoginPayload); ok {
		base.Metadata = map[string]any{
			"login_method":  payload.LoginMethod,
			"session_start": payload.SessionStart,
		}
		if payload.UserAgent != nil {
			base.Metadata["user_agent"] = *payload.UserAgent
		}
		if payload.Email != nil {
			base.Metadata["email"] = *payload.Email
		}
		if payload.Role != nil {
			base.Metadata["role"] = *payload.Role
		}
	}

	return base, nil
}

func (ep *EventProcessor) transformUserLogoutEvent(base models.AnalyticsRecord, event streaming.Event) (models.AnalyticsRecord, error) {
	base.Category = "user_activity"
	base.Action = "logout"

	if payload, ok := event.Payload.(streaming.UserLogoutPayload); ok {
		base.Value = float64Ptr(float64(payload.SessionDuration))
		base.Metadata = map[string]any{
			"session_duration_ms": payload.SessionDuration,
			"logout_reason":       payload.LogoutReason,
		}
	}

	return base, nil
}

func (ep *EventProcessor) transformQuizStartEvent(base models.AnalyticsRecord, event streaming.Event) (models.AnalyticsRecord, error) {
	base.Category = "quiz_activity"
	base.Action = "quiz_started"

	if payload, ok := event.Payload.(streaming.QuizSessionStartedPayload); ok {
		base.QuizID = &payload.QuizID
		sessionID := payload.SessionID
		base.SessionID = &sessionID
		base.Value = float64Ptr(float64(payload.TotalQuestions))

		base.Metadata = map[string]any{
			"session_code":    payload.SessionCode,
			"total_questions": payload.TotalQuestions,
			"max_score":       payload.MaxScore,
		}
		if payload.TimeLimit != nil {
			base.Metadata["time_limit_seconds"] = *payload.TimeLimit
		}
	}

	return base, nil
}

func (ep *EventProcessor) transformQuizSubmitEvent(base models.AnalyticsRecord, event streaming.Event) (models.AnalyticsRecord, error) {
	base.Category = "quiz_activity"
	base.Action = "quiz_completed"

	if payload, ok := event.Payload.(streaming.QuizSessionCompletedPayload); ok {
		base.QuizID = &payload.QuizID
		sessionID := payload.SessionID
		base.SessionID = &sessionID
		base.Value = &payload.TotalScore

		base.Metadata = map[string]any{
			"total_score":           payload.TotalScore,
			"max_score":             payload.MaxScore,
			"completion_time_ms":    payload.CompletionTimeMS,
			"questions_correct":     payload.QuestionsCorrect,
			"questions_answered":    payload.QuestionsAnswered,
			"questions_skipped":     payload.QuestionsSkipped,
			"average_response_time": payload.AverageResponseTime,
		}
	}

	return base, nil
}

func (ep *EventProcessor) transformQuizAnswerEvent(base models.AnalyticsRecord, event streaming.Event) (models.AnalyticsRecord, error) {
	base.Category = "quiz_activity"
	base.Action = "answer_submitted"

	if payload, ok := event.Payload.(streaming.QuizAnswerSubmittedPayload); ok {
		base.QuizID = &payload.QuizID
		base.QuestionID = &payload.QuestionID
		sessionID := payload.SessionID
		base.SessionID = &sessionID
		base.Value = float64Ptr(float64(payload.ResponseTimeMS))

		base.Metadata = map[string]any{
			"question_sequence": payload.QuestionSequence,
			"response_time_ms":  payload.ResponseTimeMS,
		}

		if payload.IsCorrect != nil {
			base.Metadata["is_correct"] = *payload.IsCorrect
		}
		if payload.AnswerChanges != nil {
			base.Metadata["answer_changes"] = *payload.AnswerChanges
		}
		if payload.Points != nil {
			base.Metadata["points"] = *payload.Points
		}
	}

	return base, nil
}

func (ep *EventProcessor) transformQuizQuestionEvent(base models.AnalyticsRecord, event streaming.Event) (models.AnalyticsRecord, error) {
	base.Category = "quiz_activity"
	base.Action = "question_shown"

	if payload, ok := event.Payload.(streaming.QuizQuestionShownPayload); ok {
		base.QuizID = &payload.QuizID
		base.QuestionID = &payload.QuestionID
		sessionID := payload.SessionID
		base.SessionID = &sessionID
		base.Value = float64Ptr(float64(payload.QuestionSequence))

		base.Metadata = map[string]any{
			"question_sequence": payload.QuestionSequence,
			"question_type":     payload.QuestionType,
		}
		if payload.TimeLimit != nil {
			base.Metadata["time_limit_seconds"] = *payload.TimeLimit
		}
	}

	return base, nil
}

func (ep *EventProcessor) transformAppInteractionEvent(base models.AnalyticsRecord, event streaming.Event) (models.AnalyticsRecord, error) {
	base.Category = "app_activity"
	base.Action = "interaction"

	if payload, ok := event.Payload.(streaming.AppInteractionPayload); ok {
		base.Metadata = map[string]any{
			"interaction_type": payload.InteractionType,
			"screen_name":      payload.ScreenName,
		}
		if payload.ElementClicked != nil {
			base.Metadata["element_clicked"] = *payload.ElementClicked
		}
		if payload.TimeSpentMS != nil {
			base.Value = float64Ptr(float64(*payload.TimeSpentMS))
			base.Metadata["time_spent_ms"] = *payload.TimeSpentMS
		}
		if payload.Action != nil {
			base.Metadata["action"] = *payload.Action
		}
	}

	return base, nil
}

func (ep *EventProcessor) transformAppNavigationEvent(base models.AnalyticsRecord, event streaming.Event) (models.AnalyticsRecord, error) {
	base.Category = "app_activity"
	base.Action = "navigation"

	if payload, ok := event.Payload.(streaming.AppNavigationPayload); ok {
		base.Metadata = map[string]any{
			"from_screen":     payload.FromScreen,
			"to_screen":       payload.ToScreen,
			"navigation_type": payload.NavigationType,
		}
		if payload.TimeSpentMS != nil {
			base.Value = float64Ptr(float64(*payload.TimeSpentMS))
			base.Metadata["time_spent_ms"] = *payload.TimeSpentMS
		}
	}

	return base, nil
}

func (ep *EventProcessor) transformSystemEvent(base models.AnalyticsRecord, event streaming.Event) (models.AnalyticsRecord, error) {
	base.Category = "system"
	base.Action = "system_event"

	// Add basic metadata
	base.Metadata = map[string]any{
		"event_type": event.Type.String(),
	}

	return base, nil
}

func (ep *EventProcessor) updateMetrics(ctx context.Context, records []models.AnalyticsRecord) error {
	// Update various metrics based on the processed records

	// User activity metrics
	if err := ep.updateUserActivityMetrics(ctx, records); err != nil {
		return err
	}

	// Quiz performance metrics
	if err := ep.updateQuizMetrics(ctx, records); err != nil {
		return err
	}

	// School-level metrics
	if err := ep.updateSchoolMetrics(ctx, records); err != nil {
		return err
	}

	return nil
}

func (ep *EventProcessor) updateUserActivityMetrics(ctx context.Context, records []models.AnalyticsRecord) error {
	userMetrics := make(map[string]*models.UserActivityMetric)

	for _, record := range records {
		if record.Category != "user_activity" {
			continue
		}

		dateKey := record.Timestamp.Format("2006-01-02")
		userKey := fmt.Sprintf("%s_%s", record.UserID.String(), dateKey)

		metric, exists := userMetrics[userKey]
		if !exists {
			metric = &models.UserActivityMetric{
				UserID:    record.UserID,
				SchoolID:  record.SchoolID,
				Date:      record.Timestamp.Truncate(24 * time.Hour),
				UpdatedAt: time.Now().UTC(),
			}
			userMetrics[userKey] = metric
		}

		switch record.Action {
		case "login":
			metric.LoginCount++
		case "logout":
			if record.Value != nil {
				// Convert session duration from ms to minutes
				sessionMinutes := int(*record.Value / 60000)
				metric.SessionTime += sessionMinutes
			}
		}
	}

	// Convert map to slice
	metrics := make([]models.UserActivityMetric, 0, len(userMetrics))
	for _, metric := range userMetrics {
		metrics = append(metrics, *metric)
	}

	return ep.clickhouseRepo.UpsertUserActivityMetrics(ctx, metrics)
}

func (ep *EventProcessor) updateQuizMetrics(ctx context.Context, records []models.AnalyticsRecord) error {
	quizMetrics := make(map[string]*models.QuizMetric)

	for _, record := range records {
		if record.Category != "quiz_activity" || record.QuizID == nil {
			continue
		}

		dateKey := record.Timestamp.Format("2006-01-02")
		quizKey := fmt.Sprintf("%s_%s", record.QuizID.String(), dateKey)

		metric, exists := quizMetrics[quizKey]
		if !exists {
			metric = &models.QuizMetric{
				QuizID:    *record.QuizID,
				SchoolID:  record.SchoolID,
				Date:      record.Timestamp.Truncate(24 * time.Hour),
				UpdatedAt: time.Now().UTC(),
			}
			quizMetrics[quizKey] = metric
		}

		switch record.Action {
		case "quiz_started":
			metric.ParticipantCount++
		case "quiz_completed":
			if record.Value != nil {
				metric.AverageScore = *record.Value
			}
		case "answer_submitted":
			if record.Value != nil {
				metric.AverageResponseTime = int(*record.Value)
			}
		}
	}

	// Convert map to slice
	metrics := make([]models.QuizMetric, 0, len(quizMetrics))
	for _, metric := range quizMetrics {
		metrics = append(metrics, *metric)
	}

	return ep.clickhouseRepo.UpsertQuizMetrics(ctx, metrics)
}

func (ep *EventProcessor) updateSchoolMetrics(ctx context.Context, records []models.AnalyticsRecord) error {
	schoolMetrics := make(map[string]*models.SchoolMetric)
	activeUsers := make(map[string]map[string]bool) // schoolKey -> userID -> exists

	for _, record := range records {
		dateKey := record.Timestamp.Format("2006-01-02")
		schoolKey := fmt.Sprintf("%s_%s", record.SchoolID.String(), dateKey)

		metric, exists := schoolMetrics[schoolKey]
		if !exists {
			metric = &models.SchoolMetric{
				SchoolID:  record.SchoolID,
				Date:      record.Timestamp.Truncate(24 * time.Hour),
				UpdatedAt: time.Now().UTC(),
			}
			schoolMetrics[schoolKey] = metric
			activeUsers[schoolKey] = make(map[string]bool)
		}

		metric.TotalEvents++

		// Track unique active users
		activeUsers[schoolKey][record.UserID.String()] = true

		if record.Category == "quiz_activity" && record.Action == "quiz_started" {
			metric.TotalQuizzes++
		}
	}

	// Update active user counts
	for schoolKey, users := range activeUsers {
		if metric, exists := schoolMetrics[schoolKey]; exists {
			metric.ActiveUsers = len(users)
		}
	}

	// Convert map to slice
	metrics := make([]models.SchoolMetric, 0, len(schoolMetrics))
	for _, metric := range schoolMetrics {
		metrics = append(metrics, *metric)
	}

	return ep.clickhouseRepo.UpsertSchoolMetrics(ctx, metrics)
}

// Helper function to create float64 pointer
func float64Ptr(f float64) *float64 {
	return &f
}
