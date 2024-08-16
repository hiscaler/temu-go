package temu

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

type goodsSizeChartClassService service

type GoodsSizeChartClassQueryParams struct {
	CatId int `json:"catId,omitempty"` // 叶子类目id，通过bg.goods.cats.get获取
}

func (m GoodsSizeChartClassQueryParams) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.CatId, validation.Required.Error("类目 ID 不能为空。")),
	)
}

// All 查询尺码表模板
func (s *goodsSizeChartClassService) All(params GoodsSizeChartClassQueryParams) (items []entity.GoodsSizeChartClass, err error) {
	if err = params.Validate(); err != nil {
		return
	}

	var result = struct {
		normal.Response
		Result struct {
			SizeSpecClassCat []entity.GoodsSizeChartClass `json:"sizeSpecClassCat"` // 尺码分类对象
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().SetBody(params).SetResult(&result).Post("bg.goods.sizecharts.class.get")
	if err == nil {
		err = parseResponse(resp, result.Response)
	}
	if err != nil {
		return
	}

	return result.Result.SizeSpecClassCat, nil
}
