package streaming

import (
	"time"

	"github.com/google/uuid"
)

type EventType string

// Quiz event types
const (
	QuizSessionStarted   EventType = "quiz.session.started"
	QuizQuestionShown    EventType = "quiz.question.shown"
	QuizAnswerSubmitted  EventType = "quiz.answer.submitted"
	QuizSessionCompleted EventType = "quiz.session.completed"
	QuizSessionAbandoned EventType = "quiz.session.abandoned"
	QuizSessionPaused    EventType = "quiz.session.paused"
	QuizSessionResumed   EventType = "quiz.session.resumed"
)

// User event types
const (
	UserLogin      EventType = "user.login"
	UserLogout     EventType = "user.logout"
	AppInteraction EventType = "app.interaction"
	AppNavigation  EventType = "app.navigation"
	AppFocusChange EventType = "app.focus.change"
	AppBackground  EventType = "app.background"
	AppForeground  EventType = "app.foreground"
)

// System event types
const (
	APIRequest     EventType = "api.request"
	APIResponse    EventType = "api.response"
	ErrorOccurred  EventType = "error.occurred"
	SystemStartup  EventType = "system.startup"
	SystemShutdown EventType = "system.shutdown"
)

func (e EventType) String() string {
	return string(e)
}

// AppType
type AppType string

const (
	AppTypeWhite AppType = "whiteboard"
	AppTypeNote  AppType = "notebook"
)

func (a AppType) String() string {
	return string(a)
}

// Event - a generic analytics event
type Event struct {
	ID          uuid.UUID      `json:"event_id"`
	Type        EventType      `json:"event_type"`
	Timestamp   time.Time      `json:"timestamp"`
	UserID      uuid.UUID      `json:"user_id"`
	SchoolID    uuid.UUID      `json:"school_id"`
	ClassroomID *uuid.UUID     `json:"classroom_id,omitempty"`
	AppType     AppType        `json:"app_type"`
	Payload     map[string]any `json:"payload"`
	Metadata    Metadata       `json:"metadata"`
}

// Metadata - technical metadata about the event
type Metadata struct {
	AppVersion  string  `json:"app_version"`
	DeviceType  string  `json:"device_type"`
	DeviceID    string  `json:"device_id"`
	NetworkType *string `json:"network_type,omitempty"`
	ClientTime  *string `json:"client_time,omitempty"`
	SessionID   *string `json:"session_id,omitempty"`
	IPAddress   *string `json:"ip_address,omitempty"`
}

func (e *Event) Validate() error {
	if e.ID == uuid.Nil {
		return ErrMissingEventID
	}
	if e.Type == "" {
		return ErrMissingEventType
	}
	if e.UserID == uuid.Nil {
		return ErrMissingUserID
	}
	if e.SchoolID == uuid.Nil {
		return ErrMissingSchoolID
	}
	if e.AppType != AppTypeWhite && e.AppType != AppTypeNote {
		return ErrInvalidAppType
	}
	if e.Metadata.AppVersion == "" {
		return ErrMissingAppVersion
	}
	if e.Metadata.DeviceType == "" {
		return ErrMissingDeviceType
	}
	if e.Metadata.DeviceID == "" {
		return ErrMissingDeviceID
	}
	return nil
}

func (e *Event) IsQuizEvent() bool {
	switch e.Type {
	case QuizSessionStarted, QuizQuestionShown, QuizAnswerSubmitted,
		QuizSessionCompleted, QuizSessionAbandoned, QuizSessionPaused, QuizSessionResumed:
		return true
	default:
		return false
	}
}

func (e *Event) IsUserEvent() bool {
	switch e.Type {
	case UserLogin, UserLogout, AppInteraction, AppNavigation,
		AppFocusChange, AppBackground, AppForeground:
		return true
	default:
		return false
	}
}

func (e *Event) IsSystemEvent() bool {
	switch e.Type {
	case APIRequest, APIResponse, ErrorOccurred, SystemStartup, SystemShutdown:
		return true
	default:
		return false
	}
}

func (e *Event) GetTopic() string {
	return GetTopicForEventType(e.Type)
}
