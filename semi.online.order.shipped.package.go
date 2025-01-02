package temu

import (
	"context"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

// 已发货包裹服务
type semiOnlineOrderShippedPackageService service

type SemiPlatformLogisticsShippedPackageQueryParams struct {
	normal.ParameterWithPager
	PageNumber        int      `json:"pageNumber"`        // 第几页
	ParentOrderSnList []string `json:"parentOrderSnList"` // PO 单号列表
	OrderSnList       []string `json:"orderSnList"`       // O 单号列表
}

func (m SemiPlatformLogisticsShippedPackageQueryParams) validate() error {
	return nil
}

// Confirm 确认包裹发货接口（bg.logistics.shipped.package.confirm）
// https://seller.kuajingmaihuo.com/sop/view/144659541206936016#92SpUJ
func (s semiOnlineOrderShippedPackageService) Confirm(ctx context.Context, params SemiPlatformLogisticsShippedPackageQueryParams) (items []entity.SemiPlatformLogisticsUnshippedPackage, total, totalPages int, isLastPage bool, err error) {
	params.TidyPager()
	params.PageNumber = params.Pager.Page
	if err = params.validate(); err != nil {
		err = invalidInput(err)
		return
	}

	var result = struct {
		normal.Response
		Result struct {
			TotalItemNum     int                                            `json:"totalItemNum"`
			UnshippedPackage []entity.SemiPlatformLogisticsUnshippedPackage `json:"unshippedPackage"` // 待确认包裹
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(params).
		SetResult(&result).
		Post("bg.logistics.shipped.package.confirm")
	if err = recheckError(resp, result.Response, err); err != nil {
		return
	}

	items = result.Result.UnshippedPackage
	total, totalPages, isLastPage = parseResponseTotal(params.Page, params.PageSize, result.Result.TotalItemNum)
	return
}
