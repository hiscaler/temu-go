package entity

type GoodsLifeCycle struct {
	SkcList []struct {
		SkcId        int64 `json:"skcId"`        // 货品 skcId
		SelectStatus int   `json:"selectStatus"` // 选品状态
		SkuList      struct {
			SkuId int64 `json:"skuId"`
		} `json:"skuList"` // sku 列表
		ApplyJitStatus  int  // 申诉 JIT 的状态(1：可申请、3：不可申请)
		SuggestCloseJit bool `json:"suggestCloseJit"` // 是否建议关闭 JIT 按钮
	} `json:"skcList"` // SKC 列表
	ProductId int64 `json:"productId"` // 货品 ID
}
