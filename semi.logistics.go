package temu

import (
	"context"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

// 半托管物流服务
type semiLogisticsService service

// Companies 物流商查询
// https://seller.kuajingmaihuo.com/sop/view/867739977041685428#KLthO6
func (s semiLogisticsService) Companies(ctx context.Context, regionId int) ([]entity.SemiLogisticsCompany, error) {
	var result = struct {
		normal.Response
		Result struct {
			Items []entity.SemiLogisticsCompany `json:"items"` // 物流商列表
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(map[string]int{"regionId": regionId}).
		SetResult(&result).
		Post("bg.logistics.companies.get")
	if err = recheckError(resp, result.Response, err); err != nil {
		return nil, err
	}

	return result.Result.Items, nil
}
