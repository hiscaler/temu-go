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

	"github.com/hiscaler/temu-go/config"
	"github.com/hiscaler/temu-go/entity"
)

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
				if request["type"] != "test.call" || request["data_type"] != "JSON" || request["access_token"] != redacted || request["app_key"] != "public-key" || request["sign"] == redacted || request["sign"] == "" {
					t.Fatalf("unexpected request: %#v", request)
				}
				result := response["result"].(map[string]any)
				if response["success"] != true || result["categoryId"] != float64(123) || result["categoryName"] != "Clothes" || result["address"] != redacted {
					t.Fatalf("unexpected response: %#v", response)
				}
			}
		})
	}
}

func TestRedactLogJSON(t *testing.T) {
	input := map[string]any{
		"accessToken": "token-value", "SIGN": "signed", "app_key": "key-value",
		"nested":   []any{map[string]any{"receiverPhone": "private-phone", "goodsName": "T-shirt", "qty": 3}},
		"embedded": `{"email":"private@example.com","sku":"ABC"}`,
		"imageUrl": "https://example.com/image?token=private-url#secret",
		"message":  "bad token-value",
	}
	output := redactLogJSON(input, "token-value")
	for _, secret := range []string{"token-value", "private-phone", "private@example.com", "private-url", "#secret"} {
		if strings.Contains(output, secret) {
			t.Fatalf("secret %q in output: %s", secret, output)
		}
	}
	for _, expected := range []string{"signed", "key-value", "T-shirt", `"qty": 3`, "ABC", "https://example.com/image?REDACTED"} {
		if !strings.Contains(output, expected) {
			t.Fatalf("missing %q in output: %s", expected, output)
		}
	}
	if got := redactLogJSON([]byte("not JSON secret")); strings.Contains(got, "secret") {
		t.Fatalf("non-JSON body leaked: %s", got)
	}
	largeID := `{"productId":9223372036854775807,"nested":"{\"skuId\":9223372036854775807}"}`
	if got := redactLogJSON([]byte(largeID)); strings.Count(got, "9223372036854775807") != 2 {
		t.Fatalf("large IDs must retain their exact values: %s", got)
	}
}

func TestRestyLoggerRedactsSensitiveURLQuery(t *testing.T) {
	var output bytes.Buffer
	logger := &logger{
		l:       slog.New(slog.NewJSONHandler(&output, &slog.HandlerOptions{Level: slog.LevelDebug})),
		secrets: []string{"secret-token", "secret-app"},
	}
	logger.Debugf("%s", "/router?app_key=public&sign=signature&access_token=secret-token&phone=private-phone")
	for _, secret := range []string{"secret-token", "secret-app", "private-phone"} {
		if strings.Contains(output.String(), secret) {
			t.Fatalf("sensitive value in log: %s", output.String())
		}
	}
	if !strings.Contains(output.String(), "app_key=public") || !strings.Contains(output.String(), "sign=signature") {
		t.Fatalf("missing public signing fields: %s", output.String())
	}
}

func TestTransportErrorIncludesRedactedRequestInDebug(t *testing.T) {
	var output bytes.Buffer
	client := NewClient(config.Config{
		Env: "prod", LogLevel: "debug", AccessToken: "secret-token",
		Logger: slog.New(slog.NewJSONHandler(&output, &slog.HandlerOptions{Level: slog.LevelDebug})),
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
		Logger:  slog.New(slog.NewJSONHandler(&output, &slog.HandlerOptions{Level: slog.LevelDebug})),
		Proxies: config.RegionEnvUrls{entity.ChinaRegion: {Prod: server.URL}},
	})
	_, err := client.Services.Goods.service.httpClient.R().SetBody(map[string]int{"categoryId": 42}).Post("test.call")
	if err != nil || !strings.Contains(output.String(), "STATUS       : 400 Bad Request") || !strings.Contains(output.String(), "invalid category") || !strings.Contains(output.String(), "public-key") || strings.Contains(output.String(), "secret-token") {
		t.Fatalf("unexpected error or logs: %v, %s", err, output.String())
	}
}
