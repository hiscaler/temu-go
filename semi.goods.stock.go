package temu

import (
	"context"
	"errors"
	"strings"

	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

// semiGoodsStockService 半托管商品库存服务
type semiGoodsStockService service

type SemiGoodsStockQueryRequest struct {
	Language     string   `json:"language"`
	OutSkuSnList []string `json:"outSkuSnList"` // 外部 SKU SN 列表
	SkuIdList    []int64  `json:"skuIdList"`    // SKU ID 列表
	GoodsId      int64    `json:"goodsId"`      // 商品 ID
}

func (m SemiGoodsStockQueryRequest) validate() error {
	return nil

}

// Query 查询半托管商品库存
// https://partner-us.temu.com/documentation?menu_code=fb16b05f7a904765aac4af3a24b87d4a&sub_menu_code=f14a3f28b654441b80f90e76a0a77c6e
func (s semiGoodsStockService) Query(ctx context.Context, request SemiGoodsStockQueryRequest) ([]entity.SemiGoodsStock, error) {
	if err := request.validate(); err != nil {
		return nil, err
	}

	var result = struct {
		normal.Response
		Result struct {
			StockList []entity.SemiGoodsStock `json:"stockList"` // 库存列表
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(request).
		SetResult(&result).
		Post("temu.goods.stock.query")

	if err = recheckError(resp, result.Response, err); err != nil {
		return nil, err
	}

	return result.Result.StockList, nil
}

type SemiGoodsStockUpdateRequest struct {
	GoodsId            int64 `json:"goodsId"`   // 商品 ID
	StockType          int   `json:"stockType"` // Stock type. Enumeration values: 0 - self ordinary stock, 1 - self Pre-sale stock; default is self ordinary stock if not filled.
	SkuStockChangeList []struct {
		SkuId     int64 `json:"skuId"`
		StockDiff int   `json:"stockDiff"`
	} `json:"skuStockChangeList"`
	SkuStockTargetList []struct {
		SkuId       int64 `json:"skuId"`
		StockTarget int   `json:"stockTarget"`
	} `json:"skuStockTargetList"`
	RequestUniqueKey string `json:"requestUniqueKey"`
}

func (m SemiGoodsStockUpdateRequest) validate() error {
	return nil

}

// Update 更新半托管商品库存
// https://partner-us.temu.com/documentation?menu_code=fb16b05f7a904765aac4af3a24b87d4a&sub_menu_code=429ffa60b265451d9421cd5a2004eeef
func (s semiGoodsStockService) Update(ctx context.Context, req SemiGoodsStockUpdateRequest) (*entity.SemiGoodsStockUpdateResult, bool, error) {
	if err := req.validate(); err != nil {
		return nil, false, err
	}

	var result = struct {
		normal.Response
		entity.SemiGoodsStockUpdateResult
		OperateResult bool `json:"operateResult"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(map[string]SemiGoodsStockUpdateRequest{"request": req}).
		SetResult(&result).
		Post("bg.local.goods.stock.edit")

	if err = recheckError(resp, result.Response, err); err != nil {
		return nil, false, err
	}

	if !result.OperateResult {
		reasons := make([]string, 0)
		for _, item := range result.SkuStockEditStatusInfoList {
			if !item.StockEditStatus {
				reasons = append(reasons, item.ErrorMsg)
			}
		}
		if len(reasons) > 0 {
			reasons = append(reasons, "更新库存失败，原因未知")
		}
		return nil, false, errors.New(strings.Join(reasons, "; "))
	}
	return &result.SemiGoodsStockUpdateResult, true, nil
}
