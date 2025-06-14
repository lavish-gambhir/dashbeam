package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTClaims represents the claims we expect in mobile app JWTs
type JWTClaims struct {
	UserID      uuid.UUID  `json:"user_id"`
	Email       string     `json:"email"`
	Name        string     `json:"name"`
	Role        string     `json:"role"` // student, teacher
	SchoolID    uuid.UUID  `json:"school_id"`
	ClassroomID *uuid.UUID `json:"classroom_id,omitempty"`
	AppType     string     `json:"app_type"` // whiteboard, notebook

	// Standard JWT claims
	jwt.RegisteredClaims
}

// UserContext represents authenticated user context
type UserContext struct {
	UserID      uuid.UUID  `json:"user_id"`
	Email       string     `json:"email"`
	Name        string     `json:"name"`
	Role        string     `json:"role"`
	SchoolID    uuid.UUID  `json:"school_id"`
	ClassroomID *uuid.UUID `json:"classroom_id,omitempty"`
	AppType     string     `json:"app_type"`
	IssuedAt    time.Time  `json:"issued_at"`
	ExpiresAt   time.Time  `json:"expires_at"`
}

// DashboardUser represents a dashboard user (sys-admin)
type DashboardUser struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	Username     string     `json:"username" db:"username"`
	PasswordHash string     `json:"-" db:"password_hash"` // Never serialize password
	FullName     string     `json:"full_name" db:"full_name"`
	Email        string     `json:"email" db:"email"`
	IsActive     bool       `json:"is_active" db:"is_active"`
	LastLoginAt  *time.Time `json:"last_login_at" db:"last_login_at"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

// DashboardSession represents an active dashboard session
type DashboardSession struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Username  string    `json:"username"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// IsValid checks if the user context is valid
func (uc *UserContext) IsValid() bool {
	return uc.UserID != uuid.Nil &&
		uc.SchoolID != uuid.Nil &&
		uc.Email != "" &&
		uc.Role != "" &&
		time.Now().Before(uc.ExpiresAt)
}

// HasClassroomAccess checks if user has access to a classroom
func (uc *UserContext) HasClassroomAccess(classroomID uuid.UUID) bool {
	if uc.ClassroomID == nil {
		return false
	}
	return *uc.ClassroomID == classroomID
}

// IsTeacher checks if the user is a teacher
func (uc *UserContext) IsTeacher() bool {
	return uc.Role == "teacher"
}

// IsStudent checks if the user is a student
func (uc *UserContext) IsStudent() bool {
	return uc.Role == "student"
}
