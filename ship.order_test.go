package temu

import (
	"github.com/hiscaler/gox/jsonx"
	"github.com/hiscaler/temu-go/entity"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v4"
	"math/rand"
	"testing"
)

func TestShipOrderService_Query(t *testing.T) {
	// status := entity.PurchaseOrderStatusMerchantReceived
	params := ShipOrderQueryParams{
		DeliveryOrderSnList: []string{"FH2411182674018"},
		Status:              null.IntFrom(1),
		// SubPurchaseOrderSnList: []string{"WB2409161373926"},
		// IsCustomProduct: true,
		// IsPrintBoxMark: IntPtr(1),
		// Status:             IntPtr(0),
		// SubWarehouseIdList: []int64{438429773460},
	}
	params.Page = 1
	params.PageSize = 1
	_, _, _, _, err := temuClient.Services.ShipOrder.Query(ctx, params)
	assert.Nilf(t, err, "Services.ShipOrder.Query: %s", jsonx.ToJson(params, "{}"))
}

func TestShipOrderService_Create(t *testing.T) {
	deliveryAddress, err := temuClient.Services.Mall.DeliveryAddress.One(ctx, 5441063557369)
	assert.Nil(t, err, "Query mall deliveryAddress")

	subPurchaseOrderSn := "WB2408182975602"
	addr, err := temuClient.Services.ShipOrder.ReceiveAddress.One(ctx, subPurchaseOrderSn)
	assert.Nilf(t, err, "Services.ShipOrder.ReceiveAddress.One(ctx, %s)", subPurchaseOrderSn)
	receiveAddress := addr.ReceiveAddressInfo

	shipOrderStaging, err := temuClient.Services.ShipOrder.Staging.One(ctx, subPurchaseOrderSn)
	assert.Nil(t, err, "Query shop order staging")

	shipOrderCreateRequestDeliveryOrder := ShipOrderCreateRequestDeliveryOrder{
		DeliveryOrderCreateInfos: make([]ShipOrderCreateRequestOrderInfo, 0),
		ReceiveAddressInfo: &entity.ReceiveAddress{
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
		SubWarehouseId: null.IntFrom(shipOrderStaging.SubPurchaseOrderBasicVO.SubWarehouseId),
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

// TestShipOrderService_SimpleCreate 只上传部分基础数据（备货单号、发货地址 ID），其他部分系统自动填充
func TestShipOrderService_SimpleCreate(t *testing.T) {
	deliveryAddress, err := temuClient.Services.Mall.DeliveryAddress.One(ctx, 5441063557369)
	assert.Nil(t, err, "Query mall deliveryAddress")
	req := ShipOrderCreateRequest{
		DeliveryOrderCreateGroupList: []ShipOrderCreateRequestDeliveryOrder{},
	}
	for _, purchaseOrderNumber := range []string{"WB2408182975602"} {
		req.DeliveryOrderCreateGroupList = append(req.DeliveryOrderCreateGroupList, ShipOrderCreateRequestDeliveryOrder{
			DeliveryOrderCreateInfos: []ShipOrderCreateRequestOrderInfo{
				{
					SubPurchaseOrderSn: purchaseOrderNumber,
					DeliveryAddressId:  deliveryAddress.ID,
				},
			},
		})
	}
	_, err = temuClient.Services.ShipOrder.Create(ctx, req)
	assert.Nilf(t, err, "temuClient.Services.ShipOrder.Create(ctx, %s)", jsonx.ToJson(req, "{}"))
	if err != nil {
		t.Logf("error: %s", err.Error())
	}
}

func TestShipOrderService_Cancel(t *testing.T) {
	status := entity.ShipOrderStatusWaitingPacking
	shipOrders, _, _, _, err := temuClient.Services.ShipOrder.Query(ctx, ShipOrderQueryParams{
		Status:   null.IntFrom(int64(status)),
		SortType: null.IntFrom(0),
	})
	assert.Nil(t, err, "temuClient.Services.ShipOrder.Query(ctx, {})")
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
