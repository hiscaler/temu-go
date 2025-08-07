package entity

// GoodsLifeCycle 货品生命周期状态
type GoodsLifeCycle struct {
	ProductId int64 `json:"productId"` // 货品 ID
	SkcList   []struct {
		SkcId int64 `json:"skcId"` // 货品 skcId
		// 选品状态
		// 0："已弃用"
		// 1："待平台选品"
		// 2："待上传生产资料"
		// 3："待寄样"
		// 4："寄样中"
		// 5："待平台审版"
		// 6："审版不合格"
		// 7："平台核价中"
		// 8："待修改生产资料"
		// 9："核价未通过"
		// 10："待下首单"
		// 11："已下首单"
		// 12："已加入站点"
		// 13："已下架"
		// 14："待卖家修改"
		// 15："已修改"
		// 16："服饰可加色"
		// 17："已终止"
		SelectStatus int `json:"selectStatus"` // 选品状态
		SkuList      []struct {
			SkuId int64 `json:"skuId"` // SKU ID
		} `json:"skuList"` // sku 列表
		ApplyJitStatus  int  `json:"applyJitStatus"`  // 申诉 JIT 的状态(1：可申请、3：不可申请)
		SuggestCloseJit bool `json:"suggestCloseJit"` // 是否建议关闭 JIT 按钮
	} `json:"skcList"` // SKC 列表
}
