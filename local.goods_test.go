package temu

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/hiscaler/temu-go/normal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLocalGoodsService_Query 验证本本商品查询请求、分页和商品信息解析。
func TestLocalGoodsService_Query(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/bg.local.goods.list.query", r.URL.Path)

		var body map[string]any
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, float64(10), body["pageSize"])
		assert.Equal(t, "mug", body["searchText"])
		assert.Equal(t, float64(2), body["pageNo"])

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"success": true,
			"result": {
				"pageNo": 2,
				"total": 25,
				"goodsList": [{
					"goodsId": 101,
					"outGoodsSn": "MUG-001",
					"goodsName": "Custom Mug",
					"price": "9.99",
					"currency": "USD",
					"quantity": 12,
					"catId": 200,
					"skuIdList": [1001],
					"skuInfoList": [{
						"skuId": 1001,
						"stock": 12,
						"skuSn": "MUG-001-WH"
					}]
				}]
			}
		}`))
	}))
	defer server.Close()

	service := localGoodsService{service: service{httpClient: resty.New().SetBaseURL(server.URL)}}
	items, total, totalPages, isLastPage, err := service.Query(context.Background(), LocalGoodsQueryParams{
		SearchText: "mug",
		ParameterWithPager: normal.ParameterWithPager{
			Pager: normal.Pager{Page: 2, PageSize: 10},
		},
	})

	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, int64(101), items[0].GoodsId)
	assert.Equal(t, "MUG-001", items[0].OutGoodsSn)
	assert.Equal(t, "Custom Mug", items[0].GoodsName)
	assert.Equal(t, "9.99", items[0].Price)
	assert.Equal(t, "USD", items[0].Currency)
	assert.Equal(t, 12, items[0].Quantity)
	require.Len(t, items[0].SkuInfoList, 1)
	assert.Equal(t, int64(1001), items[0].SkuInfoList[0].SkuId)
	assert.Equal(t, 25, total)
	assert.Equal(t, 3, totalPages)
	assert.False(t, isLastPage)
}

// TestLocalGoodsService_Query_Real 使用 config/config.json 中的凭据查询 Temu 本本商品。
func TestLocalGoodsService_Query_Real(t *testing.T) {
	items, total, totalPages, isLastPage, err := temuClient.Services.Local.Goods.Query(ctx, LocalGoodsQueryParams{
		ParameterWithPager: normal.ParameterWithPager{
			Pager: normal.Pager{Page: 1, PageSize: 10},
		},
	})

	require.NoError(t, err)
	assert.LessOrEqual(t, len(items), 10)
	assert.GreaterOrEqual(t, total, 0)
	assert.GreaterOrEqual(t, totalPages, 0)
	assert.Equal(t, totalPages <= 1, isLastPage)
}
