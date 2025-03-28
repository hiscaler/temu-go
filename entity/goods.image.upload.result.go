package entity

// GoodsImageUploadResult 上传货品图片返回结果
type GoodsImageUploadResult struct {
	ImageUrl string   `json:"imageUrl"` // 原图链接
	Url      string   `json:"url"`      // 单张 AI 裁图链接
	Urls     []string `json:"urls"`     // 多张 AI 裁图链接
}
