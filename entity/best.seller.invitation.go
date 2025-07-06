package entity

type BestSellerInvitation struct {
	InvitationId   int      `json:"invitationId"`   //  招标单 ID
	InvitationName string   `json:"invitationName"` // 招标单名称
	SiteId         int      `json:"siteId"`         // 站点 ID
	SiteName       string   `json:"siteName"`       // 站点名称
	CatIdList      []int    `json:"catIdList"`      // 分类 ID
	CatList        []string `json:"catList"`        // 分类名称
	ImageList      []string `json:"imageList"`      // 图片列表
}
