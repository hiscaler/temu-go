package entity

type AdStat struct {
	Val int `json:"val"`
}

// Ad 广告
type Ad struct {
	ProductId          int `json:"productId"`
	AdPhase            int `json:"adPhase"` // 广告阶段：0：一阶段，学习期；1：二阶段，平稳期
	SiteStatusInfoList []struct {
		ForbidReason string   `json:"forbidReason"`
		SiteNameList []string `json:"siteNameList"`
		AdShowStatus int      `json:"adShowStatus"`
	} `json:"siteStatusInfoList"`
	ReportsSummaryDTO AdReportSummary `json:"reportsSummaryDTO"` // 整体报表信息
	Roas              int             `json:"roas"`              // 目标广告投资回报率
	AdShowStatus      int             `json:"adShowStatus"`      // 广告状态：0：no balance；1：today budget 0；2：goods sold out；3：goods offline；4：goods under review；5：review rejected；6：promotion limited；7：pause；8：promoting；9：del；10：not creat；11：low traffic；12：low traffic soft roas
	Budget            int             `json:"budget"`            // 广告日预算
}

// AdLog 广告操作日志
type AdLog struct {
	ChangeInfo       string `json:"changeInfo"`       // 修改详细内容
	EventType        string `json:"eventType"`        // 修改类型：目前有新增，更新，删除三种类型
	UpdateSellerName string `json:"updateSellerName"` // 商家名称
	UpdatedAt        string `json:"updatedAt"`        // 修改时间
}

type AdReportSummary struct {
	ClkCntAll      AdStat `json:"clkCntAll"`      // 点击量
	OrderPayCntAll AdStat `json:"orderPayCntAll"` // 订单量
	AdSpendAll     AdStat `json:"adSpendAll"`     // 总花费
	AcosAll        AdStat `json:"acosAll"`        // 广告费比
	CtrAll         AdStat `json:"ctrAll"`         // 点击率
	ImprCntAll     AdStat `json:"imprCntAll"`     // 曝光量
	OrderPayAmtAll AdStat `json:"orderPayAmtAll"` // 申报价销售额
	CartCntAll     AdStat `json:"cartCntAll"`     // 加购数
	RoasAll        AdStat `json:"roasAll"`        // 广告投资回报率
}

// AdReport 广告报告
type AdReport struct {
	Ts          int64  `json:"ts"`          // 时间段开始
	Ctr         AdStat `json:"ctr"`         // 点击率
	CartCnt     AdStat `json:"cartCnt"`     // 加购数
	ClkCnt      AdStat `json:"clkCnt"`      // 点击量
	OrderPayAmt AdStat `json:"orderPayAmt"` // 申报价销售额
	OrderPayCnt AdStat `json:"orderPayCnt"` // 订单量
	Roas        AdStat `json:"roas"`        // 广告投资回报率
	Acos        AdStat `json:"acos"`        // 广告费比
	Adspend     AdStat `json:"adSpend"`     // 总花费
	ImprCnt     AdStat `json:"imprCnt"`     // 曝光量
}
