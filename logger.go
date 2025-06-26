package temu

import (
	"fmt"
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

func (l *logger) Errorf(format string, v ...interface{}) {
	l.l.Error(fmt.Sprintf(format, v...))
}

func (l *logger) Warnf(format string, v ...interface{}) {
	l.l.Warn(fmt.Sprintf(format, v...))
}

func (l *logger) Infof(format string, v ...interface{}) {
	l.l.Info(fmt.Sprintf(format, v...))
}

func (l *logger) Debugf(format string, v ...interface{}) {
	l.l.Debug(fmt.Sprintf(format, v...))
}
