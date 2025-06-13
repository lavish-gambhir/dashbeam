package logger

import (
	"log/slog"
	"os"

	"github.com/lavish-gambhir/dashbeam/internal/config"
)

func NewSlogger(cnf *config.AppConfig) *slog.Logger {
	level := slog.LevelDebug
	if cnf.Env == config.Production {
		level = slog.LevelInfo
	}
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	return slog.New(handler)
}
