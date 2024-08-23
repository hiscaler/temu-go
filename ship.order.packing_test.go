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
	_, err := temuClient.Services.ShipOrderPacking.Match(req)
	assert.Nilf(t, err, "temuClient.Services.ShipOrderPacking.Match(%s)", jsonx.ToJson(req, "{}"))
}

func TestShipOrderPackingService_Send(t *testing.T) {
	// 发货地址
	addresses, err := temuClient.Services.MallAddress.All()
	assert.Nilf(t, err, "temuClient.Services.MallAddress.All(): error")
	assert.Equal(t, true, len(addresses) > 0, "temuClient.Services.MallAddress.All(): results")
	address := addresses[0]

	// 快递公司
	companies, err := temuClient.Services.Logistics.Companies()
	assert.Nilf(t, err, "temuClient.Services.Logistics.Companies(): error")
	assert.Equal(t, true, len(companies) > 0, "temuClient.Services.Logistics.Companies(): results")
	company := companies[0]

	params := ShipOrderQueryParams{
		Status:         entity.ShipOrderStatusWaitingPacking, // 这条件没用？
		IsPrintBoxMark: 1,
	}
	params.PageSize = 100
	items, _, _, _, err := temuClient.Services.ShipOrder.All(params)
	assert.Nilf(t, err, "temuClient.Services.ShipOrder.All(%s)", jsonx.ToJson(params, "{}"))
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
	assert.Equalf(t, exists, true, "temuClient.Services.ShipOrder.All(%s)", jsonx.ToJson(params, "{}"))
	if exists {

		// 必须打印箱唛
		if !shipOrder.IsPrintBoxMark {
			_, err = temuClient.Services.Barcode.BoxMark(shipOrder.DeliveryOrderSn)
			assert.Nilf(t, err, "temuClient.Services.Barcode.BoxMark(%s)", shipOrder.DeliveryOrderSn)
		}

		req := ShipOrderPackingSendRequest{
			DeliveryAddressId:   address.ID,
			DeliveryOrderSnList: []string{shipOrder.DeliveryOrderSn},
			DeliverMethod:       entity.DeliveryMethodPlatformRecommendation,
			ThirdPartyDeliveryInfo: &ShipOrderPackingSendRequestPlatformRecommendationDeliveryInformation{
				ExpressCompanyId:          company.ShipId,
				TmsChannelId:              0,
				ExpressCompanyName:        company.ShipName,
				StandbyExpress:            false,
				ExpressDeliverySn:         shipOrder.ExpressDeliverySn,
				PredictTotalPackageWeight: shipOrder.PredictTotalPackageWeight,
				ExpectPickUpGoodsTime:     int(time.Now().Unix()) + 100,
				ExpressPackageNum:         len(shipOrder.PackageList),
				MinChargeAmount:           0.01,
				MaxChargeAmount:           0.02,
				PredictId:                 123, // ?
			},
		}
		if req.ThirdPartyDeliveryInfo.PredictTotalPackageWeight < 1000 {
			req.ThirdPartyDeliveryInfo.PredictTotalPackageWeight = 1000
		}
		_, err = temuClient.Services.ShipOrderPacking.Send(req)
		fmt.Println(fmt.Errorf("sssssssssss %#v", err))
		assert.Nilf(t, err, "temuClient.Services.ShipOrderPacking.Match(%s)", jsonx.ToJson(req, "{}"))
	}
}
