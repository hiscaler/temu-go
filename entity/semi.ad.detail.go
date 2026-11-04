package entity

type AdSummaryValue struct {
	Val int `json:"val"` // 指标数值
}

// AdSummary 广告指标按总量、广告和净值维度汇总的数据。
type AdSummary struct {
	Total    AdSummaryValue `json:"total"`    // 总量（文档：total number of records）
	Ad       AdSummaryValue `json:"ad"`       // 广告指标值
	NetTotal AdSummaryValue `json:"netTotal"` // 净总指标值
	NetAd    AdSummaryValue `json:"netAd"`    // 净广告指标值
}

// SemiAdDetail 商品广告活动详情。
type SemiAdDetail struct {
	Summary struct {
		Ctr             AdSummary `json:"ctr"`             // 点击率（CTR）
		CartCnt         AdSummary `json:"cartCnt"`         // 加购次数
		ClkCnt          AdSummary `json:"clkCnt"`          // 点击次数
		OrderPayAmt     AdSummary `json:"orderPayAmt"`     // 基础价格销售额
		Spend           AdSummary `json:"spend"`           // 广告消耗
		OrderPayCnt     AdSummary `json:"orderPayCnt"`     // 订单数
		Roas            AdSummary `json:"roas"`            // 广告支出回报率（ROAS）
		Acos            AdSummary `json:"acos"`            // 广告销售成本率（ACOS）
		TransactionCost AdSummary `json:"transactionCost"` // 成交成本
		GoodsNum        AdSummary `json:"goodsNum"`        // 商品数量
		ImprCnt         AdSummary `json:"imprCnt"`         // 曝光次数
		Cvr             AdSummary `json:"cvr"`             // 转化率（CVR）
	} `json:"summary"`
	AdPhase            int `json:"adPhase"` // 广告阶段：0 首阶段（学习期）；1 第二阶段（稳定期）
	SiteStatusInfoList []struct {
		ForbidReason string   `json:"forbidReason"` // 禁止投放原因
		SiteNameList []string `json:"siteNameList"` // 站点名称列表
		AdShowStatus int      `json:"adShowStatus"` // 按站点汇总的广告状态，取值见 SemiAdDetail.AdShowStatus
	} `json:"siteStatusInfoList"`
	GoodsId           int `json:"goodsId"` // 商品 ID
	ReportsSummaryDTO struct {
		ClkCntAll      AdSummaryValue `json:"clkCntAll"`      // 点击次数
		OrderPayCntAll AdSummaryValue `json:"orderPayCntAll"` // 订单数
		AdSpendAll     AdSummaryValue `json:"adSpendAll"`     // 广告消耗
		AcosAll        AdSummaryValue `json:"acosAll"`        // 广告销售成本率（ACOS）
		CtrAll         AdSummaryValue `json:"ctrAll"`         // 点击率（CTR）
		ImprCntAll     AdSummaryValue `json:"imprCntAll"`     // 曝光次数
		OrderPayAmtAll AdSummaryValue `json:"orderPayAmtAll"` // 基础价格销售额
		CartCntAll     AdSummaryValue `json:"cartCntAll"`     // 加购次数
		RoasAll        AdSummaryValue `json:"roasAll"`        // 广告支出回报率（ROAS）
	} `json:"reportsSummaryDTO"`
	AdShowStatusList []struct {
		ForbidReason string `json:"forbidReason"` // 禁止投放原因
		AdShowStatus int    `json:"adShowStatus"` // 广告状态，最多返回 4 项，取值见 SemiAdDetail.AdShowStatus
	} `json:"adShowStatusList"`
	Roas         int `json:"roas"`         // 广告支出回报率（ROAS）
	AdShowStatus int `json:"adShowStatus"` // 广告状态：0 无余额；1 当日预算为 0；2 商品售罄；3 商品下架；4 商品审核中；5 审核拒绝；6 推广受限；7 暂停；8 推广中；9 已删除；10 未创建；11 流量低；12 低流量低 ROAS；13 审核通过；14 达到月度上限；15 已禁止
	Budget       int `json:"budget"`       // 广告预算
}
