package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/lavish-gambhir/dashbeam/pkg/apperr"
	"github.com/lavish-gambhir/dashbeam/shared/database/postgres"
	"github.com/lavish-gambhir/dashbeam/shared/models"
)

type QuizRepository struct {
	db *postgres.DB
}

func NewQuizRepository(db *postgres.DB) *QuizRepository {
	return &QuizRepository{
		db: db,
	}
}

func (r *QuizRepository) CreateOrUpdateParticipant(ctx context.Context, participant *models.QuizParticipant) error {
	existing, err := r.GetParticipant(ctx, participant.SessionID, participant.UserID)
	if err != nil && !apperr.Is(err, apperr.DBRecordNotFound) {
		return apperr.Wrapf(err, apperr.DBQueryFailed, "failed to check existing participant: session=%s, user=%s", participant.SessionID, participant.UserID)
	}

	if existing != nil {
		return r.updateParticipant(ctx, participant)
	}

	return r.createParticipant(ctx, participant)
}

func (r *QuizRepository) createParticipant(ctx context.Context, participant *models.QuizParticipant) error {
	query := `
		INSERT INTO quiz_participants (
			id, session_id, user_id, joined_at, started_at, submitted_at, disconnected_at,
			total_score, max_possible_score, completion_percentage, total_time_seconds,
			questions_answered, questions_correct, questions_skipped, average_response_time_ms,
			fastest_response_time_ms, slowest_response_time_ms, total_interactions,
			answer_changes_count, focus_loss_count, status
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21
		)`

	_, err := r.db.Conn(ctx).Exec(ctx, query,
		participant.ID,
		participant.SessionID,
		participant.UserID,
		participant.JoinedAt,
		participant.StartedAt,
		participant.SubmittedAt,
		participant.DisconnectedAt,
		participant.TotalScore,
		participant.MaxPossibleScore,
		participant.CompletionPercentage,
		participant.TotalTimeSeconds,
		participant.QuestionsAnswered,
		participant.QuestionsCorrect,
		participant.QuestionsSkipped,
		participant.AverageResponseTimeMS,
		participant.FastestResponseTimeMS,
		participant.SlowestResponseTimeMS,
		participant.TotalInteractions,
		participant.AnswerChangesCount,
		participant.FocusLossCount,
		participant.Status,
	)

	if err != nil {
		return apperr.Wrapf(err, apperr.DBQueryFailed, "failed to create participant: session=%s, user=%s", participant.SessionID, participant.UserID)
	}

	return nil
}

func (r *QuizRepository) updateParticipant(ctx context.Context, participant *models.QuizParticipant) error {
	query := `
		UPDATE quiz_participants SET
			started_at = $3,
			submitted_at = $4,
			disconnected_at = $5,
			total_score = $6,
			max_possible_score = $7,
			completion_percentage = $8,
			total_time_seconds = $9,
			questions_answered = $10,
			questions_correct = $11,
			questions_skipped = $12,
			average_response_time_ms = $13,
			fastest_response_time_ms = $14,
			slowest_response_time_ms = $15,
			total_interactions = $16,
			answer_changes_count = $17,
			focus_loss_count = $18,
			status = $19
		WHERE session_id = $1 AND user_id = $2`

	result, err := r.db.Conn(ctx).Exec(ctx, query,
		participant.SessionID,
		participant.UserID,
		participant.StartedAt,
		participant.SubmittedAt,
		participant.DisconnectedAt,
		participant.TotalScore,
		participant.MaxPossibleScore,
		participant.CompletionPercentage,
		participant.TotalTimeSeconds,
		participant.QuestionsAnswered,
		participant.QuestionsCorrect,
		participant.QuestionsSkipped,
		participant.AverageResponseTimeMS,
		participant.FastestResponseTimeMS,
		participant.SlowestResponseTimeMS,
		participant.TotalInteractions,
		participant.AnswerChangesCount,
		participant.FocusLossCount,
		participant.Status,
	)

	if err != nil {
		return apperr.Wrapf(err, apperr.DBQueryFailed, "failed to update participant: session=%s, user=%s", participant.SessionID, participant.UserID)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return apperr.Newf(apperr.DBRecordNotFound, "participant not found: session=%s, user=%s", participant.SessionID, participant.UserID)
	}

	return nil
}

func (r *QuizRepository) GetParticipant(ctx context.Context, sessionID, userID string) (*models.QuizParticipant, error) {
	query := `
		SELECT
			id, session_id, user_id, joined_at, started_at, submitted_at, disconnected_at,
			total_score, max_possible_score, completion_percentage, total_time_seconds,
			questions_answered, questions_correct, questions_skipped, average_response_time_ms,
			fastest_response_time_ms, slowest_response_time_ms, total_interactions,
			answer_changes_count, focus_loss_count, status
		FROM quiz_participants
		WHERE session_id = $1 AND user_id = $2`

	var participant models.QuizParticipant
	err := r.db.Conn(ctx).QueryRow(ctx, query, sessionID, userID).Scan(
		&participant.ID,
		&participant.SessionID,
		&participant.UserID,
		&participant.JoinedAt,
		&participant.StartedAt,
		&participant.SubmittedAt,
		&participant.DisconnectedAt,
		&participant.TotalScore,
		&participant.MaxPossibleScore,
		&participant.CompletionPercentage,
		&participant.TotalTimeSeconds,
		&participant.QuestionsAnswered,
		&participant.QuestionsCorrect,
		&participant.QuestionsSkipped,
		&participant.AverageResponseTimeMS,
		&participant.FastestResponseTimeMS,
		&participant.SlowestResponseTimeMS,
		&participant.TotalInteractions,
		&participant.AnswerChangesCount,
		&participant.FocusLossCount,
		&participant.Status,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperr.Newf(apperr.DBRecordNotFound, "participant not found: session=%s, user=%s", sessionID, userID)
		}
		return nil, apperr.Wrapf(err, apperr.DBQueryFailed, "failed to get participant: session=%s, user=%s", sessionID, userID)
	}

	return &participant, nil
}

func (r *QuizRepository) UpdateParticipantProgress(ctx context.Context, sessionID, userID string, score float64, questionsAnswered, questionsCorrect int, responseTimeMS int) error {
	participant, err := r.GetParticipant(ctx, sessionID, userID)
	if err != nil {
		if apperr.Is(err, apperr.DBRecordNotFound) {
			newParticipant := &models.QuizParticipant{
				ID:                    uuid.New().String(),
				SessionID:             sessionID,
				UserID:                userID,
				JoinedAt:              time.Now().UTC(),
				TotalScore:            score,
				QuestionsAnswered:     questionsAnswered,
				QuestionsCorrect:      questionsCorrect,
				AverageResponseTimeMS: responseTimeMS,
				FastestResponseTimeMS: &responseTimeMS,
				SlowestResponseTimeMS: &responseTimeMS,
				Status:                string(models.ParticipantStatusActive),
			}
			return r.createParticipant(ctx, newParticipant)
		}
		return err
	}

	participant.TotalScore += score
	participant.QuestionsAnswered += questionsAnswered
	participant.QuestionsCorrect += questionsCorrect

	if participant.QuestionsAnswered > 0 {
		totalResponseTime := participant.AverageResponseTimeMS * (participant.QuestionsAnswered - questionsAnswered)
		totalResponseTime += responseTimeMS
		participant.AverageResponseTimeMS = totalResponseTime / participant.QuestionsAnswered
	}

	if participant.FastestResponseTimeMS == nil || responseTimeMS < *participant.FastestResponseTimeMS {
		participant.FastestResponseTimeMS = &responseTimeMS
	}
	if participant.SlowestResponseTimeMS == nil || responseTimeMS > *participant.SlowestResponseTimeMS {
		participant.SlowestResponseTimeMS = &responseTimeMS
	}

	participant.Status = string(models.ParticipantStatusActive)

	return r.updateParticipant(ctx, participant)
}

func (r *QuizRepository) CompleteParticipantSession(ctx context.Context, sessionID, userID string, finalScore, maxScore float64, completionTimeSeconds int) error {
	participant, err := r.GetParticipant(ctx, sessionID, userID)
	if err != nil {
		return apperr.Wrapf(err, apperr.DBQueryFailed, "failed to get participant for completion: session=%s, user=%s", sessionID, userID)
	}

	now := time.Now().UTC()
	participant.SubmittedAt = &now
	participant.TotalScore = finalScore
	participant.MaxPossibleScore = maxScore
	participant.TotalTimeSeconds = completionTimeSeconds
	participant.Status = string(models.ParticipantStatusCompleted)

	if maxScore > 0 {
		participant.CompletionPercentage = (finalScore / maxScore) * 100
	}

	return r.updateParticipant(ctx, participant)
}

func (r *QuizRepository) GetParticipantsBySession(ctx context.Context, sessionID string) ([]*models.QuizParticipant, error) {
	query := `
		SELECT
			id, session_id, user_id, joined_at, started_at, submitted_at, disconnected_at,
			total_score, max_possible_score, completion_percentage, total_time_seconds,
			questions_answered, questions_correct, questions_skipped, average_response_time_ms,
			fastest_response_time_ms, slowest_response_time_ms, total_interactions,
			answer_changes_count, focus_loss_count, status
		FROM quiz_participants
		WHERE session_id = $1
		ORDER BY total_score DESC, joined_at ASC`

	rows, err := r.db.Conn(ctx).Query(ctx, query, sessionID)
	if err != nil {
		return nil, apperr.Wrapf(err, apperr.DBQueryFailed, "failed to get participants for session: %s", sessionID)
	}
	defer rows.Close()

	var participants []*models.QuizParticipant
	for rows.Next() {
		var participant models.QuizParticipant
		err := rows.Scan(
			&participant.ID,
			&participant.SessionID,
			&participant.UserID,
			&participant.JoinedAt,
			&participant.StartedAt,
			&participant.SubmittedAt,
			&participant.DisconnectedAt,
			&participant.TotalScore,
			&participant.MaxPossibleScore,
			&participant.CompletionPercentage,
			&participant.TotalTimeSeconds,
			&participant.QuestionsAnswered,
			&participant.QuestionsCorrect,
			&participant.QuestionsSkipped,
			&participant.AverageResponseTimeMS,
			&participant.FastestResponseTimeMS,
			&participant.SlowestResponseTimeMS,
			&participant.TotalInteractions,
			&participant.AnswerChangesCount,
			&participant.FocusLossCount,
			&participant.Status,
		)
		if err != nil {
			return nil, apperr.Wrapf(err, apperr.DBQueryFailed, "failed to scan participant row for session: %s", sessionID)
		}
		participants = append(participants, &participant)
	}

	if err := rows.Err(); err != nil {
		return nil, apperr.Wrapf(err, apperr.DBQueryFailed, "error iterating participant rows for session: %s", sessionID)
	}

	return participants, nil
}

func (r *QuizRepository) MarkParticipantDisconnected(ctx context.Context, sessionID, userID string) error {
	query := `
		UPDATE quiz_participants SET
			disconnected_at = $3,
			status = $4
		WHERE session_id = $1 AND user_id = $2`

	now := time.Now().UTC()
	result, err := r.db.Conn(ctx).Exec(ctx, query, sessionID, userID, now, string(models.ParticipantStatusDisconnected))

	if err != nil {
		return apperr.Wrapf(err, apperr.DBQueryFailed, "failed to mark participant disconnected: session=%s, user=%s", sessionID, userID)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return apperr.Newf(apperr.DBRecordNotFound, "participant not found: session=%s, user=%s", sessionID, userID)
	}

	return nil
}
