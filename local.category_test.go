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

// TestLocalCategoryService_Query 验证本本类目查询请求与响应解析。
func TestLocalCategoryService_Query(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/bg.local.goods.cats.get", r.URL.Path)

		var body map[string]any
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, "en-US", body["language"])
		assert.Equal(t, float64(100), body["parentCatId"])
		assert.NotContains(t, body, "Language")
		assert.NotContains(t, body, "ParentCatId")

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"success": true,
			"result": {
				"items": [{
					"catId": 101,
					"availableStatus": 1,
					"level": 2,
					"catName": "Home & Kitchen",
					"secondHandCategory": false,
					"catType": 2,
					"leaf": true,
					"parentId": 100,
					"expandCatType": 0
				}]
			}
		}`))
	}))
	defer server.Close()

	service := localCategoryService{httpClient: resty.New().SetBaseURL(server.URL)}
	items, err := service.Query(context.Background(), LocalCategoryQueryParams{
		Language:    "en-US",
		ParentCatId: 100,
	})

	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, 101, items[0].CatId)
	assert.Equal(t, 1, items[0].AvailableStatus)
	assert.Equal(t, 2, items[0].Level)
	assert.Equal(t, "Home & Kitchen", items[0].CatName)
	assert.True(t, items[0].Leaf)
	assert.Equal(t, 100, items[0].ParentId)
}

// TestLocalCategoryService_Query_APIError 验证平台错误会返回给调用方。
func TestLocalCategoryService_Query_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success": false, "errorCode": 10001, "errorMsg": "invalid category"}`))
	}))
	defer server.Close()

	service := localCategoryService{httpClient: resty.New().SetBaseURL(server.URL)}
	items, err := service.Query(context.Background(), LocalCategoryQueryParams{})

	assert.Error(t, err)
	assert.Nil(t, items)
}

// TestLocalCategoryService_Query_Real 使用 config/config.json 中的凭据查询 Temu 根类目。
func TestLocalCategoryService_Query_Real(t *testing.T) {
	items, err := temuClient.Services.Local.Category.Query(ctx, LocalCategoryQueryParams{
		//Language: "zh-CN",
		ParentCatId: 1,
	})

	require.NoError(t, err)
	require.NotEmpty(t, items)
	assert.Positive(t, items[0].CatId)
	assert.NotEmpty(t, items[0].CatName)
}
