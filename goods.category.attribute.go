package temu

import (
	"context"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

// 商品属性列表服务
type goodsCategoryAttributeService service

// Query 按类目查询货品属性
func (s goodsCategoryAttributeService) Query(ctx context.Context, categoryId int64) ([]entity.GoodsCategoryAttribute, error) {
	var result = struct {
		normal.Response
		Result struct {
			InputMaxSpecNum      int                             `json:"inputMaxSpecNum"`      // 模板允许的最大的自定义规格数量
			ChooseAllQualifySpec bool                            `json:"chooseAllQualifySpec"` // 限定规格是否必须全选
			SingleSpecValueNum   int                             `json:"singleSpecValueNum"`   // 单个自定义规格值上限
			Properties           []entity.GoodsCategoryAttribute `json:"properties"`           // 模板属性
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetResult(&result).
		Post("bg.goods.attrs.get")
	if err = recheckError(resp, result.Response, err); err != nil {
		return nil, err
	}

	return result.Result.Properties, nil
}
