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
	address, err := temuClient.Services.MallAddressService.One(5441063557369)
	assert.Nil(t, err, "Query mall address")

	subPurchaseOrderSn := "WB2408173170013"
	shipOrderStaging, err := temuClient.Services.ShipOrderStaging.One(subPurchaseOrderSn)
	assert.Nil(t, err, "Query shop order staging")

	req := ShipOrderCreateRequest{}
	oi := ShipOrderCreateRequestDeliveryOrder{
		ReceiveAddressInfo: ShipOrderCreateRequestReceiveAddress{
			ProvinceName:  address.ProvinceName,
			ProvinceCode:  address.ProvinceCode,
			CityName:      address.CityName,
			CityCode:      address.CityCode,
			DistrictName:  address.DistrictName,
			DistrictCode:  address.DistrictCode,
			ReceiverName:  "",
			DetailAddress: address.Address,
			Phone:         "",
		},
		SubWarehouseId: shipOrderStaging.SubPurchaseOrderBasicVO.SubWarehouseID,
		DeliveryOrderCreateInfos: []ShipOrderCreateRequestOrderInfo{
			{
				DeliverOrderDetailInfos: []ShipOrderCreateRequestOrderDetailInfo{
					{
						DeliverSkuNum: 0,
						ProductSkuId:  0,
					},
				},
				SubPurchaseOrderSn: subPurchaseOrderSn,
				PackageInfos: []ShipOrderCreateRequestOrderPackage{
					{
						PackageDetailSaveInfos: []ShipOrderCreateRequestPackageInfo{
							{
								SkuNum:       0,
								ProductSkuId: 0,
							},
						},
					},
				},
				DeliveryAddressId: address.ID,
			},
		},
	}
	req.DeliveryOrderCreateGroupList = []ShipOrderCreateRequestDeliveryOrder{
		oi,
	}
	_, err = temuClient.Services.ShipOrder.Create(req)
	assert.Nilf(t, err, "temuClient.Services.ShipOrder.Create(%s)", jsonx.ToJson(req, "{}"))
	if err != nil {
		t.Logf("error: %s", err.Error())
	}
}
