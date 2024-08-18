package temu

import (
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

type mallService service

// IsSemiManaged 是否为半托管店铺
// https://seller.kuajingmaihuo.com/sop/view/634117628601810731#uJ0fSb
func (s shipOrderStagingService) IsSemiManaged() (bool, error) {
	var result = struct {
		normal.Response
		Result struct {
			SemiManagedMall bool `json:"semiManagedMall"`
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetResult(&result).
		Post("bg.mall.info.get")
	if err == nil {
		err = parseResponse(resp, result.Response)
	}
	if err != nil {
		return false, err
	}

	return result.Result.SemiManagedMall, nil
}

// Permission 查询店铺权限
// https://seller.kuajingmaihuo.com/sop/view/634117628601810731#uJ0fSb
func (s shipOrderStagingService) Permission() (p entity.MallPermission, err error) {
	var result = struct {
		normal.Response
		Result entity.MallPermission `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetResult(&result).
		Post("bg.open.accesstoken.info.get")
	if err == nil {
		err = parseResponse(resp, result.Response)
	}
	if err != nil {
		return
	}

	return result.Result, nil
}
