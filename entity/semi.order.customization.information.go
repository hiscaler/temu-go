package entity

import (
	"errors"
	"slices"

	"github.com/goccy/go-json"
	"github.com/hiscaler/gox/nullx"
	"github.com/spf13/cast"
	"gopkg.in/guregu/null.v4"
)

// SemiOrderCustomizationInformationPreview 半托订单定制信息预览
type SemiOrderCustomizationInformationPreview struct {
	// type of preview item, enum values:
	// - 1: overall preview image(If the product does not have a customized area configured, it represents the effect image uploaded by the merchant)
	// - 3: user uploaded image
	// - 4: customized text
	PreviewType        int         `json:"previewType"`        // 预览类型
	CustomizedAreaId   null.String `json:"customizedAreaId"`   // 定制区域 ID. This field will only be returned when templateType=3, previewType=3 or 4
	ImageUrl           null.String `json:"imageUrl"`           // Image URL
	ImageDownloadUrl   null.String `json:"imageDownloadUrl"`   // 图片下载地址
	ImageDownloadError null.String `json:"imageDownloadError"` // 图片下载错误
	CustomizedText     null.String `json:"customizedText"`     // 定制文本
}

// SemiOrderCustomizationInformation 半托订单定制信息
type SemiOrderCustomizationInformation struct {
	OrderSn string `json:"orderSn"` // OrderSn corresponding to customized information
	// Customization template type when user created customized information, return null when there is no template for the product, enum values:
	// - 1: only image
	// - 2: only text
	// - 3: text and image
	TemplateType   null.Int                                   `json:"templateType"`
	PreviewList    []SemiOrderCustomizationInformationPreview `json:"previewList"`    // Graphic customization preview information, this field will only be returned when customizedType=2
	CustomizedData CustomizedData                             `json:"customizedData"` // Graphic customization content, in json format, this field will only be returned when customizedType=2
	CustomizedText null.String                                `json:"customizedText"` // Customization text, this field will only be returned when customizedType=1
	TemplateId     int                                        `json:"templateId"`     // Customization template ID when user created customized information, return null when there is no template for the product
	// Customized type, enum values:
	// - 1: pure text customization, no customized templates
	// - 2: customized graphics and text, with customized templates available
	CustomizedType int `json:"customizedType"`
}

// 定制信息解析

const (
	image = 1
	text  = 2
)

type Dimension struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}
type Color struct {
	Red   int `json:"red"`
	Green int `json:"green"`
	Blue  int `json:"blue"`
}

type CustomizedSurfaceImageElement struct {
	Type      int    `json:"type"`
	Require   bool   `json:"require"`
	MaxSize   int    `json:"maxSize"`
	MaxWidth  int    `json:"maxWidth"`
	MaxHeight int    `json:"maxHeight"`
	Path      string `json:"path"`
	RIndex    int    `json:"rIndex"`
	ImageUrl  string `json:"imageUrl"`
}

type CustomizedSurfaceTextElement struct {
	Type              int    `json:"type"`
	Require           bool   `json:"require"`
	LengthLimit       int    `json:"lengthLimit"`
	TextAlign         int    `json:"textAlign"`
	Path              string `json:"path"`
	RIndex            int    `json:"rIndex"`
	Color             Color  `json:"color"`
	Text              string `json:"text"`
	UserPlacementData struct {
		FontSize float64  `json:"fontSize"`
		Position Position `json:"position"`
	} `json:"userPlacementData"`
}

type CustomizedSurfaceTextRegion struct {
	Elements  []CustomizedSurfaceTextElement `json:"elements"`
	Dimension Dimension                      `json:"dimension"`
	Position  Position                       `json:"position"`
}

type CustomizedSurfaceImageRegion struct {
	Elements []any `json:"elements"`
}

type CustomizedSurface struct {
	Regions   []CustomizedSurfaceImageRegion `json:"regions"`
	BaseImage struct {
		ImageUrl  string    `json:"imageUrl"`
		Dimension Dimension `json:"dimension"`
	} `json:"baseImage"`
	MaskImage struct {
		ImageUrl  string    `json:"imageUrl"`
		Dimension Dimension `json:"dimension"`
	} `json:"maskImage"`
}

type CustomizedSurfaces struct {
	Surfaces []CustomizedSurface `json:"surfaces"`
}

type CustomizedData string

type ParseResult struct {
	Region       int         `json:"region"`       // 区块
	PreviewImage null.String `json:"previewImage"` // 预览图
	Type         string      `json:"type"`         // 类型
	Text         null.String `json:"text"`         // 定制文本
	Image        null.String `json:"image"`        // 定制图片
	Error        null.String `json:"error"`        // 错误信息
	ExpireTime   int64       `json:"expireTime"`   // 过期时间
}

func (cd CustomizedData) Parse() (prs []ParseResult, err error) {
	if cd == "" {
		return prs, errors.New("customizedData is empty")
	}

	var css CustomizedSurfaces
	if err = json.Unmarshal([]byte(cd), &css); err != nil {
		return
	}

	prs = make([]ParseResult, 0)
	for _, surface := range css.Surfaces {
		for _, region := range surface.Regions {
			for _, element := range region.Elements {
				var t map[string]any
				t, err = cast.ToStringMapE(element)
				if err != nil {
					continue
				}

				elementTypeValue, ok := t["type"]
				if !ok {
					continue
				}

				elementType := cast.ToInt(elementTypeValue)
				if !slices.Contains([]int{text, image}, elementType) {
					continue
				}

				var b []byte
				b, err = json.Marshal(element)
				if err != nil {
					continue
				}

				switch elementType {
				case text:
					var textElement CustomizedSurfaceTextElement
					err = json.Unmarshal(b, &textElement)
					if err != nil {
						continue
					}
					prs = append(prs, ParseResult{
						Region:       textElement.RIndex,
						PreviewImage: nullx.StringFrom(surface.MaskImage.ImageUrl),
						Type:         "text",
						Text:         null.NewString(textElement.Text, true),
					})

				case image:
					var imageElement CustomizedSurfaceImageElement
					err = json.Unmarshal(b, &imageElement)
					if err != nil {
						continue
					}
					prs = append(prs, ParseResult{
						Region:       imageElement.RIndex,
						PreviewImage: nullx.StringFrom(surface.MaskImage.ImageUrl),
						Type:         "image",
						Image:        null.NewString(imageElement.ImageUrl, true),
					})

				default:
					continue
				}
			}
		}
	}

	return prs, nil
}
