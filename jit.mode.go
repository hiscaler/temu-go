package temu

import (
	"context"
	"github.com/hiscaler/temu-go/normal"
)

type jitModeService service

// Activate 打开JIT（bg.jitmode.activate）
// https://seller.kuajingmaihuo.com/sop/view/706628248275137588#b4ikfi
// 全托管JIT开通：全托管的SKC开通JIT模式，进行虚拟库存售卖，关联查询bg.product.search，SKC出参满足applyJitStatus=1时，可开通JIT模式，进行虚拟库存售卖
func (s jitModeService) Activate(ctx context.Context, productId, productSkcId int64) (bool, error) {
	var result = struct {
		normal.Response
		Result any `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(map[string]int64{"productId": productId, "productSkcId": productSkcId}).
		SetResult(&result).
		Post("bg.jitmode.activate")
	if err = recheckError(resp, result.Response, err); err != nil {
		return false, err
	}

	return true, nil
}
