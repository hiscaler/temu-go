package temu

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

// 图片服务
type pictureService service

type PictureCompressionRequest struct {
	Urls []string `json:"urls"` // 需要压缩的图片链接
}

func (m PictureCompressionRequest) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Urls,
			validation.Required.Error("待压缩的图片链接地址不能为空"),
			validation.Each(is.URL.Error("图片链接格式不正确")),
		),
	)
}

// Compression 高清图片压缩处理
// https://partner.kuajingmaihuo.com/document?cataId=875198836203&docId=877312019388
func (s pictureService) Compression(ctx context.Context, params PictureCompressionRequest) ([]entity.PictureCompressionResult, error) {
	if err := params.validate(); err != nil {
		return nil, invalidInput(err)
	}

	var result struct {
		normal.Response
		Result struct {
			Results []entity.PictureCompressionResult `json:"results"`
		} `json:"result"`
	}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(params).
		SetResult(&result).
		Post("bg.picturecompression.get")
	if err = recheckError(resp, result.Response, err); err != nil {
		return nil, err
	}

	return result.Result.Results, nil
}
