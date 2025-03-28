package temu

import (
    "context"
    validation "github.com/go-ozzo/ozzo-validation/v4"
    "github.com/hiscaler/temu-go/entity"
    "github.com/hiscaler/temu-go/normal"
)

// goodsPriceService 供货价/核价/调价服务
type goodsPriceService service

// 货品品牌

type GoodsPriceQueryParams struct {
    ProductSkuIds []int64 `json:"productSkuIds"` // 货品 sku ID
}

func (m GoodsPriceQueryParams) validate() error {
    return validation.ValidateStruct(&m,
        validation.Field(&m.ProductSkuIds, validation.Required.Error("货品 sku ID 列表不能为空")),
    )
}

// Query 货品供货价查询
// https://partner.kuajingmaihuo.com/document?cataId=875198836203&docId=901410718805
func (s goodsPriceService) Query(ctx context.Context, params GoodsPriceQueryParams) (items []entity.ProductSkuSupplierPrice, err error) {
    if err = params.validate(); err != nil {
        return items, invalidInput(err)
    }

    var result = struct {
        normal.Response
        Result struct {
            ProductSkuSupplierPriceList []entity.ProductSkuSupplierPrice `json:"productSkuSupplierPriceList"`
        } `json:"result"`
    }{}
    resp, err := s.httpClient.R().
        SetContext(ctx).
        SetBody(params).
        SetResult(&result).
        Post("bg.goods.price.list.get")
    if err = recheckError(resp, result.Response, err); err != nil {
        return
    }

    return result.Result.ProductSkuSupplierPriceList, nil
}

// 分页查询核价单
// https://partner.kuajingmaihuo.com/document?cataId=875198836203&docId=899321422992

type GoodsPriceReviewQueryParams struct {
    normal.ParameterWithPager
    IdLt            int64 `json:"idLt"`            // id 范围查询最大值
    IdGt            int64 `json:"idGt"`            // id 范围查询最小值
    OrderStatusList []int `json:"orderStatusList"` // 核价单状态列表. 可选值含义说明:[0:待核价;1:待供应商确认;2:核价通过;3:核价驳回;4:废弃;5:价格同步中;]
}

func (m GoodsPriceReviewQueryParams) validate() error {
    return nil
}

func (s goodsPriceService) Reviews(ctx context.Context, params GoodsPriceReviewQueryParams) (items []entity.GoodsReviewSamplePrice, err error) {
    if err = params.validate(); err != nil {
        return items, invalidInput(err)
    }

    var result = struct {
        normal.Response
        Result struct {
            Total                 int                             `json:"total"`
            ReviewSamplePriceList []entity.GoodsReviewSamplePrice `json:"reviewSamplePriceList"` // 核价单 sku 及建议价
        } `json:"result"`
    }{}
    resp, err := s.httpClient.R().
        SetContext(ctx).
        SetBody(params).
        SetResult(&result).
        Post("bg.price.review.page.query")
    if err = recheckError(resp, result.Response, err); err != nil {
        return
    }

    return result.Result.ReviewSamplePriceList, nil
}
