package temu

import (
	"fmt"
	"github.com/hiscaler/gox/jsonx"
	"github.com/hiscaler/temu-go/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_goodsService_Query(t *testing.T) {
	params := GoodsQueryParams{
		// ProductSkcIds: []int64{2646847407},
		SkuExtCodes:    []string{"8502937482"},
		ProductSkcIds:  []int64{7469668867},
		CreatedAtStart: "2024-11-18 12:00:00",
		CreatedAtEnd:   "2024-11-18 23:59:59",
	}
	params.Page = 1
	params.PageSize = 2
	items, _, _, _, err := temuClient.Services.Goods.Query(ctx, params)
	assert.Equalf(t, nil, err, "Services.Goods.Query(ctx, %s)", jsonx.ToPrettyJson(params))
	_ = items
	if len(items) != 0 {
		item := items[0]
		var sales entity.Goods
		// 根据商品 SKC ID 查询
		sales, err = temuClient.Services.Goods.One(ctx, item.ProductSkcId)
		assert.Equalf(t, nil, err, "Services.Goods.One(ctx, %d)", item.ProductSkcId)
		assert.Equalf(t, item, sales, "Services.Goods.One(ctx, %d)", item.ProductSkcId)
	}
}

func Test_goodsService_Detail(t *testing.T) {
	var productId int64 = 141911679
	detail, err := temuClient.Services.Goods.Detail(ctx, productId)
	assert.Equalf(t, nil, err, "Services.Goods.One(ctx, %d)", productId)
	assert.Equalf(t, detail.ProductId, productId, "Services.Goods.One(ctx, %d)", productId)
}

func Test_goodsService_Create(t *testing.T) {
	createRequest := GoodsCreateRequest{
		Cat1Id:                     9711,
		Cat2Id:                     9712,
		Cat3Id:                     10018,
		Cat4Id:                     10023,
		Cat5Id:                     10024,
		Cat6Id:                     0,
		Cat7Id:                     0,
		Cat8Id:                     0,
		Cat9Id:                     0,
		Cat10Id:                    0,
		ProductWarehouseRouteReq:   nil,
		ProductI18nReqs:            nil,
		ProductName:                "Ultra Absorbent & Soft Cotton Hand Towels",
		ProductCarouseVideoReqList: nil,
		ProductCustomReq:           nil,
		CarouselImageUrls: []string{
			"https://img.kwcdn.com/product/open/ebcea5a3ffdc4209a7865b6074460ca2-goods.jpeg",
		},
		//CarouselImageI18nReqs:        nil,
		ProductOuterPackageImageReqs: nil,
		MaterialImgUrl:               "https://img.kwcdn.com/product/open/ebcea5a3ffdc4209a7865b6074460ca2-goods.jpeg",
		ProductPropertyReqs: []GoodsCreateProductProperty{
			{
				TemplatePid:      201806,
				Pid:              97,
				RefPid:           185,
				PropName:         "特征",
				Vid:              3294,
				PropValue:        "防褪色",
				ValueUnit:        "",
				ValueExtendInfo:  "",
				NumberInputValue: "",
			},
			{
				TemplatePid:      201811,
				Pid:              4,
				RefPid:           20,
				PropName:         "护理说明",
				Vid:              2207,
				PropValue:        "只能手洗",
				ValueUnit:        "",
				ValueExtendInfo:  "",
				NumberInputValue: "",
			},
			{
				TemplatePid:      201817,
				Pid:              112,
				RefPid:           131,
				PropName:         "形状",
				Vid:              2466,
				PropValue:        "圆形",
				ValueUnit:        "",
				ValueExtendInfo:  "",
				NumberInputValue: "",
			},
			{
				TemplatePid:      201831,
				Pid:              89,
				RefPid:           121,
				PropName:         "材料",
				Vid:              2197,
				PropValue:        "涤纶",
				ValueUnit:        "",
				ValueExtendInfo:  "",
				NumberInputValue: "",
			},
			{
				TemplatePid:      201836,
				Pid:              176,
				RefPid:           541,
				PropName:         "毛巾主题",
				Vid:              3975,
				PropValue:        "字符",
				ValueUnit:        "",
				ValueExtendInfo:  "",
				NumberInputValue: "",
			},
			{
				TemplatePid:      201854,
				Pid:              3,
				RefPid:           19,
				PropName:         "风格",
				Vid:              136,
				PropValue:        "复古",
				ValueUnit:        "",
				ValueExtendInfo:  "",
				NumberInputValue: "",
			},
			{
				TemplatePid:      825348,
				Pid:              2,
				RefPid:           2021,
				PropName:         "封面材质",
				Vid:              74,
				PropValue:        "腈纶",
				ValueUnit:        "",
				ValueExtendInfo:  "",
				NumberInputValue: "100",
			},
			{
				TemplatePid:      961436,
				Pid:              1224,
				RefPid:           1192,
				PropName:         "织造方式",
				Vid:              29810,
				PropValue:        "梭织",
				ValueUnit:        "",
				ValueExtendInfo:  "",
				NumberInputValue: "",
			},
		},
		ProductSpecPropertyReqs: []GoodsCreateProductSpecProperty{
			{
				TemplatePid:      0,
				Pid:              0,
				RefPid:           0,
				Vid:              0,
				PropName:         "颜色",
				PropValue:        "红色",
				ParentSpecId:     1001,
				ParentSpecName:   "颜色",
				SpecId:           2,
				SpecName:         "红色",
				ValueGroupId:     0,
				ValueGroupName:   "",
				ValueUnit:        "",
				NumberInputValue: "",
				ValueExtendInfo:  "",
			},
		},
		ProductWhExtAttrReq: GoodsCreateProductWhExtAttr{
			OuterGoodsUrl: "",
			ProductOrigin: struct {
				CountryShortName string `json:"countryShortName"`
			}{
				"CN",
			},
		},
		ProductSkcReqs: []GoodsCreateProductSkc{
			{
				PreviewImgUrls:                  []string{"https://img.kwcdn.com/product/open/ebcea5a3ffdc4209a7865b6074460ca2-goods.jpeg"},
				ProductSkcCarouselImageI18nReqs: []ProductImageUrl{},
				ColorImageUrl:                   "https://img.kwcdn.com/product/open/ebcea5a3ffdc4209a7865b6074460ca2-goods.jpeg",
				MainProductSkuSpecReqs: []entity.Specification{
					{
						SpecId:         2,
						SpecName:       "红色",
						ParentSpecId:   10001,
						ParentSpecName: "颜色",
					},
				},
				ExtCode:     "test111",
				IsBasePlate: 1,
				ProductSkuReqs: []GoodsCreateProductSku{
					{
						ThumbUrl: "https://img.kwcdn.com/product/open/ebcea5a3ffdc4209a7865b6074460ca2-goods.jpeg",
						// todo recheck
						//ProductSkuThumbUrlI18nReqs: []ProductImageUrl{
						//	{
						//		ImgUrlList: []string{"https://img.kwcdn.com/product/open/ebcea5a3ffdc4209a7865b6074460ca2-goods.jpeg"},
						//		Language:   "US",
						//	},
						//},
						CurrencyType: "CNY",
						ProductSkuSpecReqs: []entity.Specification{
							{
								SpecId:         2,
								SpecName:       "红色",
								ParentSpecId:   10001,
								ParentSpecName: "颜色",
							},
						},
						SupplierPrice:               100,
						SiteSupplierPrices:          []GoodsCreateProductSkuSiteSupplierPrice{},
						ProductSkuStockQuantityReq:  nil,
						ProductSkuMultiPackReq:      nil,
						ProductSkuSuggestedPriceReq: nil,
						ProductSkuWhExtAttrReq: &GoodsCreateProductSkuWhExtAttr{
							ProductSkuWeightReq:         GoodsCreateProductSkuWhExtAttrSensitiveLimitProductSkuWeight{Value: 100},
							ProductSkuSameReferPriceReq: GoodsCreateProductSkuWhExtAttrSameReferPrice{},
							ProductSkuSensitiveLimitReq: nil,
							ProductSkuVolumeReq: GoodsCreateProductSkuWhExtAttrVolume{
								Len:    100,
								Width:  100,
								Height: 100,
							},
							ProductSkuSensitiveAttrReq: GoodsCreateProductSkuWhExtAttrSensitiveAttr{
								IsSensitive:   0,
								SensitiveList: nil,
							},
							ProductSkuBarCodeReqs: nil,
							ExtCode:               "",
						},
						ExtCode: "extcode1",
					},
				},
			},
		},
		SizeTemplateIds:     nil,
		GoodsModelReqs:      nil,
		ShowSizeTemplateIds: nil,
		ProductOuterPackageReq: &GoodsCreateProductOuterPackage{
			PackageShape: 0,
			PackageType:  2,
		},
		ProductGuideFileReqs:     nil,
		GoodsLayerDecorationReqs: nil,
		PersonalizationSwitch:    0,
		//ProductSemiManagedReq: &GoodsCreateProductSemiManaged{
		//	BindSiteIds: []int{100},
		//},
		//ProductShipmentReq: &GoodsCreateProductShipment{
		//	FreightTemplateId:   "",
		//	ShipmentLimitSecond: 259200,
		//},
		AddProductChannelType:  1,
		MaterialMultiLanguages: nil,
	}

	// 调用创建商品接口
	result, err := temuClient.Services.Goods.Create(ctx, createRequest)
	if err != nil {
		fmt.Printf("创建商品错误: %s\n", err.Error())
		return
	}

	fmt.Println(jsonx.ToPrettyJson(result))
}
