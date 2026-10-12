// Package logger implements a new *slog.Logger struct.
package logger

import (
	"log/slog"
	"os"
)

func New() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}
