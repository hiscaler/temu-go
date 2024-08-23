package temu

import (
	"fmt"
	"github.com/hiscaler/gox/jsonx"
	"github.com/hiscaler/temu-go/entity"
	"github.com/stretchr/testify/assert"
	"math/rand"
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
		Status:         entity.ShipOrderStatusWaitingPacking,
		IsPrintBoxMark: 1,
	}
	items, _, _, _, err := temuClient.Services.ShipOrder.All(params)
	assert.Nilf(t, err, "temuClient.Services.ShipOrder.All(%s)", jsonx.ToJson(params, "{}"))
	n := len(items)
	if n != 0 {
		shipOrder := items[rand.Intn(n)]
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
				ExpectPickUpGoodsTime:     int(time.Now().Unix()),
				ExpressPackageNum:         len(shipOrder.PackageList),
				MinChargeAmount:           0.01,
				MaxChargeAmount:           0.02,
				PredictId:                 123, // ?
			},
		}
		_, err = temuClient.Services.ShipOrderPacking.Send(req)
		fmt.Println(fmt.Errorf("sssssssssss %#v", err))
		assert.Nilf(t, err, "temuClient.Services.ShipOrderPacking.Match(%s)", jsonx.ToJson(req, "{}"))
	}
}
