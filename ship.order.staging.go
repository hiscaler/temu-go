package temu

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

// 发货台服务
type shipOrderStagingService service

type StagingQueryParams struct {
	normal.Parameter
	SubPurchaseOrderSnList []string `json:"subPurchaseOrderSnList,omitempty"` // 子采购单号列表
	SkcExtCode             []string `json:"skcExtCode,omitempty"`             // 货号列表
	ProductSkcIdList       []string `json:"productSkcIdList,omitempty"`       // skcId列表
	SettlementType         int      `json:"settlementType,omitempty"`         // 结算类型 0-非vmi 1-vmi
	IsFirstOrder           bool     `json:"isFirstOrder,omitempty"`           // 是否首单
	UrgencyType            bool     `json:"urgencyType,omitempty"`            // 是否是紧急发货单，0-普通 1-急采
	IsJit                  bool     `json:"isJit"`                            // 是否是jit，true:jit
	Page                   int      `json:"pageNo"`                           // 页号， 从1开始
	PageSize               int      `json:"pageSize"`                         // 每页记录数不能为空
	PurchaseStockType      int      `json:"purchaseStockType,omitempty"`      // 备货类型 0-普通备货 1-jit备货
	IsCustomProduct        int      `json:"isCustomProduct,omitempty"`        // 是否为定制品
	SubWarehouseId         int      `json:"subWarehouseId,omitempty"`         // 收货子仓
	InventoryRegion        []int    `json:"inventoryRegion,omitempty"`        // DOMESTIC(1, "国内备货"), OVERSEAS(2, "海外备货"), BOUNDED_WAREHOUSE(3, "保税仓备货"),
	// Request struct {
	//
	// } `json:"request"`
}

func (m StagingQueryParams) Validate() error {
	return nil
	// return validation.ValidateStruct(&m,
	// 	validation.Field(&m.Request, validation.When(m.Request != nil, validation.By(func(value interface{}) error {
	//
	// 		return nil
	// 	}))),
	// )
}

// All List all staging orders
// https://seller.kuajingmaihuo.com/sop/view/889973754324016047#NOA03y
func (s shipOrderStagingService) All(params StagingQueryParams) (items []entity.ShipOrderStaging, total, totalPages int, isLastPage bool, err error) {
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 10
	} else if params.PageSize > 500 {
		params.PageSize = 500
	}
	if err = params.Validate(); err != nil {
		return
	}

	var result = struct {
		normal.Response
		Result struct {
			Total int                       `json:"total"`
			List  []entity.ShipOrderStaging `json:"list"`
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().SetBody(params).SetResult(&result).Post("bg.shiporder.staging.get")
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

// One 搜索单个发货台数据
func (s shipOrderStagingService) One(subPurchaseOrderSn string) (item entity.ShipOrderStaging, err error) {
	params := StagingQueryParams{
		Page:                   1,
		PageSize:               10,
		SubPurchaseOrderSnList: []string{subPurchaseOrderSn},
	}

	items, _, _, _, err := s.All(params)
	if err != nil {
		return
	}
	if len(items) == 0 {
		return item, ErrNotFound
	}

	return items[0], nil
}

// 加入发货台

type ShipOrderStagingAddInfo struct {
	SubPurchaseOrderSn  string `json:"subPurchaseOrderSn"`  // 采购子单号（订单号）
	DeliveryAddressType int    `json:"deliveryAddressType"` // 发货地址类型，1-内地，2-香港，内地主体（店铺货币选择CNY，默认入参1，其余主体选择2）
}

type ShipOrderStagingAddRequest struct {
	normal.Parameter
	JoinInfoList []ShipOrderStagingAddInfo `json:"joinInfoList"`
}

func (m ShipOrderStagingAddRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.JoinInfoList, validation.Required.Error("发货数据不能为空。")),
	)
}

// Add 加入发货台
// https://seller.kuajingmaihuo.com/sop/view/889973754324016047#YSg2AE
func (s shipOrderStagingService) Add(req ShipOrderStagingAddRequest) (ok bool, err error) {
	if err = req.Validate(); err != nil {
		return
	}

	var result = struct {
		normal.Response
		Result any `json:"result"`
	}{}
	resp, err := s.httpClient.R().SetBody(req).SetResult(&result).Post("bg.shiporder.staging.add")
	if err == nil {
		err = parseResponse(resp, result.Response)
	}
	ok = err == nil
	return
}
