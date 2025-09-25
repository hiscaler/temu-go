package temu

import (
	"context"

	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

// 商品畅销服务
type goodsTopSellingService service

// SoldOut 批量查询爆款售罄商品（temu.goods.topselling.soldout.get）
// https://seller.kuajingmaihuo.com/sop/view/750197804480663142#1cxFkn
func (s goodsTopSellingService) SoldOut(ctx context.Context) (items []entity.GoodsTopSellingSoldOut, err error) {
	var result = struct {
		normal.Response
		Result struct {
			SellOutProducts []entity.GoodsTopSellingSoldOut `json:"sellOutProducts"`
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetResult(&result).
		Post("temu.goods.topselling.soldout.get")
	if err = recheckError(resp, result.Response, err); err != nil {
		return
	}

	items = result.Result.SellOutProducts
	return
}
