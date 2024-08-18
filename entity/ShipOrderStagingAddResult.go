package entity

import "gopkg.in/guregu/null.v4"

// ShipOrderStagingAddResult 加入发货台结果
type ShipOrderStagingAddResult struct {
	SubPurchaseOrderSn string      `json:"purchase_order_sn"` // 采购子单号
	Success            bool        `json:"success"`           // 是否成功
	ErrorCode          null.Int    `json:"error_code"`        // 错误代码
	ErrorMessage       null.String `json:"error_message"`     // 错误消息
}
