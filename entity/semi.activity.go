package entity

// SemiActivity 半托活动
type SemiActivity struct {
	ActivityId        int    `json:"activityId"`        // 活动唯一标识
	ActivityEndTime   int    `json:"activityEndTime"`   // 活动结束时间，Unix 秒级时间戳
	ActivityName      string `json:"activityName"`      // 活动名称
	ActivityStatus    int    `json:"activityStatus"`    // 活动状态：1 未开始；2 进行中；3 已结束
	ActivityStartTime int    `json:"activityStartTime"` // 活动开始时间，Unix 秒级时间戳
	ActivityType      int    `json:"activityType"`      // 活动类型：2 闪购；13 高级大促；27 清仓；100 官方大促
	IsJoinedActivity  bool   `json:"isJoinedActivity"`  // 是否已参加该活动
}

// SemiActivityEnrollResult 半托活动报名结果
type SemiActivityEnrollResult struct {
	DraftId         int    `json:"draftId"`         // 商品报名活动草稿 ID；预处理成功时生成
	OperationStatus int    `json:"operationStatus"` // 操作状态：10 处理中；20 成功；30 失败
	FailReason      string `json:"failReason"`      // 操作失败原因
	GoodsId         int    `json:"goodsId"`         // 商品 ID
}

// SemiActivityCandidateResult 半托活动候选商品查询结果
type SemiActivityCandidateResult struct {
	GoodsList    []SemiActivityCandidateGoods `json:"goodsList"`    // 可报名商品列表
	ActivityInfo SemiActivity                 `json:"activityInfo"` // 活动信息
	Total        int                          `json:"total"`        // 可报名商品总数
}

// SemiActivityCandidateGoods 半托活动候选商品
type SemiActivityCandidateGoods struct {
	RecommendActivityQuantity int                        `json:"recommendActivityQuantity"` // 建议的活动报名数量
	SkuList                   []SemiActivityCandidateSku `json:"skuList"`                   // 可报名 SKU 列表及建议活动供货价
	IsAllowedEnroll           bool                       `json:"isAllowedEnroll"`           // 是否允许报名
	GoodsId                   int64                      `json:"goodsId"`                   // 商品 ID
}

// SemiActivityCandidateSku 半托活动候选 SKU
type SemiActivityCandidateSku struct {
	RecommendActivitySupplierPrice float64 `json:"recommendActivitySupplierPrice"` // 建议活动供货价
	SkuId                          int64   `json:"skuId"`                          // SKU ID
}
