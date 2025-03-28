package temu

import (
	"context"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

type mallService struct {
	service
	DeliveryAddress mallDeliveryAddressService
}

// Type 店铺类型
// https://seller.kuajingmaihuo.com/sop/view/634117628601810731#uJ0fSb
func (s mallService) Type(ctx context.Context) (entity.MallType, error) {
	var result = struct {
		normal.Response
		Result entity.MallType `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetResult(&result).
		Post("bg.mall.info.get")
	if err = recheckError(resp, result.Response, err); err != nil {
		return entity.MallType{}, err
	}

	return result.Result, nil
}

// Permission 查询店铺权限
// https://seller.kuajingmaihuo.com/sop/view/634117628601810731#3tCaqU
func (s mallService) Permission(ctx context.Context) (p entity.MallPermission, err error) {
	var result = struct {
		normal.Response
		Result entity.MallPermission `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetResult(&result).
		Post("bg.open.accesstoken.info.get")
	if err = recheckError(resp, result.Response, err); err != nil {
		return
	}

	return result.Result, nil
}

// AccessToken 获取 Access Token
// https://seller.kuajingmaihuo.com/sop/view/634117628601810731#ov9GUf
func (s mallService) AccessToken(ctx context.Context, currentAccessToken, code string) (at entity.AccessToken, err error) {
	var result = struct {
		normal.Response
		Result entity.AccessToken `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(map[string]string{"accessToken": currentAccessToken, "code": code}).
		SetResult(&result).
		Post("bg.open.accesstoken.create")
	if err = recheckError(resp, result.Response, err); err != nil {
		return at, err
	}

	return result.Result, nil
}
