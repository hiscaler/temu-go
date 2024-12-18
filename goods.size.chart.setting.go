package temu

import (
	"context"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

// 商品尺码设置服务
type goodsSizeChartSettingService service

// View 查询尺码模板规则
func (s *goodsSizeChartSettingService) View(ctx context.Context, categoryId int64) (data entity.GoodsSizeChartSetting, err error) {
	var result = struct {
		normal.Response
		Result entity.GoodsSizeChartSetting `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(map[string]int64{"catId": categoryId}).
		SetResult(&result).
		Post("bg.goods.sizecharts.settings.get")
	if err = recheckError(resp, result.Response, err); err != nil {
		return
	}

	return result.Result, nil
}
