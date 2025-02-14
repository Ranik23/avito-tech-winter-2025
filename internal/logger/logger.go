package logger

import (
	"context"
	"log/slog"
	"os"
	"github.com/lmittmann/tint"
)

type Logger struct {
	slog.Logger
}


func NewLogger(level string) *Logger {
	var lvl slog.Level
	switch level {
	case "info":
		lvl = slog.LevelInfo
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	case "debug":
		lvl = slog.LevelDebug
	default:
		lvl = slog.LevelInfo
	}

	return &Logger{
		Logger: *slog.New(tint.NewHandler(os.Stdout, &tint.Options{
			Level: lvl,
		})),
	}
}

func (l *Logger) InfoWithAttrs(ctx context.Context, level slog.Level, message string, attrs ...slog.Attr) {
	l.Logger.LogAttrs(ctx, level, message, attrs...)
}
