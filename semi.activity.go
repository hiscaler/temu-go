package temu

import (
	"context"

	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

// semiActivityService 半托活动服务
type semiActivityService service

type SemiActivityQueryParams struct{}

func (m SemiActivityQueryParams) validate() error {
	return nil
}

// Query 活动查询
// https://partner-us.temu.com/documentation?menu_code=fb16b05f7a904765aac4af3a24b87d4a&sub_menu_code=05820fed7179430c8e353905692d51b6
func (s semiActivityService) Query(ctx context.Context, params SemiActivityQueryParams) (items []entity.SemiActivity, err error) {
	if err = params.validate(); err != nil {
		return nil, invalidInput(err)
	}

	var result = struct {
		normal.Response
		Result struct {
			Total        int64                 `json:"total"`
			ActivityList []entity.SemiActivity `json:"activityList"`
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(params).
		SetResult(&result).
		Post("bg.promotion.activity.query")
	if err = recheckError(resp, result.Response, err); err != nil {
		return
	}

	return result.Result.ActivityList, nil
}

// SemiActivityEnrollRequest 半托活动报名请求
type SemiActivityEnrollRequest struct {
	ActivityId  int64 `json:"activityId"`
	EnrollGoods struct {
		GoodsId          int64  `json:"goodsId"`
		ActivityQuantity int    `json:"activityQuantity"` // activity quantity, the quantity you set for participating activities. This is independent of product quantity
		TraceCode        string `json:"traceCode,omitempty"`
		EnrollSkuList    []struct {
			SkuId                 int64   `json:"skuId"`
			ActivitySupplierPrice float64 `json:"activitySupplierPrice"` // the base price you set for activities, the items in the activity quantity are sold at the activity price and settled at the activity base price.
		} `json:"enrollSkuList"` // 	sku information required for registering activity
	} `json:"enrollGoods"`
}

func (s SemiActivityEnrollRequest) validate() error {
	return nil
}

// Enroll 半托活动报名
// https://partner-us.temu.com/documentation?menu_code=fb16b05f7a904765aac4af3a24b87d4a&sub_menu_code=27a87ec9d0d94273a48096c050f17854
func (s semiActivityService) Enroll(ctx context.Context, req SemiActivityEnrollRequest) (*entity.SemiActivityEnrollResult, error) {
	if err := req.validate(); err != nil {
		return nil, invalidInput(err)
	}

	var result = struct {
		normal.Response
		Result entity.SemiActivityEnrollResult `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(req).
		SetResult(&result).
		Post("bg.promotion.activity.goods.enroll")
	if err = recheckError(resp, result.Response, err); err != nil {
		return nil, err
	}

	return &result.Result, nil
}
