package entity

// SemiOrderMoney 金额及币种，Amount 为币种的最小单位（美元为美分）。
type SemiOrderMoney struct {
	Currency string `json:"currency"` // 币种
	Amount   *int64 `json:"amount"`   // 币种最小单位的金额，可能为空
}

// SemiOrderAmountParent 父订单金额。
type SemiOrderAmountParent struct {
	ParentOrderSn             string          `json:"parentOrderSn"`             // 父订单号
	BasePriceTotal            *SemiOrderMoney `json:"basePriceTotal"`            // 基础价格总额
	ShippingAmountTotal       *SemiOrderMoney `json:"shippingAmountTotal"`       // 运费总额
	TaxTotal                  *SemiOrderMoney `json:"taxTotal"`                  // 消费者支付的税费总额
	DiscountFromSeller        *SemiOrderMoney `json:"discountFromSeller"`        // 卖家承担的优惠金额
	DiscountFromTEMU          *SemiOrderMoney `json:"discountFromTEMU"`          // Temu 承担的优惠金额
	RefundsTotal              *SemiOrderMoney `json:"refundsTotal"`              // 退款总额
	EstimatedRevenueDeduction *SemiOrderMoney `json:"estimatedRevenueDeduction"` // 预计收入扣减金额
	EstimatedRevenue          *SemiOrderMoney `json:"estimatedRevenue"`          // 预计收入
	RetailPriceTotal          *SemiOrderMoney `json:"retailPriceTotal"`          // 商品促销售价总额
	CustomerPaid              *SemiOrderMoney `json:"customerPaid"`              // 消费者支付总额
	TotalDiscount             *SemiOrderMoney `json:"totalDiscount"`             // 优惠总额
	SignOnDelivery            *SemiOrderMoney `json:"signOnDelivery"`            // 签收服务费，未购买时为空
}

// SemiOrderAmountItem 子订单金额。
type SemiOrderAmountItem struct {
	OrderSn            string          `json:"orderSn"`            // 子订单号
	Quantity           int             `json:"quantity"`           // 原始商品数量
	BasePrice          *SemiOrderMoney `json:"basePrice"`          // 基础价格
	UnitBasePrice      *SemiOrderMoney `json:"unitBasePrice"`      // 商品基础单价
	UnitRetailPrice    *SemiOrderMoney `json:"unitRetailPrice"`    // 商品零售单价
	DiscountFromSeller *SemiOrderMoney `json:"discountFromSeller"` // 卖家承担的商品优惠金额
	DiscountFromTEMU   *SemiOrderMoney `json:"discountFromTEMU"`   // Temu 承担的商品优惠金额
	ShipAmountTotal    *SemiOrderMoney `json:"shipAmountTotal"`    // 商品运费总额
}

// SemiOrderAmount 父订单和子订单金额查询结果。
type SemiOrderAmount struct {
	ParentOrderMap *SemiOrderAmountParent `json:"parentOrderMap"` // 父订单金额信息
	OrderList      []SemiOrderAmountItem  `json:"orderList"`      // 子订单金额列表
	Warning        []string               `json:"warning"`        // 接口返回的警告信息
}
