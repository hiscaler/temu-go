package config

import (
	"log/slog"
	"strings"

	"github.com/go-resty/resty/v2"
)

// IsDebug reports whether the effective call log level is debug.
func (c Config) IsDebug() bool {
	return strings.EqualFold(c.LogLevel, "debug") || (c.LogLevel == "" && c.Debug)
}

// RestyLogger routes auxiliary HTTP client diagnostics to the configured logger.
func (c Config) RestyLogger() *RestyLogger {
	level := strings.ToLower(c.LogLevel)
	if level == "" && c.Debug {
		level = "debug"
	}
	logger := c.Logger
	if logger == nil {
		logger = slog.Default()
	}
	return &RestyLogger{logger: logger, level: level}
}

// RestyLogger implements Resty's logger interface without exposing raw diagnostics.
type RestyLogger struct {
	logger *slog.Logger
	level  string
}

func (l *RestyLogger) Errorf(_ string, _ ...interface{}) {
	l.logger.Error("auxiliary HTTP client error")
}

func (l *RestyLogger) Warnf(_ string, _ ...interface{}) {
	if l.level == "debug" || l.level == "info" || l.level == "warn" {
		l.logger.Warn("auxiliary HTTP client warning")
	}
}

func (l *RestyLogger) Debugf(_ string, _ ...interface{}) {
	if l.level == "debug" {
		l.logger.Debug("auxiliary HTTP client debug")
	}
}

// ConfigureAuxiliaryHTTPLogging records safe summaries for file requests.
func (c Config) ConfigureAuxiliaryHTTPLogging(client *resty.Client) *resty.Client {
	logger := c.RestyLogger()
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
