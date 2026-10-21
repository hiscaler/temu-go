package entity

type AdSummaryValue struct {
	Val int `json:"val"`
}
type AdSummary struct {
	Total    AdSummaryValue `json:"total"`
	Ad       AdSummaryValue `json:"ad"`
	NetTotal AdSummaryValue `json:"netTotal"`
	NetAd    AdSummaryValue `json:"netAd"`
}
type SemiAdDetail struct {
	Summary struct {
		Ctr             AdSummary `json:"ctr"`
		CartCnt         AdSummary `json:"cartCnt"`
		ClkCnt          AdSummary `json:"clkCnt"`
		OrderPayAmt     AdSummary `json:"orderPayAmt"`
		Spend           AdSummary `json:"spend"`
		OrderPayCnt     AdSummary `json:"orderPayCnt"`
		Roas            AdSummary `json:"roas"`
		Acos            AdSummary `json:"acos"`
		TransactionCost AdSummary `json:"transactionCost"`
		GoodsNum        AdSummary `json:"goodsNum"`
		ImprCnt         AdSummary `json:"imprCnt"`
		Cvr             AdSummary `json:"cvr"`
	} `json:"summary"`
	AdPhase            int `json:"adPhase"`
	SiteStatusInfoList []struct {
		ForbidReason string   `json:"forbidReason"`
		SiteNameList []string `json:"siteNameList"`
		AdShowStatus int      `json:"adShowStatus"`
	} `json:"siteStatusInfoList"`
	GoodsId           int `json:"goodsId"`
	ReportsSummaryDTO struct {
		ClkCntAll      AdSummaryValue `json:"clkCntAll"`
		OrderPayCntAll AdSummaryValue `json:"orderPayCntAll"`
		AdSpendAll     AdSummaryValue `json:"adSpendAll"`
		AcosAll        AdSummaryValue `json:"acosAll"`
		CtrAll         AdSummaryValue `json:"ctrAll"`
		ImprCntAll     AdSummaryValue `json:"imprCntAll"`
		OrderPayAmtAll AdSummaryValue `json:"orderPayAmtAll"`
		CartCntAll     AdSummaryValue `json:"cartCntAll"`
		RoasAll        AdSummaryValue `json:"roasAll"`
	} `json:"reportsSummaryDTO"`
	AdShowStatusList []struct {
		ForbidReason string `json:"forbidReason"`
		AdShowStatus int    `json:"adShowStatus"`
	} `json:"adShowStatusList"`
	Roas         int `json:"roas"`
	AdShowStatus int `json:"adShowStatus"`
	Budget       int `json:"budget"`
}
