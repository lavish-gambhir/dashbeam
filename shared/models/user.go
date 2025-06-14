package models

import (
	"time"
)

type User struct {
	ID                     string    `json:"id" db:"id"`
	Email                  string    `json:"email" db:"email"`
	Name                   string    `json:"name" db:"name"`
	Role                   string    `json:"role" db:"role"`
	SchoolID               string    `json:"school_id" db:"school_id"`
	FirstSeenAt            time.Time `json:"first_seen_at" db:"first_seen_at"`
	LastSeenAt             time.Time `json:"last_seen_at" db:"last_seen_at"`
	TotalQuizSessions      int       `json:"total_quiz_sessions" db:"total_quiz_sessions"`
	TotalQuestionsAnswered int       `json:"total_questions_answered" db:"total_questions_answered"`
	AverageResponseTimeMS  int       `json:"average_response_time_ms" db:"average_response_time_ms"`
	CreatedAt              time.Time `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time `json:"updated_at" db:"updated_at"`
}

type UserRole string

const (
	UserRoleStudent     UserRole = "student"
	UserRoleTeacher     UserRole = "teacher"
	UserRoleSystemAdmin UserRole = "system_admin"
)
