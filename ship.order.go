package temu

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

// 发货单服务
type shipOrderService service

type ShipOrderQueryParams struct {
	normal.ParameterWithPager
	DeliveryOrderSnList    []string `json:"deliveryOrderSnList,omitempty"`    // 发货单号列表
	SubPurchaseOrderSnList []string `json:"subPurchaseOrderSnList,omitempty"` // 子采购单号列表
	ExpressDeliverySnList  []string `json:"expressDeliverySnList,omitempty"`  // 快递单号列表
	SkcExtCodeList         []string `json:"skcExtCodeList,omitempty"`         // 货号列表
	ProductSkcIdList       []int    `json:"productSkcIdList,omitempty"`       // skcId列表
	SubWarehouseIdList     []int    `json:"subWarehouseIdList,omitempty"`     // 收货子仓列表
	DeliverTimeFrom        int      `json:"deliverTimeFrom,omitempty"`        // 发货时间-开始时间
	DeliverTimeTo          int      `json:"deliverTimeTo,omitempty"`          // 发货时间-结束时间
	SettlementType         int      `json:"settlementType,omitempty"`         // 结算类型 0-非vmi 1-vmi
	IsFirstOrder           bool     `json:"isFirstOrder,omitempty"`           // 是否首单
	UrgencyType            bool     `json:"urgencyType,omitempty"`            // 是否是紧急发货单，0-普通 1-急采
	IsJit                  bool     `json:"isJit,omitempty"`                  // 是否是jit，true:jit
	PurchaseStockType      int      `json:"purchaseStockType,omitempty"`      // 备货类型 0-普通备货 1-jit备货
	IsCustomProduct        int      `json:"isCustomProduct,omitempty"`        // 是否为定制品
	SubWarehouseId         int      `json:"subWarehouseId,omitempty"`         // 收货子仓
	InventoryRegion        []int    `json:"inventoryRegion,omitempty"`        // DOMESTIC(1, "国内备货"), OVERSEAS(2, "海外备货"), BOUNDED_WAREHOUSE(3, "保税仓备货"),
}

func (m ShipOrderQueryParams) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.SettlementType, validation.When(!validation.IsEmpty(m.SettlementType), validation.In(entity.SettlementTypeVMI, entity.SettlementTypeNotVMI).Error("无效的结算类型。"))),
		validation.Field(&m.UrgencyType, validation.When(!validation.IsEmpty(m.UrgencyType), validation.In(entity.ShipOrderTypeNormal, entity.ShipOrderTypeUrgency).Error("无效的加急类型。"))),
		validation.Field(&m.PurchaseStockType, validation.When(!validation.IsEmpty(m.PurchaseStockType), validation.In(entity.PurchaseStockTypeNormal, entity.PurchaseStockTypeJIT).Error("无效的备货类型。"))),
	)
}

// All 查询发货单 V2
// https://seller.kuajingmaihuo.com/sop/view/889973754324016047#B7c51j
func (s shipOrderService) All(params ShipOrderQueryParams) (items []entity.ShipOrder, total, totalPages int, isLastPage bool, err error) {
	params.TidyPager()
	if err = params.Validate(); err != nil {
		return
	}

	var result = struct {
		normal.Response
		Result struct {
			Total int                `json:"total"`
			List  []entity.ShipOrder `json:"list"`
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetBody(params).
		SetResult(&result).
		Post("bg.shiporderv2.get")
	if err == nil {
		err = parseResponse(resp, result.Response)
	}
	if err != nil {
		return
	}

	items = result.Result.List
	total, totalPages, isLastPage = parseResponseTotal(params.Page, params.PageSize, result.Result.Total)

	return
}

// 创建发货单
// https://seller.kuajingmaihuo.com/sop/view/889973754324016047#HqGnA0

type ShipOrderCreateRequestOrderDetailInfo struct {
	DeliverSkuNum int `json:"deliverSkuNum"` // 发货sku数目
	ProductSkuId  int `json:"productSkuId"`  // 定制品，传定制品id；非定制品，传货品 skuId
}

// ShipOrderCreateRequestOrderPackage 包裹信息
type ShipOrderCreateRequestOrderPackage struct {
	PackageDetailSaveInfos []ShipOrderCreateRequestPackageInfo `json:"packageDetailSaveInfos"` // 包裹明细
}

type ShipOrderCreateRequestPackageInfo struct {
	SkuNum       int `json:"skuNum"`       // 发货 sku 数目
	ProductSkuId int `json:"productSkuId"` // skuId
}

type ShipOrderCreateRequestOrderInfo struct {
	DeliverOrderDetailInfos []ShipOrderCreateRequestOrderDetailInfo `json:"deliverOrderDetailInfos"` // 采购单创建信息列表
	SubPurchaseOrderSn      string                                  `json:"subPurchaseOrderSn"`      // 采购子单号
	PackageInfos            []ShipOrderCreateRequestOrderPackage    `json:"packageInfos"`            //	包裹信息列表
	DeliveryAddressId       int                                     `json:"deliveryAddressId"`       // 发货地址 ID
}

// ShipOrderCreateRequestReceiveAddress 收货地址
type ShipOrderCreateRequestReceiveAddress struct {
	ProvinceName  string `json:"provinceName,omitempty"`  // 省
	ProvinceCode  int    `json:"provinceCode,omitempty"`  // 省份编码
	CityName      string `json:"cityName,omitempty"`      // 市
	CityCode      int    `json:"cityCode,omitempty"`      // 市编码
	DistrictName  string `json:"districtName,omitempty"`  // 区
	DistrictCode  int    `json:"districtCode,omitempty"`  // 区编码
	ReceiverName  string `json:"receiverName,omitempty"`  // 收货人
	DetailAddress string `json:"detailAddress,omitempty"` // 详细地址
	Phone         string `json:"phone,omitempty"`         // 联系电话
}

type ShipOrderCreateRequestDeliveryOrder struct {
	DeliveryOrderCreateInfos []ShipOrderCreateRequestOrderInfo    `json:"deliveryOrderCreateInfos"` // 发货单创建组列表
	ReceiveAddressInfo       ShipOrderCreateRequestReceiveAddress `json:"receiveAddressInfo"`       // 收货地址
	SubWarehouseId           int                                  `json:"subWarehouseId"`           // 子仓 ID
}

type ShipOrderCreateRequest struct {
	normal.Parameter
	DeliveryOrderCreateGroupList []ShipOrderCreateRequestDeliveryOrder `json:"deliveryOrderCreateGroupList"` // 发货单创建组列表
}

func (m ShipOrderCreateRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.DeliveryOrderCreateGroupList, validation.Required.Error("发货单创建组列表不能为空。")),
	)
}

// Create 创建发货单接口 V3
func (s shipOrderService) Create(req ShipOrderCreateRequest) (ok bool, err error) {
	if err = req.Validate(); err != nil {
		return
	}

	var result = struct {
		normal.Response
		Result any `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetBody(req).
		SetResult(&result).
		Post("bg.shiporderv3.create")
	if err == nil {
		err = parseResponse(resp, result.Response)
	}
	ok = err == nil

	return
}

// 取消发货单

type ShipOrderCancelRequest struct {
	normal.Parameter
	DeliveryOrderSn int `json:"deliveryOrderSn"` // 发货单 ID
}

func (m ShipOrderCancelRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.DeliveryOrderSn, validation.Required.Error("发货单 ID 不能为空。")),
	)
}

// Cancel 取消发货单
// https://seller.kuajingmaihuo.com/sop/view/889973754324016047#UywT8E
func (s shipOrderService) Cancel(deliveryOrderSn int) (ok bool, err error) {
	req := ShipOrderCancelRequest{DeliveryOrderSn: deliveryOrderSn}
	if err = req.Validate(); err != nil {
		return
	}

	var result = struct {
		normal.Response
		Result any `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetBody(req).
		SetResult(&result).
		Post("bg.shiporder.cancel")
	if err == nil {
		err = parseResponse(resp, result.Response)
	}
	if err != nil {
		return
	}

	return true, nil
}
