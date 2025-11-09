package p2p

import (
	"log/slog"
	"os"
)

type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
	Debug(msg string, args ...any)
	With(args ...any) Logger
}

type SlogLogger struct {
	l *slog.Logger
}

func NewSlogLogger(level slog.Level) *SlogLogger {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	return &SlogLogger{l: slog.New(handler)}
}

func (s *SlogLogger) Info(msg string, args ...any) {
	s.l.Info(msg, args...)
}

func (s *SlogLogger) Error(msg string, args ...any) {
	s.l.Error(msg, args...)
}

func (s *SlogLogger) Debug(msg string, args ...any) {
	s.l.Debug(msg, args...)
}

func (s *SlogLogger) With(args ...any) Logger {
	return &SlogLogger{l: s.l.With(args...)}
}
