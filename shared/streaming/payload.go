package streaming

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type EventPayload interface {
	Type() string
	Validate() error
}

type QuizSessionStartedPayload struct {
	QuizID         uuid.UUID `json:"quiz_id"`
	SessionID      uuid.UUID `json:"session_id"`
	SessionCode    string    `json:"session_code"`
	TotalQuestions int       `json:"total_questions"`
	TimeLimit      *int      `json:"time_limit_seconds,omitempty"`
	MaxScore       float64   `json:"max_score"`
}

func (q QuizSessionStartedPayload) Type() string { return QuizSessionStarted.String() }
func (p QuizSessionStartedPayload) Validate() error {
	if p.QuizID == uuid.Nil {
		return ErrInvalidPayload
	}
	if p.SessionID == uuid.Nil {
		return ErrInvalidPayload
	}
	if p.SessionCode == "" {
		return ErrInvalidPayload
	}
	if p.TotalQuestions <= 0 {
		return ErrInvalidPayload
	}
	return nil
}

type QuizQuestionShownPayload struct {
	QuizID           uuid.UUID `json:"quiz_id"`
	SessionID        uuid.UUID `json:"session_id"`
	QuestionID       uuid.UUID `json:"question_id"`
	QuestionSequence int       `json:"question_sequence"`
	QuestionType     string    `json:"question_type"`
	TimeLimit        *int      `json:"time_limit_seconds,omitempty"`
}

func (q QuizQuestionShownPayload) Type() string { return QuizQuestionShown.String() }
func (p QuizQuestionShownPayload) Validate() error {
	if p.QuizID == uuid.Nil || p.SessionID == uuid.Nil || p.QuestionID == uuid.Nil {
		return ErrInvalidPayload
	}
	if p.QuestionSequence <= 0 {
		return ErrInvalidPayload
	}
	return nil
}

type QuizAnswerSubmittedPayload struct {
	QuizID           uuid.UUID       `json:"quiz_id"`
	SessionID        uuid.UUID       `json:"session_id"`
	QuestionID       uuid.UUID       `json:"question_id"`
	QuestionSequence int             `json:"question_sequence"`
	Answer           json.RawMessage `json:"answer"`
	IsCorrect        *bool           `json:"is_correct,omitempty"`
	ResponseTimeMS   int             `json:"response_time_ms"`
	AnswerChanges    *int            `json:"answer_changes,omitempty"`
	Points           *float64        `json:"points,omitempty"`
}

func (q QuizAnswerSubmittedPayload) Type() string { return QuizAnswerSubmitted.String() }
func (p QuizAnswerSubmittedPayload) Validate() error {
	if p.QuizID == uuid.Nil || p.SessionID == uuid.Nil || p.QuestionID == uuid.Nil {
		return ErrInvalidPayload
	}
	if p.QuestionSequence <= 0 {
		return ErrInvalidPayload
	}
	if p.ResponseTimeMS < 0 {
		return ErrInvalidPayload
	}
	return nil
}

type QuizSessionCompletedPayload struct {
	QuizID              uuid.UUID `json:"quiz_id"`
	SessionID           uuid.UUID `json:"session_id"`
	TotalScore          float64   `json:"total_score"`
	MaxScore            float64   `json:"max_score"`
	CompletionTimeMS    int       `json:"completion_time_ms"`
	QuestionsCorrect    int       `json:"questions_correct"`
	QuestionsAnswered   int       `json:"questions_answered"`
	QuestionsSkipped    int       `json:"questions_skipped"`
	AverageResponseTime int       `json:"average_response_time_ms"`
}

func (p QuizSessionCompletedPayload) Type() string { return QuizSessionCompleted.String() }
func (p QuizSessionCompletedPayload) Validate() error {
	if p.QuizID == uuid.Nil || p.SessionID == uuid.Nil {
		return ErrInvalidPayload
	}
	if p.MaxScore < 0 || p.TotalScore < 0 {
		return ErrInvalidPayload
	}
	if p.CompletionTimeMS < 0 {
		return ErrInvalidPayload
	}
	return nil
}

// User Event Payloads

type UserLoginPayload struct {
	LoginMethod   string     `json:"login_method"`
	SessionStart  bool       `json:"session_start"`
	PreviousLogin *time.Time `json:"previous_login_timestamp,omitempty"`
	UserAgent     *string    `json:"user_agent,omitempty"`
	IPAddress     *string    `json:"ip_address,omitempty"`
	Email         *string    `json:"email,omitempty"`
	Name          *string    `json:"name,omitempty"`
	Role          *string    `json:"role,omitempty"`
}

func (p UserLoginPayload) Type() string { return UserLogin.String() }
func (p UserLoginPayload) Validate() error {
	if p.LoginMethod == "" {
		return ErrInvalidPayload
	}
	return nil
}

type UserLogoutPayload struct {
	SessionDuration int    `json:"session_duration_ms"`
	LogoutReason    string `json:"logout_reason"` // e.g. manual, timeout, forced
}

func (p UserLogoutPayload) Type() string { return UserLogout.String() }
func (p UserLogoutPayload) Validate() error {
	if p.SessionDuration < 0 {
		return ErrInvalidPayload
	}
	return nil
}

type AppInteractionPayload struct {
	InteractionType string  `json:"interaction_type"`
	ScreenName      string  `json:"screen_name"`
	ElementClicked  *string `json:"element_clicked,omitempty"`
	TimeSpentMS     *int    `json:"time_spent_ms,omitempty"`
	ElementData     *string `json:"element_data,omitempty"`
	Action          *string `json:"action,omitempty"` // eg. tap, swipe, long_press, etc.
}

func (p AppInteractionPayload) Type() string { return AppInteraction.String() }
func (p AppInteractionPayload) Validate() error {
	if p.InteractionType == "" || p.ScreenName == "" {
		return ErrInvalidPayload
	}
	return nil
}

type AppNavigationPayload struct {
	FromScreen     string  `json:"from_screen"`
	ToScreen       string  `json:"to_screen"`
	NavigationType string  `json:"navigation_type"` // eg. push, pop, replace
	TimeSpentMS    *int    `json:"time_spent_ms,omitempty"`
	NavigationData *string `json:"navigation_data,omitempty"`
}

func (p AppNavigationPayload) Type() string { return AppNavigation.String() }
func (p AppNavigationPayload) Validate() error {
	if p.FromScreen == "" || p.ToScreen == "" || p.NavigationType == "" {
		return ErrInvalidPayload
	}
	return nil
}

// System Event Payloads

// APIRequestPayload contains data for API request events
type APIRequestPayload struct {
	Method       string  `json:"method"`
	Endpoint     string  `json:"endpoint"`
	StatusCode   *int    `json:"status_code,omitempty"`
	ResponseTime *int    `json:"response_time_ms,omitempty"`
	ErrorCode    *string `json:"error_code,omitempty"`
	RequestSize  *int    `json:"request_size_bytes,omitempty"`
	ResponseSize *int    `json:"response_size_bytes,omitempty"`
}

func (p APIRequestPayload) Type() string { return APIRequest.String() }
func (p APIRequestPayload) Validate() error {
	if p.Method == "" || p.Endpoint == "" {
		return ErrInvalidPayload
	}
	return nil
}

type ErrorOccurredPayload struct {
	ErrorType    string  `json:"error_type"`
	ErrorMessage string  `json:"error_message"`
	ErrorCode    *string `json:"error_code,omitempty"`
	StackTrace   *string `json:"stack_trace,omitempty"`
	Context      *string `json:"context,omitempty"`
	Severity     string  `json:"severity"` // low, medium, high, critical
}

func (p ErrorOccurredPayload) Type() string { return ErrorOccurred.String() }
func (p ErrorOccurredPayload) Validate() error {
	if p.ErrorType == "" || p.ErrorMessage == "" || p.Severity == "" {
		return ErrInvalidPayload
	}
	return nil
}

// PayloadFromMap converts a map to a specific payload type based on event type
func PayloadFromMap(eventType EventType, data map[string]any) (EventPayload, error) {
	if data == nil {
		return nil, fmt.Errorf("payload data is empty")
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload data: %w", err)
	}

	switch eventType {
	case QuizSessionStarted:
		var payload QuizSessionStartedPayload
		if err := json.Unmarshal(jsonData, &payload); err != nil {
			return nil, fmt.Errorf("failed to unmarshal QuizSessionStartedPayload: %w", err)
		}
		if err := payload.Validate(); err != nil {
			return nil, fmt.Errorf("QuizSessionStartedPayload validation failed: %w", err)
		}
		return payload, payload.Validate()

	case QuizQuestionShown:
		var payload QuizQuestionShownPayload
		if err := json.Unmarshal(jsonData, &payload); err != nil {
			return nil, fmt.Errorf("failed to unmarshal QuizQuestionShownPayload: %w", err)
		}
		if err := payload.Validate(); err != nil {
			return nil, fmt.Errorf("QuizQuestionShownPayload validation failed: %w", err)
		}
		return payload, payload.Validate()

	case QuizAnswerSubmitted:
		var payload QuizAnswerSubmittedPayload
		if err := json.Unmarshal(jsonData, &payload); err != nil {
			return nil, fmt.Errorf("failed to unmarshal QuizAnswerSubmittedPayload: %w", err)
		}
		if err := payload.Validate(); err != nil {
			return nil, fmt.Errorf("QuizAnswerSubmittedPayload validation failed: %w", err)
		}
		return payload, payload.Validate()

	case QuizSessionCompleted:
		var payload QuizSessionCompletedPayload
		if err := json.Unmarshal(jsonData, &payload); err != nil {
			return nil, fmt.Errorf("failed to unmarshal QuizSessionCompletedPayload: %w", err)
		}
		if err := payload.Validate(); err != nil {
			return nil, fmt.Errorf("QuizSessionCompletedPayload validation failed: %w", err)
		}
		return payload, payload.Validate()

	case UserLogin:
		var payload UserLoginPayload
		if err := json.Unmarshal(jsonData, &payload); err != nil {
			return nil, fmt.Errorf("failed to unmarshal UserLoginPayload: %w", err)
		}
		if err := payload.Validate(); err != nil {
			return nil, fmt.Errorf("UserLoginPayload validation failed: %w", err)
		}
		return payload, payload.Validate()

	case UserLogout:
		var payload UserLogoutPayload
		if err := json.Unmarshal(jsonData, &payload); err != nil {
			return nil, fmt.Errorf("failed to unmarshal UserLogoutPayload: %w", err)
		}
		if err := payload.Validate(); err != nil {
			return nil, fmt.Errorf("UserLogoutPayload validation failed: %w", err)
		}
		return payload, payload.Validate()

	case AppInteraction:
		var payload AppInteractionPayload
		if err := json.Unmarshal(jsonData, &payload); err != nil {
			return nil, fmt.Errorf("failed to unmarshal AppInteractionPayload: %w", err)
		}
		if err := payload.Validate(); err != nil {
			return nil, fmt.Errorf("AppInteractionPayload validation failed: %w", err)
		}
		return payload, payload.Validate()

	case AppNavigation:
		var payload AppNavigationPayload
		if err := json.Unmarshal(jsonData, &payload); err != nil {
			return nil, fmt.Errorf("failed to unmarshal AppNavigationPayload: %w", err)
		}
		if err := payload.Validate(); err != nil {
			return nil, fmt.Errorf("AppNaviationPayload validation failed: %w", err)
		}
		return payload, payload.Validate()

	default:
		return nil, fmt.Errorf("unknown event type: %s", eventType)
	}
}
