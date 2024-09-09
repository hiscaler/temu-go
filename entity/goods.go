package entity

type Goods struct {
	ProductProperties []struct {
		Vid              int         `json:"vid"`
		ValueUnit        string      `json:"valueUnit"`
		Language         interface{} `json:"language"`
		Pid              int         `json:"pid"`
		TemplatePid      int         `json:"templatePid"`
		NumberInputValue string      `json:"numberInputValue"`
		PropValue        string      `json:"propValue"`
		ValueExtendInfo  string      `json:"valueExtendInfo"`
		PropName         string      `json:"propName"`
		RefPid           int         `json:"refPid"`
	} `json:"productProperties"`
	ProductId      int `json:"productId"`
	ProductJitMode struct {
		QuickSellAgtSignStatus interface{} `json:"quickSellAgtSignStatus"`
		MatchJitMode           bool        `json:"matchJitMode"`
	} `json:"productJitMode"`
	ProductSkuSummaries []struct {
		ProductSkuId        int64  `json:"productSkuId"`
		ExtCode             string `json:"extCode"`
		ProductSkuWhExtAttr struct {
			ProductSkuWeight struct {
				Value int `json:"value"`
			} `json:"productSkuWeight"`
			ProductSkuWmsVolume     interface{} `json:"productSkuWmsVolume"`
			ProductSkuBarCodes      interface{} `json:"productSkuBarCodes"`
			ProductSkuSubSellMode   interface{} `json:"productSkuSubSellMode"`
			ProductSkuSensitiveAttr struct {
				SensitiveTypes []interface{} `json:"sensitiveTypes"`
				IsSensitive    int           `json:"isSensitive"`
			} `json:"productSkuSensitiveAttr"`
			ProductSkuFragileLabels    interface{} `json:"productSkuFragileLabels"`
			ProductSkuNewSensitiveAttr struct {
				Force2NormalTypes interface{} `json:"force2NormalTypes"`
				SensitiveList     []int       `json:"sensitiveList"`
				IsForce2Normal    bool        `json:"isForce2Normal"`
			} `json:"productSkuNewSensitiveAttr"`
			ProductSkuVolumeLabel interface{} `json:"productSkuVolumeLabel"`
			ProductSkuWmsWeight   interface{} `json:"productSkuWmsWeight"`
			ProductSkuVolume      struct {
				Len    int `json:"len"`
				Width  int `json:"width"`
				Height int `json:"height"`
			} `json:"productSkuVolume"`
			ProductSkuSensitiveLimit interface{} `json:"productSkuSensitiveLimit"`
			ProductSkuWmsVolumeLabel interface{} `json:"productSkuWmsVolumeLabel"`
		} `json:"productSkuWhExtAttr"`
		VirtualStock       interface{} `json:"virtualStock"`
		ProductSkuSpecList []struct {
			SpecId         int    `json:"specId"`
			ParentSpecName string `json:"parentSpecName"`
			SpecName       string `json:"specName"`
			ParentSpecId   int    `json:"parentSpecId"`
		} `json:"productSkuSpecList"`
	} `json:"productSkuSummaries"`
	ProductName              string      `json:"productName"`
	CreatedAt                int64       `json:"createdAt"`
	ProductSemiManaged       interface{} `json:"productSemiManaged"`
	IsSupportPersonalization bool        `json:"isSupportPersonalization"`
	ExtCode                  string      `json:"extCode"`
	LeafCat                  struct {
		CatId   int    `json:"catId"`
		CatName string `json:"catName"`
	} `json:"leafCat"`
	Categories struct {
		Cat8 struct {
			CatId   int    `json:"catId"`
			CatName string `json:"catName"`
		} `json:"cat8"`
		Cat9 struct {
			CatId   int    `json:"catId"`
			CatName string `json:"catName"`
		} `json:"cat9"`
		Cat6 struct {
			CatId   int    `json:"catId"`
			CatName string `json:"catName"`
		} `json:"cat6"`
		Cat7 struct {
			CatId   int    `json:"catId"`
			CatName string `json:"catName"`
		} `json:"cat7"`
		Cat4 struct {
			CatId   int    `json:"catId"`
			CatName string `json:"catName"`
		} `json:"cat4"`
		Cat5 struct {
			CatId   int    `json:"catId"`
			CatName string `json:"catName"`
		} `json:"cat5"`
		Cat2 struct {
			CatId   int    `json:"catId"`
			CatName string `json:"catName"`
		} `json:"cat2"`
		Cat3 struct {
			CatId   int    `json:"catId"`
			CatName string `json:"catName"`
		} `json:"cat3"`
		Cat10 struct {
			CatId   int    `json:"catId"`
			CatName string `json:"catName"`
		} `json:"cat10"`
		Cat1 struct {
			CatId   int    `json:"catId"`
			CatName string `json:"catName"`
		} `json:"cat1"`
		LeafCat interface{} `json:"leafCat"`
	} `json:"categories"`
	ProductSkcId    int    `json:"productSkcId"`
	MatchSkcJitMode bool   `json:"matchSkcJitMode"`
	MainImageUrl    string `json:"mainImageUrl"`
}
