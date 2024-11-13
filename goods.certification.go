package temu

import "github.com/hiscaler/temu-go/entity"

type goodsCertificationService service

type GoodsCertificationQueryRequest struct {
	CertTypeList []int  `json:"certTypeList"` // 资质类型 id 列表
	ProductId    int64  `json:"productId"`    // 货品 id
	Language     string `json:"language"`     // 语言
}

func (s goodsCertificationService) Query(request GoodsCertificationQueryRequest) (certification []entity.GoodsCertification, err error) {
	return
}
