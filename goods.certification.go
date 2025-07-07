package temu

import (
	"context"
	"errors"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
	"gopkg.in/guregu/null.v4"
)

// 商品资质服务
type goodsCertificationService service

type GoodsCertificationQueryParams struct {
	CertTypeList []int  `json:"certTypeList,omitempty"` // 资质类型 ID 列表
	ProductId    int64  `json:"productId"`              // 货品 ID
	Language     string `json:"language,omitempty"`     // 语言
}

func (m GoodsCertificationQueryParams) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.CertTypeList,
			validation.When(len(m.CertTypeList) != 0, validation.By(func(value interface{}) error {
				types, ok := value.([]int)
				if !ok {
					return errors.New("无效的资质类型")
				}
				for _, typ := range types {
					if typ < 0 || typ > 303 {
						return fmt.Errorf("无效的资质类型 %d", typ)
					}
				}
				return nil
			}))),
		validation.Field(&m.ProductId, validation.Required.Error("货品 ID 不能为空")),
	)
}

// Query 批量查询货品资质信息
// https://seller.kuajingmaihuo.com/sop/view/649320516224723675#Oq8dC9
func (s goodsCertificationService) Query(ctx context.Context, params GoodsCertificationQueryParams) (certifications []entity.GoodsCertification, err error) {
	if err = params.validate(); err != nil {
		return certifications, invalidInput(err)
	}

	var result = struct {
		normal.Response
		Result struct {
			ProductCertList []entity.GoodsCertification `json:"productCertList"`
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(params).
		SetResult(&result).
		Post("bg.arbok.open.product.cert.query")
	if err = recheckError(resp, result.Response, err); err != nil {
		return
	}

	certifications = result.Result.ProductCertList
	return
}

// 查询资质要上传的内容
// https://seller.kuajingmaihuo.com/sop/view/649320516224723675#5mZ1dI

type GoodsCertificationNeedUploadItemRequest struct {
	CertType  int   `json:"certType"`  // 资质类型
	ProductId int64 `json:"productId"` // 货品 ID
}

func (m GoodsCertificationNeedUploadItemRequest) validate() error {
	return nil
}

func (s goodsCertificationService) QueryNeedUploadItems(ctx context.Context, request GoodsCertificationNeedUploadItemRequest) (items []entity.GoodsCertificationNeedUploadItem, err error) {
	if err = request.validate(); err != nil {
		return items, invalidInput(err)
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
	if err = recheckError(resp, result.Response, err); err != nil {
		return
	}

	return result.Result.CertNeedUploadItems, nil
}

// 上传文件接口（bg.arbok.open.upload.uploadFile）
// https://seller.kuajingmaihuo.com/sop/view/649320516224723675#sFvgAq

type GoodsCertificationUploadFileRequest struct {
	Base64File string `json:"base64File"` // 支持格式有：jpg/jpeg、png,pdf格式，注意入参图片必须转码为base64编码
	FileName   string `json:"fileName"`   // 文件名，主要用来辨别格式，名字不采用，仅支持jpg/jpeg、png,pdf格式
	IsRealPic  bool   `json:"isRealPic"`  // 是否实拍图，是的话结果连接不需要签名就可以展示，不是的话就需要签名才能展示
}

func (m GoodsCertificationUploadFileRequest) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Base64File, validation.Required.Error("上传文件内容不能为空"), is.Base64.Error("无效的 Base64 内容")),
		validation.Field(&m.FileName, validation.Required.Error("上传文件名称不能为空")),
	)
}

func (s goodsCertificationService) UploadFile(ctx context.Context, request GoodsCertificationUploadFileRequest) (string, error) {
	if err := request.validate(); err != nil {
		return "", invalidInput(err)
	}

	var result = struct {
		normal.Response
		Result string `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(request).
		SetResult(&result).
		Post("bg.arbok.open.upload.uploadFile")
	if err = recheckError(resp, result.Response, err); err != nil {
		return "", err
	}

	return result.Result, nil
}

// 提交资质内容
// https://seller.kuajingmaihuo.com/sop/view/649320516224723675#dCjnye

type GoodsCertificationUploadFile struct {
	FileName      string   `json:"fileName"`                // 文件名称
	FileUrl       string   `json:"fileUrl"`                 // 文件 Url
	ShowCustomer  bool     `json:"showCustomer"`            // 是否同意向消费者展示
	ExpireTime    null.Int `json:"expireTime,omitempty"`    // 失效时间
	EffectiveTime null.Int `json:"effectiveTime,omitempty"` // 生效时间
}

func (m GoodsCertificationUploadFile) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.FileName, validation.Required.Error("文件名称不能为空")),
		validation.Field(&m.FileUrl, validation.Required.Error("文件 URL 不能为空")),
	)
}

type GoodsCertificationUploadRequest struct {
	ProductId             int64 `json:"productId"` // 货品 ID
	ProductCatCertReqList []struct {
		CertType           int                            `json:"certType"`           // 资质类型
		AuthCode           string                         `json:"authCode"`           // 资质编号
		ProductCertFiles   []GoodsCertificationUploadFile `json:"productCertFiles"`   // 货品资质文件列表
		InspectReportFiles []GoodsCertificationUploadFile `json:"inspectReportFiles"` // 检测报告文件列表
		RealPictures       []struct {
			ImageUrl string `json:"imageUrl"` // 图片 Url
		} `json:"realPictures"` // 实物图列表
	} `json:"productCatCertReqList"` // 货品类目资质请求列表，仅传需上传的资质类型，审核通过/审核中的勿传
}

func (m GoodsCertificationUploadRequest) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.ProductId, validation.Required.Error("货品 ID 不能为空")),
		validation.Field(&m.ProductCatCertReqList,
			validation.Required.Error("货品类目资质请求列表不能为空"),
		),
	)
}

func (s goodsCertificationService) Upload(ctx context.Context, request GoodsCertificationUploadRequest) (bool, error) {
	if err := request.validate(); err != nil {
		return false, invalidInput(err)
	}

	var result = struct {
		normal.Response
		Result any `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(request).
		SetResult(&result).
		Post("bg.arbok.open.cert.uploadProductCert")
	if err = recheckError(resp, result.Response, err); err != nil {
		return false, err
	}

	return true, nil
}
