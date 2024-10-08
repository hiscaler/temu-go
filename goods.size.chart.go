package temu

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

type goodsSizeChartService service

type GoodsSizeChartQueryParams struct {
	Page     int `json:"offset"`             // 锚点（第一页传0）
	PageSize int `json:"pageSize,omitempty"` // 页面大小
	CatId    int `json:"catId,omitempty"`    // 类目 ID
}

func (m GoodsSizeChartQueryParams) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Page, validation.Required.Error("页码不能为空。")),
		validation.Field(&m.PageSize, validation.Required.Error("页面大小不能为空。")),
	)
}

// All 查询尺码表模板
func (s *goodsSizeChartService) All(ctx context.Context, params GoodsSizeChartQueryParams) (items []entity.GoodsSizeChart, err error) {
	if params.Page < 0 {
		params.Page = 0
	}
	if params.PageSize <= 0 {
		params.PageSize = 10
	} else if params.PageSize > 500 {
		params.PageSize = 500
	}
	if err = params.Validate(); err != nil {
		return
	}

	var result = struct {
		normal.Response
		Result struct {
			labelCodePageResult struct {
				TotalCount         int                     `json:"totalCount"`         // 总数
				SizeSpecDataVOList []entity.GoodsSizeChart `json:"sizeSpecDataVOList"` // 列表数据
			}
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(params).
		SetResult(&result).
		Post("bg.goods.sizecharts.get")
	if err == nil {
		err = parseResponse(resp, result.Response)
	}
	if err != nil {
		return
	}

	return result.Result.labelCodePageResult.SizeSpecDataVOList, nil
}

// Create 生成尺码表模板
// https://seller.kuajingmaihuo.com/sop/view/415794628056821162#n0Wlda
func (s *goodsSizeChartService) Create(ctx context.Context, businessId int) (tempBusinessId int, err error) {
	var result = struct {
		normal.Response
		Result struct {
			TempBusinessId int `json:"tempBusinessId"` // 临时模板ID
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(map[string]int{"tempBusinessId": businessId}).
		SetResult(&result).
		Post("bg.goods.sizecharts.template.create")
	if err == nil {
		err = parseResponse(resp, result.Response)
	}
	if err != nil {
		return
	}
	return result.Result.TempBusinessId, nil
}
