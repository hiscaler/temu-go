package temu

import (
	"fmt"
	"log/slog"
	"strings"
)

type Logger interface {
	Errorf(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Debugf(format string, v ...interface{})
}

func createLogger() *logger {
	return &logger{l: slog.Default()}
}

type logger struct {
	l *slog.Logger
}

var _ Logger = (*logger)(nil)

func (l *logger) Errorf(msg string, args ...interface{}) {
	if strings.Contains(msg, "%") {
		l.l.Error(fmt.Sprintf(msg, args...))
		return
	}
	l.l.Error(msg, args...)
}

func (l *logger) Warnf(msg string, args ...interface{}) {
	if strings.Contains(msg, "%") {
		l.l.Warn(fmt.Sprintf(msg, args...))
		return
	}
	l.l.Warn(msg, args...)
}

func (l *logger) Infof(msg string, args ...interface{}) {
	if strings.Contains(msg, "%") {
		l.l.Info(fmt.Sprintf(msg, args...))
		return
	}
	l.l.Info(msg, args...)
}

func (l *logger) Debugf(msg string, args ...interface{}) {
	if strings.Contains(msg, "%") {
		l.l.Debug(fmt.Sprintf(msg, args...))
		return
	}
	l.l.Debug(msg, args...)
}
