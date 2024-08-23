package temu

import (
	"github.com/hiscaler/gox/jsonx"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShipOrderService_All(t *testing.T) {
	params := ShipOrderQueryParams{}
	params.Page = 1
	params.PageSize = 10
	_, _, _, _, err := temuClient.Services.ShipOrder.All(params)
	assert.Nilf(t, err, "Services.ShipOrder.All: %s", jsonx.ToJson(params, "{}"))
}

func TestShipOrderService_Create(t *testing.T) {
	deliveryAddress, err := temuClient.Services.MallAddress.One(5441063557369)
	assert.Nil(t, err, "Query mall deliveryAddress")

	subPurchaseOrderSn := "WB2408182975602"
	addr, err := temuClient.Services.ShipOrderReceiveAddress.One(subPurchaseOrderSn)
	assert.Nilf(t, err, "Services.ShipOrderReceiveAddress.One(%s)", subPurchaseOrderSn)
	receiveAddress := addr.ReceiveAddressInfo

	shipOrderStaging, err := temuClient.Services.ShipOrderStaging.One(subPurchaseOrderSn)
	assert.Nil(t, err, "Query shop order staging")

	shipOrderCreateRequestDeliveryOrder := ShipOrderCreateRequestDeliveryOrder{
		DeliveryOrderCreateInfos: make([]ShipOrderCreateRequestOrderInfo, 0),
		ReceiveAddressInfo: ShipOrderCreateRequestReceiveAddress{
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
		SubWarehouseId: shipOrderStaging.SubPurchaseOrderBasicVO.SubWarehouseID,
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
			ProductSkuId:  v.ProductSkuID,
		})
		deliveryOrderCreateInfo.PackageInfos = append(deliveryOrderCreateInfo.PackageInfos, ShipOrderCreateRequestOrderPackage{
			PackageDetailSaveInfos: []ShipOrderCreateRequestPackageInfo{
				{
					SkuNum:       v.ProductSkuPurchaseQuantity,
					ProductSkuId: v.ProductSkuID,
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
	_, err = temuClient.Services.ShipOrder.Create(req)
	assert.Nilf(t, err, "temuClient.Services.ShipOrder.Create(%s)", jsonx.ToJson(req, "{}"))
	if err != nil {
		t.Logf("error: %s", err.Error())
	}
}
