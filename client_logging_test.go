package temu

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	restyslog "github.com/hiscaler/resty-slog"
	"github.com/hiscaler/temu-go/config"
	"github.com/hiscaler/temu-go/entity"
)

func testLogMask() *restyslog.LogMask {
	return &restyslog.LogMask{
		Request:  []string{"access_token", "app_secret", "phone", "Authorization"},
		Response: []string{"address", "Set-Cookie"},
	}
}

func TestRequestLogLevelsAndRedaction(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Set-Cookie", "private-cookie")
		_, _ = w.Write([]byte(`{"success":true,"result":{"address":"private-address","categoryId":123,"categoryName":"Clothes"}}`))
	}))
	t.Cleanup(server.Close)
	for _, tc := range []struct {
		level string
		info  bool
		debug bool
	}{
		{"error", false, false}, {"warn", false, false}, {"info", true, false}, {"debug", true, true},
	} {
		t.Run(tc.level, func(t *testing.T) {
			t.Parallel()
			var output bytes.Buffer
			logger := slog.New(slog.NewJSONHandler(&output, &slog.HandlerOptions{Level: slog.LevelDebug}))
			client := NewClient(config.Config{
				Env: "prod", LogLevel: tc.level, Logger: logger, AppKey: "public-key", AppSecret: "secret-app", AccessToken: "secret-token",
				LogMask: testLogMask(),
				Proxies: config.RegionEnvUrls{entity.ChinaRegion: {Prod: server.URL}},
			})
			_, err := client.Services.Goods.service.httpClient.R().
				SetHeader("Authorization", "private-header-token").
				SetBody(map[string]string{"phone": "private-phone"}).Post("test.call")
			if err != nil {
				t.Fatal(err)
			}
			logs := output.String()
			if strings.Contains(logs, `"msg":"call"`) != tc.info || strings.Contains(logs, "~~~ REQUEST ~~~") != tc.debug {
				t.Fatalf("unexpected logs for %s: %s", tc.level, logs)
			}
			for _, secret := range []string{"secret-token", "secret-app", "private-address", "private-phone", "private-cookie", "private-header-token"} {
				if strings.Contains(logs, secret) {
					t.Fatalf("sensitive data in logs: %s", logs)
				}
			}
			if tc.debug {
				var details string
				for _, line := range strings.Split(strings.TrimSpace(logs), "\n") {
					var item map[string]any
					if err := json.Unmarshal([]byte(line), &item); err != nil {
						t.Fatal(err)
					}
					if message, ok := item["msg"].(string); ok && strings.Contains(message, "~~~ REQUEST ~~~") {
						details = message
					}
				}
				if !strings.Contains(details, "~~~ RESPONSE ~~~") || !strings.Contains(details, "HEADERS:") || !strings.Contains(details, "STATUS       : 200 OK") {
					t.Fatalf("missing Resty HTTP details: %s", details)
				}
				_, requestPart, found := strings.Cut(details, "BODY   :\n")
				if !found {
					t.Fatalf("missing request body: %s", details)
				}
				requestJSON, responsePart, found := strings.Cut(requestPart, "\n------------------------------------------------------------------------------\n")
				if !found {
					t.Fatalf("missing response section: %s", details)
				}
				_, responsePart, found = strings.Cut(responsePart, "BODY         :\n")
				if !found {
					t.Fatalf("missing response body: %s", details)
				}
				responseJSON, _, found := strings.Cut(responsePart, "\n==============================================================================")
				if !found {
					t.Fatalf("missing response boundary: %s", details)
				}
				var request, response map[string]any
				if err := json.Unmarshal([]byte(requestJSON), &request); err != nil {
					t.Fatal(err)
				}
				if err := json.Unmarshal([]byte(responseJSON), &response); err != nil {
					t.Fatal(err)
				}
				sign, signIsString := request["sign"].(string)
				if request["type"] != "test.call" || request["data_type"] != "JSON" || request["access_token"] != strings.Repeat("*", len("secret-token")) || request["app_key"] != "public-key" || !signIsString || sign == "" || strings.Contains(sign, "*") {
					t.Fatalf("unexpected request: %#v", request)
				}
				result := response["result"].(map[string]any)
				if response["success"] != true || result["categoryId"] != float64(123) || result["categoryName"] != "Clothes" || result["address"] != strings.Repeat("*", len("private-address")) {
					t.Fatalf("unexpected response: %#v", response)
				}
			}
		})
	}
}

func TestTransportErrorIncludesRedactedRequestInDebug(t *testing.T) {
	var output bytes.Buffer
	client := NewClient(config.Config{
		Env: "prod", LogLevel: "debug", AccessToken: "secret-token",
		LogMask: testLogMask(),
		Logger:  slog.New(slog.NewJSONHandler(&output, &slog.HandlerOptions{Level: slog.LevelDebug})),
	})
	transport, err := client.Services.Goods.service.httpClient.Transport()
	if err != nil {
		t.Fatal(err)
	}
	transport.DialContext = func(_ context.Context, _, _ string) (net.Conn, error) {
		return nil, errors.New("network unavailable")
	}
	_, err = client.Services.Goods.service.httpClient.R().SetBody(map[string]any{"categoryId": 42}).Post("test.call")
	if err == nil || !strings.Contains(output.String(), "~~~ REQUEST ~~~") || !strings.Contains(output.String(), "categoryId") || strings.Contains(output.String(), "secret-token") {
		t.Fatalf("unexpected error or logs: %v, %s", err, output.String())
	}
}

func TestBusinessErrorLoggedAtErrorLevel(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"success":false,"errorCode":123,"errorMsg":"private-error"}`))
	}))
	defer server.Close()
	var output bytes.Buffer
	client := NewClient(config.Config{
		Env: "prod", LogLevel: "error",
		Logger:  slog.New(slog.NewJSONHandler(&output, &slog.HandlerOptions{Level: slog.LevelDebug})),
		Proxies: config.RegionEnvUrls{entity.ChinaRegion: {Prod: server.URL}},
	})
	_, err := client.Services.Goods.service.httpClient.R().Post("test.call")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), `"error_code":123`) || strings.Contains(output.String(), "private-error") {
		t.Fatalf("unexpected error log: %s", output.String())
	}
}

func TestBusinessErrorDebugIncludesResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"success":false,"errorCode":123,"errorMsg":"invalid category"}`))
	}))
	defer server.Close()
	var output bytes.Buffer
	client := NewClient(config.Config{
		Env: "prod", LogLevel: "debug",
		Logger:  slog.New(slog.NewJSONHandler(&output, &slog.HandlerOptions{Level: slog.LevelDebug})),
		Proxies: config.RegionEnvUrls{entity.ChinaRegion: {Prod: server.URL}},
	})
	_, err := client.Services.Goods.service.httpClient.R().SetBody(map[string]int{"categoryId": 42}).Post("test.call")
	if err != nil || !strings.Contains(output.String(), "invalid category") || !strings.Contains(output.String(), "categoryId") {
		t.Fatalf("unexpected error or logs: %v, %s", err, output.String())
	}
}

func TestHTTPErrorDebugIncludesRestyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"success":false,"errorCode":123,"errorMsg":"invalid category"}`))
	}))
	defer server.Close()
	var output bytes.Buffer
	client := NewClient(config.Config{
		Env: "prod", LogLevel: "debug", AppKey: "public-key", AppSecret: "secret-app", AccessToken: "secret-token",
		LogMask: testLogMask(),
		Logger:  slog.New(slog.NewJSONHandler(&output, &slog.HandlerOptions{Level: slog.LevelDebug})),
		Proxies: config.RegionEnvUrls{entity.ChinaRegion: {Prod: server.URL}},
	})
	_, err := client.Services.Goods.service.httpClient.R().SetBody(map[string]int{"categoryId": 42}).Post("test.call")
	if err != nil || !strings.Contains(output.String(), "STATUS       : 400 Bad Request") || !strings.Contains(output.String(), "invalid category") || !strings.Contains(output.String(), "public-key") || strings.Contains(output.String(), "secret-token") {
		t.Fatalf("unexpected error or logs: %v, %s", err, output.String())
	}
}

func TestCallerMaskPolicyAndUnchangedHTTP(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Error(err)
		}
		if body["access_token"] != "secret-token" || r.Header.Get("Authorization") != "header-secret" {
			t.Errorf("request changed: %#v, %s", body, r.Header.Get("Authorization"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"result":{"phone":"123456","items":[{"name":"Alice"}]}}`))
	}))
	defer server.Close()

	for _, tc := range []struct {
		name   string
		mask   *restyslog.LogMask
		want   []string
		absent []string
	}{
		{"no rules", nil, []string{"secret-token", "header-secret", "123456", "Alice"}, nil},
		{"independent rules", &restyslog.LogMask{Keep: 1, Request: []string{"access_token", "Authorization", "private.value"}, Response: []string{"result.items.name"}}, []string{"s**********n", "h***********t", "p*****e", "A***e", "123456", "public-key"}, []string{"secret-token", "header-secret", "Alice"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var output bytes.Buffer
			client := NewClient(config.Config{
				Env: "prod", LogLevel: "debug", LogMask: tc.mask,
				AppKey: "public-key", AppSecret: "secret-app", AccessToken: "secret-token",
				Logger:  slog.New(slog.NewJSONHandler(&output, &slog.HandlerOptions{Level: slog.LevelDebug})),
				Proxies: config.RegionEnvUrls{entity.ChinaRegion: {Prod: server.URL}},
			})
			resp, err := client.Services.Goods.service.httpClient.R().SetHeader("Authorization", "header-secret").
				SetBody(map[string]any{"private": map[string]any{"value": "private"}}).Post("test.call")
			if err != nil || !strings.Contains(resp.String(), "Alice") {
				t.Fatalf("unexpected response: %v, %v", resp, err)
			}
			for _, value := range tc.want {
				if !strings.Contains(output.String(), value) {
					t.Errorf("missing %q: %s", value, output.String())
				}
			}
			for _, value := range tc.absent {
				if strings.Contains(output.String(), value) {
					t.Errorf("unexpected %q: %s", value, output.String())
				}
			}
		})
	}
}

func TestInvalidMaskDisablesDebugAndLegacyDebug(t *testing.T) {
	for _, tc := range []struct {
		name    string
		cfg     config.Config
		debug   bool
		invalid bool
	}{
		{"legacy debug", config.Config{Debug: true}, true, false},
		{"explicit level overrides debug", config.Config{Debug: true, LogLevel: "info"}, false, false},
		{"invalid mask", config.Config{LogLevel: "debug", LogMask: &restyslog.LogMask{Request: []string{"a..b"}}}, false, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var output bytes.Buffer
			tc.cfg.Logger = slog.New(slog.NewJSONHandler(&output, &slog.HandlerOptions{Level: slog.LevelDebug}))
			client := NewClient(tc.cfg)
			if client.Debug != tc.debug || client.Services.Goods.service.httpClient.Debug != tc.debug {
				t.Fatalf("unexpected debug state: %v", client.Debug)
			}
			if strings.Contains(output.String(), "invalid HTTP logging configuration") != tc.invalid {
				t.Fatalf("unexpected log: %s", output.String())
			}
		})
	}
}

func TestMaskConfigurationFromJSON(t *testing.T) {
	var cfg config.Config
	if err := json.Unmarshal([]byte(`{"log_level":"debug","log_mask":{"keep":1,"request":["access_token"],"response":["result.phone"]}}`), &cfg); err != nil {
		t.Fatal(err)
	}
	if cfg.LogLevel != "debug" || cfg.LogMask == nil || cfg.LogMask.Keep != 1 || len(cfg.LogMask.Request) != 1 || cfg.LogMask.Request[0] != "access_token" || len(cfg.LogMask.Response) != 1 || cfg.LogMask.Response[0] != "result.phone" {
		t.Fatalf("unexpected configuration: %#v", cfg)
	}
}
