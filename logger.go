package temu

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"unicode"
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

func scopedLogger(base *slog.Logger, level string, debug bool) *slog.Logger {
	if base == nil {
		base = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}
	threshold := slog.LevelError
	switch strings.ToLower(level) {
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

type Logger interface {
	Errorf(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Debugf(format string, v ...interface{})
}

func createLogger(debug bool) *logger {
	if debug {
		return &logger{l: slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))}
	}
	return &logger{l: slog.Default()}
}

type logger struct {
	l       *slog.Logger
	secrets []string
}

var _ Logger = (*logger)(nil)

// isf Is valid Xprint format
func isf(format string) bool {
	fnIsFlag := func(c byte) bool {
		switch c {
		case '+', '-', '#', ' ', '0':
			return true
		}
		return false
	}

	fnIsValidVerb := func(c byte) bool {
		validVerbs := "vT%tbcdoqxXUeEfFgGspw"
		return strings.IndexByte(validVerbs, c) >= 0
	}
	s := format
	for {
		i := strings.IndexByte(s, '%')
		if i < 0 {
			break
		}
		s = s[i+1:]

		// 处理 %%
		if len(s) > 0 && s[0] == '%' {
			s = s[1:]
			continue
		}

		// 解析标志
		for len(s) > 0 && fnIsFlag(s[0]) {
			s = s[1:]
		}

		// 解析宽度
		if len(s) > 0 && s[0] == '*' {
			s = s[1:]
		} else {
			for len(s) > 0 && unicode.IsDigit(rune(s[0])) {
				s = s[1:]
			}
		}

		// 解析精度
		if len(s) > 0 && s[0] == '.' {
			s = s[1:]
			if len(s) > 0 && s[0] == '*' {
				s = s[1:]
			} else {
				for len(s) > 0 && unicode.IsDigit(rune(s[0])) {
					s = s[1:]
				}
			}
		}

		// 检查格式动词
		if len(s) == 0 {
			return false
		}

		verb := s[0]
		if !fnIsValidVerb(verb) {
			return false
		}

		s = s[1:]
	}
	return true
}

func (l *logger) Errorf(msg string, args ...interface{}) {
	l.l.Error("HTTP client error")
}

func (l *logger) Warnf(msg string, args ...interface{}) {
	l.l.Warn("HTTP client warning")
}

func (l *logger) Infof(msg string, args ...interface{}) {
	if isf(msg) {
		l.l.Info(fmt.Sprintf(msg, args...))
		return
	}
	l.l.Info(msg, args...)
}

func (l *logger) Debugf(msg string, args ...interface{}) {
	message := fmt.Sprintf(msg, args...)
	for _, secret := range l.secrets {
		if secret != "" {
			message = strings.ReplaceAll(message, secret, redacted)
		}
	}
	message = sensitiveLogQuery.ReplaceAllString(message, "${1}"+redacted)
	l.l.Debug(message)
}
