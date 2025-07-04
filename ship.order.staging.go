package temu

import (
	"context"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/gox/jsonx"
	"github.com/hiscaler/gox/nullx"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
	"github.com/hiscaler/temu-go/validators/is"
	"gopkg.in/guregu/null.v4"
	"strings"
)

// 发货台服务
type shipOrderStagingService service

type ShipOrderStagingQueryParams struct {
	normal.ParameterWithPager
	SubPurchaseOrderSnList []string  `json:"subPurchaseOrderSnList,omitempty"` // 子采购单号列表
	SkcExtCode             []string  `json:"skcExtCode,omitempty"`             // 货号列表
	ProductSkcIdList       []int64   `json:"productSkcIdList,omitempty"`       // skcId列表
	SettlementType         null.Int  `json:"settlementType,omitempty"`         // 结算类型 0-非vmi 1-vmi
	IsFirstOrder           null.Bool `json:"isFirstOrder,omitempty"`           // 是否首单
	UrgencyType            null.Int  `json:"urgencyType,omitempty"`            // 是否是紧急发货单，0-普通 1-急采
	IsJit                  null.Bool `json:"isJit,omitempty"`                  // 是否是jit，true:jit
	PurchaseStockType      null.Int  `json:"purchaseStockType,omitempty"`      // 备货类型 0-普通备货 1-jit备货
	IsCustomProduct        null.Bool `json:"isCustomProduct,omitempty"`        // 是否为定制品
	SubWarehouseId         int64     `json:"subWarehouseId,omitempty"`         // 收货子仓
	InventoryRegion        []int     `json:"inventoryRegion,omitempty"`        // DOMESTIC(1, "国内备货"), OVERSEAS(2, "海外备货"), BOUNDED_WAREHOUSE(3, "保税仓备货"),
	OrderType              null.Int  `json:"orderType,omitempty"`              // 订单类型（1：普通备货单、2：JIT 备货单、3：定制备货单）此参数为扩展参数，用于简化备货类型查询处理
}

func (m ShipOrderStagingQueryParams) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.SubPurchaseOrderSnList, validation.Each(validation.By(is.PurchaseOrderNumber()))),
		validation.Field(&m.SettlementType, validation.When(m.SettlementType.Valid,
			validation.By(func(value interface{}) error {
				v, ok := value.(null.Int)
				if !ok {
					return errors.New("无效的结算类型")
				}

				return validation.Validate(int(v.Int64), validation.In(entity.SettlementTypeVMI, entity.SettlementTypeNotVMI).Error("无效的结算类型"))
			}),
		)),
		validation.Field(&m.PurchaseStockType, validation.When(m.PurchaseStockType.Valid,
			validation.By(func(value interface{}) error {
				v, ok := value.(null.Int)
				if !ok {
					return errors.New("无效的备货类型")
				}

				return validation.Validate(int(v.Int64), validation.In(entity.PurchaseStockTypeNormal, entity.PurchaseStockTypeJIT).Error("无效的备货类型"))
			}),
		)),
	)
}

// Query Queries staging orders
// https://seller.kuajingmaihuo.com/sop/view/889973754324016047#NOA03y
func (s shipOrderStagingService) Query(ctx context.Context, params ShipOrderStagingQueryParams) (items []entity.ShipOrderStaging, total, totalPages int, isLastPage bool, err error) {
	params.TidyPager()
	if params.OrderType.Valid {
		switch params.OrderType.Int64 {
		case entity.OrderTypeNormal:
			params.IsCustomProduct = null.BoolFrom(false)
			params.IsJit = null.BoolFrom(false)

		case entity.OrderTypeJIT:
			params.IsCustomProduct = null.BoolFrom(false)
			params.IsJit = null.BoolFrom(true)

		case entity.OrderTypeCustomized:
			params.IsCustomProduct = null.BoolFrom(true)
			params.IsJit = null.BoolFrom(false)
		}
		params.OrderType = null.NewInt(0, false)
	}
	if err = params.validate(); err != nil {
		err = invalidInput(err)
		return
	}

	var result = struct {
		normal.Response
		Result struct {
			Total int                       `json:"total"`
			List  []entity.ShipOrderStaging `json:"list"`
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(params).
		SetResult(&result).
		Post("bg.shiporder.staging.get")
	if err = recheckError(resp, result.Response, err); err != nil {
		return
	}

	items = result.Result.List
	for i, item := range items {
		var orderType null.Int
		if item.SubPurchaseOrderBasicVO.IsCustomProduct {
			orderType = null.IntFrom(int64(entity.OrderTypeCustomized))
		} else if item.SubPurchaseOrderBasicVO.PurchaseStockType == entity.PurchaseStockTypeJIT {
			orderType = null.IntFrom(int64(entity.OrderTypeJIT))
		} else {
			orderType = null.IntFrom(int64(entity.OrderTypeNormal))
		}
		items[i].OrderType = orderType
	}
	total, totalPages, isLastPage = parseResponseTotal(params.Page, params.PageSize, result.Result.Total)
	return
}

// One 根据备货单号搜索发货台数据
func (s shipOrderStagingService) One(ctx context.Context, purchaseOrderNumber string) (item entity.ShipOrderStaging, err error) {
	err = validation.Validate(purchaseOrderNumber, validation.By(is.PurchaseOrderNumber()))
	if err != nil {
		return item, invalidInput(err)
	}

	items, _, _, _, err := s.Query(ctx, ShipOrderStagingQueryParams{
		SubPurchaseOrderSnList: []string{purchaseOrderNumber},
	})
	if err != nil {
		return
	}

	purchaseOrderNumber = strings.ToLower(purchaseOrderNumber)
	for _, v := range items {
		if strings.ToLower(v.SubPurchaseOrderBasicVO.SubPurchaseOrderSn) == purchaseOrderNumber {
			return v, nil
		}
	}
	return item, ErrNotFound
}

// 加入发货台

type ShipOrderStagingAddInfo struct {
	SubPurchaseOrderSn  string `json:"subPurchaseOrderSn"`  // 采购子单号（订单号）
	DeliveryAddressType int    `json:"deliveryAddressType"` // 发货地址类型，1-内地，2-香港，内地主体（店铺货币选择CNY，默认入参1，其余主体选择2）
}

func (m ShipOrderStagingAddInfo) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.SubPurchaseOrderSn,
			validation.Required.Error("备货单号不能为空"),
			validation.By(is.PurchaseOrderNumber()),
		),
		validation.Field(&m.DeliveryAddressType,
			validation.
				In(entity.DeliveryAddressTypeChineseMainland, entity.DeliveryAddressTypeChineseHongKong).
				Error("无效的发货地址类型"),
		),
	)
}

type ShipOrderStagingAddRequest struct {
	normal.Parameter
	JoinInfoList []ShipOrderStagingAddInfo `json:"joinInfoList"` // 加入发货台的信息列表
}

func (m ShipOrderStagingAddRequest) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.JoinInfoList,
			validation.Required.Error("发货数据不能为空"),
			validation.Each(validation.By(func(value interface{}) error {
				v, ok := value.(ShipOrderStagingAddInfo)
				if !ok {
					return errors.New("无效的发货数据")
				}
				return v.validate()
			})),
		),
	)
}

// Add 加入发货台
// https://seller.kuajingmaihuo.com/sop/view/889973754324016047#YSg2AE
func (s shipOrderStagingService) Add(ctx context.Context, req ShipOrderStagingAddRequest) (ok bool, results []entity.Result, err error) {
	if err = req.validate(); err != nil {
		return ok, results, invalidInput(err)
	}

	results = make([]entity.Result, len(req.JoinInfoList))
	for i, info := range req.JoinInfoList {
		results[i] = entity.Result{
			Key:     info.SubPurchaseOrderSn,
			Success: true,
		}
	}

	type joinError struct {
		JoinErrorSubPurchaseOrderSn string            `json:"joinErrorSubPurchaseOrderSn"` // 加入发货台失败的发货单号
		ExtraInfoMap                map[string]string `json:"extraInfoMap"`                // 附加信息字段
		ErrorCode                   int               `json:"errorCode"`                   // 错误码
		ErrorMsg                    string            `json:"errorMsg"`                    // 错误消息
	}
	var result = struct {
		normal.Response
		Result struct {
			JoinErrorList             []joinError `json:"joinErrorList"`
			ExistJoinErrorSubPurchase bool        `json:"existJoinErrorSubPurchase"`
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(req).
		SetResult(&result).
		Post("bg.shiporder.staging.add")
	if err = recheckError(resp, result.Response, err); err != nil {
		return
	}

	ok = !result.Result.ExistJoinErrorSubPurchase
	if !ok {
		kvJoinErrors := make(map[string]joinError, len(result.Result.JoinErrorList))
		for _, v := range result.Result.JoinErrorList {
			kvJoinErrors[strings.ToLower(v.JoinErrorSubPurchaseOrderSn)] = v
		}
		for i, r := range results {
			joinErr, exists := kvJoinErrors[strings.ToLower(r.Key)]
			if !exists {
				continue
			}

			r.Success = false
			r.Code = null.IntFrom(int64(joinErr.ErrorCode))
			errMessage := joinErr.ErrorMsg
			if len(joinErr.ExtraInfoMap) != 0 {
				errMessage += fmt.Sprintf("(%s)", jsonx.ToJson(joinErr.ExtraInfoMap, "{}"))
			}
			r.Error = nullx.StringFrom(errMessage)
			results[i] = r
		}
	}
	return
}
