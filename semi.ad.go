package temu

import (
	"context"
	"errors"

	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

type semiAdService service

type SemiAdQueryParams struct {
	GoodsList             []int64 `url:"goodsList"`
	SelectedAdManagedType int     `url:"selectedAdManagedType"`
}

func (m SemiAdQueryParams) validate() error {
	return nil
}

// Query 广告查询
// https://partner-us.temu.com/documentation?menu_code=fb16b05f7a904765aac4af3a24b87d4a&sub_menu_code=66db5438c37446f49c122829489ac6d4
func (s semiAdService) Query(ctx context.Context, params SemiAdQueryParams) (items []entity.SemiAdDetail, err error) {
	if err = params.validate(); err != nil {
		return nil, invalidInput(err)
	}

	var result = struct {
		normal.Response
		Result struct {
			AdsDetail []entity.SemiAdDetail `json:"adsDetail"`
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(params).
		SetResult(&result).
		Post("temu.searchrec.ad.detail.query")
	if err = recheckError(resp, result.Response, err); err != nil {
		return
	}

	return result.Result.AdsDetail, nil
}

type SemiAdCreateRequest struct {
	GoodsId  int64 `json:"goodsId"`  // Goods id
	Roas     int   `json:"roas"`     // Target roas, multiply the actual value by 10,000
	RoasType int   `json:"roasType"` // 0: ad, 1: overAll roas
	Budget   int64 `json:"budget"`   // Advertising daily budget amount, if there is no limit, use -1
}

func (m SemiAdCreateRequest) validate() error {
	return nil
}

type SemiAdCreateResponse struct {
	AlreadyCreatedGoodsNum int `json:"alreadyCreatedGoodsNum"`
	CreateGoodsFailMap     struct {
		Key   string `json:"$key"`
		Value string `json:"$value"`
	} `json:"createGoodsFailMap"`
	SuccessGoodsIdLists     []int `json:"successGoodsIdLists"`
	SuccessCreateProductNum int   `json:"successCreateProductNum"`
	CreateGoodsFailObjList  []struct {
		Reason  string `json:"reason"`
		GoodsId int    `json:"goodsId"`
		Success bool   `json:"success"`
	} `json:"createGoodsFailObjList"`
}

// Create 广告创建
// https://partner-us.temu.com/documentation?menu_code=fb16b05f7a904765aac4af3a24b87d4a&sub_menu_code=7bc9231776304158a895e41a816b7805
func (s semiAdService) Create(ctx context.Context, reqs ...SemiAdCreateRequest) (*SemiAdCreateResponse, error) {
	for _, req := range reqs {
		if err := req.validate(); err != nil {
			return nil, invalidInput(err)
		}
	}

	var result = struct {
		normal.Response
		Result SemiAdCreateResponse `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(map[string][]SemiAdCreateRequest{"createAdReqs": reqs}).
		SetResult(&result).
		Post("temu.searchrec.ad.create")
	if err = recheckError(resp, result.Response, err); err != nil {
		return nil, err
	}

	return &result.Result, nil
}

type SemiAdUpdateRequest struct {
	GoodsId int64 `json:"goodsId"` // Goods id
	Roas    int   `json:"roas"`    // Target roas, multiply the actual value by 10,000
	Budget  int64 `json:"budget"`  // Advertising daily budget amount, if there is no limit, use -1
	Status  int   `json:"status"`  // Modification type: 1:delete, 2:pause, 3:open, 4:modify budget, 5:modify roas
}

func (m SemiAdUpdateRequest) validate() error {
	return nil

}

// Update 广告更新
// https://partner-us.temu.com/documentation?menu_code=fb16b05f7a904765aac4af3a24b87d4a&sub_menu_code=0b7140898262428eb8a4b28609112651
func (s semiAdService) Update(ctx context.Context, req SemiAdUpdateRequest) (bool, error) {
	if err := req.validate(); err != nil {
		return false, invalidInput(err)
	}

	var result = struct {
		normal.Response
		Result struct {
			SuccessModifyProductNum int `json:"successModifyProductNum"`
			ModifyGoodsRespList     []struct {
				GoodsId int64  `json:"goodsId"`
				Reason  string `json:"reason"`
				Success bool   `json:"success"`
			} `json:"modifyGoodsRespList"`
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(map[string]SemiAdUpdateRequest{"modifyAdDTO": req}).
		SetResult(&result).
		Post("temu.searchrec.ad.modify")
	if err = recheckError(resp, result.Response, err); err != nil {
		return false, err
	}

	unknownError := errors.New("未知错误")
	for _, goodsResp := range result.Result.ModifyGoodsRespList {
		if goodsResp.GoodsId == req.GoodsId {
			if goodsResp.Success {
				return true, nil
			} else {
				if goodsResp.Reason != "" {
					return false, errors.New(goodsResp.Reason)
				} else {
					return false, unknownError
				}
			}
		}
	}

	return false, unknownError
}
