package temu

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
	"github.com/hiscaler/temu-go/validators/is"
	"gopkg.in/guregu/null.v4"
	"strconv"
	"time"
)

// 商品数据服务
type goodsService service

type GoodsQueryParams struct {
	normal.ParameterWithPager
	Page                   int      `json:"page"`                             // 页码
	Cat1Id                 int      `json:"cat1Id,omitempty"`                 // 一级分类 ID
	Cat2Id                 int      `json:"cat2Id,omitempty"`                 // 二级分类 ID
	Cat3Id                 int      `json:"cat3Id,omitempty"`                 // 三级分类 ID
	Cat4Id                 int      `json:"cat4Id,omitempty"`                 // 四级分类 ID
	Cat5Id                 int      `json:"cat5Id,omitempty"`                 // 五级分类 ID
	Cat6Id                 int      `json:"cat6Id,omitempty"`                 // 六级分类 ID
	Cat7Id                 int      `json:"cat7Id,omitempty"`                 // 七级分类 ID
	Cat8Id                 int      `json:"cat8Id,omitempty"`                 // 八级分类 ID
	Cat9Id                 int      `json:"cat9Id,omitempty"`                 // 九级分类 ID
	Cat10Id                int      `json:"cat10Id,omitempty"`                // 十级分类 ID
	SkcExtCode             string   `json:"skcExtCode,omitempty"`             // 货品 SKC 外部编码
	SupportPersonalization int      `json:"supportPersonalization,omitempty"` // 是否支持定制品模板
	BindSiteIds            []int    `json:"bindSiteIds,omitempty"`            // 经营站点
	ProductName            string   `json:"productName,omitempty"`            // 货品名称
	ProductSkcIds          []int64  `json:"productSkcIds,omitempty"`          // SKC 列表
	QuickSellAgtSignStatus null.Int `json:"quickSellAgtSignStatus,omitempty"` // 快速售卖协议签署状态 0-未签署 1-已签署
	MatchJitMode           null.Int `json:"matchJitMode,omitempty"`           // 是否命中 JIT 模式
	SkcSiteStatus          null.Int `json:"skcSiteStatus,omitempty"`          // skc 加站点状态 (0: 未加入站点, 1: 已加入站点)
	CreatedAtStart         string   `json:"createdAtStart,omitempty"`         // 创建时间开始，毫秒级时间戳
	CreatedAtEnd           string   `json:"createdAtEnd,omitempty"`           // 创建时间结束，毫秒级时间戳
}

func (m GoodsQueryParams) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.CreatedAtStart,
			validation.When(m.CreatedAtStart != "" || m.CreatedAtEnd != "", validation.By(is.TimeRange(m.CreatedAtStart, m.CreatedAtEnd, time.DateOnly))),
		),
	)
}

// Query 货品列表查询
// https://seller.kuajingmaihuo.com/sop/view/750197804480663142#SjadVR
func (s goodsService) Query(ctx context.Context, params GoodsQueryParams) (items []entity.Goods, total, totalPages int, isLastPage bool, err error) {
	params.TidyPager()
	if params.Page <= 0 {
		params.Page = 1
	}
	if err = params.validate(); err != nil {
		return
	}

	if params.CreatedAtStart != "" && params.CreatedAtEnd != "" {
		t, _ := time.ParseInLocation(time.DateTime, params.CreatedAtStart+" 00:00:00", time.Local)
		params.CreatedAtStart = strconv.Itoa(int(t.UnixMilli()))
		t, _ = time.ParseInLocation(time.DateTime, params.CreatedAtEnd+" 23:59:59", time.Local)
		params.CreatedAtEnd = strconv.Itoa(int(t.UnixMilli()))
	}
	var result = struct {
		normal.Response
		Result struct {
			Data       []entity.Goods `json:"data"`
			TotalCount int            `json:"totalCount"`
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(params).
		SetResult(&result).
		Post("bg.goods.list.get")
	if err = recheckError(resp, result.Response, err); err != nil {
		return
	}

	items = result.Result.Data
	total, totalPages, isLastPage = parseResponseTotal(params.Page, params.PageSize, result.Result.TotalCount)
	return
}

// One 根据商品 SKC ID 查询
func (s goodsService) One(ctx context.Context, productSkcId int64) (item entity.Goods, err error) {
	items, _, _, _, err := s.Query(ctx, GoodsQueryParams{ProductSkcIds: []int64{productSkcId}})
	if err != nil {
		return
	}

	for _, v := range items {
		if v.ProductSkcId == productSkcId {
			return v, nil
		}
	}

	return item, ErrNotFound
}
