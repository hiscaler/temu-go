package temu

import (
	"context"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/helpers"
	"github.com/hiscaler/temu-go/normal"
	"github.com/hiscaler/temu-go/validators/is"
)

// 订单服务（半托管专属，必须在 US/EU 网关调用）
type semiOrderService service

type OrderQueryParams struct {
	normal.ParameterWithPager
	PageNumber int `json:"pageNumber"` // 第几页
	// 父单状态，默认查全部枚举值如下
	// 0-全部
	// 1-”PENDING“，挂起中
	// 2-"UN_SHIPPING"，待发货
	// 3-"CANCELED",已取消
	// 4-"SHIPPED"，已发货
	// 5-“RECEIPTED”,已签收
	// 备注：
	// 本本订单还存在如下状态
	// 41-部分取消
	// 51-部分签收
	ParentOrderStatus         int      `json:"parentOrderStatus"`         // 父单状态
	CreateBefore              string   `json:"createBefore"`              // 父单创建时间结束查询时间，格式是秒时间戳。查询时间需要同时入参开始和结束时间才生效
	CreateAfter               string   `json:"createAfter"`               // 父单创建时间开始查询时间，格式是秒时间戳
	ParentOrderSnList         []string `json:"parentOrderSnList"`         // 父单号列表，单次请求最多 20 个
	ExpectShipLatestTimeStart string   `json:"expectShipLatestTimeStart"` // 期望最晚发货时间开始查询时间，格式是秒时间戳
	ExpectShipLatestTimeEnd   string   `json:"expectShipLatestTimeEnd"`   // 期望最晚发货时间结束查询时间，格式是秒时间戳。查询时间需要同时入参开始和结束时间才生效
	UpdateAtStart             string   `json:"updateAtStart"`             // 订单更新时间开始查询时间，格式是秒时间戳
	UpdateAtEnd               string   `json:"updateAtEnd"`               // 订单更新时间结束查询时间，格式是秒时间戳。查询时间需要同时入参开始和结束时间才生效
	RegionId                  int      `json:"regionId"`                  // 区域ID，美国-211
	// 子订单履约类型，具体枚举值如下：
	// 1. 数组只传入fulfillBySeller，只返回卖家履约子订单列表
	// 2. 数组只传入fulfillByCooperativeWarehouse，只返回合作仓履约子订单列表
	// 3. 数组传入fulfillBySeller和fulfillByCooperativeWarehouse，返回卖家履约子订单列表+合作仓履约子订单列表
	// 4. fulfillmentTypeList不传或者传了为空，默认返回卖家履约子订单列表
	FulfillmentTypeList []string `json:"fulfillmentTypeList"` // 子订单履约类型
}

func (m OrderQueryParams) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.CreateBefore,
			validation.When(m.CreateBefore != "" || m.CreateAfter != "", validation.By(is.TimeRange(m.CreateBefore, m.CreateAfter, time.DateTime))),
		),
		validation.Field(&m.ExpectShipLatestTimeStart,
			validation.When(m.ExpectShipLatestTimeStart != "" || m.ExpectShipLatestTimeEnd != "", validation.By(is.TimeRange(m.ExpectShipLatestTimeStart, m.ExpectShipLatestTimeEnd, time.DateTime))),
		),
		validation.Field(&m.UpdateAtStart,
			validation.When(m.UpdateAtStart != "" || m.UpdateAtEnd != "", validation.By(is.TimeRange(m.UpdateAtStart, m.UpdateAtEnd, time.DateTime))),
		),
	)
}

// Query 订单列表查询接口
// https://seller.kuajingmaihuo.com/sop/view/867739977041685428#r2WKrz
func (s semiOrderService) Query(ctx context.Context, params OrderQueryParams) (items []entity.PageItem, total, totalPages int, isLastPage bool, err error) {
	params.TidyPager()
	params.PageNumber = params.Pager.Page
	if err = params.validate(); err != nil {
		err = invalidInput(err)
		return
	}

	if params.CreateBefore != "" && params.CreateAfter != "" {
		if start, end, e := helpers.StrTime2UnixMilli(params.CreateBefore, params.CreateAfter); e == nil {
			params.CreateBefore = start
			params.CreateAfter = end
		}
	}

	if params.ExpectShipLatestTimeStart != "" && params.ExpectShipLatestTimeEnd != "" {
		if start, end, e := helpers.StrTime2UnixMilli(params.ExpectShipLatestTimeStart, params.ExpectShipLatestTimeEnd); e == nil {
			params.ExpectShipLatestTimeStart = start
			params.ExpectShipLatestTimeEnd = end
		}
	}

	if params.UpdateAtStart != "" && params.UpdateAtEnd != "" {
		if start, end, e := helpers.StrTime2UnixMilli(params.UpdateAtStart, params.UpdateAtEnd); e == nil {
			params.UpdateAtStart = start
			params.UpdateAtEnd = end
		}
	}

	var result = struct {
		normal.Response
		Result struct {
			Result entity.OrderResult `json:"result"`
		} `json:"result"`
	}{}

	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(params).
		SetResult(&result).
		Post("bg.order.list.get")
	if err = recheckError(resp, result.Response, err); err != nil {
		return
	}

	items = result.Result.Result.PageItems
	 
	total, totalPages, isLastPage = parseResponseTotal(params.Page, params.PageSize, result.Result.Result.TotalItemNum)
	return
}
