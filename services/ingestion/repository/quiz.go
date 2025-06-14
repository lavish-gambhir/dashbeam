package repository

import (
	"context"

	"github.com/lavish-gambhir/dashbeam/shared/models"
)

type Quiz interface {
	// CreateOrUpdateParticipant creates or updates a quiz participant
	CreateOrUpdateParticipant(ctx context.Context, participant *models.QuizParticipant) error

	// GetParticipant retrieves a quiz participant
	GetParticipant(ctx context.Context, sessionID, userID string) (*models.QuizParticipant, error)

	// UpdateParticipantProgress updates participant progress during quiz
	UpdateParticipantProgress(ctx context.Context, sessionID, userID string, score float64, questionsAnswered, questionsCorrect int, responseTimeMS int) error

	// CompleteParticipantSession marks a participant's session as completed
	CompleteParticipantSession(ctx context.Context, sessionID, userID string, finalScore, maxScore float64, completionTimeSeconds int) error
}
