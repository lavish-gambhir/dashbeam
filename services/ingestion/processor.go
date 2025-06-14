package ingestion

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/lavish-gambhir/dashbeam/pkg/apperr"
	"github.com/lavish-gambhir/dashbeam/shared/models"
	"github.com/lavish-gambhir/dashbeam/shared/streaming"
)

// processBatchEvents processes multiple events in a batch
func (h *handler) processBatchEvents(ctx context.Context, events []streaming.Event) ([]string, error) {
	var eventIDs []string

	for i, event := range events {
		if event.ID == uuid.Nil {
			id, err := uuid.NewV7()
			if err != nil {
				return nil, apperr.Wrapf(err, apperr.Internal, "failed to generate event ID for batch item %d", i)
			}
			event.ID = id
			events[i] = event
		}

		if event.Timestamp.IsZero() {
			event.Timestamp = time.Now().UTC()
			events[i] = event
		}

		if err := event.Validate(); err != nil {
			return nil, apperr.Wrapf(err, apperr.ValidationFailed, "validation failed for batch item %d", i)
		}

		if err := event.Payload.Validate(); err != nil {
			return nil, apperr.Wrapf(err, apperr.ValidationFailed, "payload validation failed for batch item %d", i)
		}

		if err := h.processOperationalData(ctx, event); err != nil {
			h.logger.Warn("failed to update operational data", "event_id", event.ID.String(), "error", err)
			// do not fail: Continue processing other events
		}

		// Publish event to message queue
		topic := streaming.GetTopicForEventType(event.Type)
		if err := h.messageQueue.Publish(ctx, topic, event); err != nil {
			h.logger.Error("failed to publish event", "event_id", event.ID.String(), "topic", topic, "error", err)
			// do not fail: continue processing other events for batch
		}

		eventIDs = append(eventIDs, event.ID.String())
	}

	return eventIDs, nil
}

func (h *handler) processSingleEvent(ctx context.Context, event streaming.Event) (string, error) {
	if event.ID == uuid.Nil {
		id, err := uuid.NewV7()
		if err != nil {
			return "", apperr.Wrap(err, apperr.Internal, "failed to generate event ID")
		}
		event.ID = id
	}

	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now().UTC()
	}

	if err := event.Validate(); err != nil {
		return "", apperr.Wrap(err, apperr.ValidationFailed, "event validation failed")
	}
	if err := event.Payload.Validate(); err != nil {
		return "", apperr.Wrapf(err, apperr.ValidationFailed, "%s: %s", "event payload validation failed for event", event.ID.String())
	}

	// Process operational data updates
	if err := h.processOperationalData(ctx, event); err != nil {
		h.logger.Warn("failed to update operational data", "event_id", event.ID.String(), "error", err)
		// Continue with event publishing
	}

	// Publish event to message queue
	topic := streaming.GetTopicForEventType(event.Type)
	if err := h.messageQueue.Publish(ctx, topic, event); err != nil {
		return "", apperr.Wrapf(err, apperr.Internal, "failed to publish event to topic %s", topic)
	}

	return event.ID.String(), nil
}

// processOperationalData updates operational database based on event type
func (h *handler) processOperationalData(ctx context.Context, event streaming.Event) error {
	switch event.Type {
	case streaming.QuizAnswerSubmitted:
		return h.processQuizAnswerSubmitted(ctx, event)
	case streaming.QuizSessionCompleted:
		return h.processQuizSessionCompleted(ctx, event)
	case streaming.UserLogin:
		return h.processUserLogin(ctx, event)
	default:
		// No operational data update needed for this event type
		return nil
	}
}

func (h *handler) processQuizAnswerSubmitted(ctx context.Context, event streaming.Event) error {
	payload, ok := event.Payload.(streaming.QuizAnswerSubmittedPayload)
	if !ok {
		return apperr.New(apperr.ValidationFailed, "invalid payload type for quiz answer submitted event")
	}

	score := 0.0
	questionsCorrect := 0
	if payload.IsCorrect != nil && *payload.IsCorrect {
		score = 1.0
		questionsCorrect = 1
	}

	return h.quizRepo.UpdateParticipantProgress(
		ctx,
		payload.SessionID.String(),
		event.UserID.String(),
		score,
		1, // questionsAnswered
		questionsCorrect,
		payload.ResponseTimeMS,
	)
}

func (h *handler) processQuizSessionCompleted(ctx context.Context, event streaming.Event) error {
	payload, ok := event.Payload.(streaming.QuizSessionCompletedPayload)
	if !ok {
		return apperr.New(apperr.ValidationFailed, "invalid payload type for quiz session completed event")
	}

	return h.quizRepo.CompleteParticipantSession(
		ctx,
		payload.SessionID.String(),
		event.UserID.String(),
		payload.TotalScore,
		payload.MaxScore,
		payload.CompletionTimeMS/1000, // Convert to seconds
	)
}

func (h *handler) processUserLogin(ctx context.Context, event streaming.Event) error {
	payload, ok := event.Payload.(streaming.UserLoginPayload)
	if !ok {
		return apperr.New(apperr.ValidationFailed, "invalid payload type for user login event")
	}

	user := &models.User{
		ID:         event.UserID.String(),
		SchoolID:   event.SchoolID.String(),
		LastSeenAt: event.Timestamp,
	}

	if payload.Email != nil {
		user.Email = *payload.Email
	}
	if payload.Name != nil {
		user.Name = *payload.Name
	}
	if payload.Role != nil {
		user.Role = *payload.Role
	} else {
		user.Role = string(models.UserRoleStudent) // Default
	}

	// Try to get existing user
	existingUser, err := h.userRepo.GetUserByID(ctx, user.ID)
	if err != nil && !apperr.Is(err, apperr.DBRecordNotFound) {
		return err
	}

	if existingUser == nil {
		user.FirstSeenAt = event.Timestamp
		user.CreatedAt = time.Now().UTC()
		user.UpdatedAt = time.Now().UTC()
		return h.userRepo.CreateUser(ctx, user)
	}

	return h.userRepo.UpdateUserLastSeen(ctx, user.ID)
}
