package temu

import (
	"fmt"
	"github.com/hiscaler/gox/jsonx"
	"github.com/hiscaler/temu-go/entity"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestShipOrderPackingService_Match(t *testing.T) {
	req := ShipOrderPackingMatchRequest{
		DeliveryOrderSnList: []string{"FH2408231977953"},
	}
	_, err := temuClient.Services.ShipOrderPacking.Match(ctx, req)
	assert.Nilf(t, err, "temuClient.Services.ShipOrderPacking.Match(ctx, %s)", jsonx.ToJson(req, "{}"))
}

func TestShipOrderPackingService_Send(t *testing.T) {
	// 发货地址
	addresses, err := temuClient.Services.MallAddress.All(ctx)
	assert.Nilf(t, err, "temuClient.Services.MallAddress.All(ctx): error")
	assert.Equal(t, true, len(addresses) > 0, "temuClient.Services.MallAddress.All(ctx): results")
	address := addresses[0]

	// 快递公司
	companies, err := temuClient.Services.Logistics.Companies(ctx)
	assert.Nilf(t, err, "temuClient.Services.Logistics.Companies(ctx): error")
	assert.Equal(t, true, len(companies) > 0, "temuClient.Services.Logistics.Companies(ctx): results")
	company := companies[0]

	status := entity.ShipOrderStatusWaitingPacking
	params := ShipOrderQueryParams{
		Status:         IntPtr(status),
		IsPrintBoxMark: IntPtr(1),
	}
	params.PageSize = 100
	items, _, _, _, err := temuClient.Services.ShipOrder.All(ctx, params)
	assert.Nilf(t, err, "temuClient.Services.ShipOrder.All(ctx, %s)", jsonx.ToJson(params, "{}"))
	exists := false
	var shipOrder entity.ShipOrder
	for _, v := range items {
		// 发货单状态异常，存在非待发货状态的发货单，请刷新页面重试
		if v.Status == entity.ShipOrderStatusWaitingPacking {
			exists = true
			shipOrder = v
			break
		}
	}
	assert.Equalf(t, exists, true, "temuClient.Services.ShipOrder.All(ctx, %s)", jsonx.ToJson(params, "{}"))
	if exists {

		// 必须打印箱唛
		if !shipOrder.IsPrintBoxMark {
			_, err = temuClient.Services.Barcode.BoxMark(ctx, shipOrder.DeliveryOrderSn)
			assert.Nilf(t, err, "temuClient.Services.Barcode.BoxMark(ctx, %s)", shipOrder.DeliveryOrderSn)
		}

		d, _ := time.Parse(time.DateTime, "2024-09-01 18:00:00")
		req := ShipOrderPackingSendRequest{
			DeliveryAddressId:   address.ID,
			DeliveryOrderSnList: []string{shipOrder.DeliveryOrderSn},
			DeliverMethod:       IntPtr(entity.DeliveryMethodPlatformRecommendation),
			ThirdPartyDeliveryInfo: &ShipOrderPackingSendRequestPlatformRecommendationDeliveryInformation{
				ExpressCompanyId:          company.ShipId,
				TmsChannelId:              0,
				ExpressCompanyName:        company.ShipName,
				StandbyExpress:            false,
				ExpressDeliverySn:         shipOrder.ExpressDeliverySn,
				PredictTotalPackageWeight: shipOrder.PredictTotalPackageWeight,
				ExpectPickUpGoodsTime:     int(d.Unix()),
				ExpressPackageNum:         len(shipOrder.PackageList),
				MinChargeAmount:           0.01,
				MaxChargeAmount:           0.02,
				PredictId:                 123, // ?
			},
		}
		if req.ThirdPartyDeliveryInfo.PredictTotalPackageWeight < 1000 {
			req.ThirdPartyDeliveryInfo.PredictTotalPackageWeight = 1000
		}
		_, err = temuClient.Services.ShipOrderPacking.Send(ctx, req)
		fmt.Println(fmt.Errorf("sssssssssss %#v", err))
		assert.Nilf(t, err, "temuClient.Services.ShipOrderPacking.Match(ctx, %s)", jsonx.ToJson(req, "{}"))
	}
}
