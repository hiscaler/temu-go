package temu

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
	"gopkg.in/guregu/null.v4"
)

// ad 广告服务
type adService service

type ADQueryParams struct {
	GoodsInfoList []int `json:"productId"`
}

func (m ADQueryParams) validate() error {
	return nil
}

// Query 广告投放查询接口
// https://agentpartner.temu.com/document?cataId=875198836203&docId=929736716892
func (s adService) Query(ctx context.Context, params ADQueryParams) ([]entity.Ad, error) {
	if err := params.validate(); err != nil {
		return nil, invalidInput(err)
	}

	var result = struct {
		normal.Response
		Result struct {
			AdsDetail []entity.Ad `json:"adsDetail"`
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(params).
		SetResult(&result).
		Post("bg.glo.searchrec.ad.detail.query")
	if err = recheckError(resp, result.Response, err); err != nil {
		return nil, err
	}

	return result.Result.AdsDetail, nil
}

type AdCreateRequestItem struct {
	ProductId int64 `json:"productId"` // 货品 ID
	Roas      int   `json:"roas"`      // 目标广告投资回报率，按照实际值乘 10000
	Budget    int   `json:"budget"`    // 广告日预算金额，不限制则传 -1
}
type AdCreateRequest struct {
	CreateAdReqs []AdCreateRequestItem `json:"createAdReqs"` // 创建广告请求
}

func (m AdCreateRequest) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.CreateAdReqs, validation.Required.Error("创建广告请求不能为空")),
	)
}

type AdCreateResult struct {
	ProductId int64       `json:"productId"`
	Success   bool        `json:"success"`
	Message   null.String `json:"message"`
}

// Create 创建广告
func (s adService) Create(ctx context.Context, request AdCreateRequest) ([]AdCreateResult, error) {
	if err := request.validate(); err != nil {
		return nil, invalidInput(err)
	}

	var result = struct {
		normal.Response
		Result struct {
			SuccessCreateProductNum int `json:"successCreateProductNum"` // 创建成功商品数量
			CreateGoodsFailObjList  []struct {
				ProductId int64       `json:"productId"`
				Success   bool        `json:"success"`
				Reason    null.String `json:"reason"`
			} `json:"createGoodsFailObjList"`
			SuccessProductIdLists []int64           `json:"successProductIdLists"`
			CreateProductFailMap  map[string]string `json:"createProductFailMap"`
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(request).
		SetResult(&result).
		Post("bg.glo.searchrec.ad.create")
	if err = recheckError(resp, result.Response, err); err != nil {
		return nil, err
	}

	results := make([]AdCreateResult, len(result.Result.CreateGoodsFailObjList))
	for i, item := range result.Result.CreateGoodsFailObjList {
		results[i] = AdCreateResult{
			ProductId: item.ProductId,
			Success:   item.Success,
			Message:   item.Reason,
		}
	}

	return results, nil
}
