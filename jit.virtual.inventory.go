package temu

import (
	"context"
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

// 虚拟库存
type jitVirtualInventoryService service

// Query 虚拟库存查询接口
// https://seller.kuajingmaihuo.com/sop/view/706628248275137588#ag3EtD
func (s jitVirtualInventoryService) Query(ctx context.Context, productSkcId int64) (items []entity.JitProductVirtualInventory, err error) {
	var result = struct {
		normal.Response
		Result struct {
			Total               int                                 `json:"total"`               // 总数
			ProductSkuStockList []entity.JitProductVirtualInventory `json:"productSkuStockList"` // 货品 SKU 库存列表
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(map[string]int64{"productSkcId": productSkcId}).
		SetResult(&result).
		Post("bg.virtualinventoryjit.get")
	if err = recheckError(resp, result.Response, err); err != nil {
		return
	}

	return result.Result.ProductSkuStockList, nil
}

// 虚拟库存编辑接口（bg.virtualinventoryjit.edit）
// 全托管JIT库存限制：调整后虚拟库存数量必须 ≥ skuId在Temu仓库中的实物库存数量

type SkuVirtualStockChangeRequest struct {
	ProductSkuId     int64 `json:"productSkuId"`     // 货品 SKU ID
	VirtualStockDiff int   `json:"virtualStockDiff"` // 虚拟库存(含商家自管库存)变更，大于0代表增加，小于0代表减少
}

func (m SkuVirtualStockChangeRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.ProductSkuId, validation.Required.Error("货品 SKU 不能为空。")),
		validation.Field(&m.VirtualStockDiff,
			validation.Required.Error("库存变更数量不能为空。"),
			validation.By(func(value interface{}) error {
				qty, ok := value.(int)
				if !ok || qty == 0 {
					return errors.New("无效的库存变更数量。")
				}
				return nil
			}),
		),
	)
}

type VirtualInventoryJitEditRequest struct {
	ProductSkcId              int64                          `json:"productSkcId"`              // 货品SKC ID
	SkuVirtualStockChangeList []SkuVirtualStockChangeRequest `json:"skuVirtualStockChangeList"` // 虚拟库存模式下使用，虚拟库存调整信息.
}

func (m VirtualInventoryJitEditRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.ProductSkcId, validation.Required.Error("货品 SKC 不能为空。")),
		validation.Field(&m.SkuVirtualStockChangeList, validation.Required.Error("虚拟库存不能为空。")),
	)
}

// Edit 虚拟库存编辑接口
// https://seller.kuajingmaihuo.com/sop/view/706628248275137588#hALnFd
func (s jitVirtualInventoryService) Edit(ctx context.Context, request VirtualInventoryJitEditRequest) (ok bool, err error) {
	if err = request.Validate(); err != nil {
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
		Post("bg.virtualinventoryjit.edit")
	if err = recheckError(resp, result.Response, err); err != nil {
		return
	}

	return true, nil
}
