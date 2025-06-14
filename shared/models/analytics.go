package models

import (
	"time"

	"github.com/google/uuid"
)

// AnalyticsRecord represents a processed event record for ClickHouse
type AnalyticsRecord struct {
	EventID     uuid.UUID      `json:"event_id" ch:"event_id"`
	EventType   string         `json:"event_type" ch:"event_type"`
	Category    string         `json:"category" ch:"category"`
	Action      string         `json:"action" ch:"action"`
	UserID      uuid.UUID      `json:"user_id" ch:"user_id"`
	SchoolID    uuid.UUID      `json:"school_id" ch:"school_id"`
	SessionID   *uuid.UUID     `json:"session_id,omitempty" ch:"session_id"`
	QuizID      *uuid.UUID     `json:"quiz_id,omitempty" ch:"quiz_id"`
	QuestionID  *uuid.UUID     `json:"question_id,omitempty" ch:"question_id"`
	Value       *float64       `json:"value,omitempty" ch:"value"`
	Metadata    map[string]any `json:"metadata" ch:"metadata"`
	Timestamp   time.Time      `json:"timestamp" ch:"timestamp"`
	ProcessedAt time.Time      `json:"processed_at" ch:"processed_at"`
}

// UserActivityMetric represents aggregated user activity data
type UserActivityMetric struct {
	UserID      uuid.UUID `json:"user_id" ch:"user_id"`
	SchoolID    uuid.UUID `json:"school_id" ch:"school_id"`
	Date        time.Time `json:"date" ch:"date"`
	LoginCount  int       `json:"login_count" ch:"login_count"`
	SessionTime int       `json:"session_time_minutes" ch:"session_time_minutes"`
	QuizCount   int       `json:"quiz_count" ch:"quiz_count"`
	UpdatedAt   time.Time `json:"updated_at" ch:"updated_at"`
}

// QuizMetric represents aggregated quiz performance data
type QuizMetric struct {
	QuizID              uuid.UUID `json:"quiz_id" ch:"quiz_id"`
	SchoolID            uuid.UUID `json:"school_id" ch:"school_id"`
	Date                time.Time `json:"date" ch:"date"`
	ParticipantCount    int       `json:"participant_count" ch:"participant_count"`
	CompletionRate      float64   `json:"completion_rate" ch:"completion_rate"`
	AverageScore        float64   `json:"average_score" ch:"average_score"`
	AverageResponseTime int       `json:"average_response_time_ms" ch:"average_response_time_ms"`
	UpdatedAt           time.Time `json:"updated_at" ch:"updated_at"`
}

// SchoolMetric represents aggregated school-level data
type SchoolMetric struct {
	SchoolID     uuid.UUID `json:"school_id" ch:"school_id"`
	Date         time.Time `json:"date" ch:"date"`
	ActiveUsers  int       `json:"active_users" ch:"active_users"`
	TotalQuizzes int       `json:"total_quizzes" ch:"total_quizzes"`
	TotalEvents  int       `json:"total_events" ch:"total_events"`
	UpdatedAt    time.Time `json:"updated_at" ch:"updated_at"`
}
