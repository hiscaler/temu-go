package temu

import (
	"context"
	"github.com/hiscaler/temu-go/normal"
)

// 商品规格服务
type goodsSpecificationService service

type GoodsSpecificationCreateRequest struct {
	ParentSpecId int `json:"parentSpecId"` // 父规格 id
	// 限制的特殊字符枚举如下
	// - emoji表情符号
	// - }、{、!、<、@、>、？、$、^、…、€、®、™、©、£、†、o、â、¥、¢、‡
	SpecName string `json:"specName"` // 子规格名称
}

func (m GoodsSpecificationCreateRequest) validate() error {
	return nil
}

// Create 生成规格（bg.goods.spec.create）
//
//	https://seller.kuajingmaihuo.com/sop/view/728777295758127187#MOa2Iu
func (s goodsSpecificationService) Create(ctx context.Context, request GoodsSpecificationCreateRequest) (int, error) {
	if err := request.validate(); err != nil {
		return 0, err
	}
	var result = struct {
		normal.Response
		Result struct {
			SpecId int `json:"specId"`
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(request).
		SetResult(&result).
		Post("bg.goods.spec.create")
	if err = recheckError(resp, result.Response, err); err != nil {
		return 0, err
	}

	return result.Result.SpecId, nil
}
