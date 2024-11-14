package temu

import (
	"context"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

type goodsCertificationService service

type GoodsCertificationQueryRequest struct {
	CertTypeList []int  `json:"certTypeList,omitempty"` // 资质类型 ID 列表
	ProductId    int64  `json:"productId"`              // 货品 ID
	Language     string `json:"language,omitempty"`     // 语言
}

func (m GoodsCertificationQueryRequest) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.CertTypeList,
			validation.When(len(m.CertTypeList) != 0, validation.By(func(value interface{}) error {
				types, ok := value.([]int)
				if !ok {
					return errors.New("无效的资质类型 ID。")
				}
				for _, typ := range types {
					if typ < 0 || typ > 303 {
						return fmt.Errorf("无效的资质类型 ID: %d。", typ)
					}
				}
				return nil
			}))),
		validation.Field(&m.ProductId, validation.Required.Error("货品 ID 不能为空。")),
	)
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
