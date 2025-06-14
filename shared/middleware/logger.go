package middleware

import (
	"log"
	"log/slog"
	"net/http"
	"time"

	sharedcontext "github.com/lavish-gambhir/dashbeam/shared/context"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"go.opentelemetry.io/otel/trace"
)

type responseWriter struct {
	http.ResponseWriter
	status int
	length int
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.status == 0 {
		rw.status = http.StatusOK
	}
	n, err := rw.ResponseWriter.Write(b)
	rw.length += n
	return n, err
}

func Logging(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID, err := gonanoid.New()
			if err != nil {
				log.Fatal(err)
			}

			ctx := r.Context()
			spanCtx := trace.SpanContextFromContext(ctx)
			traceID := spanCtx.TraceID().String()

			// TODO: Set up `TraceProvider` for tracing.

			// Add request ID to the context
			r = r.WithContext(sharedcontext.WithRequestID(ctx, requestID))
			// Add trace ID to the context
			r = r.WithContext(sharedcontext.WithTraceID(ctx, traceID))

			// Create a custom response writer to capture status code
			rw := &responseWriter{ResponseWriter: w}

			// Pre-request logging
			logger.LogAttrs(ctx, slog.LevelInfo, "Request started",
				slog.String("method", r.Method),
				slog.String("url", r.URL.String()),
				slog.String("requestID", requestID),
				slog.String("traceID", traceID),
				slog.String("userAgent", r.UserAgent()),
				slog.String("remoteAddr", r.RemoteAddr),
			)

			start := time.Now()

			next.ServeHTTP(rw, r)

			duration := time.Since(start)

			logger.LogAttrs(ctx, slog.LevelInfo, "Request completed",
				slog.String("method", r.Method),
				slog.String("url", r.URL.String()),
				slog.String("requestID", requestID),
				slog.String("traceID", traceID),
				slog.Int("status", rw.status),
				slog.Int("responseSize", rw.length),
				slog.Duration("duration", duration),
			)
		})
	}
}
