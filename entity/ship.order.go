package entity

import "gopkg.in/guregu/null.v4"

// ShipOrder 发货单
type ShipOrder struct {
	ReceiveSkcNum                int      `json:"receiveSkcNum"`                // 实收 skc 数目
	ExpressPackageNum            int      `json:"expressPackageNum"`            // 交接给快递公司的包裹数量
	LatestFeedbackStatus         int      `json:"latestFeedbackStatus"`         // 物流异常反馈最新反馈状态
	ExpectLatestPickTime         null.Int `json:"expectLatestPickTime"`         // 要求最晚揽收时间
	DeliveryOrderCancelLeftTime  null.Int `json:"deliveryOrderCancelLeftTime"`  // 发货单超时取消剩余时间,单位毫秒. 只针对非加急发货单,加急发货单(urgencyType=1)该字段为null
	ExpressDeliverySn            string   `json:"expressDeliverySn"`            // 快递单号
	DeliveryAddressId            null.Int `json:"deliveryAddressId"`            // 发货地址 id
	ExpressWeightFeedbackStatus  int      `json:"expressWeightFeedbackStatus"`  // 运单计费重量异常状态. 可选值含义说明:[0:未定义（数据库默认值）或无异常;1:异常待确认;2:已提交异常反馈，待物流商处理;3:物流商处理完成;4:平台介入处理中;5:平台处理完成;6:卖家已确认;7:卖家超期自动确认;8:物流商介入处理，卖家确认或超时自动确认;9:结算消息驱动卖家确认;10:无需公示;11:结算物流单计算重量查询失败;12:结算理论计费重拦截;13:SKU重量体积拦截;]
	ExpressRejectStatus          null.Int `json:"expressRejectStatus"`          // 物流单拒收状态 0-无拒收信息 1-存在拒收，待物流商处理 2-存在拒收，物流商已处理. 可选值含义说明:[0:无拒收信息;1:存在拒收，待物流商处理;2:存在拒收，物流商已处理;]
	PackageReceiveInfoVOList     any      `json:"packageReceiveInfoVOList"`     // 包裹收货信息（包裹收货时间）
	TaxWarehouseApplyOperateType int      `json:"taxWarehouseApplyOperateType"` // 入保税仓申请操作类型 0-不可操作 1-可申请 2-可查看
	ProductSkcId                 int64    `json:"productSkcId"`                 // skcId
	SkcExtCode                   string   `json:"skcExtCode"`                   // skc 货号信息
	InboundTime                  null.Int `json:"inboundTime"`                  // 发货单入库时间
	SubWarehouseId               int64    `json:"subWarehouseId"`               // 子仓 id
	PackageList                  []struct {
		SkcNum    int    `json:"skcNum"`    // skc数量
		PackageSn string `json:"packageSn"` // 包裹号
	} `json:"packageList"` //
	InventoryRegion             int            `json:"inventoryRegion"`             // 备货类型
	DeliverPackageNum           int            `json:"deliverPackageNum"`           // 实发包裹数
	DriverName                  string         `json:"driverName"`                  // 司机姓名
	SubPurchaseOrderSn          string         `json:"subPurchaseOrderSn"`          // 采购子单号
	ExpressCompanyId            int64          `json:"expressCompanyId"`            // 快递公司id
	DefectiveSkcNum             int            `json:"defectiveSkcNum"`             // 次品skc数目
	Status                      int            `json:"status"`                      // 状态
	ExpectPickUpGoodsTime       int64          `json:"expectPickUpGoodsTime"`       // 预约取货时间
	PredictTotalPackageWeight   int64          `json:"predictTotalPackageWeight"`   // 预估总包裹重量，单位g
	SupplierId                  int64          `json:"supplierId"`                  // 供应商id
	IsDisplayCourier            bool           `json:"isDisplayCourier"`            // 是否可以展示快递小哥联系方式，部分快递未接入
	IsCustomProduct             bool           `json:"isCustomProduct"`             // 是否为定制品 false-非定制品 true-定制品
	DeliveryMethod              int            `json:"deliveryMethod"`              // 发货方式. 可选值含义说明:[0:无;1:自送;2:公司指定物流;3:第三方物流;]
	ExpressWeightFeedbackTip    string         `json:"expressWeightFeedbackTip"`    // 运单计费重量异常提示文案 运单重量异常，待确认 || 物流商已回复重量异常，待确认. 可选值含义说明:[0:未定义（数据库默认值）或无异常;1:异常待确认;2:已提交异常反馈，待物流商处理;3:物流商处理完成;4:平台介入处理中;5:平台处理完成;6:卖家已确认;7:卖家超期自动确认;8:物流商介入处理，卖家确认或超时自动确认;9:结算消息驱动卖家确认;10:无需公示;11:结算物流单计算重量查询失败;12:结算理论计费重拦截;13:SKU重量体积拦截;]
	ExceptionFeedBackTotalCount null.Int       `json:"exceptionFeedBackTotalCount"` // 异常反馈总记录数
	OtherDeliveryPackageNum     int            `json:"otherDeliveryPackageNum"`     // 其他发货单的包裹数目
	PurchaseStockType           int            `json:"purchaseStockType"`           // 备货类型（0-普通备货 1-jit备货）
	IfCanOperateDeliver         bool           `json:"ifCanOperateDeliver"`         // 是否可以操作发货
	ReceivePackageNum           int            `json:"receivePackageNum"`           // 实收包裹数
	IsPrintBoxMark              bool           `json:"isPrintBoxMark"`              // 是否打印箱唛
	ExpressCompany              string         `json:"expressCompany"`              // 快递公司名称
	IsClothCategory             bool           `json:"isClothCategory"`             // 是否服饰类目
	DeliveryOrderSn             string         `json:"deliveryOrderSn"`             // 发货单号
	DeliverTime                 null.Int       `json:"deliverTime"`                 // 发货单发货时间
	UrgencyType                 int            `json:"urgencyType"`                 // 是否是紧急发货单（0-普通 1-急采）
	ExpressBatchSn              string         `json:"expressBatchSn"`              // 发货批次号
	ReceiveAddressInfo          ReceiveAddress `json:"receiveAddressInfo"`          // 收货仓详细地址
	PlateNumber                 string         `json:"plateNumber"`                 // 车牌号
	ReceiveTime                 null.Int       `json:"receiveTime"`                 // 发货单收货时间
	PackageDetailList           []struct {
		ProductSkuId         int64       `json:"productSkuId"`         // skuId
		ProductOriginalSkuId null.Int    `json:"productOriginalSkuId"` // 原 skuId
		PersonalText         null.String `json:"personalText"`         // 定制内容
		SkuNum               int         `json:"skuNum"`               // sku 数量
	} `json:"packageDetailList"` //	包裹详情列表
	SubPurchaseOrderBasicVO struct {
		SupplierId         int64       `json:"supplierId"`         // 供应商 id
		IsCustomProduct    bool        `json:"isCustomProduct"`    // 是否为定制品
		ProductSkcPicture  string      `json:"productSkcPicture"`  // 货品图片
		IsFirst            bool        `json:"isFirst"`            // 是否首单
		PurchaseStockType  int         `json:"purchaseStockType"`  // 备货类型（0-普通备货 1-jit备货）
		IsClothCategory    bool        `json:"isClothCategory"`    // 是否服饰类目
		ProductSkcId       int64       `json:"productSkcId"`       // skcId
		SettlementType     int         `json:"settlementType"`     // 结算类型（0-非vmi 1-vmi）
		SkcExtCode         string      `json:"skcExtCode"`         // 货号
		SubWarehouseId     null.Int    `json:"subWarehouseId"`     // 子仓 id
		UrgencyType        int         `json:"urgencyType"`        // 是否是紧急发货单(0-普通 1-急采)
		FragileTag         bool        `json:"fragileTag"`         // 易碎品打标
		PurchaseQuantity   int         `json:"purchaseQuantity"`   // 下单数量
		SubWarehouseName   null.String `json:"subWarehouseName"`   // 子仓名称
		PurchaseTime       int64       `json:"purchaseTime"`       // 下单时间（毫秒）
		SubPurchaseOrderSn string      `json:"subPurchaseOrderSn"` // 采购子单号
	} `json:"subPurchaseOrderBasicVO"` // 采购单信息
	SubWarehouseName        string   `json:"subWarehouseName"`        // 子仓名称
	PurchaseTime            int64    `json:"purchaseTime"`            // 下单时间（时间戳：毫秒）
	SkcPurchaseNum          int      `json:"skcPurchaseNum"`          // 下单数量
	DeliverSkcNum           int      `json:"deliverSkcNum"`           // 实发skc数目
	DeliveryOrderCreateTime int64    `json:"deliveryOrderCreateTime"` // 发货单创建时间
	OrderType               null.Int `json:"orderType"`               // 自处理后的发货单类型（1：普通、2：紧急、3：定制）
}
