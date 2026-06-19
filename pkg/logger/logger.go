package logger

import (
	"log/slog"
	"os"
)

type Logger struct {
	*slog.Logger
}

func InitLogger(env string) *Logger {
	var handler slog.Handler

	if env == "production" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	} else {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	}

	baseLogger := slog.New(handler)

	slog.SetDefault(baseLogger)

	return &Logger{Logger: baseLogger}
}
