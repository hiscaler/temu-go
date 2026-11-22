package httplog

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"
)

func TestScopedLoggerLevelsAndAttributes(t *testing.T) {
	for _, tc := range []struct {
		name      string
		callLevel string
		baseLevel slog.Level
		wantDebug bool
		wantInfo  bool
	}{
		{"debug", "debug", slog.LevelDebug, true, true},
		{"info", "info", slog.LevelDebug, false, true},
		{"handler filters debug", "debug", slog.LevelError, false, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var output bytes.Buffer
			base := slog.New(slog.NewTextHandler(&output, &slog.HandlerOptions{Level: tc.baseLevel})).With("connector", "shop")
			logger := ScopedLogger(base, tc.callLevel, false).WithGroup("request").With("endpoint", "test")
			logger.Debug("debug-message")
			logger.Info("info-message")
			logger.Error("error-message")
			logs := output.String()
			if strings.Contains(logs, "debug-message") != tc.wantDebug || strings.Contains(logs, "info-message") != tc.wantInfo {
				t.Fatalf("unexpected levels: %s", logs)
			}
			for _, value := range []string{"error-message", "connector=shop", "request.endpoint=test"} {
				if !strings.Contains(logs, value) {
					t.Fatalf("missing %q: %s", value, logs)
				}
			}
		})
	}
}
