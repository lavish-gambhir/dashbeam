package repository

import (
	"context"

	"github.com/lavish-gambhir/dashbeam/shared/models"
)

// DashboardRepository defines the interface for dashboard user operations
type DashboardRepository interface {
	// GetUserByUsername retrieves a dashboard user by username
	GetUserByUsername(ctx context.Context, username string) (*models.DashboardUser, error)

	// UpdateLastLogin updates the user's last login timestamp
	UpdateLastLogin(ctx context.Context, userID string) error

	// CreateUser creates a new dashboard user (for initial setup)
	CreateUser(ctx context.Context, user *models.DashboardUser) error
}
