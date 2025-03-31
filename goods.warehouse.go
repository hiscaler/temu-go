package temu

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
	"github.com/hiscaler/temu-go/validators/is"
)

// 商品仓库服务
type goodsWarehouseService service

type GoodsWarehouseOpenApiUser struct {
	SupplierId int `json:"supplierId"` // 供应商 ID
}

func (m GoodsWarehouseOpenApiUser) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.SupplierId, validation.Required.Error("供应商 ID 不能为空")),
	)
}

type GoodsWarehouseQueryParams struct {
	SiteIdList  []int                     `json:"siteIdList"`  // 站点列表
	OpenApiUser GoodsWarehouseOpenApiUser `json:"openApiUser"` // 用户信息
}

func (m GoodsWarehouseQueryParams) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.SiteIdList,
			validation.Required.Error("站点列表不能为空"),
			validation.By(is.SiteIds(entity.SiteIds)),
		),
		validation.Field(&m.OpenApiUser,
			validation.Required.Error("用户信息不能为空"),
		),
	)
}

// Query 根据站点查询可绑定的发货仓库信息接口（半托管专属）
// https://seller.kuajingmaihuo.com/sop/view/867739977041685428#hpIdAo
func (s goodsWarehouseService) Query(ctx context.Context, params GoodsWarehouseQueryParams) (items []entity.SiteWarehouses, err error) {
	if err = params.validate(); err != nil {
		return items, invalidInput(err)
	}

	var result = struct {
		normal.Response
		Result struct {
			WarehouseDTOList []entity.SiteWarehouses `json:"warehouseDTOList"` // 可选发货仓列表
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(params).
		SetResult(&result).
		Post("bg.goods.warehouse.list.get")
	if err = recheckError(resp, result.Response, err); err != nil {
		return items, err
	}

	return result.Result.WarehouseDTOList, nil
}
