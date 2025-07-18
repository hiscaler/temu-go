package temu

import (
	"log/slog"
)

type Logger interface {
	Errorf(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Debugf(format string, v ...interface{})
}

func createLogger() *logger {
	ll := slog.Default()
	l := &logger{l: ll}
	return l
}

type logger struct {
	l *slog.Logger
}

var _ Logger = (*logger)(nil)

func (l *logger) Errorf(msg string, args ...interface{}) {
	l.l.Error(msg, args...)
}

func (l *logger) Warnf(msg string, args ...interface{}) {
	l.l.Warn(msg, args...)
}

func (l *logger) Infof(msg string, args ...interface{}) {
	l.l.Info(msg, args...)
}

func (l *logger) Debugf(msg string, args ...interface{}) {
	l.l.Debug(msg, args...)
}
