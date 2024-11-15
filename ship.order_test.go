package temu

import (
	"github.com/hiscaler/gox/jsonx"
	"github.com/hiscaler/temu-go/entity"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestShipOrderService_All(t *testing.T) {
	status := entity.PurchaseOrderStatusWaitingMerchantReceive
	params := ShipOrderQueryParams{
		SubPurchaseOrderSnList: []string{"WB2408232992138"},
		// IsCustomProduct: true,
		IsPrintBoxMark: IntPtr(0),
		Status:         IntPtr(status),
	}
	params.Page = 1
	params.PageSize = 22
	_, _, _, _, err := temuClient.Services.ShipOrder.All(ctx, params)
	assert.Nilf(t, err, "Services.ShipOrder.All: %s", jsonx.ToJson(params, "{}"))
}

func TestShipOrderService_Create(t *testing.T) {
	deliveryAddress, err := temuClient.Services.MallAddress.One(ctx, 5441063557369)
	assert.Nil(t, err, "Query mall deliveryAddress")

	subPurchaseOrderSn := "WB2408182975602"
	addr, err := temuClient.Services.ShipOrderReceiveAddress.One(ctx, subPurchaseOrderSn)
	assert.Nilf(t, err, "Services.ShipOrderReceiveAddress.One(ctx, %s)", subPurchaseOrderSn)
	receiveAddress := addr.ReceiveAddressInfo

	shipOrderStaging, err := temuClient.Services.ShipOrderStaging.One(ctx, subPurchaseOrderSn)
	assert.Nil(t, err, "Query shop order staging")

	shipOrderCreateRequestDeliveryOrder := ShipOrderCreateRequestDeliveryOrder{
		DeliveryOrderCreateInfos: make([]ShipOrderCreateRequestOrderInfo, 0),
		ReceiveAddressInfo: entity.ReceiveAddress{
			ProvinceName:  receiveAddress.ProvinceName,
			ProvinceCode:  receiveAddress.ProvinceCode,
			CityName:      receiveAddress.CityName,
			CityCode:      receiveAddress.CityCode,
			DistrictName:  receiveAddress.DistrictName,
			DistrictCode:  receiveAddress.DistrictCode,
			ReceiverName:  receiveAddress.ReceiverName,
			DetailAddress: receiveAddress.DetailAddress,
			Phone:         receiveAddress.Phone,
		},
		SubWarehouseId: shipOrderStaging.SubPurchaseOrderBasicVO.SubWarehouseId,
	}

	deliveryOrderCreateInfo := ShipOrderCreateRequestOrderInfo{
		DeliverOrderDetailInfos: make([]ShipOrderCreateRequestOrderDetailInfo, 0),
		SubPurchaseOrderSn:      subPurchaseOrderSn,
		PackageInfos:            make([]ShipOrderCreateRequestOrderPackage, 0),
		DeliveryAddressId:       deliveryAddress.ID,
	}
	for _, v := range shipOrderStaging.OrderDetailVOList {
		deliveryOrderCreateInfo.DeliverOrderDetailInfos = append(deliveryOrderCreateInfo.DeliverOrderDetailInfos, ShipOrderCreateRequestOrderDetailInfo{
			DeliverSkuNum: v.ProductSkuPurchaseQuantity,
			ProductSkuId:  v.ProductSkuId,
		})
		deliveryOrderCreateInfo.PackageInfos = append(deliveryOrderCreateInfo.PackageInfos, ShipOrderCreateRequestOrderPackage{
			PackageDetailSaveInfos: []ShipOrderCreateRequestPackageInfo{
				{
					SkuNum:       v.ProductSkuPurchaseQuantity,
					ProductSkuId: v.ProductSkuId,
				},
			},
		})
	}
	shipOrderCreateRequestDeliveryOrder.DeliveryOrderCreateInfos = append(shipOrderCreateRequestDeliveryOrder.DeliveryOrderCreateInfos, deliveryOrderCreateInfo)

	req := ShipOrderCreateRequest{
		DeliveryOrderCreateGroupList: []ShipOrderCreateRequestDeliveryOrder{
			shipOrderCreateRequestDeliveryOrder,
		},
	}
	_, err = temuClient.Services.ShipOrder.Create(ctx, req)
	assert.Nilf(t, err, "temuClient.Services.ShipOrder.Create(ctx, %s)", jsonx.ToJson(req, "{}"))
	if err != nil {
		t.Logf("error: %s", err.Error())
	}
}

func TestShipOrderService_Cancel(t *testing.T) {
	status := entity.ShipOrderStatusWaitingPacking
	shipOrders, _, _, _, err := temuClient.Services.ShipOrder.All(ctx, ShipOrderQueryParams{
		Status:   IntPtr(status),
		SortType: IntPtr(0),
	})
	assert.Nil(t, err, "temuClient.Services.ShipOrder.All(ctx, {})")
	n := len(shipOrders)
	if n != 0 {
		shipOrder := shipOrders[rand.Intn(n)]
		deliveryOrderSn := shipOrder.DeliveryOrderSn
		deliveryOrderSn = "FH2408271533488"
		var ok bool
		ok, err = temuClient.Services.ShipOrder.Cancel(ctx, deliveryOrderSn)
		assert.Nilf(t, err, "temuClient.Services.ShipOrder.Cancel(ctx, %s", deliveryOrderSn)
		assert.Truef(t, ok, "temuClient.Services.ShipOrder.Cancel(ctx, %s) ok value", deliveryOrderSn)
	}

}
