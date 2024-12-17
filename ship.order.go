package temu

import (
	"context"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/gox/inx"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
	"github.com/hiscaler/temu-go/validators/is"
	"gopkg.in/guregu/null.v4"
	"maps"
	"slices"
	"strings"
)

// 发货单服务
type shipOrderService struct {
	service
	Staging        shipOrderStagingService        // 发货台
	ReceiveAddress shipOrderReceiveAddressService // 收货地址
	Packing        shipOrderPackingService        // 装箱发货
	Package        shipOrderPackageService        // 包裹
}

type ShipOrderQueryParams struct {
	normal.ParameterWithPager
	DeliveryOrderSnList      []string  `json:"deliveryOrderSnList,omitempty"`      // 发货单号列表
	SubPurchaseOrderSnList   []string  `json:"subPurchaseOrderSnList,omitempty"`   // 子采购单号列表
	ExpressDeliverySnList    []string  `json:"expressDeliverySnList,omitempty"`    // 快递单号列表
	SkcExtCodeList           []string  `json:"skcExtCodeList,omitempty"`           // 货号列表
	ProductSkcIdList         []int64   `json:"productSkcIdList,omitempty"`         // skcId 列表
	SubWarehouseIdList       []int64   `json:"subWarehouseIdList,omitempty"`       // 收货子仓列表
	DeliverTimeFrom          int64     `json:"deliverTimeFrom,omitempty"`          // 发货时间-开始时间
	DeliverTimeTo            int64     `json:"deliverTimeTo,omitempty"`            // 发货时间-结束时间
	Status                   null.Int  `json:"status,omitempty"`                   // 发货单状态，0：待装箱发货，1：待仓库收货，2：已收货，3：已入库，4：已退货，5：已取消，6：部分收货，查询发货批次时仅支持查询发货单状态=1
	UrgencyType              null.Int  `json:"urgencyType,omitempty"`              // 是否是紧急发货单，0-普通，1-急采
	IsCustomProduct          null.Bool `json:"isCustomProduct,omitempty"`          // 是否为定制品
	IsVim                    null.Int  `json:"isVmi,omitempty"`                    // 是否是vmi，0-非VMI，1-VMI
	IsJit                    null.Bool `json:"isJit,omitempty"`                    // 是否是jit，true:jit
	LatestFeedbackStatusList []int     `json:"latestFeedbackStatusList,omitempty"` // 最新反馈状态列表，0-当前无异常，1-已提交，2-物流商处理中，3-已撤销，4-已反馈
	SortType                 null.Int  `json:"sortType,omitempty"`                 // 排序类型，0-创建时间最新在上，1-要求发货时间较早在上，2-按照仓库名称排序
	InventoryRegion          []int     `json:"inventoryRegion,omitempty"`          // 发货区域，1-国内备货，2-海外备货，3-保税仓备货
	IsPrintBoxMark           null.Int  `json:"isPrintBoxMark,omitempty"`           // 是否已打印商品打包标签，0-未打印，1-已打印
	TargetReceiveAddress     string    `json:"targetReceiveAddress,omitempty"`     // 筛选项-收货地址（精准匹配）
	TargetDeliveryAddress    string    `json:"targetDeliveryAddress,omitempty"`    // 筛选项-发货地址（精准匹配）
	OrderType                null.Int  `json:"orderType,omitempty"`                // 订单类型（1：普通备货单、2：JIT 备货单、3：定制备货单）此参数为扩展参数，用于简化备货类型查询处理
}

func (m ShipOrderQueryParams) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.UrgencyType, validation.When(
			m.UrgencyType.Valid,
			validation.By(func(value interface{}) error {
				v, ok := value.(null.Int)
				if !ok {
					return errors.New("无效的加急类型")
				}
				return validation.Validate(int(v.Int64), validation.In(entity.UrgencyTypeNormal, entity.UrgencyTypeUrgency).Error("无效的加急类型"))
			}),
		)),
		validation.Field(&m.IsVim, validation.When(
			m.IsVim.Valid,
			validation.By(func(value interface{}) error {
				v, ok := value.(null.Int)
				if !ok {
					return errors.New("无效的 VMI")
				}
				return validation.Validate(int(v.Int64), validation.In(entity.FalseNumber, entity.TrueNumber).Error("无效的 VMI"))
			}),
		)),
		validation.Field(&m.OrderType,
			validation.When(m.OrderType.Valid, validation.By(func(value interface{}) error {
				v, ok := value.(null.Int)
				if !ok {
					return errors.New("无效的发货单类型")
				}

				return validation.Validate(int(v.Int64), validation.In(entity.OrderTypeNormal, entity.OrderTypeJIT, entity.OrderTypeCustomized).Error("无效的发货单类型"))
			})),
		),
	)
}

// Query 查询发货单 V2
// https://seller.kuajingmaihuo.com/sop/view/889973754324016047#B7c51j
func (s shipOrderService) Query(ctx context.Context, params ShipOrderQueryParams) (items []entity.ShipOrder, total, totalPages int, isLastPage bool, err error) {
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
			Total int                `json:"total"`
			List  []entity.ShipOrder `json:"list"`
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(params).
		SetResult(&result).
		Post("bg.shiporderv2.get")
	if err = recheckError(resp, result.Response, err); err != nil {
		return
	}

	items = result.Result.List
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
	total, totalPages, isLastPage = parseResponseTotal(params.Page, params.PageSize, result.Result.Total)

	return
}

// 创建发货单
// https://seller.kuajingmaihuo.com/sop/view/889973754324016047#HqGnA0

type ShipOrderCreateRequestOrderDetailInfo struct {
	ProductSkuId  int64 `json:"productSkuId"`  // 定制品，传定制品id；非定制品，传货品 skuId
	DeliverSkuNum int   `json:"deliverSkuNum"` // 发货sku数目
}

// ShipOrderCreateRequestOrderPackage 包裹信息
type ShipOrderCreateRequestOrderPackage struct {
	PackageDetailSaveInfos []ShipOrderCreateRequestPackageInfo `json:"packageDetailSaveInfos"` // 包裹明细
}

type ShipOrderCreateRequestPackageInfo struct {
	ProductSkuId int64 `json:"productSkuId"` // skuId
	SkuNum       int   `json:"skuNum"`       // 发货 sku 数目
}

type ShipOrderCreateRequestOrderInfo struct {
	SubPurchaseOrderSn      string                                  `json:"subPurchaseOrderSn"`      // 采购子单号
	DeliverOrderDetailInfos []ShipOrderCreateRequestOrderDetailInfo `json:"deliverOrderDetailInfos"` // 采购单创建信息列表
	PackageInfos            []ShipOrderCreateRequestOrderPackage    `json:"packageInfos"`            //	包裹信息列表
	DeliveryAddressId       int64                                   `json:"deliveryAddressId"`       // 发货地址 ID
}

func (m ShipOrderCreateRequestOrderInfo) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.SubPurchaseOrderSn, validation.By(is.PurchaseOrderNumber())),
		validation.Field(&m.DeliveryAddressId, validation.Required.Error("发货地址不能为空")),
	)
}

type ShipOrderCreateRequestDeliveryOrder struct {
	DeliveryOrderCreateInfos []ShipOrderCreateRequestOrderInfo `json:"deliveryOrderCreateInfos"`     // 发货单创建组列表
	ReceiveAddressInfo       entity.ReceiveAddress             `json:"receiveAddressInfo,omitempty"` // 收货地址
	SubWarehouseId           int64                             `json:"subWarehouseId,omitempty"`     // 子仓 ID
}

func (m ShipOrderCreateRequestDeliveryOrder) validate(ctx context.Context, s shipOrderService) error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.DeliveryOrderCreateInfos,
			validation.Required.Error("发货单创建组列表不能为空"),
			validation.By(func(value interface{}) error {
				errorMessages := make([]string, 0)
				requests, ok := value.([]ShipOrderCreateRequestOrderInfo)
				if !ok {
					return errors.New("无效的发货单创建组列表")
				}
				purchaseOrderNumbers := make([]string, 0)
				for _, request := range requests {
					err := request.validate()
					if err != nil {
						errorMessages = append(errorMessages, err.Error())
						continue
					}

					if inx.StringIn(request.SubPurchaseOrderSn, purchaseOrderNumbers...) {
						errorMessages = append(errorMessages, fmt.Sprintf("备货单 %s 重复", request.SubPurchaseOrderSn))
						continue
					}
					purchaseOrderNumbers = append(purchaseOrderNumbers, request.SubPurchaseOrderSn)
				}
				if len(errorMessages) != 0 {
					return errors.New(strings.Join(errorMessages, "; "))
				}

				stagingOrders := make([]entity.ShipOrderStaging, 0, len(purchaseOrderNumbers))
				pageSize := 20
				for chunkPurchaseOrderNumbers := range slices.Chunk(purchaseOrderNumbers, pageSize) {
					shipOrderStagingQueryParams := ShipOrderStagingQueryParams{
						SubPurchaseOrderSnList: chunkPurchaseOrderNumbers,
					}
					shipOrderStagingQueryParams.Page = 1
					shipOrderStagingQueryParams.PageSize = pageSize
					rawStagingOrders, _, _, _, err := s.Staging.Query(ctx, shipOrderStagingQueryParams)
					if err != nil {
						return err
					} else if len(rawStagingOrders) != 0 {
						stagingOrders = append(stagingOrders, rawStagingOrders...)
					}
				}

				kvNumberStagingOrder := make(map[string]entity.ShipOrderStaging, len(stagingOrders))
				if len(stagingOrders) != 0 {
					warehouseIdName := make(map[int64]string, 0)
					for _, order := range stagingOrders {
						purchaseOrder := order.SubPurchaseOrderBasicVO
						m.SubWarehouseId = purchaseOrder.SubWarehouseId
						m.ReceiveAddressInfo = purchaseOrder.ReceiveAddressInfo
						warehouseIdName[m.SubWarehouseId] = purchaseOrder.SubWarehouseName
						kvNumberStagingOrder[strings.ToLower(purchaseOrder.SubPurchaseOrderSn)] = order
					}
					switch len(warehouseIdName) {
					case 0:
						errorMessages = append(errorMessages, "无法获取收货仓库")
					case 1:
					default:
						names := slices.Collect(maps.Values(warehouseIdName))
						errorMessages = append(errorMessages, fmt.Sprintf("存在多个收货仓库: %s", strings.Join(names, ", ")))
					}
				}
				// 验证 DeliverOrderDetailInfos, PackageInfos 数据
				// DeliverOrderDetailInfos, PackageInfos 为空表示全部数据创建发货单，如果指定的话则只有指定的数据会创建发货单
				for i, request := range requests {
					purchaseOrderNumber := request.SubPurchaseOrderSn
					shipOrderStaging, exists := kvNumberStagingOrder[strings.ToLower(purchaseOrderNumber)]
					if !exists {
						errorMessages = append(errorMessages, fmt.Sprintf("%s 在发货台中不存在", purchaseOrderNumber))
						continue
					}

					kvSkuIdQuantity := make(map[int64]int, len(shipOrderStaging.OrderDetailVOList)) // 默认发货的 sku 和数量（skuId: quantity）
					for _, v := range shipOrderStaging.OrderDetailVOList {
						skuId := v.ProductSkuId
						if shipOrderStaging.SubPurchaseOrderBasicVO.IsCustomProduct {
							skuId = v.ProductOriginalSkuId
						}
						kvSkuIdQuantity[skuId] = v.SkuDeliveryQuantityMaxLimit
					}
					if len(request.DeliverOrderDetailInfos) == 0 && len(request.PackageInfos) == 0 {
						// 用户未主动添加发货信息，默认将所有可发货的数据加进来
						// 用户要不全部提供，要不全部不提供由系统处理，不能只添加部分信息
						for skuId, quantity := range kvSkuIdQuantity {
							request.DeliverOrderDetailInfos = append(request.DeliverOrderDetailInfos, ShipOrderCreateRequestOrderDetailInfo{
								ProductSkuId:  skuId,
								DeliverSkuNum: quantity,
							})
							request.PackageInfos = append(request.PackageInfos, ShipOrderCreateRequestOrderPackage{
								PackageDetailSaveInfos: []ShipOrderCreateRequestPackageInfo{
									{
										ProductSkuId: skuId,
										SkuNum:       quantity,
									},
								},
							})
						}
						m.DeliveryOrderCreateInfos[i] = request // 补充发货单创建组列表数据
					} else {
						// 验证用户主动提交的信息
						if len(request.DeliverOrderDetailInfos) == 0 {
							errorMessages = append(errorMessages, fmt.Sprintf("%s 采购单创建信息不能为空", purchaseOrderNumber))
						}
						for _, v := range request.DeliverOrderDetailInfos {
							var qty int
							if qty, exists = kvSkuIdQuantity[v.ProductSkuId]; !exists {
								errorMessages = append(errorMessages, fmt.Sprintf("%s SKU %d 不存在", purchaseOrderNumber, v.ProductSkuId))
								continue
							}
							deliveryQty := v.DeliverSkuNum
							if deliveryQty <= 0 {
								errorMessages = append(errorMessages, fmt.Sprintf("%s SKU %d 发货数量不能小于零", purchaseOrderNumber, v.ProductSkuId))
							}
							if deliveryQty > qty {
								errorMessages = append(errorMessages, fmt.Sprintf("%s SKU %d 发货数量不能超过 %d", purchaseOrderNumber, v.ProductSkuId, qty))
							}
						}

						if len(request.PackageInfos) == 0 {
							errorMessages = append(errorMessages, fmt.Sprintf("%s 包裹信息列表不能为空", purchaseOrderNumber))
						}
						for _, packageInfo := range request.PackageInfos {
							for _, v := range packageInfo.PackageDetailSaveInfos {
								var qty int
								if qty, exists = kvSkuIdQuantity[v.ProductSkuId]; !exists {
									errorMessages = append(errorMessages, fmt.Sprintf("%s SKU %d 不存在", purchaseOrderNumber, v.ProductSkuId))
									continue
								}
								deliveryQty := v.SkuNum
								if deliveryQty <= 0 {
									errorMessages = append(errorMessages, fmt.Sprintf("%s SKU %d 包裹发货数量不能小于零", purchaseOrderNumber, v.ProductSkuId))
								}
								if deliveryQty > qty {
									errorMessages = append(errorMessages, fmt.Sprintf("%s SKU %d 包裹发货数量不能超过 %d", purchaseOrderNumber, v.ProductSkuId, qty))
								}
							}
						}
					}
				}

				if len(errorMessages) != 0 {
					return errors.New(strings.Join(errorMessages, "; "))
				}
				return nil
			})),
		// validation.Field(&m.SubWarehouseId, validation.Required.Error("收货仓库不能为空")),
	)
}

type ShipOrderCreateRequest struct {
	normal.Parameter
	DeliveryOrderCreateGroupList []ShipOrderCreateRequestDeliveryOrder `json:"deliveryOrderCreateGroupList"` // 发货单创建组列表
}

func (m ShipOrderCreateRequest) validate(ctx context.Context, s shipOrderService) error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.DeliveryOrderCreateGroupList,
			validation.Required.Error("发货单创建组列表不能为空"),
			validation.Each(validation.By(func(value interface{}) error {
				v, ok := value.(ShipOrderCreateRequestDeliveryOrder)
				if !ok {
					return errors.New("无效的发货单创建数据")
				}
				return v.validate(ctx, s)
			})),
		),
	)
}

// Create 创建发货单接口 V3
func (s shipOrderService) Create(ctx context.Context, req ShipOrderCreateRequest) (ok bool, err error) {
	if err = req.validate(ctx, s); err != nil {
		err = invalidInput(err)
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
	if err = recheckError(resp, result.Response, err); err != nil {
		return
	}

	return true, nil
}

// 取消发货单

type ShipOrderCancelRequest struct {
	normal.Parameter
	DeliveryOrderSn string `json:"deliveryOrderSn"` // 发货单 ID
}

func (m ShipOrderCancelRequest) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.DeliveryOrderSn,
			validation.Required.Error("发货单号不能为空"),
			validation.By(is.ShipOrderNumber()),
		),
	)
}

// Cancel 取消发货单
// https://seller.kuajingmaihuo.com/sop/view/889973754324016047#UywT8E
func (s shipOrderService) Cancel(ctx context.Context, shipOrderNumber string) (ok bool, err error) {
	req := ShipOrderCancelRequest{DeliveryOrderSn: shipOrderNumber}
	if err = req.validate(); err != nil {
		err = invalidInput(err)
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
	if err = recheckError(resp, result.Response, err); err != nil {
		return
	}

	return true, nil
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
	if err = recheckError(resp, result.Response, err); err != nil {
		return
	}

	return result.Result, nil
}
