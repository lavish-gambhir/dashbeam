package ingestion

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/lavish-gambhir/dashbeam/pkg/apperr"
	"github.com/lavish-gambhir/dashbeam/pkg/utils"
	"github.com/lavish-gambhir/dashbeam/services/ingestion/repository"
	sharedcontext "github.com/lavish-gambhir/dashbeam/shared/context"
	"github.com/lavish-gambhir/dashbeam/shared/streaming"
)

type handler struct {
	userRepo     repository.User
	quizRepo     repository.Quiz
	messageQueue streaming.MessageQueue
	logger       *slog.Logger
	maxBatchSize int
}

func NewHandler(userRepo repository.User, quizRepo repository.Quiz, messageQueue streaming.MessageQueue, logger *slog.Logger, maxBatchSize int) *handler {
	log := logger.With("handler", "ingestion.handler")
	return &handler{
		userRepo:     userRepo,
		quizRepo:     quizRepo,
		messageQueue: messageQueue,
		logger:       log,
		maxBatchSize: maxBatchSize,
	}
}

func (h *handler) handleBatchEvents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqID, _ := sharedcontext.GetRequestID(ctx)
	logger := h.logger.With("fn", "handleBatchEvents").With("requestID", reqID)
	if r.Method != http.MethodPost {
		utils.WriteJSONError(w, apperr.New(apperr.BadRequest, "method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	userContext, ok := sharedcontext.GetUserContext(ctx)
	if !ok {
		logger.Error("user context not found - middleware not applied correctly")
		utils.WriteJSONError(w, apperr.New(apperr.Unauthorized, "authentication context missing"), http.StatusUnauthorized)
		return
	}
	logger = logger.With("userID", userContext.UserID.String()).With("schoolID", userContext.SchoolID.String())

	var req BatchEventsRequest

	if err := utils.FromJson(r.Body, &req); err != nil {
		utils.WriteJSONError(w, apperr.Wrap(err, apperr.BadRequest, "invalid request body"), http.StatusBadRequest)
		return
	}

	if len(req.Events) == 0 {
		utils.WriteJSONError(w, apperr.New(apperr.BadRequest, "no events provided"), http.StatusBadRequest)
		return
	}

	if len(req.Events) > h.maxBatchSize {
		utils.WriteJSONError(w, apperr.Newf(apperr.BadRequest, "batch size %d exceeds maximum %d", len(req.Events), h.maxBatchSize), http.StatusRequestEntityTooLarge)
		return
	}

	logger.Info("processing batch events", slog.Int("count", len(req.Events)))
	eventIDs, err := h.processBatchEvents(ctx, req.Events)
	if err != nil {
		h.logger.Error("failed to process batch events", slog.Any("error", err), slog.Int("count", len(req.Events)))
		utils.WriteJSONError(w, err, http.StatusInternalServerError)
		return
	}

	utils.WriteJSONSuccess(w, EventResponse{
		Status:    "success",
		EventIDs:  eventIDs,
		Processed: len(req.Events),
		Timestamp: time.Now().UTC(),
	})
}

func (h *handler) handleQuizEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteJSONError(w, apperr.New(apperr.BadRequest, "method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	var req SingleEventRequest

	if err := utils.FromJson(r.Body, &req); err != nil {
		utils.WriteJSONError(w, apperr.Wrap(err, apperr.BadRequest, "invalid request body"), http.StatusBadRequest)
		return
	}

	if !req.Event.IsQuizEvent() {
		utils.WriteJSONError(w, apperr.New(apperr.BadRequest, "not a quiz event"), http.StatusBadRequest)
		return
	}

	eventID, err := h.processSingleEvent(ctx, req.Event)
	if err != nil {
		h.logger.Error("failed to process quiz event", slog.Any("error", err), slog.String("event_type", req.Event.Type.String()))
		utils.WriteJSONError(w, err, http.StatusInternalServerError)
		return
	}

	utils.WriteJSONSuccess(w, EventResponse{
		Status:    "success",
		EventIDs:  []string{eventID},
		Processed: 1,
		Timestamp: time.Now().UTC(),
	})
}

func (h *handler) handleUserEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteJSONError(w, apperr.New(apperr.BadRequest, "method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	var req SingleEventRequest

	if err := utils.FromJson(r.Body, &req); err != nil {
		utils.WriteJSONError(w, apperr.Wrap(err, apperr.BadRequest, "invalid request body"), http.StatusBadRequest)
		return
	}

	if !req.Event.IsUserEvent() {
		utils.WriteJSONError(w, apperr.New(apperr.BadRequest, "not a user event"), http.StatusBadRequest)
		return
	}

	eventID, err := h.processSingleEvent(ctx, req.Event)
	if err != nil {
		h.logger.Error("failed to process user event", slog.Any("error", err), slog.String("event_type", req.Event.Type.String()))
		utils.WriteJSONError(w, err, http.StatusInternalServerError)
		return
	}

	utils.WriteJSONSuccess(w, EventResponse{
		Status:    "success",
		EventIDs:  []string{eventID},
		Processed: 1,
		Timestamp: time.Now().UTC(),
	})
}

func (h *handler) handleSystemEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteJSONError(w, apperr.New(apperr.BadRequest, "method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	var req SingleEventRequest

	if err := utils.FromJson(r.Body, &req); err != nil {
		utils.WriteJSONError(w, apperr.Wrap(err, apperr.BadRequest, "invalid request body"), http.StatusBadRequest)
		return
	}

	if !req.Event.IsSystemEvent() {
		utils.WriteJSONError(w, apperr.New(apperr.BadRequest, "not a system event"), http.StatusBadRequest)
		return
	}

	eventID, err := h.processSingleEvent(ctx, req.Event)
	if err != nil {
		h.logger.Error("failed to process system event", slog.Any("error", err), slog.String("event_type", req.Event.Type.String()))
		utils.WriteJSONError(w, err, http.StatusInternalServerError)
		return
	}

	utils.WriteJSONSuccess(w, EventResponse{
		Status:    "success",
		EventIDs:  []string{eventID},
		Processed: 1,
		Timestamp: time.Now().UTC(),
	})
}
