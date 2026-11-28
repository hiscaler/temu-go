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

func TestSemiActivityService_QueryAndCandidateLastPage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, float64(3), body["pageNumber"])
		assert.NotContains(t, body, "pageNo")
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/bg.promotion.activity.query" {
			_, _ = w.Write([]byte(`{"success":true,"result":{"total":25,"activityList":[{"activityId":123}]}}`))
			return
		}
		_, _ = w.Write([]byte(`{"success":true,"result":{"total":25,"goodsList":[{"goodsId":456}]}}`))
	}))
	defer server.Close()
	service := semiActivityService{httpClient: resty.New().SetBaseURL(server.URL)}
	pager := normal.Pager{Page: 3, PageSize: 10}
	activities, total, pages, last, err := service.Query(context.Background(), SemiActivityQueryParams{Pager: pager, ActivityType: 2})
	require.NoError(t, err)
	assert.Len(t, activities, 1)
	assert.Equal(t, 25, total)
	assert.Equal(t, 3, pages)
	assert.True(t, last)
	goods, total, pages, last, err := service.Candidate(context.Background(), SemiActivityCandidateQueryParams{Pager: pager, ActivityId: 123})
	require.NoError(t, err)
	assert.Len(t, goods, 1)
	assert.Equal(t, 25, total)
	assert.Equal(t, 3, pages)
	assert.True(t, last)
}

func TestSemiActivityService_Query_Real(t *testing.T) {
	activityTypeText := os.Getenv("TEMU_SEMI_ACTIVITY_TYPE")
	if activityTypeText == "" {
		t.Skip("set TEMU_SEMI_ACTIVITY_TYPE to run the real Temu API test")
	}

	activityType, err := strconv.Atoi(activityTypeText)
	require.NoError(t, err)
	require.Contains(t, []int{2, 13, 27, 100}, activityType)

	const pageSize = 10
	items, total, totalPages, isLastPage, err := temuClient.Services.SemiManaged.Activity.Query(ctx, SemiActivityQueryParams{
		Pager:        normal.Pager{Page: 1, PageSize: pageSize},
		ActivityType: activityType,
	})

	require.NoError(t, err)
	assert.LessOrEqual(t, len(items), pageSize)
	assert.GreaterOrEqual(t, total, len(items))
	assert.Equal(t, (total+pageSize-1)/pageSize, totalPages)
	assert.Equal(t, totalPages <= 1, isLastPage)
	for _, item := range items {
		t.Logf("activity_id=%d name=%q type=%d status=%d", item.ActivityId, item.ActivityName, item.ActivityType, item.ActivityStatus)
	}
	t.Logf("total=%d total_pages=%d is_last_page=%t", total, totalPages, isLastPage)
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
