package temu

import (
	"context"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
	"github.com/hiscaler/temu-go/validators/is"
	"gopkg.in/guregu/null.v4"
	"strings"
)

// 采购单（备货单）服务
type purchaseOrderService service

// 查询采购单列表 V2

type PurchaseOrderQueryParams struct {
	normal.ParameterWithPager
	SettlementType                  null.Int  `json:"settlementType,omitempty"`                  // 结算类型 0-非vmi(采购) 1-vmi(备货)
	UrgencyType                     null.Int  `json:"urgencyType,omitempty"`                     // 是否紧急 0-否 1-是
	StatusList                      []int     `json:"statusList,omitempty"`                      // 订单状态 0-待接单；1-已接单，待发货；2-已送货；3-已收货；4-已拒收；5-已验收，全部退回；6-已验收；7-已入库；8-作废；9-已超时；10-已取消
	SubPurchaseOrderSnList          []string  `json:"subPurchaseOrderSnList,omitempty"`          // 订单号（采购子单号）
	ProductSnList                   []string  `json:"productSnList,omitempty"`                   // 货号列表
	ProductSkcIdList                []int64   `json:"productSkcIdList,omitempty"`                // skc 列表
	PurchaseTimeFrom                int64     `json:"purchaseTimeFrom,omitempty"`                // 下单时间-开始：毫秒
	PurchaseTimeTo                  int64     `json:"purchaseTimeTo,omitempty"`                  // 下单时间-结束：毫秒
	DeliverOrderSnList              []string  `json:"deliverOrderSnList,omitempty"`              // 发货单号列表
	IsDelayDeliver                  null.Bool `json:"isDelayDeliver,omitempty"`                  // 是否延迟发货
	IsDelayArrival                  null.Bool `json:"isDelayArrival,omitempty"`                  // 是否延迟到货
	ExpectLatestDeliverTimeFrom     int64     `json:"expectLatestDeliverTimeFrom,omitempty"`     // 要求最晚发货时间-开始（时间戳 单位：毫秒）
	ExpectLatestDeliverTimeTo       int64     `json:"expectLatestDeliverTimeTo,omitempty"`       // 要求最晚发货时间-结束（时间戳 单位：毫秒）
	ExpectLatestArrivalTimeFrom     int64     `json:"expectLatestArrivalTimeFrom,omitempty"`     // 要求最晚到达时间-开始（时间戳 单位：毫秒）
	ExpectLatestArrivalTimeTo       int64     `json:"expectLatestArrivalTimeTo,omitempty"`       // 要求最晚到达时间-结束（时间戳 单位：毫秒）
	PurchaseStockType               null.Int  `json:"purchaseStockType,omitempty"`               // 是否是JIT备货， 0-普通，1-JIT备货
	IsFirst                         null.Bool `json:"isFirst,omitempty"`                         // 是否首单 0-否 1-是
	IsCustomGoods                   null.Bool `json:"isCustomGoods,omitempty"`                   // 是否为定制品
	OriginalPurchaseOrderSnList     []string  `json:"originalPurchaseOrderSnList,omitempty"`     // 母订单号列表
	DeliverOrArrivalDelayStatusList []int     `json:"deliverOrArrivalDelayStatusList,omitempty"` // 发货或者到货逾期状态 101-发货即将逾期，102-发货已逾期，201-到货即将逾期，202-到货已逾期
	TodayCanDeliver                 null.Bool `json:"todayCanDeliver,omitempty"`                 // 是否今日可发货
	SkuLackSnapshot                 null.Int  `json:"skuLackSnapshot,omitempty"`                 // 创单时是否存在缺货sku，0-不存在 1-存在
	QcReject                        null.Int  `json:"qcReject,omitempty"`                        // 创单时是否存在质检不合格sku，0-不存在 1-存在
	QcOption                        null.Int  `json:"qcOption,omitempty"`                        // 是否存在质检不合格的sku，10-是，20-否
	SourceList                      []int     `json:"sourceList,omitempty"`                      // 下单来源，0-运营下单，1-卖家下单，9999-平台下单
	IsSystemAutoPurchaseSource      null.Bool `json:"isSystemAutoPurchaseSource,omitempty"`      // 是否系统下单 是-系统自动下单 否-其他
	LackOrSoldOutTagList            []int     `json:"lackOrSoldOutTagList,omitempty"`            // 标签：1-含缺货SKU；2-含售罄SKU
	IsTodayPlatformPurchase         null.Bool `json:"isTodayPlatformPurchase,omitempty"`         // 是否今日平台下单
	JoinDeliveryPlatform            null.Bool `json:"joinDeliveryPlatform,omitempty"`            // 是否加入了发货台
	OrderType                       null.Int  `json:"orderType,omitempty"`                       // 订单类型（1：普通备货单、2：JIT 备货单、3：定制备货单）此参数为扩展参数，用于简化备货类型查询处理
}

func (m PurchaseOrderQueryParams) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.SettlementType,
			validation.When(m.SettlementType.Valid,
				validation.By(func(value interface{}) error {
					v, ok := value.(null.Int)
					if !ok {
						return errors.New("无效的结算类型")
					}

					return validation.Validate(int(v.Int64), validation.In(entity.SettlementTypeNotVMI, entity.SettlementTypeVMI).Error("无效的结算类型"))
				}),
			),
		),
		validation.Field(&m.UrgencyType,
			validation.When(m.UrgencyType.Valid,
				validation.By(func(value interface{}) error {
					v, ok := value.(null.Int)
					if !ok {
						return errors.New("无效的紧急类型")
					}

					return validation.Validate(int(v.Int64), validation.In(entity.FalseNumber, entity.TrueNumber).Error("无效的紧急类型"))
				}),
			),
		),
		validation.Field(&m.StatusList,
			validation.When(
				len(m.StatusList) != 0,
				validation.Each(validation.In(
					entity.PurchaseOrderStatusWaitingMerchantReceive,
					entity.PurchaseOrderStatusMerchantReceived,
					entity.PurchaseOrderStatusMerchantSend,
					entity.PurchaseOrderStatusPlatformReceived,
					entity.PurchaseOrderStatusPlatformRejected,
					entity.PurchaseOrderStatusPlatformReturned,
					entity.PurchaseOrderStatusPlatformApproved,
					entity.PurchaseOrderStatusPlatformPutInStorage,
					entity.PurchaseOrderStatusDiscard,
					entity.PurchaseOrderStatusTimeout,
					entity.PurchaseOrderStatusCancel,
				).Error("无效的状态")),
			),
		),
		validation.Field(&m.SubPurchaseOrderSnList, validation.When(len(m.SubPurchaseOrderSnList) != 0, validation.Each(validation.By(is.PurchaseOrderNumber())))),
		validation.Field(&m.PurchaseStockType,
			validation.When(m.PurchaseStockType.Valid,
				validation.By(func(value interface{}) error {
					v, ok := value.(null.Int)
					if !ok {
						return errors.New("无效的备货类型")
					}

					return validation.Validate(int(v.Int64), validation.In(entity.PurchaseStockTypeNormal, entity.PurchaseStockTypeJIT).Error("无效的备货类型"))
				}),
			),
		),
		validation.Field(&m.QcReject,
			validation.When(m.PurchaseStockType.Valid,
				validation.By(func(value interface{}) error {
					v, ok := value.(null.Int)
					if !ok {
						return errors.New("无效的创单时质检是否合格查询值")
					}

					return validation.Validate(int(v.Int64), validation.In(entity.FalseNumber, entity.TrueNumber).Error("无效的创单时质检是否合格查询值"))
				}),
			),
		),
		validation.Field(&m.SourceList,
			validation.When(
				len(m.SourceList) > 0,
				validation.Each(validation.In(
					entity.PurchaseOrderSourceOperationalStaff,
					entity.PurchaseOrderSourceSeller,
					entity.PurchaseOrderSourcePlatform,
				).Error("无效的下单来源")),
			),
		),
		validation.Field(&m.OrderType,
			validation.When(m.OrderType.Valid, validation.By(func(value interface{}) error {
				v, ok := value.(null.Int)
				if !ok {
					return errors.New("无效的备货单类型")
				}

				return validation.Validate(int(v.Int64), validation.In(entity.OrderTypeNormal, entity.OrderTypeJIT, entity.OrderTypeCustomized).Error("无效的备货单类型"))
			})),
		),
	)
}

// Query 查询采购单列表 V2
// https://seller.kuajingmaihuo.com/sop/view/889973754324016047#Ip0Gso
func (s purchaseOrderService) Query(ctx context.Context, params PurchaseOrderQueryParams) (items []entity.PurchaseOrder, stat entity.PurchaseOrderStatistic, err error) {
	params.TidyPager()
	if params.OrderType.Valid {
		switch params.OrderType.Int64 {
		case entity.OrderTypeNormal:
			params.IsCustomGoods = null.BoolFrom(false)
			params.PurchaseStockType = null.IntFrom(entity.PurchaseStockTypeNormal)

		case entity.OrderTypeJIT:
			params.IsCustomGoods = null.BoolFrom(false)
			params.PurchaseStockType = null.IntFrom(entity.PurchaseStockTypeJIT)

		case entity.OrderTypeCustomized:
			params.IsCustomGoods = null.BoolFrom(true)
			params.PurchaseStockType = null.NewInt(0, false)
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
			entity.PurchaseOrderStatistic
			SubOrderForSupplierList []entity.PurchaseOrder `json:"subOrderForSupplierList"` // 订单信息
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(params).
		SetResult(&result).
		Post("bg.purchaseorderv2.get")
	if err = recheckError(resp, result.Response, err); err != nil {
		return
	}

	items = result.Result.SubOrderForSupplierList
	for i, item := range items {
		orderType := 0 // Unknown
		if item.IsCustomProduct {
			orderType = entity.OrderTypeCustomized
		} else if item.PurchaseStockType == entity.PurchaseStockTypeJIT {
			orderType = entity.OrderTypeJIT
		} else {
			orderType = entity.OrderTypeNormal
		}
		items[i].OrderType = orderType
	}

	return items, result.Result.PurchaseOrderStatistic, nil
}

// One 根据子采购单或者母采购单号查询采购单数据
func (s purchaseOrderService) One(ctx context.Context, number string) (item entity.PurchaseOrder, err error) {
	if len(number) <= 2 {
		err = ErrInvalidParameters
		return
	}

	prefix := strings.ToLower(number[0:2])
	if prefix != "wp" && prefix != "wb" {
		err = ErrInvalidParameters
		return
	}

	isPurchaseOrder := prefix == "wb"
	params := PurchaseOrderQueryParams{}
	if isPurchaseOrder {
		params.SubPurchaseOrderSnList = []string{number}
	} else {
		params.OriginalPurchaseOrderSnList = []string{number}
	}
	items, _, err := s.Query(ctx, params)
	if err != nil {
		return
	}

	for _, order := range items {
		if (isPurchaseOrder && strings.EqualFold(order.SubPurchaseOrderSn, number)) ||
			(!isPurchaseOrder && strings.EqualFold(order.OriginalPurchaseOrderSn, number)) {
			return order, nil
		}
	}

	return item, ErrNotFound
}

// 申请备货

type PurchaseOrderApplyDetail struct {
	ProductSkuId               int64 `json:"productSkuId"`               // skuId
	ProductSkuPurchaseQuantity int   `json:"productSkuPurchaseQuantity"` // 申请备货数量
}

func (m PurchaseOrderApplyDetail) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.ProductSkuId, validation.Required.Error("SKU 不能为空")),
		validation.Field(&m.ProductSkuPurchaseQuantity, validation.Min(1).Error("备货数量不能小于 {.min}")),
	)
}

type PurchaseOrderApplyRequest struct {
	normal.Parameter
	ProductSkcId            int64                    `json:"productSkcId"`                      // skcId
	ExpectLatestDeliverTime null.Int                 `json:"expectLatestDeliverTime,omitempty"` // 最晚发货时间
	ExpectLatestArrivalTime null.Int                 `json:"expectLatestArrivalTime,omitempty"` // 最晚送达时间
	PurchaseDetailList      PurchaseOrderApplyDetail `json:"purchaseDetailList"`                // 详情
}

func (m PurchaseOrderApplyRequest) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.ProductSkcId, validation.Required.Error("SKC 不能为空")),
		validation.Field(&m.ExpectLatestDeliverTime, validation.When(m.ExpectLatestDeliverTime.Valid), validation.By(is.Millisecond())),
		validation.Field(&m.ExpectLatestArrivalTime,
			validation.When(m.ExpectLatestArrivalTime.Valid),
			validation.By(is.Millisecond()),
			validation.When(m.ExpectLatestDeliverTime.Valid, validation.By(func(value interface{}) error {
				v, _ := value.(int64)
				if v > m.ExpectLatestArrivalTime.Int64 {
					return errors.New("最晚送达时间不能小于最晚发货时间")
				}
				return nil
			})),
		),
		validation.Field(&m.PurchaseDetailList, validation.Required.Error("备货详情不能为空")),
	)
}

// Apply 申请备货
// https://seller.kuajingmaihuo.com/sop/view/889973754324016047#nsjLx8
func (s purchaseOrderService) Apply(ctx context.Context, request PurchaseOrderApplyRequest) (ok bool, err error) {
	if err = request.validate(); err != nil {
		err = invalidInput(err)
		return
	}
	var result = struct {
		normal.Response
		Result any `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(request).
		SetResult(&result).
		Post("bg.purchaseorder.apply")
	if err = recheckError(resp, result.Response, err); err != nil {
		return
	}

	return true, nil
}

// 修改备货单下单数量（bg.purchaseorder.edit）
// https://seller.kuajingmaihuo.com/sop/view/889973754324016047#YT2bPD

type PurchaseOrderEditItem struct {
	ProductSkuId               int64 `json:"productSkuId"`               // 货品 SKU ID
	ProductSkuPurchaseQuantity int   `json:"productSkuPurchaseQuantity"` // 货品 SKU 下单数量
}

type PurchaseOrderEditRequest struct {
	SubPurchaseOrderSn string                  `json:"subPurchaseOrderSn"` // 采购子单号（订单号） 支持修改（待创建）普通备货单的备货数量 备货数量仅支持向下修改
	PurchaseDetailList []PurchaseOrderEditItem `json:"purchaseDetailList"` // 采购详情列表
}

func (m PurchaseOrderEditRequest) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.SubPurchaseOrderSn,
			validation.Required.Error("备货单号不能为空"),
			validation.By(is.PurchaseOrderNumber()),
		),
		validation.Field(&m.PurchaseDetailList,
			validation.Required.Error("待修改备货单详情不能为空"),
			validation.Each(validation.WithContext(func(ctx context.Context, value interface{}) error {
				item, ok := value.(PurchaseOrderEditItem)
				if !ok {
					return errors.New("无效的备货单详情")
				}
				return validation.ValidateStruct(&item,
					validation.Field(&item.ProductSkuId, validation.Required.Error("SKU 不能为空")),
					validation.Field(&item.ProductSkuPurchaseQuantity,
						validation.Min(1).Error("修改数量不能小于 {.min}"),
					),
				)
			})),
		),
	)
}

func (s purchaseOrderService) Edit(ctx context.Context, request PurchaseOrderEditRequest) (ok bool, err error) {
	if err = request.validate(); err != nil {
		return false, invalidInput(err)
	}

	var result = struct {
		normal.Response
		Result any `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(request).
		SetResult(&result).
		Post("bg.purchaseorder.edit")
	if err = recheckError(resp, result.Response, err); err != nil {
		return
	}

	return true, nil
}

// 批量取消待接单的备货单（bg.purchaseorder.cancel）

func (s purchaseOrderService) Cancel(ctx context.Context, rawPurchaseOrderNumbers ...string) (results []entity.Result, err error) {
	err = validation.Validate(rawPurchaseOrderNumbers,
		validation.Required.Error("备货单号不能为空"),
		validation.Each(validation.By(is.PurchaseOrderNumber())),
	)
	if err != nil {
		return results, invalidInput(err)
	}

	results = make([]entity.Result, len(rawPurchaseOrderNumbers))
	numbers := make([]string, 0)
	for i, number := range rawPurchaseOrderNumbers {
		result := entity.Result{Key: number}
		number = strings.TrimSpace(number)
		err = validation.Validate(number, validation.By(is.PurchaseOrderNumber()))
		if err != nil {
			result.Success = false
			result.Error = null.StringFrom(err.Error())
		} else {
			result.Success = true // Default is true? Must check API response result and reset it.
			numbers = append(numbers, number)
		}
		results[i] = result
	}

	var result = struct {
		normal.Response
		Result struct {
			IsSuccess       bool  `json:"isSuccess"`
			SuccessInfoList []any `json:"successInfoList"`
			ErrorInfoList   []any `json:"errorInfoList"`
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(map[string][]string{"subPurchaseOrderSnList": numbers}).
		SetResult(&result).
		Post("bg.purchaseorder.cancel")
	err = recheckError(resp, result.Response, err)
	if err != nil || result.Result.IsSuccess {
		return
	}
	// todo Must check results

	return
}
