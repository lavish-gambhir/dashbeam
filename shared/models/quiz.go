package models

import (
	"time"
)

type QuizParticipant struct {
	ID                    string     `json:"id" db:"id"`
	SessionID             string     `json:"session_id" db:"session_id"`
	UserID                string     `json:"user_id" db:"user_id"`
	JoinedAt              time.Time  `json:"joined_at" db:"joined_at"`
	StartedAt             *time.Time `json:"started_at" db:"started_at"`
	SubmittedAt           *time.Time `json:"submitted_at" db:"submitted_at"`
	DisconnectedAt        *time.Time `json:"disconnected_at" db:"disconnected_at"`
	TotalScore            float64    `json:"total_score" db:"total_score"`
	MaxPossibleScore      float64    `json:"max_possible_score" db:"max_possible_score"`
	CompletionPercentage  float64    `json:"completion_percentage" db:"completion_percentage"`
	TotalTimeSeconds      int        `json:"total_time_seconds" db:"total_time_seconds"`
	QuestionsAnswered     int        `json:"questions_answered" db:"questions_answered"`
	QuestionsCorrect      int        `json:"questions_correct" db:"questions_correct"`
	QuestionsSkipped      int        `json:"questions_skipped" db:"questions_skipped"`
	AverageResponseTimeMS int        `json:"average_response_time_ms" db:"average_response_time_ms"`
	FastestResponseTimeMS *int       `json:"fastest_response_time_ms" db:"fastest_response_time_ms"`
	SlowestResponseTimeMS *int       `json:"slowest_response_time_ms" db:"slowest_response_time_ms"`
	TotalInteractions     int        `json:"total_interactions" db:"total_interactions"`
	AnswerChangesCount    int        `json:"answer_changes_count" db:"answer_changes_count"`
	FocusLossCount        int        `json:"focus_loss_count" db:"focus_loss_count"`
	Status                string     `json:"status" db:"status"`
}

type ParticipantStatus string

const (
	ParticipantStatusJoined       ParticipantStatus = "joined"
	ParticipantStatusActive       ParticipantStatus = "active"
	ParticipantStatusCompleted    ParticipantStatus = "completed"
	ParticipantStatusDisconnected ParticipantStatus = "disconnected"
	ParticipantStatusAbandoned    ParticipantStatus = "abandoned"
)
