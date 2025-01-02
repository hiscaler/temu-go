package entity

type OrderResult struct {
	TotalItemNum int        `json:"totalItemNum"`
	PageItems    []PageItem `json:"pageItems"`
}

type PageItem struct {
	ParentOrderMap ParentOrderMap `json:"parentOrderMap"`
	OrderList      []Order        `json:"orderList"`
}

type ParentOrderMap struct {
	ParentOrderLabel             []Label  `json:"parentOrderLabel"`
	ParentOrderSn                string   `json:"parentOrderSn"`
	ParentOrderStatus            int      `json:"parentOrderStatus"`
	ParentOrderTime              int      `json:"parentOrderTime"`
	ParentOrderPendingFinishTime int      `json:"parentOrderPendingFinishTime"`
	ExpectShipLatestTime         int      `json:"expectShipLatestTime"`
	ParentShippingTime           int      `json:"parentShippingTime"`
	FulfillmentWarning           []string `json:"fulfillmentWarning"`
}

type Label struct {
	Name  string `json:"nameupdate"`
	Value int    `json:"valueupdate"`
}

type Order struct {
	OrderSn                         string    `json:"orderSn"`
	Quantity                        int       `json:"quantity"`
	OriginalOrderQuantity           int       `json:"originalOrderQuantity"`
	CanceledQuantityBeforeShipment  int       `json:"canceledQuantityBeforeShipment"`
	InventoryDeductionWarehouseId   string    `json:"inventoryDeductionWarehouseId"`
	InventoryDeductionWarehouseName string    `json:"inventoryDeductionWarehouseName"`
	OrderLabel                      []Label   `json:"orderLabel"`
	Spec                            string    `json:"spec"`
	ThumbUrl                        string    `json:"thumbUrl"`
	OrderStatus                     int       `json:"orderStatus"`
	FulfillmentType                 string    `json:"fulfillmentType"`
	GoodsName                       string    `json:"goodsName"`
	ProductList                     []Product `json:"productList"`
	RegionId                        int       `json:"regionId"`
	SiteId                          int       `json:"siteId"`
}

type Product struct {
	ProductSkuId int    `json:"productSkuId"`
	SoldFactor   int    `json:"soldFactor"`
	ProductId    int64  `json:"productId"`
	ExtCode      string `json:"extCode"`
}
