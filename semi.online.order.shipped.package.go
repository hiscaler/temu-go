package temu

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

// 已发货包裹服务
type semiOnlineOrderShippedPackageService service

type SemiOnlineOrderPlatformLogisticsShippedPackageRequest struct {
	PackageSendInfoList []struct {
		PackageSn      string `json:"packageSn"`      // 包裹号
		TrackingNumber string `json:"trackingNumber"` // 运单号
		PackageDetail  []struct {
			ParentOrderSn string `json:"parentOrderSn"` // 父单号
			OrderSn       string `json:"orderSn"`       // 子单号
			Quantity      int    `json:"quantity"`      // 发货件数
		} `json:"packageDetail"` // 发货包裹详情
	} `json:"packageSendInfoList"` // 确认发货包裹列表
}

func (m SemiOnlineOrderPlatformLogisticsShippedPackageRequest) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.PackageSendInfoList, validation.Required.Error("确认发货包裹列表不能为空")),
		// todo 更多的数据验证
	)
}

// Confirm 确认包裹发货接口（bg.logistics.shipped.package.confirm）
// https://seller.kuajingmaihuo.com/sop/view/144659541206936016#92SpUJ
func (s semiOnlineOrderShippedPackageService) Confirm(ctx context.Context, params SemiOnlineOrderPlatformLogisticsShippedPackageRequest) (items []entity.SemiOnlineOrderShippedPackageConfirmResult, err error) {
	if err = params.validate(); err != nil {
		return items, invalidInput(err)
	}

	var result = struct {
		normal.Response
		Result struct {
			WarningMessage []entity.SemiOnlineOrderShippedPackageConfirmResult `json:"warningMessage"` // 提醒信息列表
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(params).
		SetResult(&result).
		Post("bg.logistics.shipped.package.confirm")
	if err = recheckError(resp, result.Response, err); err != nil {
		return
	}

	return result.Result.WarningMessage, nil
}
