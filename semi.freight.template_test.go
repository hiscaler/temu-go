package temu

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/hiscaler/temu-go/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSemiFreightTemplateService_Templates(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/bg.glo.logistics.template.get", r.URL.Path)

		var body map[string][]int
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, []int{101, 202}, body["siteIds"])

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"success": true,
			"result": {
				"freightTemplates": [{
					"freightTemplateId": "123",
					"templateName": "US Standard Shipping"
				}]
			}
		}`))
	}))
	defer server.Close()

	service := semiFreightTemplateService{httpClient: resty.New().SetBaseURL(server.URL)}
	items, err := service.Templates(context.Background(), 101, 202)

	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "123", items[0].FreightTemplateId)
	assert.Equal(t, "US Standard Shipping", items[0].TemplateName)
}

func TestSemiFreightTemplateService_Templates_WithoutSiteIDs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.NotContains(t, body, "siteIds")

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success": true, "result": {"freightTemplates": []}}`))
	}))
	defer server.Close()

	service := semiFreightTemplateService{httpClient: resty.New().SetBaseURL(server.URL)}
	items, err := service.Templates(context.Background())

	require.NoError(t, err)
	assert.Empty(t, items)
}

func TestSemiFreightTemplateService_BaseURL(t *testing.T) {
	assert.Equal(
		t,
		"https://openapi-b-partner.temu.com/openapi/router",
		baseUrl("bg.glo.logistics.template.get", entity.ChinaRegion, prodEnv, nil),
	)
}
