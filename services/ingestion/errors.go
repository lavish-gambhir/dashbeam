package ingestion

import "errors"

var (
	ErrMissingEventID   = errors.New("event_id is required")
	ErrMissingEventType = errors.New("event_type is required")
	ErrMissingUserID    = errors.New("user_id is required")
	ErrMissingSchoolID  = errors.New("school_id is required")
	ErrInvalidAppType   = errors.New("app_type must be 'white' or 'note'")
	ErrInvalidEventType = errors.New("invalid event_type")
	ErrBatchTooLarge    = errors.New("batch size exceeds maximum limit")
	ErrEventValidation  = errors.New("event validation failed")
	ErrUserCreation     = errors.New("failed to create/update user")
	ErrEventPublishing  = errors.New("failed to publish event")
	ErrInternalServer   = errors.New("internal server error")
)
