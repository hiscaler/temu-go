package entity

type AdSummary struct {
	Total struct {
		Val int `json:"val"`
	} `json:"total"`
	Ad struct {
		Val int `json:"val"`
	} `json:"ad"`
	NetTotal struct {
		Val int `json:"val"`
	} `json:"netTotal"`
	NetAd struct {
		Val int `json:"val"`
	} `json:"netAd"`
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
		ClkCntAll struct {
			Val int `json:"val"`
		} `json:"clkCntAll"`
		OrderPayCntAll struct {
			Val int `json:"val"`
		} `json:"orderPayCntAll"`
		AdSpendAll struct {
			Val int `json:"val"`
		} `json:"adSpendAll"`
		AcosAll struct {
			Val int `json:"val"`
		} `json:"acosAll"`
		CtrAll struct {
			Val int `json:"val"`
		} `json:"ctrAll"`
		ImprCntAll struct {
			Val int `json:"val"`
		} `json:"imprCntAll"`
		OrderPayAmtAll struct {
			Val int `json:"val"`
		} `json:"orderPayAmtAll"`
		CartCntAll struct {
			Val int `json:"val"`
		} `json:"cartCntAll"`
		RoasAll struct {
			Val int `json:"val"`
		} `json:"roasAll"`
	} `json:"reportsSummaryDTO"`
	AdShowStatusList []struct {
		ForbidReason string `json:"forbidReason"`
		AdShowStatus int    `json:"adShowStatus"`
	} `json:"adShowStatusList"`
	Roas         int `json:"roas"`
	AdShowStatus int `json:"adShowStatus"`
	Budget       int `json:"budget"`
}
