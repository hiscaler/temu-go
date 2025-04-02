package entity

// PictureCompressionResult 图片压缩结果
type PictureCompressionResult struct {
	Size      int64  `json:"size"`
	OriginUrl string `json:"originUrl"`
	Width     int    `json:"width"`
	ResultUrl string `json:"resultUrl"`
	Height    int    `json:"height"`
}
