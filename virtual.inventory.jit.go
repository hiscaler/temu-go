package temu

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

// 虚拟库存
type virtualInventoryJitService service

// View 虚拟库存查询接口
// https://seller.kuajingmaihuo.com/sop/view/706628248275137588#ag3EtD
func (s virtualInventoryJitService) View(productSkcId int) (items []entity.VirtualInventoryJit, err error) {
	var result = struct {
		normal.Response
		Result struct {
			Total               int                          `json:"total"`               // 总数
			ProductSkuStockList []entity.VirtualInventoryJit `json:"productSkuStockList"` // 订单信息
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetBody(map[string]int{"productSkcId": productSkcId}).
		SetResult(&result).
		Post("bg.virtualinventoryjit.get")
	if err == nil {
		err = parseResponse(resp, result.Response)
	}
	if err != nil {
		return
	}

	return result.Result.ProductSkuStockList, nil
}

// 虚拟库存编辑接口（bg.virtualinventoryjit.edit）
// 全托管JIT库存限制：调整后虚拟库存数量必须 ≥ skuId在Temu仓库中的实物库存数量

type VirtualInventoryJitEditRequest struct {
	ProductSkcId              int `json:"productSkcId"` // 货品SKC ID
	SkuVirtualStockChangeList struct {
		VirtualStockDiff int `json:"virtualStockDiff"` // 虚拟库存(含商家自管库存)变更，大于0代表增加，小于0代表减少
		ProductSkuId     int `json:"productSkuId"`     // 货品 SKU ID.
	} `json:"skuVirtualStockChangeList"` // 虚拟库存模式下使用，虚拟库存调整信息.
}

func (m VirtualInventoryJitEditRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.ProductSkcId, validation.Required.Error("货品SKC ID不能为空。")),
		validation.Field(&m.SkuVirtualStockChangeList, validation.Required.Error("虚拟库存不能为空。")),
	)
}

// Edit 虚拟库存编辑接口
// https://seller.kuajingmaihuo.com/sop/view/706628248275137588#hALnFd
func (s virtualInventoryJitService) Edit(request VirtualInventoryJitEditRequest) (ok bool, err error) {
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
		Post("bg.virtualinventoryjit.edit")
	if err == nil {
		err = parseResponse(resp, result.Response)
	}

	ok = err == nil
	return
}
