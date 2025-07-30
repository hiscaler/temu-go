package temu

import (
	"fmt"
	"log/slog"
	"strings"
	"unicode"
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
	if isf(msg) {
		l.l.Error(fmt.Sprintf(msg, args...))
		return
	}
	l.l.Error(msg, args...)
}

func (l *logger) Warnf(msg string, args ...interface{}) {
	if isf(msg) {
		l.l.Warn(fmt.Sprintf(msg, args...))
		return
	}
	l.l.Warn(msg, args...)
}

func (l *logger) Infof(msg string, args ...interface{}) {
	if isf(msg) {
		l.l.Info(fmt.Sprintf(msg, args...))
		return
	}
	l.l.Info(msg, args...)
}

func (l *logger) Debugf(msg string, args ...interface{}) {
	if isf(msg) {
		l.l.Debug(fmt.Sprintf(msg, args...))
		return
	}
	l.l.Debug(msg, args...)
}
