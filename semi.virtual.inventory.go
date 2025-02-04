package temu

import (
	"context"
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
	"gopkg.in/guregu/null.v4"
)

// 半托管虚拟库存服务
type semiVirtualInventoryService service

type SemiVirtualInventoryQueryParams struct {
	ProductSkcId int64 `json:"productSkcId"` // 货品 SKC ID
}

func (m SemiVirtualInventoryQueryParams) validate() error {
	return validation.ValidateStruct(&m, validation.Field(&m.ProductSkcId, validation.Required.Error("货品 SKC ID 不能为空")))
}

// Query 查询商品虚拟库存
// https://seller.kuajingmaihuo.com/sop/view/867739977041685428#hm9Qgt
func (s *semiVirtualInventoryService) Query(ctx context.Context, params SemiVirtualInventoryQueryParams) (items []entity.SemiVirtualInventory, err error) {
	if err = params.validate(); err != nil {
		err = invalidInput(err)
		return nil, err
	}

	var result = struct {
		normal.Response
		Result struct {
			Total               int                           `json:"total"`
			ProductSkuStockList []entity.SemiVirtualInventory `json:"productSkuStockList"`
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(params).
		SetResult(&result).
		Post("bg.goods.quantity.get")
	if err = recheckError(resp, result.Response, err); err != nil {
		return nil, err
	}

	return result.Result.ProductSkuStockList, nil
}

type SemiVirtualInventoryChangeItem struct {
	ProductSkuId          int64       `json:"productSkuId"`                    // 货品 SKU ID
	StockDiff             null.Int    `json:"stockDiff,omitempty"`             // 虚拟库存变更(通过1-增减变更时：虚拟库存(含商家自管库存)变更，大于0代表在现有库存基础上增加，小于0代表在现有库存基础上减少)
	TargetStockAvailable  null.Int    `json:"targetStockAvailable,omitempty"`  // 覆盖变更目标库存值(通过2-覆盖变更时：覆盖变更目标库存值（填多少，则变更后库存则为多少，不能为负数))
	WarehouseId           null.String `json:"warehouseId"`                     // 发货仓ID(发货仓ID-当变更方式为2时，是必填字段。) 货品SKUId维度数据，欧洲地区支持分仓库存
	CurrentShippingMode   null.Int    `json:"currentShippingMode,omitempty"`   // 当前发货模式
	CurrentStockAvailable null.Int    `json:"currentStockAvailable,omitempty"` // 当前库存件数
}

type SemiVirtualInventoryUpdateRequest struct {
	QuantityChangeMode int                              `json:"quantityChangeMode"`     // 更新库存数量方式（1-增减变更 2-覆盖变更，默认为1）
	ProductSkcId       null.Int                         `json:"productSkcId,omitempty"` // 货品 SKC ID
	SkuStockChangeList []SemiVirtualInventoryChangeItem `json:"skuStockChangeList"`     // 虚拟库存调整信息
}

func (m SemiVirtualInventoryUpdateRequest) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.QuantityChangeMode,
			validation.Required.Error("更新虚拟库存数量方式不能为空"),
			validation.In(entity.QuantityChangeModeInDecrease, entity.QuantityChangeModeReplace).Error("无效的更新库存数量方式"),
		),
		validation.Field(&m.SkuStockChangeList,
			validation.Required.Error("虚拟库存调整信息不能为空"),
			validation.Each(validation.By(func(value interface{}) error {
				v, ok := value.(SemiVirtualInventoryChangeItem)
				if !ok {
					return errors.New("无效的虚拟库存调整数据")
				}

				if !v.WarehouseId.Valid || len(v.WarehouseId.String) == 0 {
					return errors.New("虚拟库存调整仓库 ID 不能为空")
				}

				switch m.QuantityChangeMode {
				case entity.QuantityChangeModeInDecrease:
					if !v.StockDiff.Valid || v.StockDiff.Int64 == 0 {
						return errors.New("虚拟库存变更值不能为空或者零值")
					}
				case entity.QuantityChangeModeReplace:
					if !v.TargetStockAvailable.Valid || v.TargetStockAvailable.Int64 < 0 {
						return errors.New("目标库存值不能为空或者负数")
					}
				}

				return nil
			})),
		),
	)
}

// Update 更新虚拟库存
// https://seller.kuajingmaihuo.com/sop/view/867739977041685428#DMwO8O
func (s *semiVirtualInventoryService) Update(ctx context.Context, params SemiVirtualInventoryUpdateRequest) (bool, error) {
	if err := params.validate(); err != nil {
		return false, err
	}

	var result = struct {
		normal.Response
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(params).
		SetResult(&result).
		Post("bg.goods.quantity.update")

	if err = recheckError(resp, result.Response, err); err != nil {
		return false, err
	}

	return result.Success, nil
}
