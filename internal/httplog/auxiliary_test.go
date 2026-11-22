package httplog

import (
	"bytes"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/hiscaler/temu-go/config"
)

type transport func(*http.Request) (*http.Response, error)

func (f transport) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}

func TestAuxiliaryLogging(t *testing.T) {
	for _, tc := range []struct {
		name    string
		level   string
		status  int
		failure bool
		wantLog string
	}{
		{"debug success", "debug", http.StatusOK, false, "auxiliary HTTP request"},
		{"info success", "info", http.StatusOK, false, "auxiliary HTTP request"},
		{"warn success", "warn", http.StatusOK, false, ""},
		{"error success", "error", http.StatusOK, false, ""},
		{"HTTP error", "error", http.StatusBadRequest, false, "auxiliary HTTP request failed"},
		{"transport error", "error", 0, true, "auxiliary HTTP transport error"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var output bytes.Buffer
			logger := slog.New(slog.NewTextHandler(&output, &slog.HandlerOptions{Level: slog.LevelDebug}))
			client := resty.New().SetDebug(true).SetTransport(transport(func(request *http.Request) (*http.Response, error) {
				if tc.failure {
					return nil, errors.New("private-error")
				}
				return &http.Response{
					StatusCode: tc.status,
					Status:     http.StatusText(tc.status),
					Header:     http.Header{"Content-Type": {"text/plain"}},
					Body:       io.NopCloser(strings.NewReader("private-body")),
					Request:    request,
				}, nil
			}))
			Configure(client, config.Config{LogLevel: tc.level, Logger: logger})
			response, err := client.R().SetHeader("Authorization", "private-token").Get("https://example.test/?token=private-token")
			if client.Debug || (err != nil) != tc.failure {
				t.Fatalf("unexpected client state or error: debug=%v err=%v", client.Debug, err)
			}
			if !tc.failure && response.String() != "private-body" {
				t.Fatalf("response body changed: %s", response.String())
			}
			logs := output.String()
			if tc.wantLog == "" && logs != "" || tc.wantLog != "" && !strings.Contains(logs, tc.wantLog) {
				t.Fatalf("unexpected log: %s", logs)
			}
			for _, sensitive := range []string{"private-token", "private-body", "private-error"} {
				if strings.Contains(logs, sensitive) {
					t.Fatalf("unexpected sensitive value in log: %s", logs)
				}
			}
		})
	}
}

func TestIsDebug(t *testing.T) {
	for _, tc := range []struct {
		cfg  config.Config
		want bool
	}{
		{config.Config{Debug: true}, true},
		{config.Config{Debug: true, LogLevel: "info"}, false},
		{config.Config{LogLevel: " DEBUG "}, true},
		{config.Config{}, false},
	} {
		if got := IsDebug(tc.cfg); got != tc.want {
			t.Errorf("IsDebug(%+v) = %v, want %v", tc.cfg, got, tc.want)
		}
	}
}
