package auth

import (
	"time"

	"github.com/google/uuid"
	"github.com/lavish-gambhir/dashbeam/shared/models"
)

type ValidateJWTRequest struct {
	Token string `json:"token"`
}

type ValidateJWTResponse struct {
	Valid       bool                `json:"valid"`
	UserContext *models.UserContext `json:"user_context,omitempty"`
	Error       string              `json:"error,omitempty"`
}

type DashboardLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type DashboardLoginResponse struct {
	Success     bool                     `json:"success"`
	User        *models.DashboardUser    `json:"user,omitempty"`
	Session     *models.DashboardSession `json:"session,omitempty"`
	AccessToken string                   `json:"access_token,omitempty"`
	ExpiresAt   time.Time                `json:"expires_at,omitempty"`
	Error       string                   `json:"error,omitempty"`
}

type CurrentUserResponse struct {
	User    *models.DashboardUser    `json:"user,omitempty"`
	Session *models.DashboardSession `json:"session,omitempty"`
	Error   string                   `json:"error,omitempty"`
}

type LogoutResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

type dashboardClaims struct {
	UserID    uuid.UUID `json:"user_id"`
	Username  string    `json:"username"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	SessionID uuid.UUID `json:"session_id"`
	IssuedAt  time.Time `json:"iat"`
	ExpiresAt time.Time `json:"exp"`
	NotBefore time.Time `json:"nbf"`
	Issuer    string    `json:"iss"`
	Subject   string    `json:"sub"`
}
