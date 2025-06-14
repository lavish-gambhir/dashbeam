package streaming

import "errors"

// Event validation errors
var (
	ErrMissingEventID    = errors.New("event ID is required")
	ErrMissingEventType  = errors.New("event type is required")
	ErrMissingUserID     = errors.New("user ID is required")
	ErrMissingSchoolID   = errors.New("school ID is required")
	ErrInvalidAppType    = errors.New("app type must be 'whiteboard' or 'notebook'")
	ErrMissingAppVersion = errors.New("app version is required in metadata")
	ErrMissingDeviceType = errors.New("device type is required in metadata")
	ErrMissingDeviceID   = errors.New("device ID is required in metadata")
	ErrInvalidEventType  = errors.New("invalid event type")
	ErrInvalidTimestamp  = errors.New("invalid timestamp")
	ErrInvalidPayload    = errors.New("invalid event payload")
	ErrEventTooLarge     = errors.New("event payload too large")
)

// Processing errors
var (
	ErrEventProcessing = errors.New("failed to process event")
	ErrTopicNotFound   = errors.New("topic not found")
	ErrPublishFailed   = errors.New("failed to publish event")
	ErrSubscribeFailed = errors.New("failed to subscribe to topic")
	ErrHandlerFailed   = errors.New("event handler failed")
	ErrRetryExhausted  = errors.New("retry attempts exhausted")
)

// Database errors
var (
	ErrDatabaseConnection  = errors.New("database connection failed")
	ErrUserNotFound        = errors.New("user not found")
	ErrUserCreation        = errors.New("failed to create user")
	ErrUserUpdate          = errors.New("failed to update user")
	ErrQuizNotFound        = errors.New("quiz not found")
	ErrParticipantNotFound = errors.New("quiz participant not found")
)
