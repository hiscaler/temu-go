package temu

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

// 物流服务

type logisticsService service

// Companies 查询发货快递公司
// https://seller.kuajingmaihuo.com/sop/view/889973754324016047#wjtGTK
func (s logisticsService) Companies() (items []entity.LogisticsCompany, err error) {
	var result = struct {
		normal.Response
		Result struct {
			ShipList []entity.LogisticsCompany `json:"shipList"` // 快递公司列表
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetResult(&result).
		Post("bg.logistics.company.get")
	if err == nil {
		err = parseResponse(resp, result.Response)
	}
	if err != nil {
		return
	}

	return result.Result.ShipList, nil
}

// 平台推荐物流商匹配接口
// https://seller.kuajingmaihuo.com/sop/view/889973754324016047#16WiXI

type LogisticsMatchRequest struct {
	DeliveryAddressId         int  `json:"deliveryAddressId"`         // 发货地址
	PredictTotalPackageWeight int  `json:"predictTotalPackageWeight"` // 预估总包裹重量，单位g
	UrgencyType               int  `json:"urgencyType"`               // 是否是紧急发货单，0-普通 1-急采
	SubWarehouseId            int  `json:"subWarehouseId"`            // 收货子仓id
	QueryStandbyExpress       bool `json:"queryStandbyExpress"`       // 是否查询备用快递服务商, false-不查询 true-查询
	TotalPackageNum           int  `json:"totalPackageNum"`           // 包裹件数
	ReceiveAddressInfo        struct {
		DistrictCode  int    `json:"districtCode"`
		CityName      string `json:"cityName"`
		DistrictName  string `json:"districtName"`
		Phone         string `json:"phone"`
		ProvinceCode  int    `json:"provinceCode"`
		CityCode      int    `json:"cityCode"`
		ReceiverName  string `json:"receiverName"`
		DetailAddress string `json:"detailAddress"`
		ProvinceName  string `json:"provinceName"`
	} `json:"receiveAddressInfo"` // 收货地址
	DeliveryOrderSns []string `json:"deliveryOrderSns"`
}

func (m LogisticsMatchRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.PredictTotalPackageWeight, validation.Required.Error("预估总包裹重量不能为空。")),
		validation.Field(&m.TotalPackageNum, validation.Required.Error("包裹件数不能为空。")),
	)
}

func (s logisticsService) Match(request LogisticsMatchRequest) (items []entity.LogisticsMatch, err error) {
	if err = request.Validate(); err != nil {
		return
	}

	var result = struct {
		normal.Response
		Result []entity.LogisticsMatch `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetBody(request).
		SetResult(&result).
		Post("bg.shiporderv2.logisticsmatch.get")
	if err == nil {
		err = parseResponse(resp, result.Response)
	}
	if err != nil {
		return
	}

	return result.Result, nil
}
