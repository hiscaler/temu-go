package httplog

import (
	"context"
	"log/slog"
	"os"
	"strings"
)

type levelHandler struct {
	slog.Handler
	level slog.Level
}

func (h levelHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.level && h.Handler.Enabled(ctx, level)
}

func (h levelHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return levelHandler{Handler: h.Handler.WithAttrs(attrs), level: h.level}
}

func (h levelHandler) WithGroup(name string) slog.Handler {
	return levelHandler{Handler: h.Handler.WithGroup(name), level: h.level}
}

// ScopedLogger applies the call's level without overriding the underlying handler.
func ScopedLogger(base *slog.Logger, level string, debug bool) *slog.Logger {
	if base == nil {
		base = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}
	threshold := slog.LevelError
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		threshold = slog.LevelDebug
	case "info":
		threshold = slog.LevelInfo
	case "warn":
		threshold = slog.LevelWarn
	case "":
		if debug {
			threshold = slog.LevelDebug
		}
	}
	return slog.New(levelHandler{Handler: base.Handler(), level: threshold})
}
