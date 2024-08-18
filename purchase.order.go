package temu

import (
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

type purchaseOrderService service

// 采购单（备货单）服务

// 查询采购单列表 V2

type PurchaseOrderQueryParams struct {
	normal.ParameterWithPager
	SettlementType                  int      `json:"settlementType"`                            // 结算类型 0-非vmi(采购) 1-vmi(备货)
	UrgencyType                     int      `json:"urgencyType,omitempty"`                     // 是否紧急 0-否 1-是
	StatusList                      []int    `json:"statusList,omitempty"`                      // 订单状态 0-待接单；1-已接单，待发货；2-已送货；3-已收货；4-已拒收；5-已验收，全部退回；6-已验收；7-已入库；8-作废；9-已超时；10-已取消
	SubPurchaseOrderSnList          []string `json:"subPurchaseOrderSnList,omitempty"`          // 订单号（采购子单号）
	ProductSnList                   []int    `json:"productSnList,omitempty"`                   // 货号列表
	ProductSkcIdList                []int    `json:"productSkcIdList,omitempty"`                // skc列表
	PurchaseTimeFrom                int      `json:"purchaseTimeFrom,omitempty"`                // 下单时间-开始：毫秒
	PurchaseTimeTo                  int      `json:"purchaseTimeTo,omitempty"`                  // 下单时间-结束：毫秒
	DeliverOrderSnList              []string `json:"deliverOrderSnList,omitempty"`              // 发货单号列表
	IsDelayDeliver                  bool     `json:"isDelayDeliver,omitempty"`                  // 是否延迟发货
	IsDelayArrival                  bool     `json:"isDelayArrival,omitempty"`                  // 是否延迟到货
	ExpectLatestDeliverTimeFrom     int      `json:"expectLatestDeliverTimeFrom,omitempty"`     // 要求最晚发货时间-开始（时间戳 单位：毫秒）
	ExpectLatestDeliverTimeTo       int      `json:"expectLatestDeliverTimeTo,omitempty"`       // 要求最晚发货时间-结束（时间戳 单位：毫秒）
	ExpectLatestArrivalTimeFrom     int      `json:"expectLatestArrivalTimeFrom,omitempty"`     // 要求最晚到达时间-开始（时间戳 单位：毫秒）
	ExpectLatestArrivalTimeTo       int      `json:"expectLatestArrivalTimeTo,omitempty"`       // 要求最晚到达时间-结束（时间戳 单位：毫秒）
	PurchaseStockType               int      `json:"purchaseStockType,omitempty"`               // 是否是JIT备货， 0-普通，1-JIT备货
	IsFirst                         bool     `json:"isFirst,omitempty"`                         // 是否首单 0-否 1-是
	IsCustomGoods                   bool     `json:"isCustomGoods,omitempty"`                   // 是否为定制品
	OriginalPurchaseOrderSnList     []string `json:"originalPurchaseOrderSnList,omitempty"`     // 母订单号列表
	DeliverOrArrivalDelayStatusList []int    `json:"deliverOrArrivalDelayStatusList,omitempty"` // 发货或者到货逾期状态 101-发货即将逾期，102-发货已逾期，201-到货即将逾期，202-到货已逾期
	TodayCanDeliver                 bool     `json:"todayCanDeliver,omitempty"`                 // 是否今日可发货
	SkuLackSnapshot                 int      `json:"skuLackSnapshot,omitempty"`                 // 创单时是否存在缺货sku，0-不存在 1-存在
	QcReject                        int      `json:"qcReject,omitempty"`                        // 创单时是否存在质检不合格sku，0-不存在 1-存在
	QcOption                        int      `json:"qcOption,omitempty"`                        // 是否存在质检不合格的sku，10-是，20-否
	SourceList                      []int    `json:"sourceList,omitempty"`                      // 下单来源，0-运营下单，1-卖家下单，9999-平台下单
	IsSystemAutoPurchaseSource      bool     `json:"isSystemAutoPurchaseSource,omitempty"`      // 是否系统下单 是-系统自动下单 否-其他
	LackOrSoldOutTagList            []int    `json:"lackOrSoldOutTagList,omitempty"`            // 标签：1-含缺货SKU；2-含售罄SKU
	IsTodayPlatformPurchase         bool     `json:"isTodayPlatformPurchase,omitempty"`         // 是否今日平台下单
	JoinDeliveryPlatform            bool     `json:"joinDeliveryPlatform,omitempty"`            // 是否加入发货台
}

func (m PurchaseOrderQueryParams) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Page, validation.Required.Error("页码不能为空。")),
		validation.Field(&m.PageSize, validation.Required.Error("页面大小不能为空。")),
		validation.Field(&m.SettlementType,
			validation.Required.Error("结算类型不能为空。"),
			validation.In(entity.SettlementTypeNotVMI, entity.SettlementTypeVMI).Error("无效的结算类型。"),
		),
		validation.Field(&m.SourceList,
			validation.When(len(m.SourceList) > 0, validation.By(func(value interface{}) error {
				sources, ok := value.([]int)
				if !ok {
					return errors.New("无效的下单来源")
				}

				validSources := map[int]any{
					entity.PurchaseOrderSourceOperationalStaff: nil,
					entity.PurchaseOrderSourceSeller:           nil,
					entity.PurchaseOrderSourcePlatform:         nil,
				}
				for _, source := range sources {
					if _, ok = validSources[source]; !ok {
						return fmt.Errorf("无效的下单来源：%d", source)
					}
				}
				return nil
			})),
		),
	)
}

// All 查询采购单列表 V2
// https://seller.kuajingmaihuo.com/sop/view/889973754324016047#Ip0Gso
func (s purchaseOrderService) All(params PurchaseOrderQueryParams) (items []entity.PurchaseOrder, err error) {
	params.TidyPager()
	if err = params.Validate(); err != nil {
		return
	}
	var result = struct {
		normal.Response
		Result struct {
			labelCodePageResult struct {
				TotalCount int                    `json:"totalCount"` // 总数
				Data       []entity.PurchaseOrder `json:"data"`       // 结果列表
			}
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetBody(params).
		SetResult(&result).
		Post("bg.purchaseorderv2.get")
	if err == nil {
		err = parseResponse(resp, result.Response)
	}
	if err != nil {
		return
	}

	return result.Result.labelCodePageResult.Data, nil
}

// One 单个采购单查询
func (s purchaseOrderService) One(subPurchaseOrderSn string) (item entity.PurchaseOrder, err error) {
	items, err := s.All(PurchaseOrderQueryParams{SubPurchaseOrderSnList: []string{subPurchaseOrderSn}})
	if err != nil {
		return
	}
	if len(items) == 0 {
		return item, ErrNotFound
	}

	return items[0], nil
}

// 申请备货

type PurchaseOrderApplyRequest struct {
	normal.Parameter
	ProductSkcId            int `json:"productSkcId"`            // skcId
	ExpectLatestDeliverTime int `json:"expectLatestDeliverTime"` // 最晚发货时间
	ExpectLatestArrivalTime int `json:"expectLatestArrivalTime"` // 最晚送达时间
	PurchaseDetailList      struct {
		ProductSkuId               int `json:"productSkuId"`               // skuId
		ProductSkuPurchaseQuantity int `json:"productSkuPurchaseQuantity"` // 申请备货数量
	} `json:"purchaseDetailList"` // 详情
}

func (m PurchaseOrderApplyRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.ProductSkcId, validation.Required.Error("skcId不能为空。")),
		validation.Field(&m.PurchaseDetailList, validation.Required.Error("详情不能为空。")),
	)
}

// Apply 申请备货
// https://seller.kuajingmaihuo.com/sop/view/889973754324016047#nsjLx8
func (s purchaseOrderService) Apply(request PurchaseOrderApplyRequest) (ok bool, err error) {
	if err = request.Validate(); err != nil {
		return
	}
	var result = struct {
		normal.Response
		Result any `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetBody(request).
		SetResult(&result).
		Post("bg.purchaseorder.apply")
	if err == nil {
		err = parseResponse(resp, result.Response)
	}
	if err != nil {
		return
	}

	return true, nil
}
