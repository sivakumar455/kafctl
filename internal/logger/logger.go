package logger

import "log/slog"

type Logger interface {
	Info(msg string, args ...any)
	Debug(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

var log Logger

func init() {
	log = slog.Default()
}

// Exported functions for logging
func Info(msg string, args ...any) {
	log.Info(msg, args...)
}

func Debug(msg string, args ...any) {
	log.Debug(msg, args...)
}

func Warn(msg string, args ...any) {
	log.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	log.Error(msg, args...)
}
