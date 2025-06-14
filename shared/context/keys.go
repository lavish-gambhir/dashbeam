package context

import (
	"context"
)

type contextKey string

const (
	userIDKey    contextKey = "userID"
	requestIDKey contextKey = "requestID"
	traceIDKey   contextKey = "traceID"
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
