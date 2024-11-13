package temu

import (
	"context"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

type goodsCertificationService service

type GoodsCertificationQueryRequest struct {
	CertTypeList []int  `json:"certTypeList"` // 资质类型 id 列表
	ProductId    int64  `json:"productId"`    // 货品 id
	Language     string `json:"language"`     // 语言
}

func (m GoodsCertificationQueryRequest) validate() error {
	return nil
}

// Query 批量查询货品资质信息
// https://seller.kuajingmaihuo.com/sop/view/649320516224723675#Oq8dC9
func (s goodsCertificationService) Query(ctx context.Context, request GoodsCertificationQueryRequest) (certifications []entity.GoodsCertification, err error) {
	if err = request.validate(); err != nil {
		return
	}
	var result = struct {
		normal.Response
		Result struct {
			ProductCertList []entity.GoodsCertification `json:"productCertList"`
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(request).
		SetResult(&result).
		Post("bg.arbok.open.product.cert.query")
	if err == nil {
		err = parseResponse(resp, result.Response)
	}
	if err != nil {
		return
	}
	certifications = result.Result.ProductCertList
	return
}

type GoodsCertificationNeedUploadItemRequest struct {
	CertType  int   `json:"certType"`  // 资质类型
	ProductId int64 `json:"productId"` // 货品id
}

func (m GoodsCertificationNeedUploadItemRequest) validate() error {
	return nil
}

func (s goodsCertificationService) QueryNeedUploadItems(ctx context.Context, request GoodsCertificationNeedUploadItemRequest) (items []entity.GoodsCertificationNeedUploadItem, err error) {
	if err = request.validate(); err != nil {
		return
	}
	var result = struct {
		normal.Response
		Result struct {
			CertNeedUploadItems []entity.GoodsCertificationNeedUploadItem `json:"certNeedUploadItems"`
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(request).
		SetResult(&result).
		Post("bg.arbok.open.cert.queryNeedUploadItems")
	if err == nil {
		err = parseResponse(resp, result.Response)
	}
	if err != nil {
		return
	}

	items = result.Result.CertNeedUploadItems
	return
}
