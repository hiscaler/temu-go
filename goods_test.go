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
		Cat3Id:                     9809,
		Cat4Id:                     10018,
		Cat5Id:                     10023,
		Cat6Id:                     0,
		Cat7Id:                     0,
		Cat8Id:                     0,
		Cat9Id:                     0,
		Cat10Id:                    0,
		ProductWarehouseRouteReq:   nil,
		ProductI18nReqs:            []GoodsCreateProductI18n{},
		ProductName:                "Ultra Absorbent & Soft Cotton Hand Towels",
		ProductCarouseVideoReqList: nil,
		ProductCustomReq:           nil,
		CarouselImageUrls: []string{
			"https://img.cdnfe.com/product/fancy/952b1680-bad4-431d-b288-7c4e41629757.jpg",
			"https://img.cdnfe.com/product/fancy/952b1680-bad4-431d-b288-7c4e41629757.jpg",
			"https://img.cdnfe.com/product/fancy/084ba15f-0af1-4dcd-ae4e-1976b7eecead.jpg",
			"https://img.cdnfe.com/product/fancy/20d9ccbe-f159-4a56-9e8a-4e17ae6a24a4.jpg",
			"https://img.cdnfe.com/product/fancy/5504b457-98d3-435c-84ef-a68cd56cac33.jpg",
		},
		CarouselImageI18nReqs: nil,
		ProductOuterPackageImageReqs: []GoodsCreateProductOuterPackageImage{
			{
				ImageUrl: "https://kjpfs-cn.kuajingmaihuo.com/product-material-private-tag/1f29b02490/1f15de43-494a-4e6f-b241-fb9f48140c3f_1080x1080.jpeg",
			},
		},
		MaterialImgUrl: "https://img.cdnfe.com/product/fancy/69632d08-6605-42f5-ae03-f14263031160.jpg",
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
				ValueUnit:        "%",
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
				Region1ShortName string `json:"region1ShortName"`
				Region2Id        string `json:"region2Id"`
			}{Region1ShortName: "CN", Region2Id: "43000000000006"},
		},
		ProductSkcReqs: []GoodsCreateProductSkc{
			{
				PreviewImgUrls:                  []string{"https://img.cdnfe.com/product/fancy/69632d08-6605-42f5-ae03-f14263031160.jpg"},
				ProductSkcCarouselImageI18nReqs: []ProductImageUrl{},
				ColorImageUrl:                   "https://img.cdnfe.com/product/fancy/69632d08-6605-42f5-ae03-f14263031160.jpg",
				MainProductSkuSpecReqs: []entity.Specification{
					{
						SpecId:         0,
						SpecName:       "",
						ParentSpecId:   0,
						ParentSpecName: "",
					},
				},
				ExtCode:     "test111",
				IsBasePlate: 1,
				ProductSkuReqs: []GoodsCreateProductSku{
					{
						ThumbUrl: "https://img.cdnfe.com/product/fancy/69632d08-6605-42f5-ae03-f14263031160.jpg",
						// todo recheck
						// ProductSkuThumbUrlI18nReqs: []ProductImageUrl{
						//	{
						//		ImgUrlList: []string{"https://img.kwcdn.com/product/open/ebcea5a3ffdc4209a7865b6074460ca2-goods.jpeg"},
						//		Language:   "US",
						//	},
						// },
						CurrencyType: "CNY",
						ProductSkuSpecReqs: []entity.Specification{
							{
								SpecId:         2,
								SpecName:       "红色",
								ParentSpecId:   1001,
								ParentSpecName: "颜色",
							},
						},
						SupplierPrice: 100,
						// SiteSupplierPrices:         []GoodsCreateProductSkuSiteSupplierPrice{},
						ProductSkuStockQuantityReq: nil,
						ProductSkuMultiPackReq:     nil,
						ProductSkuSuggestedPriceReq: &GoodsCreateProductSkuSuggestedPrice{
							SpecialSuggestedPrice:      "NA",
							SuggestedPriceCurrencyType: "",
							SuggestedPrice:             0,
						},
						ProductSkuWhExtAttrReq: &GoodsCreateProductSkuWhExtAttr{
							ProductSkuWeightReq:         GoodsCreateProductSkuWhExtAttrSensitiveLimitProductSkuWeight{Value: 10000},
							ProductSkuSameReferPriceReq: GoodsCreateProductSkuWhExtAttrSameReferPrice{},
							ProductSkuSensitiveLimitReq: nil,
							ProductSkuVolumeReq: GoodsCreateProductSkuWhExtAttrVolume{
								Len:    100,
								Width:  100,
								Height: 100,
							},
							ProductSkuSensitiveAttrReq: GoodsCreateProductSkuWhExtAttrSensitiveAttr{
								IsSensitive:   0,
								SensitiveList: []int{},
							},
							ProductSkuBarCodeReqs: nil,
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
		// ProductSemiManagedReq: &GoodsCreateProductSemiManaged{
		//	BindSiteIds: []int{100},
		// },
		// ProductShipmentReq: &GoodsCreateProductShipment{
		//	FreightTemplateId:   "",
		//	ShipmentLimitSecond: 259200,
		// },
		AddProductChannelType:  1,
		MaterialMultiLanguages: nil,
	}

	// 调用创建商品接口
	result, err := temuClient.Services.Goods.Create(ctx, createRequest)
	assert.Equalf(t, nil, err, "Create(ctx, %v)", jsonx.ToPrettyJson(createRequest))
	if err == nil {
		fmt.Println(jsonx.ToPrettyJson(result))
	}
}
