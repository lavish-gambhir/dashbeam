package repositories

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/lavish-gambhir/dashbeam/pkg/apperr"
	"github.com/lavish-gambhir/dashbeam/shared/database/postgres"
	"github.com/lavish-gambhir/dashbeam/shared/models"
	"github.com/lavish-gambhir/dashbeam/shared/streaming"
)

type UserRepository struct {
	db *postgres.DB
}

func NewUserRepository(db *postgres.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (
			id, email, name, role, school_id, first_seen_at, last_seen_at,
			total_quiz_sessions, total_questions_answered, average_response_time_ms,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
		)`

	_, err := r.db.Conn(ctx).Exec(ctx, query,
		user.ID,
		user.Email,
		user.Name,
		user.Role,
		user.SchoolID,
		user.FirstSeenAt,
		user.LastSeenAt,
		user.TotalQuizSessions,
		user.TotalQuestionsAnswered,
		user.AverageResponseTimeMS,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		return apperr.Wrapf(err, apperr.DBQueryFailed, "failed to create user with ID: %s", user.ID)
	}

	return nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users SET
			email = $2,
			name = $3,
			role = $4,
			school_id = $5,
			last_seen_at = $6,
			total_quiz_sessions = $7,
			total_questions_answered = $8,
			average_response_time_ms = $9,
			updated_at = $10
		WHERE id = $1`

	result, err := r.db.Conn(ctx).Exec(ctx, query,
		user.ID,
		user.Email,
		user.Name,
		user.Role,
		user.SchoolID,
		user.LastSeenAt,
		user.TotalQuizSessions,
		user.TotalQuestionsAnswered,
		user.AverageResponseTimeMS,
		time.Now().UTC(),
	)

	if err != nil {
		return apperr.Wrapf(err, apperr.DBQueryFailed, "failed to update user with ID: %s", user.ID)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return apperr.New(apperr.DBRecordNotFound, "user not found")
	}

	return nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	query := `
		SELECT
			id, email, name, role, school_id, first_seen_at, last_seen_at,
			total_quiz_sessions, total_questions_answered, average_response_time_ms,
			created_at, updated_at
		FROM users
		WHERE id = $1`

	var user models.User
	err := r.db.Conn(ctx).QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Role,
		&user.SchoolID,
		&user.FirstSeenAt,
		&user.LastSeenAt,
		&user.TotalQuizSessions,
		&user.TotalQuestionsAnswered,
		&user.AverageResponseTimeMS,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperr.Newf(apperr.DBRecordNotFound, "user not found with ID: %s", id)
		}
		return nil, apperr.Wrapf(err, apperr.DBQueryFailed, "failed to get user by ID: %s", id)
	}

	return &user, nil
}

func (r *UserRepository) UpdateUserLastSeen(ctx context.Context, userID string) error {
	query := `
		UPDATE users SET
			last_seen_at = $2,
			updated_at = $2
		WHERE id = $1`

	now := time.Now().UTC()
	result, err := r.db.Conn(ctx).Exec(ctx, query, userID, now)

	if err != nil {
		return apperr.Wrapf(err, apperr.DBQueryFailed, "failed to update last seen for user: %s", userID)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return apperr.Newf(apperr.DBRecordNotFound, "user not found with ID: %s", userID)
	}

	return nil
}

func (r *UserRepository) CreateOrUpdateUserFromEvent(ctx context.Context, event streaming.Event) error {
	userID := event.UserID.String()
	schoolID := event.SchoolID.String()

	existingUser, err := r.GetUserByID(ctx, userID)
	if err != nil && !apperr.Is(err, apperr.DBRecordNotFound) {
		return apperr.Wrapf(err, apperr.DBQueryFailed, "failed to check existing user: %s", userID)
	}

	now := time.Now().UTC()

	if existingUser == nil {
		user := &models.User{
			ID:                     userID,
			Email:                  extractEmailFromPayload(event.Payload),
			Name:                   extractNameFromPayload(event.Payload),
			Role:                   extractRoleFromPayload(event.Payload),
			SchoolID:               schoolID,
			FirstSeenAt:            event.Timestamp,
			LastSeenAt:             event.Timestamp,
			TotalQuizSessions:      0,
			TotalQuestionsAnswered: 0,
			AverageResponseTimeMS:  0,
			CreatedAt:              now,
			UpdatedAt:              now,
		}

		return r.CreateUser(ctx, user)
	}

	existingUser.LastSeenAt = event.Timestamp
	existingUser.UpdatedAt = now

	return r.UpdateUser(ctx, existingUser)
}

func extractEmailFromPayload(payload map[string]any) string {
	if email, ok := payload["email"].(string); ok {
		return email
	}
	return ""
}

func extractNameFromPayload(payload map[string]any) string {
	if name, ok := payload["name"].(string); ok {
		return name
	}
	return ""
}

func extractRoleFromPayload(payload map[string]any) string {
	if role, ok := payload["role"].(string); ok {
		return role
	}
	return string(models.UserRoleStudent) // Default to student
}
