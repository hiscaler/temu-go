package temu

import (
	"github.com/hiscaler/gox/jsonx"
	"github.com/hiscaler/temu-go/entity"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v4"
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

// TestShipOrderPackingService_SendForSelf 自送发货
func TestShipOrderPackingService_SendForSelf(t *testing.T) {
	// 发货地址
	addresses, err := temuClient.Services.MallAddress.Query(ctx)
	assert.Nilf(t, err, "temuClient.Services.MallAddress.Query(ctx): error")
	assert.Equal(t, true, len(addresses) > 0, "temuClient.Services.MallAddress.Query(ctx): results")
	address := addresses[0]

	params := ShipOrderQueryParams{
		Status: null.IntFrom(entity.ShipOrderStatusWaitingPacking),
	}
	params.PageSize = 1
	items, _, _, _, err := temuClient.Services.ShipOrder.Query(ctx, params)
	assert.Nilf(t, err, "temuClient.Services.ShipOrder.Query(ctx, %s)", jsonx.ToJson(params, "{}"))
	if len(items) != 0 {
		shipOrder := items[0]
		// 必须打印箱唛
		if !shipOrder.IsPrintBoxMark {
			_, err = temuClient.Services.Goods.Barcode.BoxMark(ctx, shipOrder.DeliveryOrderSn)
			assert.Nilf(t, err, "temuClient.Services.Goods.Barcode.BoxMark(ctx, %s)", shipOrder.DeliveryOrderSn)
		}

		driverName := shipOrder.DriverName
		if driverName == "" {
			driverName = "Zhang San"
		}
		req := ShipOrderPackingSendRequest{
			DeliveryAddressId:   address.ID,
			DeliveryOrderSnList: []string{shipOrder.DeliveryOrderSn},
			DeliverMethod:       null.IntFrom(entity.DeliveryMethodSelf),
			SelfDeliveryInfo: &ShipOrderPackingSendRequestSelfDeliveryInformation{
				// DriverUid:             0,
				DriverName: driverName,
				// PlateNumber:           "",
				// DeliveryContactNumber: "",
				// DeliveryContactAreaNo: "",
				ExpressPackageNum: len(shipOrder.PackageList),
			},
		}
		_, err = temuClient.Services.ShipOrderPacking.Send(ctx, req)
		assert.Nilf(t, err, "temuClient.Services.ShipOrderPacking.Send(ctx, %s)", jsonx.ToJson(req, "{}"))
	} else {
		t.Logf("not found waitingPackage status purchase order")
	}
}

// TestShipOrderPackingService_SendForPlatformRecommendation 平台推荐物流发货
func TestShipOrderPackingService_SendForPlatformRecommendation(t *testing.T) {
	// 发货地址
	addresses, err := temuClient.Services.MallAddress.Query(ctx)
	assert.Nilf(t, err, "temuClient.Services.MallAddress.Query(ctx): error")
	assert.Equal(t, true, len(addresses) > 0, "temuClient.Services.MallAddress.Query(ctx): results")
	address := addresses[0]

	// 快递公司
	companies, err := temuClient.Services.Logistics.Companies(ctx)
	assert.Nilf(t, err, "temuClient.Services.Logistics.Companies(ctx): error")
	assert.Equal(t, true, len(companies) > 0, "temuClient.Services.Logistics.Companies(ctx): results")
	company := companies[0]

	params := ShipOrderQueryParams{
		Status: null.IntFrom(entity.ShipOrderStatusWaitingPacking),
	}
	params.PageSize = 1
	items, _, _, _, err := temuClient.Services.ShipOrder.Query(ctx, params)
	assert.Nilf(t, err, "temuClient.Services.ShipOrder.Query(ctx, %s)", jsonx.ToJson(params, "{}"))
	if len(items) != 0 {
		shipOrder := items[0]
		// 必须打印箱唛
		if !shipOrder.IsPrintBoxMark {
			_, err = temuClient.Services.Goods.Barcode.BoxMark(ctx, shipOrder.DeliveryOrderSn)
			assert.Nilf(t, err, "temuClient.Services.Goods.Barcode.BoxMark(ctx, %s)", shipOrder.DeliveryOrderSn)
		}

		d, _ := time.ParseInLocation(time.DateTime, time.Now().Format(time.DateOnly)+" 18:00:00", temuClient.TimeLocation)
		req := ShipOrderPackingSendRequest{
			DeliveryAddressId:   address.ID,
			DeliveryOrderSnList: []string{shipOrder.DeliveryOrderSn},
			DeliverMethod:       null.IntFrom(entity.DeliveryMethodPlatformRecommendation),
			ThirdPartyDeliveryInfo: &ShipOrderPackingSendRequestPlatformRecommendationDeliveryInformation{
				ExpressCompanyId:          company.ShipId,
				TmsChannelId:              0,
				ExpressCompanyName:        company.ShipName,
				StandbyExpress:            false,
				ExpressDeliverySn:         shipOrder.ExpressDeliverySn,
				PredictTotalPackageWeight: shipOrder.PredictTotalPackageWeight,
				ExpectPickUpGoodsTime:     d.UnixMilli(),
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
		assert.Nilf(t, err, "temuClient.Services.ShipOrderPacking.Send(ctx, %s)", jsonx.ToJson(req, "{}"))
	} else {
		t.Logf("not found waitingPackage status purchase order")
	}
}
