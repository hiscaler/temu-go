package temu

import (
	"context"

	"github.com/hiscaler/temu-go/normal"
)

// 商品尺码表模板服务
type goodsSizeChartTemplateService service

// Create 生成尺码表模板
// https://agentpartner.temu.com/document?cataId=875198836203&docId=877348687822
func (s *goodsSizeChartTemplateService) Create(ctx context.Context, businessId int64) (tempBusinessId int64, err error) {
	var result = struct {
		normal.Response
		Result struct {
			TempBusinessId int64 `json:"tempBusinessId"` // 临时模板 Id
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(map[string]int64{"tempBusinessId": businessId}).
		SetResult(&result).
		Post("bg.goods.sizecharts.template.create")
	if err = recheckError(resp, result.Response, err); err != nil {
		return
	}

	return result.Result.TempBusinessId, nil
}
