package temu

import (
	"context"

	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

// 待发货包裹服务
type semiOnlineOrderUnshippedPackageService service

type SemiOnlineOrderPlatformLogisticsUnshippedPackageQueryParams struct {
	normal.ParameterWithPager
	PageNumber        int      `json:"pageNumber"`        // 第几页
	ParentOrderSnList []string `json:"parentOrderSnList"` // PO 单号列表
	OrderSnList       []string `json:"orderSnList"`       // O 单号列表
}

func (m SemiOnlineOrderPlatformLogisticsUnshippedPackageQueryParams) validate() error {
	return nil
}

// Query 下 Call 成功待发货包裹列表查询接口（bg.order.unshipped.package.get）
// https://seller.kuajingmaihuo.com/sop/view/144659541206936016#9QfRAV
func (s semiOnlineOrderUnshippedPackageService) Query(ctx context.Context, params SemiOnlineOrderPlatformLogisticsUnshippedPackageQueryParams) (items []entity.SemiOnlineOrderPlatformLogisticsUnshippedPackage, total, totalPages int, isLastPage bool, err error) {
	params.PageNumber = params.TidyPager().Page
	params.OmitPage()
	if err = params.validate(); err != nil {
		err = invalidInput(err)
		return
	}

	var result = struct {
		normal.Response
		Result struct {
			TotalItemNum     int                                                       `json:"totalItemNum"`
			UnshippedPackage []entity.SemiOnlineOrderPlatformLogisticsUnshippedPackage `json:"unshippedPackage"` // 待确认包裹
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
