package utils

import (
	"log/slog"
	"os"
)

type Logger struct {
	logger *slog.Logger
}

func NewLogger() *Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	return &Logger{ logger: logger }
}

func (l *Logger) Debug(msg string, args... any) {
	l.logger.Debug(msg, args...)
}

func (l *Logger) Warn(msg string, args... any) {
	l.logger.Warn(msg, args...)
}

func (l *Logger) Error(msg string, args... any) {
	l.logger.Error(msg, args...)
}

func (l *Logger) Info(msg string, args... any) {
	l.logger.Info(msg, args...)
}