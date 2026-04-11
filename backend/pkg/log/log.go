package log

import (
	"context"
	"log/slog"
)

type Logger struct {
	ctx context.Context
}

func GetLogger(ctx context.Context) *Logger {
	return &Logger{ctx: ctx}
}

func (l *Logger) Error(msg string, err error) {
	slog.ErrorContext(l.ctx, msg, "error", err)
}

func (l *Logger) Info(msg string) {
	slog.InfoContext(l.ctx, msg)
}
