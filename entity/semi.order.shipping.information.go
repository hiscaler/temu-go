package entity

// SemiOrderShippingInformation 半托管订单收货信息
type SemiOrderShippingInformation struct {
	ReceiptName           string `json:"receiptName"`           // 收货人姓名
	Mail                  string `json:"mail"`                  // 邮箱（目前透出为虚拟邮箱）
	ReceiptAdditionalName string `json:"receiptAdditionalName"` // 片假名
	Mobile                string `json:"mobile"`                // 收货人手机号(虚拟号)，parentOrderStatus在待发货UN_SHIPPING、已发货SHIPPED）才会返回，格式为：11位手机号-4位分机号，形如："12312312312-1234"
	BackupMobile          string `json:"backupMobile"`          // 备用电话，部分国家会有返回
	RegionName1           string `json:"regionName1"`           // 区域地址1，国家，如：United States
	RegionName2           string `json:"regionName2"`           // 区域地址2，州，如NJ
	RegionName3           string `json:"regionName3"`           // 区域地址3，城市，如Jersey City
	RegionName4           string `json:"regionName4"`           // 区域地址4，日韩等国会有此区划
	AddressLine1          string `json:"addressLine1"`          // 详细地址行1，街道，如River Road
	AddressLine2          string `json:"addressLine2"`          // 详细地址行2，详细地址，如APT2601
	AddressLine3          string `json:"addressLine3"`          // 详细地址行3，部分国家会有返回
	AddressLineAll        string `json:"addressLineall"`        // 详细地址行1,详细地址行2，详细地址行3等
	PostCode              string `json:"postCode"`              // 邮编地址
}
