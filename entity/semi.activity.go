package entity

type SemiActivity struct {
	ActivityId        int    `json:"activityId"`
	ActivityEndTime   int    `json:"activityEndTime"`
	ActivityName      string `json:"activityName"`
	ActivityStatus    int    `json:"activityStatus"`
	ActivityStartTime int    `json:"activityStartTime"`
	ActivityType      int    `json:"activityType"`
	IsJoinedActivity  bool   `json:"isJoinedActivity"`
}

// SemiActivityEnrollResult 半托活动报名结果
type SemiActivityEnrollResult struct {
	DraftId         int    `json:"draftId"`
	OperationStatus int    `json:"operationStatus"`
	FailReason      string `json:"failReason"`
	GoodsId         int    `json:"goodsId"`
}

// SemiActivityCandidateResult 半托活动候选商品查询结果
type SemiActivityCandidateResult struct {
	GoodsList    []SemiActivityCandidateGoods `json:"goodsList"`
	ActivityInfo SemiActivity                 `json:"activityInfo"`
	Total        int                          `json:"total"`
}

// SemiActivityCandidateGoods 半托活动候选商品
type SemiActivityCandidateGoods struct {
	RecommendActivityQuantity int                        `json:"recommendActivityQuantity"`
	SkuList                   []SemiActivityCandidateSku `json:"skuList"`
	IsAllowedEnroll           bool                       `json:"isAllowedEnroll"`
	GoodsId                   int64                      `json:"goodsId"`
}

// SemiActivityCandidateSku 半托活动候选 SKU
type SemiActivityCandidateSku struct {
	RecommendActivitySupplierPrice float64 `json:"recommendActivitySupplierPrice"`
	SkuId                          int64   `json:"skuId"`
}
