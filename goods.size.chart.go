package temu

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

// 商品尺码服务
type goodsSizeChartService service

// GoodsSizeChartQueryParams
// Page 第一页从 0 开始
type GoodsSizeChartQueryParams struct {
	normal.ParameterWithPager
	CatId int64 `json:"catId,omitempty"` // 类目 ID
}

func (m GoodsSizeChartQueryParams) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Page, validation.Required.Error("页码不能为空")),
		validation.Field(&m.PageSize, validation.Required.Error("页面大小不能为空")),
	)
}

// Query 查询尺码表模板
func (s *goodsSizeChartService) Query(ctx context.Context, params GoodsSizeChartQueryParams) (items []entity.GoodsSizeChart, err error) {
	params.TidyPager(0)
	if err = params.validate(); err != nil {
		err = invalidInput(err)
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
	if err = recheckError(resp, result.Response, err); err != nil {
		return
	}

	return result.Result.labelCodePageResult.SizeSpecDataVOList, nil
}
