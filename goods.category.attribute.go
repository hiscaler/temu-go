package temu

import (
	"context"
	"fmt"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

// 商品属性列表服务
type goodsCategoryAttributeService service

// Query 按类目查询货品属性
// https://seller.kuajingmaihuo.com/sop/view/728777295758127187#6bz75P
func (s goodsCategoryAttributeService) Query(ctx context.Context, categoryId int64) (*entity.GoodsCategoryAttribute, error) {
	if categoryId <= 0 {
		return nil, fmt.Errorf("无效的分类：%d", categoryId)
	}

	var result = struct {
		normal.Response
		Result *entity.GoodsCategoryAttribute `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(map[string]int64{"catId": categoryId}).
		SetResult(&result).
		Post("bg.goods.attrs.get")
	if err = recheckError(resp, result.Response, err); err != nil {
		return nil, err
	}

	if result.Result == nil {
		return nil, fmt.Errorf("类目 %d 无货品属性", categoryId)
	}

	return result.Result, nil
}
