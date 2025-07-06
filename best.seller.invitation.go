package temu

import (
	"context"

	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

// best seller 招标单服务
type bestSellerInvitationService service

type BestSellerInvitationQueryParams struct {
	normal.ParameterWithPager
}

func (m BestSellerInvitationQueryParams) validate() error {
	return nil
}

// Query best seller 招标单查询
// https://partner.kuajingmaihuo.com/document?cataId=875198836203&docId=917138067169
func (s bestSellerInvitationService) Query(ctx context.Context, params BestSellerInvitationQueryParams) (items []entity.BestSellerInvitation, total, totalPages int, isLastPage bool, err error) {
	if err = params.validate(); err != nil {
		err = invalidInput(err)
		return
	}

	var result = struct {
		normal.Response
		Result struct {
			List  []entity.BestSellerInvitation `json:"list"`
			Total int                           `json:"total"`
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(params).
		SetResult(&result).
		Post("temu.best.seller.invitation.query")
	if err = recheckError(resp, result.Response, err); err != nil {
		return
	}

	items = result.Result.List
	total, totalPages, isLastPage = parseResponseTotal(params.Page, params.PageSize, result.Result.Total)
	return
}
