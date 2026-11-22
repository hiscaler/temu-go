package temu

import (
	"net/http"
	"testing"

	"github.com/hiscaler/temu-go/config"
)

func TestClientProxyConfiguration(t *testing.T) {
	for _, proxyURL := range []string{"", "socks5://127.0.0.1:1080", "http://127.0.0.1:8080"} {
		t.Run(proxyURL, func(t *testing.T) {
			client := NewClient(config.Config{ProxyURL: proxyURL})
			transport, err := client.Services.Goods.service.httpClient.Transport()
			if err != nil {
				t.Fatal(err)
			}
			if proxyURL == "" {
				if transport.Proxy != nil {
					t.Fatal("expected a direct connection")
				}
				return
			}
			if transport.Proxy == nil {
				t.Fatal("proxy is not configured")
			}
			request, err := http.NewRequest(http.MethodGet, "https://example.com", nil)
			if err != nil {
				t.Fatal(err)
			}
			proxy, err := transport.Proxy(request)
			if err != nil || proxy == nil || proxy.String() != proxyURL {
				t.Fatalf("proxy = %v, error = %v; want %s", proxy, err, proxyURL)
			}
		})
	}
}
