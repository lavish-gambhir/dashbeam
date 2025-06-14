package streaming

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	sharedutil "github.com/lavish-gambhir/dashbeam/shared"
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
	ID          uuid.UUID    `json:"event_id"`
	Type        EventType    `json:"event_type"`
	Timestamp   time.Time    `json:"timestamp"`
	UserID      uuid.UUID    `json:"user_id"`
	SchoolID    uuid.UUID    `json:"school_id"`
	ClassroomID *uuid.UUID   `json:"classroom_id,omitempty"`
	AppType     AppType      `json:"app_type"`
	Payload     EventPayload `json:"payload"`
	Metadata    Metadata     `json:"metadata"`
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

// eventJSON is used for custom JSON marshaling/unmarshaling
type eventJSON struct {
	ID          string                 `json:"event_id"`
	Type        EventType              `json:"event_type"`
	Timestamp   time.Time              `json:"timestamp"`
	UserID      string                 `json:"user_id"`
	SchoolID    string                 `json:"school_id"`
	ClassroomID *string                `json:"classroom_id,omitempty"`
	AppType     AppType                `json:"app_type"`
	Payload     map[string]interface{} `json:"payload"`
	Metadata    Metadata               `json:"metadata"`
}

func (e Event) MarshalJSON() ([]byte, error) {
	// Convert payload to map
	payloadBytes, err := json.Marshal(e.Payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	var payloadMap map[string]interface{}
	if err := json.Unmarshal(payloadBytes, &payloadMap); err != nil {
		return nil, fmt.Errorf("failed to convert payload to map: %w", err)
	}

	// Convert UUIDs to strings and prepare classroom ID
	var classroomIDStr *string
	if e.ClassroomID != nil {
		str := e.ClassroomID.String()
		classroomIDStr = &str
	}

	eventData := eventJSON{
		ID:          e.ID.String(),
		Type:        e.Type,
		Timestamp:   e.Timestamp,
		UserID:      e.UserID.String(),
		SchoolID:    e.SchoolID.String(),
		ClassroomID: classroomIDStr,
		AppType:     e.AppType,
		Payload:     payloadMap,
		Metadata:    e.Metadata,
	}

	return json.Marshal(eventData)
}

func (e *Event) UnmarshalJSON(data []byte) error {
	var eventData eventJSON
	if err := json.Unmarshal(data, &eventData); err != nil {
		return fmt.Errorf("failed to unmarshal event data: %w", err)
	}

	// Parse UUIDs
	eventID, err := sharedutil.ParseUUID(eventData.ID)
	if err != nil && eventData.ID != "" { // Allow empty event_id
		return fmt.Errorf("invalid event_id: %w", err)
	}

	userID, err := sharedutil.ParseUUID(eventData.UserID)
	if err != nil {
		return fmt.Errorf("invalid user_id: %w", err)
	}

	schoolID, err := sharedutil.ParseUUID(eventData.SchoolID)
	if err != nil {
		return fmt.Errorf("invalid school_id: %w", err)
	}

	var classroomID *uuid.UUID
	if eventData.ClassroomID != nil {
		cid, err := sharedutil.ParseUUID(*eventData.ClassroomID)
		if err != nil {
			return fmt.Errorf("invalid classroom_id: %w", err)
		}
		classroomID = &cid
	}

	// Parse payload based on event type
	payload, err := PayloadFromMap(eventData.Type, eventData.Payload)
	if err != nil {
		return fmt.Errorf("failed to parse payload for event type %s: %w", eventData.Type, err)
	}

	// Set all fields
	e.ID = eventID
	e.Type = eventData.Type
	e.Timestamp = eventData.Timestamp
	e.UserID = userID
	e.SchoolID = schoolID
	e.ClassroomID = classroomID
	e.AppType = eventData.AppType
	e.Payload = payload
	e.Metadata = eventData.Metadata

	return nil
}
