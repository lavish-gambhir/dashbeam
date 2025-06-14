package streaming

const (
	TopicQuizEvents       = "quiz-events"
	TopicUserEvents       = "user-events"
	TopicEngagementEvents = "engagement-events"
	TopicSystemEvents     = "system-events"
)

func GetTopicForEventType(eventType EventType) string {
	switch eventType {
	case QuizSessionStarted, QuizQuestionShown, QuizAnswerSubmitted, QuizSessionCompleted, QuizSessionAbandoned, QuizSessionPaused, QuizSessionResumed:
		return TopicQuizEvents
	case UserLogin, UserLogout:
		return TopicUserEvents
	case AppInteraction, AppNavigation, AppFocusChange, AppBackground, AppForeground:
		return TopicEngagementEvents
	default:
		return TopicSystemEvents
	}
}
