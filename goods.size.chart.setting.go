package temu

import (
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

type goodsSizeChartSettingService service

// View 查询尺码模板规则
func (s *goodsSizeChartSettingService) View(catId int) (data entity.GoodsSizeChartSetting, err error) {
	var result = struct {
		normal.Response
		Result entity.GoodsSizeChartSetting `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetBody(map[string]int{"catId": catId}).
		SetResult(&result).
		Post("bg.goods.sizecharts.settings.get")
	if err == nil {
		err = parseResponse(resp, result.Response)
	}
	if err != nil {
		return
	}

	return result.Result, nil
}
