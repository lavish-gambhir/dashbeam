package repository

import (
	"context"

	"github.com/lavish-gambhir/dashbeam/shared/models"
)

type User interface {
	// CreateUser creates a new user
	CreateUser(ctx context.Context, user *models.User) error

	// UpdateUser updates an existing user
	UpdateUser(ctx context.Context, user *models.User) error

	// GetUserByID retrieves a user by ID
	GetUserByID(ctx context.Context, id string) (*models.User, error)

	// UpdateUserLastSeen updates the user's last seen timestamp
	UpdateUserLastSeen(ctx context.Context, userID string) error
}
