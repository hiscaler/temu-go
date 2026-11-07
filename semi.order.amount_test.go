package temu

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/hiscaler/temu-go/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type amountRoundTripper func(*http.Request) (*http.Response, error)

// RoundTrip 使用测试函数模拟 HTTP 请求。
func (f amountRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

// TestSemiOrderService_Amount 验证请求参数及多个子订单的金额解析。
func TestSemiOrderService_Amount(t *testing.T) {
	transport := amountRoundTripper(func(r *http.Request) (*http.Response, error) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/bg.order.amount.query", r.URL.Path)
		var body map[string]string
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, "PO-211-123", body["parentOrderSn"])
		return &http.Response{StatusCode: http.StatusOK, Header: http.Header{"Content-Type": []string{"application/json"}}, Body: io.NopCloser(strings.NewReader(`{"success":true,"result":{"parentOrderMap":{"parentOrderSn":"PO-211-123","basePriceTotal":{"currency":"USD","amount":1234},"signOnDelivery":null},"orderList":[{"orderSn":"211-1","quantity":2,"unitBasePrice":{"currency":"USD","amount":617}},{"orderSn":"211-2","quantity":1,"unitBasePrice":{"currency":"USD","amount":0}}],"warning":["test"]}}`)), Request: r}, nil
	})

	service := semiOrderService{httpClient: resty.New().SetBaseURL("https://example.test").SetTransport(transport)}
	result, err := service.Amount(context.Background(), "PO-211-123")
	require.NoError(t, err)
	require.NotNil(t, result.ParentOrderMap)
	require.NotNil(t, result.ParentOrderMap.BasePriceTotal)
	require.NotNil(t, result.ParentOrderMap.BasePriceTotal.Amount)
	assert.Equal(t, int64(1234), *result.ParentOrderMap.BasePriceTotal.Amount)
	assert.Equal(t, "USD", result.ParentOrderMap.BasePriceTotal.Currency)
	assert.Nil(t, result.ParentOrderMap.SignOnDelivery)
	require.Len(t, result.OrderList, 2)
	require.NotNil(t, result.OrderList[1].UnitBasePrice.Amount)
	assert.Zero(t, *result.OrderList[1].UnitBasePrice.Amount)
	assert.Equal(t, []string{"test"}, result.Warning)
}

// TestSemiOrderService_Amount_InvalidParentOrder 验证空父单号被拒绝。
func TestSemiOrderService_Amount_InvalidParentOrder(t *testing.T) {
	service := semiOrderService{}
	result, err := service.Amount(context.Background(), "  ")
	assert.Error(t, err)
	assert.Nil(t, result)
}

// TestSemiOrderService_Amount_APIError 验证平台错误向调用方返回。
func TestSemiOrderService_Amount_APIError(t *testing.T) {
	transport := amountRoundTripper(func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(strings.NewReader(`{"success":false,"errorCode":140020008,"errorMsg":"Some amount fields are still undergoing system calculation"}`)),
			Request:    r,
		}, nil
	})
	service := semiOrderService{httpClient: resty.New().SetBaseURL("https://example.test").SetTransport(transport)}
	result, err := service.Amount(context.Background(), "PO-211-123")
	assert.Error(t, err)
	assert.Nil(t, result)
}

// TestSemiOrderAmountGateway 验证金额查询使用美国半托网关。
func TestSemiOrderAmountGateway(t *testing.T) {
	assert.Equal(t, "https://openapi-b-us.temu.com/openapi/router", baseUrl("bg.order.amount.query", entity.AmericanRegion, "prod", nil))
}
