package temu

import (
	"context"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/gox/nullx"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
	"github.com/samber/lo"
	"gopkg.in/guregu/null.v4"
	"strings"
)

// 发货台服务
type shipOrderStagingService service

type ShipOrderStagingQueryParams struct {
	normal.ParameterWithPager
	SubPurchaseOrderSnList []string `json:"subPurchaseOrderSnList,omitempty"` // 子采购单号列表
	SkcExtCode             []string `json:"skcExtCode,omitempty"`             // 货号列表
	ProductSkcIdList       []int64  `json:"productSkcIdList,omitempty"`       // skcId列表
	SettlementType         *int     `json:"settlementType,omitempty"`         // 结算类型 0-非vmi 1-vmi
	IsFirstOrder           bool     `json:"isFirstOrder,omitempty"`           // 是否首单
	UrgencyType            *int     `json:"urgencyType,omitempty"`            // 是否是紧急发货单，0-普通 1-急采
	IsJit                  bool     `json:"isJit,omitempty"`                  // 是否是jit，true:jit
	PurchaseStockType      *int     `json:"purchaseStockType,omitempty"`      // 备货类型 0-普通备货 1-jit备货
	IsCustomProduct        bool     `json:"isCustomProduct,omitempty"`        // 是否为定制品
	SubWarehouseId         int64    `json:"subWarehouseId,omitempty"`         // 收货子仓
	InventoryRegion        []int    `json:"inventoryRegion,omitempty"`        // DOMESTIC(1, "国内备货"), OVERSEAS(2, "海外备货"), BOUNDED_WAREHOUSE(3, "保税仓备货"),
}

func (m ShipOrderStagingQueryParams) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.SettlementType, validation.When(!validation.IsEmpty(m.SettlementType), validation.In(entity.SettlementTypeVMI, entity.SettlementTypeNotVMI).Error("无效的结算类型。"))),
		validation.Field(&m.PurchaseStockType, validation.When(!validation.IsEmpty(m.PurchaseStockType), validation.In(entity.PurchaseStockTypeNormal, entity.PurchaseStockTypeJIT).Error("无效的结算类型。"))),
	)
}

// All List all staging orders
// https://seller.kuajingmaihuo.com/sop/view/889973754324016047#NOA03y
func (s shipOrderStagingService) All(ctx context.Context, params ShipOrderStagingQueryParams) (items []entity.ShipOrderStaging, total, totalPages int, isLastPage bool, err error) {
	params.TidyPager()
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
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(params).
		SetResult(&result).
		Post("bg.shiporder.staging.get")
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
func (s shipOrderStagingService) One(ctx context.Context, subPurchaseOrderSn string) (item entity.ShipOrderStaging, err error) {
	items, _, _, _, err := s.All(ctx, ShipOrderStagingQueryParams{
		SubPurchaseOrderSnList: []string{subPurchaseOrderSn},
	})
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

func (m ShipOrderStagingAddInfo) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.SubPurchaseOrderSn,
			validation.Required.Error("备货单号不能为空。"),
			validation.By(func(value interface{}) error {
				number, ok := value.(string)
				if !ok || !strings.HasPrefix(strings.ToLower(number), "wb") {
					return fmt.Errorf("无效的备货单号：%v", value)
				}
				return nil
			}),
		),
		validation.Field(&m.DeliveryAddressType,
			validation.In(entity.DeliveryAddressTypeChineseMainland, entity.DeliveryAddressTypeChineseHongKong).Error("无效的发货地址类型。"),
		),
	)
}

type ShipOrderStagingAddRequest struct {
	normal.Parameter
	JoinInfoList []ShipOrderStagingAddInfo `json:"joinInfoList"` // 加入发货台的信息列表
}

func (m ShipOrderStagingAddRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.JoinInfoList, validation.Required.Error("发货数据不能为空。")),
	)
}

// Add 加入发货台
// https://seller.kuajingmaihuo.com/sop/view/889973754324016047#YSg2AE
func (s shipOrderStagingService) Add(ctx context.Context, req ShipOrderStagingAddRequest) (ok bool, results []entity.ShipOrderStagingAddResult, err error) {
	if err = req.Validate(); err != nil {
		return
	}

	for _, info := range req.JoinInfoList {
		results = append(results, entity.ShipOrderStagingAddResult{
			SubPurchaseOrderSn: info.SubPurchaseOrderSn,
			Success:            false,
			ErrorMessage:       nullx.StringFrom("Unknown"),
		})
	}

	type joinError struct {
		ExtraInfoMap                any    `json:"extraInfoMap"`
		JoinErrorSubPurchaseOrderSn string `json:"joinErrorSubPurchaseOrderSn"`
		ErrorCode                   int    `json:"errorCode"`
		ErrorMsg                    string `json:"errorMsg"`
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
	if err == nil {
		err = parseResponse(resp, result.Response)
	}
	if err != nil {
		return
	}

	if result.Result.ExistJoinErrorSubPurchase {
		lo.ForEach(results, func(r entity.ShipOrderStagingAddResult, index int) {
			value, exists := lo.Find(result.Result.JoinErrorList, func(rr joinError) bool {
				return strings.EqualFold(rr.JoinErrorSubPurchaseOrderSn, r.SubPurchaseOrderSn)
			})
			if exists {
				r.ErrorCode = null.IntFrom(int64(value.ErrorCode))
				r.ErrorMessage = nullx.StringFrom(value.ErrorMsg)
			} else {
				r.Success = true
				r.ErrorMessage = nullx.NullString()
			}
			results[index] = r
		})
		_, foundFailed := lo.Find(results, func(item entity.ShipOrderStagingAddResult) bool {
			return item.Success == false
		})
		ok = !foundFailed
		if foundFailed {
			errMessages := make([]string, 0)
			for _, addResult := range results {
				if !addResult.Success {
					errMessages = append(errMessages, fmt.Sprintf("%s: %s", addResult.SubPurchaseOrderSn, addResult.ErrorMessage.ValueOrZero()))
				}
			}
			err = errors.New(strings.Join(errMessages, "; "))
		}
	} else {
		ok = true
		lo.ForEach(results, func(item entity.ShipOrderStagingAddResult, index int) {
			item.Success = true
			item.ErrorMessage = nullx.NullString()
			results[index] = item
		})
	}
	return
}
