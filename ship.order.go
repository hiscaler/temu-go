package temu

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

// 发货单服务
type shipOrderService service

type ShipOrderQueryParams struct {
	normal.ParameterWithPager
	DeliveryOrderSnList      []string `json:"deliveryOrderSnList,omitempty"`      // 发货单号列表
	SubPurchaseOrderSnList   []string `json:"subPurchaseOrderSnList,omitempty"`   // 子采购单号列表
	ExpressDeliverySnList    []string `json:"expressDeliverySnList,omitempty"`    // 快递单号列表
	SkcExtCodeList           []string `json:"skcExtCodeList,omitempty"`           // 货号列表
	ProductSkcIdList         []int    `json:"productSkcIdList,omitempty"`         // skcId 列表
	SubWarehouseIdList       []int    `json:"subWarehouseIdList,omitempty"`       // 收货子仓列表
	DeliverTimeFrom          int      `json:"deliverTimeFrom,omitempty"`          // 发货时间-开始时间
	DeliverTimeTo            int      `json:"deliverTimeTo,omitempty"`            // 发货时间-结束时间
	Status                   *int     `json:"status,omitempty"`                   // 发货单状态，0：待装箱发货，1：待仓库收货，2：已收货，3：已入库，4：已退货，5：已取消，6：部分收货，查询发货批次时仅支持查询发货单状态=1
	UrgencyType              bool     `json:"urgencyType,omitempty"`              // 是否是紧急发货单，0-普通，1-急采
	IsCustomProduct          bool     `json:"isCustomProduct,omitempty"`          // 是否为定制品
	IsVim                    *int     `json:"isVmi,omitempty"`                    // 是否是vmi，0-非VMI，1-VMI
	IsJit                    bool     `json:"isJit,omitempty"`                    // 是否是jit，true:jit
	LatestFeedbackStatusList []int    `json:"latestFeedbackStatusList,omitempty"` // 最新反馈状态列表，0-当前无异常，1-已提交，2-物流商处理中，3-已撤销，4-已反馈
	SortType                 *int     `json:"sortType,omitempty"`                 // 排序类型，0-创建时间最新在上，1-要求发货时间较早在上，2-按照仓库名称排序
	InventoryRegion          []int    `json:"inventoryRegion,omitempty"`          // 发货区域，1-国内备货，2-海外备货，3-保税仓备货
	IsPrintBoxMark           *int     `json:"isPrintBoxMark,omitempty"`           // 是否已打印商品打包标签，0-未打印，1-已打印
	TargetReceiveAddress     string   `json:"targetReceiveAddress,omitempty"`     // 筛选项-收货地址（精准匹配）
	TargetDeliveryAddress    string   `json:"targetDeliveryAddress,omitempty"`    // 筛选项-发货地址（精准匹配）
}

func (m ShipOrderQueryParams) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.UrgencyType, validation.When(!validation.IsEmpty(m.UrgencyType), validation.In(entity.ShipOrderTypeNormal, entity.ShipOrderTypeUrgency).Error("无效的加急类型。"))),
	)
}

// All 查询发货单 V2
// https://seller.kuajingmaihuo.com/sop/view/889973754324016047#B7c51j
func (s shipOrderService) All(ctx context.Context, params ShipOrderQueryParams) (items []entity.ShipOrder, total, totalPages int, isLastPage bool, err error) {
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
		SetContext(ctx).
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

type ShipOrderCreateRequestDeliveryOrder struct {
	DeliveryOrderCreateInfos []ShipOrderCreateRequestOrderInfo `json:"deliveryOrderCreateInfos"` // 发货单创建组列表
	ReceiveAddressInfo       entity.ReceiveAddress             `json:"receiveAddressInfo"`       // 收货地址
	SubWarehouseId           int                               `json:"subWarehouseId"`           // 子仓 ID
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
func (s shipOrderService) Create(ctx context.Context, req ShipOrderCreateRequest) (ok bool, err error) {
	if err = req.Validate(); err != nil {
		return
	}

	var result = struct {
		normal.Response
		Result any `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
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
	DeliveryOrderSn string `json:"deliveryOrderSn"` // 发货单 ID
}

func (m ShipOrderCancelRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.DeliveryOrderSn, validation.Required.Error("发货单号不能为空。")),
	)
}

// Cancel 取消发货单
// https://seller.kuajingmaihuo.com/sop/view/889973754324016047#UywT8E
func (s shipOrderService) Cancel(ctx context.Context, deliveryOrderSn string) (ok bool, err error) {
	req := ShipOrderCancelRequest{DeliveryOrderSn: deliveryOrderSn}
	if err = req.Validate(); err != nil {
		return
	}

	var result = struct {
		normal.Response
		Result any `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(req).
		SetResult(&result).
		Post("bg.shiporder.cancel")
	if err == nil {
		err = parseResponse(resp, result.Response)
	}

	ok = err == nil
	return
}

// ThirdPartyLogisticsCompanies 自行委托三方物流公司查询接口
func (s shipOrderService) ThirdPartyLogisticsCompanies(ctx context.Context) (companies []entity.LogisticsExpressCompany, err error) {
	var result = struct {
		normal.Response
		Result []entity.LogisticsExpressCompany `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetResult(&result).
		Post("bg.shiporder.logistics.get")
	if err == nil {
		err = parseResponse(resp, result.Response)
	}
	if err != nil {
		return
	}

	return result.Result, nil
}
