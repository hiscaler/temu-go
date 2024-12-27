package entity

// OrderLabel 订单标签
type OrderLabel struct {
	Name  string `json:"name"`  // 标签名称
	Value string `json:"value"` // 是否有标签（0=无标签，1=有标签）
}

// ParentOrder 父订单
type ParentOrder struct {
	// 标签名称具体枚举如下
	// soon_to_be_overdue-即将逾期、past_due-已逾期、pending_buyer_cancellation-买家取消待确认订单、pending_buyer_address_change-买家改地址待确认订单
	ParentOrderLabel             OrderLabel `json:"parentOrderLabel"`             // PO 单订单状态标签
	ParentOrderSn                string     `json:"parentOrderSn"`                // 订单号
	ParentOrderStatus            int        `json:"parentOrderStatus"`            // 订单状态
	ParentOrderTime              int64      `json:"parentOrderTime"`              // 订单创建时间
	ParentOrderPendingFinishTime int64      `json:"parentOrderPendingFinishTime"` // 订单结束pending转为自发货时间
	ExpectShipLatestTime         int64      `json:"expectShipLatestTime"`         // 要求最晚发货时间
	ParentShippingTime           int64      `json:"parentShippingTime"`           // 父单发货时间
	FulfillmentWarning           []string   `json:"fulfillmentWarning"`           // 履约相关提醒: SUGGEST_SIGNATURE_ON_DELIVERY-建议发货时购买签名服务
}

// ChildOrder 子订单
type ChildOrder struct {
	OrderSn string `json:"orderSn"` // 订单号
	// 单应履约件数
	// 备注：代表商家实际需要发货件数，在订单部分取消时：应履约件数=下单件数-发货前售后件数
	Quantity                          int    `json:"quantity"`                          // 单应履约件数
	OriginalOrderQuantityNew          int    `json:"originalOrderQuantityNew"`          // 用户初始下单时的 O 单件数
	CanceledQuantityBeforeShipmentNew int    `json:"canceledQuantityBeforeShipmentNew"` // O 单发货前，用户发起部分取消的件数（用户申请且退款已受理）
	InventoryDeductionWarehouseId     string `json:"inventoryDeductionWarehouseId"`     // 库存扣减仓库 id
	InventoryDeductionWarehouseName   string `json:"inventoryDeductionWarehouseName"`   // 库存扣减仓库名称
	// 标签名称，具体枚举如下 customized_products：定制品标签	US_to_CA：美发加订单标签
	OrderLabel      OrderLabel `json:"orderLabel"`      // 子订单O单标签，内部请求异常返回为空，返回为空时请重试
	GoodsId         int64      `json:"goodsId"`         // 商品 id
	GoodsName       string     `json:"goodsName"`       // 商品名称
	SkuId           int64      `json:"skuId"`           // SKU id
	Spec            string     `json:"spec"`            // 商品信息描述
	ThumbUrl        string     `json:"thumbUrl"`        // 商品缩略图图片
	OrderStatus     int        `json:"orderStatus"`     // 订单状态，3 是已取消
	FulfillmentType string     `json:"fulfillmentType"` // 子订单履约类型 - 卖家履约订单值返回：fulfillBySeller 	- 合作仓履约订单返回：fulfillByCooperativeWarehouse
	ProductList     []struct {
		ProductSkuId int64  `json:"productSkuId"` // 货品 skuId
		SoldFactor   int    `json:"soldFactor"`   // 商品和货品数量转换系数，商品数量(quantity)乘以转换系数，代表货品数量
		ProductId    int64  `json:"productId"`    // 货品 Id
		ExtCode      string `json:"extCode"`      // 货品编码
	} `json:"productList"` // 货品信息
	RegionId int `json:"regionId"` // 区域 ID
	SiteId   int `json:"siteId"`   // 站点 ID
}

// Order 订单
type Order struct {
	ParentOrder ParentOrder  `json:"parent_order"`
	Items       []ChildOrder `json:"items"`
}
