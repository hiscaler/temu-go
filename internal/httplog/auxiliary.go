// Package httplog configures call logging and auxiliary HTTP diagnostics.
package httplog

import (
	"log/slog"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/hiscaler/temu-go/config"
)

// IsDebug reports whether the effective call log level is debug.
func IsDebug(cfg config.Config) bool {
	level := strings.TrimSpace(cfg.LogLevel)
	return strings.EqualFold(level, "debug") || (level == "" && cfg.Debug)
}

// NewLogger routes auxiliary HTTP diagnostics to the configured logger.
func NewLogger(cfg config.Config) *Logger {
	level := strings.ToLower(strings.TrimSpace(cfg.LogLevel))
	if level == "" && cfg.Debug {
		level = "debug"
	}
	logger := cfg.Logger
	if logger == nil {
		logger = slog.Default()
	}
	return &Logger{logger: logger, level: level}
}

// Logger implements Resty's logger interface without exposing raw diagnostics.
type Logger struct {
	logger *slog.Logger
	level  string
}

func (l *Logger) Errorf(_ string, _ ...interface{}) {
	l.logger.Error("auxiliary HTTP client error")
}

func (l *Logger) Warnf(_ string, _ ...interface{}) {
	if l.level == "debug" || l.level == "info" || l.level == "warn" {
		l.logger.Warn("auxiliary HTTP client warning")
	}
}

func (l *Logger) Debugf(_ string, _ ...interface{}) {
	if l.level == "debug" {
		l.logger.Debug("auxiliary HTTP client debug")
	}
}

// Configure records safe summaries for auxiliary file requests.
func Configure(client *resty.Client, cfg config.Config) *resty.Client {
	logger := NewLogger(cfg)
	return client.SetDebug(false).SetLogger(logger).
		OnAfterResponse(func(_ *resty.Client, response *resty.Response) error {
			if response.IsError() {
				logger.logger.Error("auxiliary HTTP request failed", "status", response.StatusCode())
			} else if logger.level == "info" || logger.level == "debug" {
				logger.logger.Info("auxiliary HTTP request", "status", response.StatusCode(), "duration_ms", response.Time().Milliseconds())
			}
			return nil
		}).
		OnError(func(_ *resty.Request, _ error) {
			logger.logger.Error("auxiliary HTTP transport error")
		})
}
