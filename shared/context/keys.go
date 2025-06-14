package context

import (
	"context"

	"github.com/google/uuid"
	"github.com/lavish-gambhir/dashbeam/shared/models"
)

type contextKey string

const (
	userIDKey      contextKey = "user_id"
	requestIDKey   contextKey = "request_id"
	traceIDKey     contextKey = "trace_id"
	userContextKey contextKey = "user_context"
)

func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func GetUserID(ctx context.Context) (string, bool) {
	val := ctx.Value(userIDKey)
	userID, ok := val.(string)
	return userID, ok
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

func GetRequestID(ctx context.Context) (string, bool) {
	val := ctx.Value(requestIDKey)
	requestID, ok := val.(string)
	return requestID, ok
}

func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey, traceID)
}

func GetTraceID(ctx context.Context) (string, bool) {
	val := ctx.Value(traceIDKey)
	requestID, ok := val.(string)
	return requestID, ok
}

// WithUserContext adds authenticated user context to the request context
func WithUserContext(ctx context.Context, userContext *models.UserContext) context.Context {
	return context.WithValue(ctx, userContextKey, userContext)
}

// GetUserContext retrieves the authenticated user context from request context
func GetUserContext(ctx context.Context) (*models.UserContext, bool) {
	userContext, ok := ctx.Value(userContextKey).(*models.UserContext)
	return userContext, ok
}

// GetAuthenticatedUserID is a convenience function to get the authenticated user ID from context
func GetAuthenticatedUserID(ctx context.Context) (uuid.UUID, bool) {
	userContext, ok := GetUserContext(ctx)
	if !ok || userContext == nil {
		return uuid.Nil, false
	}
	return userContext.UserID, true
}
