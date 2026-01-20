package temu

import (
	"context"
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
	"github.com/samber/lo"
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

func (m AdCreateRequestItem) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.ProductId, validation.Required.Error("货品 ID 不能为空")),
		validation.Field(&m.Roas, validation.Required.Error("目标广告投资回报率不能为空")),
		validation.Field(&m.Budget, validation.Required.Error("广告日预算金额不能为空")),
	)
}

type AdCreateRequest struct {
	CreateAdReqs []AdCreateRequestItem `json:"createAdReqs"` // 创建广告请求
}

func (m AdCreateRequest) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.CreateAdReqs,
			validation.Required.Error("创建广告请求不能为空"),
			validation.Each(validation.By(func(value interface{}) error {
				v, ok := value.(AdCreateRequestItem)
				if !ok {
					return errors.New("创建广告请求参数错误")
				}
				return v.validate()
			})),
		),
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

type AdUpdateRequestItem struct {
	ProductId int64 `json:"productId"` // 货品 ID
	Roas      int   `json:"roas"`      // 目标广告投资回报率，按照实际值乘 10000
	Budget    int   `json:"budget"`    // 广告日预算金额，不限制则传 -1
	Status    int   `json:"status"`    // 修改类型：2:暂停, 3:开启
}

func (m AdUpdateRequestItem) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.ProductId, validation.Required.Error("货品 ID 不能为空")),
		validation.Field(&m.Roas, validation.Required.Error("目标广告投资回报率不能为空")),
		validation.Field(&m.Budget, validation.Required.Error("广告日预算金额不能为空")),
		validation.Field(&m.Status, validation.Required.Error("修改类型不能为空"), validation.In(2, 3).ErrorObject(validation.NewError("422", "无效的修改类型 {{.value}}").SetParams(map[string]interface{}{"value": m.Status}))),
	)
}

type AdUpdateRequest struct {
	ModifyAdDTOs []AdUpdateRequestItem `json:"modifyAdDTOs"`
}

func (m AdUpdateRequest) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.ModifyAdDTOs,
			validation.Required.Error("修改广告请求不能为空"),
			validation.Each(validation.By(func(value interface{}) error {
				v, ok := value.(AdUpdateRequestItem)
				if !ok {
					return errors.New("修改广告请求参数错误")
				}
				return v.validate()
			})),
		),
	)
}

type AdUpdateResult = AdCreateResult

// Update 批量修改广告接口
// https://agentpartner.temu.com/document?cataId=875198836203&docId=931828782212
func (s adService) Update(ctx context.Context, request AdUpdateRequest) ([]AdUpdateResult, error) {
	if err := request.validate(); err != nil {
		return nil, invalidInput(err)
	}

	var result = struct {
		normal.Response
		Result struct {
			SuccessModifyProductNum int `json:"successModifyProductNum"` // 创建成功商品数量
			ModifyGoodsRespList     []struct {
				ProductId int64       `json:"productId"`
				Success   bool        `json:"success"`
				Reason    null.String `json:"reason"`
			} `json:"modifyGoodsRespList"`
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(request).
		SetResult(&result).
		Post("bg.glo.searchrec.ad.batch.modify")
	if err = recheckError(resp, result.Response, err); err != nil {
		return nil, err
	}

	results := make([]AdUpdateResult, len(result.Result.ModifyGoodsRespList))
	for i, item := range result.Result.ModifyGoodsRespList {
		results[i] = AdUpdateResult{
			ProductId: item.ProductId,
			Success:   item.Success,
			Message:   item.Reason,
		}
	}

	return results, nil
}

// Return on Advertising Spend [ROAS]

type AdRoasRequest struct {
	GoodsInfoList []struct {
		ProductIds []int64 `json:"productIds"` // 货品 ID
	} `json:"goodsInfoList"`
}

func (m AdRoasRequest) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.GoodsInfoList,
			validation.Required.Error("货品信息列表不能为空"),
			validation.Each(validation.By(func(value interface{}) error {
				_, ok := value.(struct {
					ProductIds []int64 `json:"productIds"` // 货品 ID
				})
				if !ok {
					return errors.New("货品信息参数错误")
				}
				return nil
			})),
		),
	)
}

type AdRoasResult struct {
	ProductId int64 `json:"productId"`
	PredList  []struct {
		Roas string `json:"roas"` // 推荐 roas(广告投资回报率)
	} `json:"predList"`
}

// Roas 广告投资回报率查询接口
// https://agentpartner.temu.com/document?cataId=875198836203&docId=929735887634
func (s adService) Roas(ctx context.Context, productId ...int64) ([]AdRoasResult, error) {
	if len(productId) == 0 {
		return nil, errors.New("请提供查询的商品 ID")
	}

	var result = struct {
		normal.Response
		Result struct {
			QueryAdBidResult []AdRoasResult `json:"queryAdBidResult"`
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(map[string]any{
			"goodsInfoList": lo.Map(productId, func(item int64, index int) struct {
				ProductId int64 `json:"product_id"`
			} {
				return struct {
					ProductId int64 `json:"product_id"`
				}{ProductId: item}
			}),
		}).
		SetResult(&result).
		Post("bg.glo.searchrec.ad.roas.pred")
	if err = recheckError(resp, result.Response, err); err != nil {
		return nil, err
	}

	return result.Result.QueryAdBidResult, nil
}

type AdLogQueryParams struct {
	ProductId null.Int `json:"productId,omitempty"` // 货品 ID
	StartTime int64    `json:"startTime"`           // 查询开始时间，毫秒级时间戳（值以当地时间0点为开始时间）
	EndTime   int64    `json:"endTime"`             // 查询结束时间，毫秒级时间戳（值以当地时间23点59分59秒999毫秒为结束时间）
}

func (m AdLogQueryParams) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.ProductId, validation.When(m.ProductId.Valid, validation.By(func(value interface{}) error {
			v, ok := value.(int64)
			if !ok {
				return errors.New("货品 ID 参数错误")
			}
			return validation.Validate(v, validation.Min(1).Error("无效的货品 ID"))
		}))),
		validation.Field(&m.StartTime, validation.Required.Error("查询开始时间不能为空")),
		validation.Field(&m.EndTime, validation.Required.Error("查询结束时间不能为空")),
	)
}

// Logs 操作日志查询接口
// https://agentpartner.temu.com/document?cataId=875198836203&docId=931830463288
func (s adService) Logs(ctx context.Context, params AdLogQueryParams) ([]entity.AdLog, error) {
	if err := params.validate(); err != nil {
		return nil, invalidInput(err)
	}

	var result = struct {
		normal.Response
		Result []entity.AdLog `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(params).
		SetResult(&result).
		Post("bg.glo.searchrec.ad.log.query")
	if err = recheckError(resp, result.Response, err); err != nil {
		return nil, err
	}

	return result.Result, nil
}

type AdProductReportQueryParams struct {
	ProductId int64 `json:"productId"` // 货品 ID
	StartTs   int64 `json:"startTs"`   // 查询结束时间，毫秒级时间戳（值以当地时间23点59分59秒999毫秒为结束时间）
	EndTs     int64 `json:"endTs"`     // 查询开始时间，毫秒级时间戳（值以当地时间0点为开始时间）
}

func (m AdProductReportQueryParams) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.ProductId, validation.Required.Error("无效的货品 ID")),
		validation.Field(&m.StartTs, validation.Required.Error("查询开始时间不能为空")),
		validation.Field(&m.EndTs, validation.Required.Error("查询结束时间不能为空")),
	)
}

// ProductReport 广告商品投放数据效果
// https://agentpartner.temu.com/document?cataId=875198836203&docId=929739237497
func (s adService) ProductReport(ctx context.Context, params AdProductReportQueryParams) ([]entity.AdReport, error) {
	if err := params.validate(); err != nil {
		return nil, invalidInput(err)
	}

	var result = struct {
		normal.Response
		Result struct {
			ReportInfo struct {
				ReportsItemList []entity.AdReport      `json:"reportsItemList"` // 分时间段报表信息，按照天级或小时级划分的报表信息（请求时间跨度大于一天按照天级划分，等于一天按照小时级划分）
				ReportsSummary  entity.AdReportSummary `json:"reportsSummary"`  // 整体报表信息
			} `json:"reportInfo"`
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(params).
		SetResult(&result).
		Post("bg.glo.searchrec.ad.reports.goods.query")
	if err = recheckError(resp, result.Response, err); err != nil {
		return nil, err
	}

	return result.Result.ReportInfo.ReportsItemList, nil
}
