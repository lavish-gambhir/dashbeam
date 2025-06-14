package repositories

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/lavish-gambhir/dashbeam/pkg/apperr"
	"github.com/lavish-gambhir/dashbeam/shared/database/postgres"
	"github.com/lavish-gambhir/dashbeam/shared/models"
)

type DashboardUserRepository struct {
	db *postgres.DB
}

func NewDashboardUserRepository(db *postgres.DB) *DashboardUserRepository {
	return &DashboardUserRepository{
		db: db,
	}
}

func (r *DashboardUserRepository) GetUserByUsername(ctx context.Context, username string) (*models.DashboardUser, error) {
	query := `
		SELECT
			id, username, password_hash, full_name, email, role, school_access, permissions,
			is_active, last_login_at, created_at, updated_at
		FROM dashboard_users
		WHERE username = $1`

	var user models.DashboardUser
	err := r.db.Conn(ctx).QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.FullName,
		&user.Email,
		&user.Role,
		&user.SchoolAccess,
		&user.Permissions,
		&user.IsActive,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperr.Newf(apperr.DBRecordNotFound, "dashboard user not found with username: %s", username)
		}
		return nil, apperr.Wrapf(err, apperr.DBQueryFailed, "failed to get dashboard user by username: %s", username)
	}

	return &user, nil
}

func (r *DashboardUserRepository) UpdateLastLogin(ctx context.Context, userID string) error {
	query := `
		UPDATE dashboard_users SET
			last_login_at = $2,
			updated_at = $2
		WHERE id = $1`

	now := time.Now().UTC()
	result, err := r.db.Conn(ctx).Exec(ctx, query, userID, now)

	if err != nil {
		return apperr.Wrapf(err, apperr.DBQueryFailed, "failed to update last login for dashboard user: %s", userID)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return apperr.Newf(apperr.DBRecordNotFound, "dashboard user not found with ID: %s", userID)
	}

	return nil
}

func (r *DashboardUserRepository) CreateUser(ctx context.Context, user *models.DashboardUser) error {
	// Hash the password if it's not already hashed
	if user.PasswordHash == "" {
		return apperr.New(apperr.BadRequest, "password hash is required")
	}

	query := `
		INSERT INTO dashboard_users (
			id, username, password_hash, full_name, email, is_active,
			last_login_at, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		)`

	_, err := r.db.Conn(ctx).Exec(ctx, query,
		user.ID,
		user.Username,
		user.PasswordHash,
		user.FullName,
		user.Email,
		user.IsActive,
		user.LastLoginAt,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		return apperr.Wrapf(err, apperr.DBQueryFailed, "failed to create dashboard user with username: %s", user.Username)
	}

	return nil
}

func (r *DashboardUserRepository) GetUserByID(ctx context.Context, userID string) (*models.DashboardUser, error) {
	query := `
		SELECT
			id, username, password_hash, full_name, email, is_active,
			last_login_at, created_at, updated_at
		FROM dashboard_users
		WHERE id = $1`

	var user models.DashboardUser
	err := r.db.Conn(ctx).QueryRow(ctx, query, userID).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.FullName,
		&user.Email,
		&user.IsActive,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperr.Newf(apperr.DBRecordNotFound, "dashboard user not found with ID: %s", userID)
		}
		return nil, apperr.Wrapf(err, apperr.DBQueryFailed, "failed to get dashboard user by ID: %s", userID)
	}

	return &user, nil
}

func (r *DashboardUserRepository) UpdatePassword(ctx context.Context, userID string, hashedPassword string) error {
	query := `
		UPDATE dashboard_users SET
			password_hash = $2,
			updated_at = $3
		WHERE id = $1`

	now := time.Now().UTC()
	result, err := r.db.Conn(ctx).Exec(ctx, query, userID, hashedPassword, now)

	if err != nil {
		return apperr.Wrapf(err, apperr.DBQueryFailed, "failed to update password for dashboard user: %s", userID)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return apperr.Newf(apperr.DBRecordNotFound, "dashboard user not found with ID: %s", userID)
	}

	return nil
}

func (r *DashboardUserRepository) SetUserActiveStatus(ctx context.Context, userID string, isActive bool) error {
	query := `
		UPDATE dashboard_users SET
			is_active = $2,
			updated_at = $3
		WHERE id = $1`

	now := time.Now().UTC()
	result, err := r.db.Conn(ctx).Exec(ctx, query, userID, isActive, now)

	if err != nil {
		return apperr.Wrapf(err, apperr.DBQueryFailed, "failed to update active status for dashboard user: %s", userID)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return apperr.Newf(apperr.DBRecordNotFound, "dashboard user not found with ID: %s", userID)
	}

	return nil
}

func (r *DashboardUserRepository) ListUsers(ctx context.Context, limit, offset int) ([]*models.DashboardUser, error) {
	query := `
		SELECT
			id, username, password_hash, full_name, email, is_active,
			last_login_at, created_at, updated_at
		FROM dashboard_users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.Conn(ctx).Query(ctx, query, limit, offset)
	if err != nil {
		return nil, apperr.Wrap(err, apperr.DBQueryFailed, "failed to list dashboard users")
	}
	defer rows.Close()

	var users []*models.DashboardUser
	for rows.Next() {
		var user models.DashboardUser
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.PasswordHash,
			&user.FullName,
			&user.Email,
			&user.IsActive,
			&user.LastLoginAt,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, apperr.Wrap(err, apperr.DBQueryFailed, "failed to scan dashboard user row")
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, apperr.Wrap(err, apperr.DBQueryFailed, "error iterating dashboard user rows")
	}

	return users, nil
}

func (r *DashboardUserRepository) HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", apperr.Wrap(err, apperr.Internal, "failed to hash password")
	}
	return string(hashedBytes), nil
}

func (r *DashboardUserRepository) CreateUserWithPassword(ctx context.Context, user *models.DashboardUser, password string) error {
	hashedPassword, err := r.HashPassword(password)
	if err != nil {
		return err
	}

	user.PasswordHash = hashedPassword

	return r.CreateUser(ctx, user)
}
