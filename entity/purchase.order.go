package entity

import "gopkg.in/guregu/null.v4"

// PurchaseOrder 采购单
type PurchaseOrder struct {
	Total                   int            `json:"total"`
	StatusNumMap            map[string]int `json:"statusNumMap"` // key:采购单状态 value:该状态下数量。 订单状态 0-待接单；1-已接单，待发货；2-已送货；3-已收货；5-已验收，全部退回；6-已验收；7-已入库；8-作废；9-已超时； 10-已取消
	TodayCanDeliverNum      int            `json:"todayCanDeliverNum"`
	DelayNumMap             map[string]int `json:"delayNumMap"`
	SubOrderForSupplierList []struct {
		OriginalPurchaseOrderSn string      `json:"originalPurchaseOrderSn"`
		Source                  int         `json:"source"`
		ProductName             string      `json:"productName"`
		FulfilmentFormStatus    interface{} `json:"fulfilmentFormStatus"`
		IsFirst                 bool        `json:"isFirst"`
		SkuQuantityDetailList   []struct {
			CurrencyType                 string      `json:"currencyType"`
			ClassName                    string      `json:"className"`
			RealReceiveAuthenticQuantity int         `json:"realReceiveAuthenticQuantity"`
			FulfilmentProductSkuId       int         `json:"fulfilmentProductSkuId"`
			CustomizationType            int         `json:"customizationType"`
			ProductSkuId                 int         `json:"productSkuId"`
			DeliverQuantity              int         `json:"deliverQuantity"`
			ThumbUrlList                 []string    `json:"thumbUrlList"`
			QcResult                     interface{} `json:"qcResult"`
			ExtCode                      string      `json:"extCode"`
			PurchaseQuantity             int         `json:"purchaseQuantity"`
		} `json:"skuQuantityDetailList"`
		DeliverInfo struct {
			ReceiveTime                      null.Int    `json:"receiveTime"`
			DeliverTime                      null.Int    `json:"deliverTime"`
			ReceiveWarehouseId               null.Int    `json:"receiveWarehouseId"`
			ReceiveWarehouseName             null.String `json:"receiveWarehouseName"`
			ExpectLatestDeliverTimeOrDefault int         `json:"expectLatestDeliverTimeOrDefault"`
			ExpectLatestArrivalTimeOrDefault int         `json:"expectLatestArrivalTimeOrDefault"`
			DeliveryOrderSn                  null.String `json:"deliveryOrderSn"`
		} `json:"deliverInfo"`
		ProductSkcId         int `json:"productSkcId"`
		ProductId            int `json:"productId"`
		HasQcBill            int `json:"hasQcBill"`
		SupplyStatus         int `json:"supplyStatus"`
		ApplyDeleteStatus    int `json:"applyDeleteStatus"`
		SkuQuantityTotalInfo struct {
			CurrencyType                 interface{} `json:"currencyType"`
			ClassName                    interface{} `json:"className"`
			RealReceiveAuthenticQuantity int         `json:"realReceiveAuthenticQuantity"`
			CustomizationType            interface{} `json:"customizationType"`
			ProductSkuId                 interface{} `json:"productSkuId"`
			DeliverQuantity              int         `json:"deliverQuantity"`
			ExtCode                      interface{} `json:"extCode"`
			PurchaseQuantity             int         `json:"purchaseQuantity"`
		} `json:"skuQuantityTotalInfo"`
		IsCanJoinDeliverPlatform bool        `json:"isCanJoinDeliverPlatform"`
		CategoryType             int         `json:"categoryType"`
		SubPurchaseOrderSn       string      `json:"subPurchaseOrderSn"`
		Status                   int         `json:"status"`
		FulfilmentFormId         interface{} `json:"fulfilmentFormId"`
		ProductSkcPicture        string      `json:"productSkcPicture"`
		LackOrSoldOutTagList     []struct {
			IsLack     bool        `json:"isLack"`
			SkuDisplay string      `json:"skuDisplay"`
			SoldOut    interface{} `json:"soldOut"`
		} `json:"lackOrSoldOutTagList"`
		QcReject          int `json:"qcReject"`
		PurchaseStockType int `json:"purchaseStockType"`
		SkuLackItemList   []struct {
			SkuDisplay string `json:"skuDisplay"`
		} `json:"skuLackItemList"`
		SkuLackSnapshot         int         `json:"skuLackSnapshot"`
		DeliveryOrderSn         *string     `json:"deliveryOrderSn"`
		SettlementType          int         `json:"settlementType"`
		UrgencyType             int         `json:"urgencyType"`
		ProductSn               *string     `json:"productSn"`
		SkuQcRejectItemList     interface{} `json:"skuQcRejectItemList"`
		DefectiveTime           interface{} `json:"defectiveTime"`
		TodayCanDeliver         bool        `json:"todayCanDeliver"`
		PurchaseTime            int         `json:"purchaseTime"`
		ApplyChangeSupplyStatus int         `json:"applyChangeSupplyStatus"`
		Category                string      `json:"category"`
	} `json:"subOrderForSupplierList"`
}
