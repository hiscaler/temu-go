package temu

import (
	"context"

	"github.com/hiscaler/temu-go/normal"
)

type semiGoodsImageService service

type SemiGoodsImageRequest struct {
	FileUrl string `json:"fileUrl"`
	CatId   int64  `json:"catId"`
	Usage   int    `json:"usage"`
}

func (m SemiGoodsImageRequest) validate() error {
	return nil
}

// Upload 商品图片上传
// https://partner-us.temu.com/documentation?menu_code=fb16b05f7a904765aac4af3a24b87d4a&sub_menu_code=37d470d8c6e149f78953311aa0b0296d
func (s semiGoodsImageService) Upload(ctx context.Context, params SemiGoodsImageRequest) ([]string, error) {
	if err := params.validate(); err != nil {
		return nil, err
	}

	var result = struct {
		normal.Response
		Result struct {
			Images []struct {
				Url string `json:"url"`
			} `json:"images"`
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(params).
		SetResult(&result).
		Post("bg.local.goods.list.query")
	if err = recheckError(resp, result.Response, err); err != nil {
		return nil, err
	}

	urls := make([]string, len(result.Result.Images))
	for i, image := range result.Result.Images {
		urls[i] = image.Url
	}
	return urls, nil
}
