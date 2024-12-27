package temu

import (
	"context"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

// 平台物流服务
type semiPlatformLogisticsService service

type SemiPlatformLogisticsUnshippedPackageQueryParams struct {
	normal.ParameterWithPager
	PageNumber        int      `json:"pageNumber"`        // 第几页
	ParentOrderSnList []string `json:"parentOrderSnList"` // PO单号列表
	OrderSnList       []string `json:"orderSnList"`       // O单号列表
}

func (m SemiPlatformLogisticsUnshippedPackageQueryParams) validate() error {
	return nil
}

// UnshippedPackages 下 Call 成功待发货包裹列表查询接口（bg.order.unshipped.package.get）
// https://seller.kuajingmaihuo.com/sop/view/144659541206936016#9QfRAV
func (s semiPlatformLogisticsService) UnshippedPackages(ctx context.Context, params SemiPlatformLogisticsUnshippedPackageQueryParams) (items []entity.SemiPlatformLogisticsUnshippedPackage, total, totalPages int, isLastPage bool, err error) {
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
		Post("bg.order.unshipped.package.get")
	if err = recheckError(resp, result.Response, err); err != nil {
		return
	}

	items = result.Result.UnshippedPackage
	total, totalPages, isLastPage = parseResponseTotal(params.Page, params.PageSize, result.Result.TotalItemNum)
	return
}