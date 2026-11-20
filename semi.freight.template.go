package temu

import (
	"context"

	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

// semiFreightTemplateService 半托运费模板服务
type semiFreightTemplateService service

// Templates 查询半托运费模板
// https://agentpartner.temu.com/document?cataId=875198836203&docId=929751463671
func (s semiFreightTemplateService) Templates(ctx context.Context, siteIds ...int) ([]entity.SemiFreightTemplate, error) {
	var result = struct {
		normal.Response
		Result struct {
			FreightTemplates []entity.SemiFreightTemplate `json:"freightTemplates"`
		} `json:"result"`
	}{}
	request := struct {
		SiteIds []int `json:"siteIds,omitempty"`
	}{
		SiteIds: siteIds,
	}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(request).
		SetResult(&result).
		Post("bg.glo.logistics.template.get")
	if err = recheckError(resp, result.Response, err); err != nil {
		return nil, err
	}

	return result.Result.FreightTemplates, nil
}
