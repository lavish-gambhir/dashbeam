package ingestion

import (
	"log/slog"
	"net/http"

	"github.com/lavish-gambhir/dashbeam/services/ingestion/repository"
	"github.com/lavish-gambhir/dashbeam/shared/streaming"
)

type Service interface {
	RegisterRoutes(mux *http.ServeMux, prefix string)
}

type service struct {
	userRepo     repository.User
	quizRepo     repository.Quiz
	messageQueue streaming.MessageQueue
	logger       *slog.Logger
	maxBatchSize int
}

func New(userRepo repository.User, quizRepo repository.Quiz, messageQueue streaming.MessageQueue, logger *slog.Logger, maxBatchSize int) Service {
	return &service{
		userRepo:     userRepo,
		quizRepo:     quizRepo,
		messageQueue: messageQueue,
		logger:       logger,
		maxBatchSize: maxBatchSize,
	}
}

// RegisterRoutes registers all ingestion service routes
func (s *service) RegisterRoutes(parentmux *http.ServeMux, prefix string) {
	h := &handler{
		userRepo:     s.userRepo,
		quizRepo:     s.quizRepo,
		messageQueue: s.messageQueue,
		logger:       s.logger,
		maxBatchSize: s.maxBatchSize,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/batch", h.handleBatchEvents)
	mux.HandleFunc("/quiz", h.handleQuizEvent)
	mux.HandleFunc("/user", h.handleUserEvent)
	mux.HandleFunc("/system", h.handleSystemEvent)
	parentmux.Handle(prefix+"/", http.StripPrefix(prefix, mux))
}
