package temu

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/hiscaler/temu-go/normal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSemiActivityService_Candidate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/bg.promotion.activity.candidate.goods.query", r.URL.Path)

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode request body: %v", err)
			return
		}
		assert.Equal(t, float64(123), body["activityId"])
		assert.Equal(t, float64(2), body["pageNumber"])
		assert.Equal(t, float64(10), body["pageSize"])
		assert.NotContains(t, body, "pageNo")

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"success": true,
			"result": {
				"goodsList": [{
					"goodsId": 456,
					"isAllowedEnroll": true,
					"recommendActivityQuantity": 3,
					"skuList": [{
						"skuId": 789,
						"recommendActivitySupplierPrice": 12.34
					}]
				}],
				"total": 25
			}
		}`))
	}))
	defer server.Close()

	service := semiActivityService{httpClient: resty.New().SetBaseURL(server.URL)}
	items, total, totalPages, isLastPage, err := service.Candidate(context.Background(), SemiActivityCandidateQueryParams{
		Pager:      normal.Pager{Page: 2, PageSize: 10},
		ActivityId: 123,
	})

	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, int64(456), items[0].GoodsId)
	assert.True(t, items[0].IsAllowedEnroll)
	assert.Equal(t, 3, items[0].RecommendActivityQuantity)
	require.Len(t, items[0].SkuList, 1)
	assert.Equal(t, int64(789), items[0].SkuList[0].SkuId)
	assert.Equal(t, 12.34, items[0].SkuList[0].RecommendActivitySupplierPrice)
	assert.Equal(t, 25, total)
	assert.Equal(t, 3, totalPages)
	assert.False(t, isLastPage)
}

func TestSemiActivityService_Candidate_InvalidActivityID(t *testing.T) {
	service := semiActivityService{}
	items, total, totalPages, isLastPage, err := service.Candidate(context.Background(), SemiActivityCandidateQueryParams{
		ActivityId: 0,
	})

	assert.Error(t, err)
	assert.Nil(t, items)
	assert.Zero(t, total)
	assert.Zero(t, totalPages)
	assert.False(t, isLastPage)
}

func TestSemiActivityService_Candidate_Real(t *testing.T) {
	activityIDText := os.Getenv("TEMU_SEMI_ACTIVITY_ID")
	if activityIDText == "" {
		t.Skip("set TEMU_SEMI_ACTIVITY_ID to run the real Temu API test")
	}

	activityID, err := strconv.ParseInt(activityIDText, 10, 64)
	require.NoError(t, err)

	items, total, totalPages, isLastPage, err := temuClient.Services.SemiManaged.Activity.Candidate(ctx, SemiActivityCandidateQueryParams{
		Pager:      normal.Pager{Page: 1, PageSize: 10},
		ActivityId: activityID,
	})

	require.NoError(t, err)
	assert.LessOrEqual(t, len(items), 10)
	assert.GreaterOrEqual(t, total, 0)
	assert.GreaterOrEqual(t, totalPages, 0)
	if totalPages > 0 {
		assert.Equal(t, total == 0 || totalPages == 1, isLastPage)
	}
}
