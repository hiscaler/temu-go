package temu

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

// 商品管理数据服务
type goodsSalesService service

type GoodsSalesQueryParams struct {
	normal.ParameterWithPager
	IsLack                    int      `json:"isLack,omitempty"`                    // 是否缺货 0-不缺货 1-缺货
	ProductSkcIdList          []int    `json:"productSkcIdList,omitempty"`          // skc列表
	MaxRemanentInventoryNum   int      `json:"maxRemanentInventoryNum,omitempty"`   // sku最大剩余库存
	OnSalesDurationOfflineLte int      `json:"onSalesDurationOfflineLte,omitempty"` // 加入站点时长小于等于
	MinRemanentInventoryNum   int      `json:"minRemanentInventoryNum,omitempty"`   // sku最小剩余库存
	SelectStatusList          []int    `json:"selectStatusList,omitempty"`          // 选品状态 10-待下首单 11-已下首单 12-已加入站点 13-已下架
	TodaySaleVolumMax         int      `json:"todaySaleVolumMax,omitempty"`         // SKC今日销量最大值
	MaxAvailableSaleDays      int      `json:"maxAvailableSaleDays,omitempty"`      // 最大可售天数
	OnSalesDurationOfflineGte int      `json:"onSalesDurationOfflineGte,omitempty"` // 加入站点时长大于等于
	SkuExtCodeList            []string `json:"skuExtCodeList,omitempty"`            // sku货号列表
	ProductName               string   `json:"productName,omitempty"`               // 货品名称
	ThirtyDaysSaleVolumMax    int      `json:"thirtyDaysSaleVolumMax,omitempty"`    // SKC近30天销量最大值
	ThirtyDaysSaleVolumMin    int      `json:"thirtyDaysSaleVolumMin,omitempty"`    // SKC近30天销量最小值
	CategoryList              []int    `json:"categoryList,omitempty"`              // 类目列表
	IsTrustManagementMall     bool     `json:"isTrustManagementMall,omitempty"`     // 是否托管店铺
	SevenDaysSaleVolumMax     int      `json:"sevenDaysSaleVolumMax,omitempty"`     // SKC近7天销量最大值
	SettlementType            int      `json:"settlementType,omitempty"`            // 结算类型 0-非vmi 1-vmi
	StockStatusList           []int    `json:"stockStatusList,omitempty"`           // 售罄状态 (0-库存充足 1-即将售罄 2-已经售罄)
	SkcExtCodeList            []int    `json:"skcExtCodeList,omitempty"`            // skc货号列表
	TodaySaleVolumMin         int      `json:"todaySaleVolumMin,omitempty"`         // SKC今日销量最小值
	SevenDaysSaleVolumMin     int      `json:"sevenDaysSaleVolumMin,omitempty"`     // SKC近7天销量最小值
	OrderByDesc               int      `json:"orderByDesc,omitempty"`               // 排序，0-升序，1-降序
	IsAdviceStock             bool     `json:"isAdviceStock,omitempty"`             // 是否建议备货
	PictureAuditStatusList    []int    `json:"pictureAuditStatusList,omitempty"`    // 图片审核状态 1-未完成；2-已完成
	IsCustomGoods             bool     `json:"isCustomGoods,omitempty"`             // 是否为定制品
	OrderByParam              string   `json:"orderByParam,omitempty"`              // 排序参数，传入后端返回的字段
	MinAvailableSaleDays      string   `json:"minAvailableSaleDays,omitempty"`      // 最小可售天数
}

func (m GoodsSalesQueryParams) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Page, validation.Required.Error("页码不能为空。")),
		validation.Field(&m.PageSize, validation.Required.Error("页数不能为空。")),
	)
}

// All 销售管理数据查询接口
// https://seller.kuajingmaihuo.com/sop/view/889973754324016047#8sCZQ8
func (s goodsSalesService) All(params GoodsSalesQueryParams) (items []entity.GoodsSales, err error) {
	params.TidyPager()
	if err = params.Validate(); err != nil {
		return
	}
	var result = struct {
		normal.Response
		Result struct {
			Total        int                 `json:"total"`        // 总数
			SubOrderList []entity.GoodsSales `json:"subOrderList"` // 订单信息
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetBody(params).
		SetResult(&result).
		Post("bg.goods.sales.get")
	if err == nil {
		err = parseResponse(resp, result.Response)
	}
	if err != nil {
		return
	}

	return result.Result.SubOrderList, nil
}

// One 根据商品 SKC ID 查询
func (s goodsSalesService) One(productSkcId int) (item entity.GoodsSales, err error) {
	items, err := s.All(GoodsSalesQueryParams{ProductSkcIdList: []int{productSkcId}})
	if err != nil {
		return
	}

	for _, sales := range items {
		if sales.ProductSkcID == productSkcId {
			return sales, nil
		}
	}

	return item, ErrNotFound
}
