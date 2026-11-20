package temu

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSemiLogisticsService_Templates(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/temu.logistics.template.get", r.URL.Path)

		var body map[string][]int
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, []int{101, 202}, body["siteIds"])

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"success": true,
			"result": {
				"freightTemplates": [{
					"freightTemplateId": 123,
					"templateName": "物流模板"
				}]
			}
		}`))
	}))
	defer server.Close()

	service := semiLogisticsService{httpClient: resty.New().SetBaseURL(server.URL)}
	items, err := service.Templates(context.Background(), 101, 202)

	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, int64(123), items[0].FreightTemplateId)
	assert.Equal(t, "物流模板", items[0].TemplateName)
}

func TestSemiLogisticsService_Templates_RequiresSiteIDs(t *testing.T) {
	items, err := (semiLogisticsService{}).Templates(context.Background())

	assert.EqualError(t, err, "请提供站点 ID")
	assert.Nil(t, items)
}
