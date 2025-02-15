package logger_test

import (
	"avito/internal/logger"
	"log/slog"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {
	
	tests := []struct {
		name string
		level string
		exp slog.Level
	}{
		{"info level", "info", slog.LevelInfo},
		{"warn level", "warn", slog.LevelWarn},
		{"error level", "error", slog.LevelError},
		{"debug level", "debug", slog.LevelDebug},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := logger.NewLogger(tt.level)
			assert.NotNil(t, log)
		})
	}
}