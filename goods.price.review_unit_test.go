package temu

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hiscaler/temu-go/config"
	"github.com/hiscaler/temu-go/entity"
)

func TestGoodsPriceReviewQueryRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request map[string]any
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Error(err)
		}
		if request["type"] != "bg.semi.price.review.page.query.order" || request["pageNo"] != float64(1) {
			t.Errorf("unexpected request type or page: %v, %v", request["type"], request["pageNo"])
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true,"result":{"total":1,"reviewSamplePriceList":[{"orderId":123,"orderStatus":1,"productSkuIdList":[456]}]}}`))
	}))
	defer server.Close()

	client := NewClient(config.Config{
		Env: "prod", AppKey: "test", AppSecret: "test", AccessToken: "test",
		Timeout: 5, VerifySSL: true, Proxies: config.RegionEnvUrls{entity.Partner: {Prod: server.URL}},
	})
	params := GoodsPriceReviewQueryParams{OrderStatusList: []int{1}}
	params.Page = 1
	params.PageSize = 10
	items, err := client.Services.Goods.Price.Review.Query(context.Background(), params)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].OrderStatus != 1 {
		t.Fatalf("unexpected reviews: %+v", items)
	}
}

func TestGoodsPriceReviewQueryRejectsInvalidStatus(t *testing.T) {
	params := GoodsPriceReviewQueryParams{OrderStatusList: []int{1, 6}}
	if err := params.validate(); err == nil {
		t.Fatal("expected invalid status error")
	}
}
