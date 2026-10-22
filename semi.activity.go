package temu

import (
	"context"
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

// semiActivityService 半托活动服务
type semiActivityService service

type SemiActivityQueryParams struct {
	normal.Pager
	PageNumber              int     `json:"pageNumber"`                        // 页码，从 1 开始，由 Pager.Page 映射
	ActivityEndTime         int64   `json:"activityEndTime,omitempty"`         // 活动结束时间，Unix 秒级时间戳
	ActivityIdList          []int64 `json:"activityIdList,omitempty"`          // 活动 ID 列表
	ActivityStatus          int     `json:"activityStatus,omitempty"`          // 活动状态：1-未开始，2-进行中，3-已结束
	ActivityStartTime       int64   `json:"activityStartTime,omitempty"`       // 活动开始时间，Unix 秒级时间戳
	ActivityType            int     `json:"activityType"`                      // 活动类型：2-闪购，13-高级大促，27-清仓，100-官方大促
	OnlyQueryJoinedActivity bool    `json:"onlyQueryJoinedActivity,omitempty"` // 是否只查询已参加的活动
}

func (m SemiActivityQueryParams) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Page,
			validation.Required.Error("页码不能为空"),
			validation.Min(1).Error("页码必须大于 0"),
		),
		validation.Field(&m.PageSize,
			validation.Required.Error("页面大小不能为空"),
			validation.Min(1).Error("页面大小必须大于 0"),
			validation.Max(entity.MaxPageSize).Error("页面大小不能超过 100"),
		),
		validation.Field(&m.ActivityType,
			validation.Required.Error("活动类型不能为空"),
			validation.In(2, 13, 27, 100).Error("无效的活动类型"),
		),
		validation.Field(&m.ActivityStatus,
			validation.When(m.ActivityStatus != 0,
				validation.In(1, 2, 3).Error("无效的活动状态"),
			),
		),
	)
}

// Query 活动查询
// https://partner-us.temu.com/documentation?menu_code=fb16b05f7a904765aac4af3a24b87d4a&sub_menu_code=05820fed7179430c8e353905692d51b6
func (s semiActivityService) Query(ctx context.Context, params SemiActivityQueryParams) (items []entity.SemiActivity, err error) {
	params.PageNumber = params.TidyPager().Page
	params.OmitPage()
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
	return validation.ValidateStruct(&s,
		validation.Field(&s.ActivityId,
			validation.Required.Error("活动 ID 不能为空"),
			validation.Min(1).Error("活动 ID 必须大于 0"),
		),
		validation.Field(&s.EnrollGoods.GoodsId,
			validation.Required.Error("商品 ID 不能为空"),
			validation.Min(1).Error("商品 ID 必须大于 0"),
		),
		validation.Field(&s.EnrollGoods.ActivityQuantity,
			validation.Required.Error("活动数量不能为空"),
			validation.Min(1).Error("活动数量必须大于 0"),
		),
		validation.Field(&s.EnrollGoods.EnrollSkuList,
			validation.Required.Error("报名 SKU 列表不能为空"),
			validation.Each(validation.By(func(value any) error {
				v, ok := value.(struct {
					SkuId                 int64   `json:"skuId"`
					ActivitySupplierPrice float64 `json:"activitySupplierPrice"`
				})
				if !ok {
					return errors.New("无效的报名 SKU 数据")
				}

				return validation.ValidateStruct(&v,
					validation.Field(&v.SkuId,
						validation.Required.Error("SKU ID 不能为空"),
						validation.Min(1).Error("SKU ID 必须大于 0"),
					),
					validation.Field(&v.ActivitySupplierPrice,
						validation.Required.Error("活动供货价不能为空"),
						validation.Min(0).Error("活动供货价不能为负数"),
					),
				)
			})),
		),
	)
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
