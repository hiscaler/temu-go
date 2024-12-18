package temu

import (
	"context"
	"github.com/hiscaler/temu-go/normal"
)

// 商品尺码表模板服务
type goodsSizeChartTemplateService service

// Create 生成尺码表模板
// https://seller.kuajingmaihuo.com/sop/view/415794628056821162#n0Wlda
func (s *goodsSizeChartTemplateService) Create(ctx context.Context, businessId int) (tempBusinessId int, err error) {
	var result = struct {
		normal.Response
		Result struct {
			TempBusinessId int `json:"tempBusinessId"` // 临时模板 Id
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(map[string]int{"tempBusinessId": businessId}).
		SetResult(&result).
		Post("bg.goods.sizecharts.template.create")
	if err = recheckError(resp, result.Response, err); err != nil {
		return
	}

	return result.Result.TempBusinessId, nil
}
