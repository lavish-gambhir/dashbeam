package logger

import (
	"log/slog"
	"os"
)

func NewSlogger(env string) *slog.Logger {
	level := slog.LevelDebug

	if env == "production" || env == "prod" {
		level = slog.LevelInfo
	}
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	return slog.New(handler)
}
