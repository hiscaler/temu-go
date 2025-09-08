package temu

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

// 半托管物流扫描单服务
type semiOrderLogisticsScanFormService service

type SemiOrderLogisticsScanFormCreateRequest struct {
	PackageSnList []string `json:"packageSnList"` // 包裹号列表
	ShipCompanyId int64    `json:"shipCompanyId"` // 发货物流公司 ID
	WarehouseId   string   `json:"warehouseId"`   // 仓库 ID
}

func (m SemiOrderLogisticsScanFormCreateRequest) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.PackageSnList,
			validation.Required.Error("扫描包裹号列表不能为空"),
		),
		validation.Field(&m.ShipCompanyId,
			validation.Required.Error("发货物流公司 ID 不能为空"),
		),
		validation.Field(&m.WarehouseId,
			validation.Required.Error("仓库 ID 不能为空"),
		),
	)
}

// Create 创建扫描单
// https://partner-us.temu.com/documentation?menu_code=fb16b05f7a904765aac4af3a24b87d4a&sub_menu_code=0ef1cf008e144cbb987771ae3a8fd99d
func (s *semiOrderLogisticsScanFormService) Create(ctx context.Context, request SemiOrderLogisticsScanFormCreateRequest) ([]entity.SemiOrderLogisticsScanFormResult, error) {
	if err := request.validate(); err != nil {
		return nil, err
	}

	var result = struct {
		normal.Response
		Result struct {
			ScanFormInfoList []entity.SemiOrderLogisticsScanFormResult `json:"scanFormInfoList"` // 扫描信息列表
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(request).
		SetResult(&result).
		Post("temu.logistics.scanform.create")

	if err = recheckError(resp, result.Response, err); err != nil {
		return nil, err
	}

	return result.Result.ScanFormInfoList, nil
}

// Document 获取扫描单文件
func (s *semiOrderLogisticsScanFormService) Document(ctx context.Context, scanFormSn string) (entity.File, error) {
	var result = struct {
		normal.Response
		Result struct {
			Url entity.SignatureUrl `json:"url"` // 面单文件
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(map[string]string{"scanFormSn": scanFormSn}).
		SetResult(&result).
		Post("temu.logistics.scanform.document.get")

	if err = recheckError(resp, result.Response, err); err != nil {
		return entity.File{}, err
	}
	return result.Result.Url.Decode(s.config)
}
