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
