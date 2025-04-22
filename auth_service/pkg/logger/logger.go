package logger

import (
	"fmt"

	"go.uber.org/zap"
)

type Logger struct {
	L *zap.Logger
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.L.Info(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.L.Error(msg, fields...)
}

func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.L.Fatal(msg, fields...)
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.L.Debug(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.L.Warn(msg, fields...)
}

func Setup() (*Logger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("failed to setup logger: %w", err)
	}

	return &Logger{L: logger}, nil
}
