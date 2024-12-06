package temu

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
	"github.com/hiscaler/temu-go/validators/is"
	"gopkg.in/guregu/null.v4"
	"strconv"
	"time"
)

// 商品数据服务
type goodsService service

type GoodsQueryParams struct {
	normal.ParameterWithPager
	Page                   int      `json:"page"`                             // 页码
	Cat1Id                 int      `json:"cat1Id,omitempty"`                 // 一级分类 ID
	Cat2Id                 int      `json:"cat2Id,omitempty"`                 // 二级分类 ID
	Cat3Id                 int      `json:"cat3Id,omitempty"`                 // 三级分类 ID
	Cat4Id                 int      `json:"cat4Id,omitempty"`                 // 四级分类 ID
	Cat5Id                 int      `json:"cat5Id,omitempty"`                 // 五级分类 ID
	Cat6Id                 int      `json:"cat6Id,omitempty"`                 // 六级分类 ID
	Cat7Id                 int      `json:"cat7Id,omitempty"`                 // 七级分类 ID
	Cat8Id                 int      `json:"cat8Id,omitempty"`                 // 八级分类 ID
	Cat9Id                 int      `json:"cat9Id,omitempty"`                 // 九级分类 ID
	Cat10Id                int      `json:"cat10Id,omitempty"`                // 十级分类 ID
	SkcExtCode             string   `json:"skcExtCode,omitempty"`             // 货品 SKC 外部编码
	SupportPersonalization int      `json:"supportPersonalization,omitempty"` // 是否支持定制品模板
	BindSiteIds            []int    `json:"bindSiteIds,omitempty"`            // 经营站点
	ProductName            string   `json:"productName,omitempty"`            // 货品名称
	ProductSkcIds          []int64  `json:"productSkcIds,omitempty"`          // SKC 列表
	QuickSellAgtSignStatus null.Int `json:"quickSellAgtSignStatus,omitempty"` // 快速售卖协议签署状态 0-未签署 1-已签署
	MatchJitMode           null.Int `json:"matchJitMode,omitempty"`           // 是否命中 JIT 模式
	SkcSiteStatus          null.Int `json:"skcSiteStatus,omitempty"`          // skc 加站点状态 (0: 未加入站点, 1: 已加入站点)
	CreatedAtStart         string   `json:"createdAtStart,omitempty"`         // 创建时间开始，毫秒级时间戳
	CreatedAtEnd           string   `json:"createdAtEnd,omitempty"`           // 创建时间结束，毫秒级时间戳
}

func (m GoodsQueryParams) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.CreatedAtStart,
			validation.When(m.CreatedAtStart != "" || m.CreatedAtEnd != "", validation.By(is.TimeRange(m.CreatedAtStart, m.CreatedAtEnd, time.DateOnly))),
		),
	)
}

// Query 货品列表查询
// https://seller.kuajingmaihuo.com/sop/view/750197804480663142#SjadVR
func (s goodsService) Query(ctx context.Context, params GoodsQueryParams) (items []entity.Goods, total, totalPages int, isLastPage bool, err error) {
	params.TidyPager()
	if params.Page <= 0 {
		params.Page = 1
	}
	if err = params.validate(); err != nil {
		err = invalidInput(err)
		return
	}

	if params.CreatedAtStart != "" && params.CreatedAtEnd != "" {
		t, _ := time.ParseInLocation(time.DateTime, params.CreatedAtStart+" 00:00:00", time.Local)
		params.CreatedAtStart = strconv.Itoa(int(t.UnixMilli()))
		t, _ = time.ParseInLocation(time.DateTime, params.CreatedAtEnd+" 23:59:59", time.Local)
		params.CreatedAtEnd = strconv.Itoa(int(t.UnixMilli()))
	}
	var result = struct {
		normal.Response
		Result struct {
			Data       []entity.Goods `json:"data"`
			TotalCount int            `json:"totalCount"`
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(params).
		SetResult(&result).
		Post("bg.goods.list.get")
	if err = recheckError(resp, result.Response, err); err != nil {
		return
	}

	items = result.Result.Data
	total, totalPages, isLastPage = parseResponseTotal(params.Page, params.PageSize, result.Result.TotalCount)
	return
}

// One 根据商品 SKC ID 查询
func (s goodsService) One(ctx context.Context, productSkcId int64) (item entity.Goods, err error) {
	items, _, _, _, err := s.Query(ctx, GoodsQueryParams{ProductSkcIds: []int64{productSkcId}})
	if err != nil {
		return
	}

	for _, v := range items {
		if v.ProductSkcId == productSkcId {
			return v, nil
		}
	}

	return item, ErrNotFound
}

// 添加货品

type GoodsCreateRequest struct {
	Cat1Id                   int `json:"cat1Id"`  // 一级类目id
	Cat2Id                   int `json:"cat2Id"`  // 二级类目id
	Cat3Id                   int `json:"cat3Id"`  // 三级类目id
	Cat4Id                   int `json:"cat4Id"`  // 四级类目id（没有的情况传 0）
	Cat5Id                   int `json:"cat5Id"`  // 五级类目id（没有的情况传 0）
	Cat6Id                   int `json:"cat6Id"`  // 六级类目id（没有的情况传 0）
	Cat7Id                   int `json:"cat7Id"`  // 七级类目id（没有的情况传 0）
	Cat8Id                   int `json:"cat8Id"`  // 八级类目id（没有的情况传 0）
	Cat9Id                   int `json:"cat9Id"`  // 九级类目id（没有的情况传 0）
	Cat10Id                  int `json:"cat10Id"` // 十级类目id（没有的情况传 0）
	ProductWarehouseRouteReq struct {
		TargetRouteList []struct {
			SiteIdList  []int  `json:"siteIdList"`  // 站点列表
			WarehouseId string `json:"warehouseId"` // 仓库ID，使用goods.warehouse.list.get查询
		} `json:"targetRouteList"` // 商品的目标路由列表
	} `json:"productWarehouseRouteReq"` // 库存仓库配置对象。	半托管发品必传，全托管店铺不需要传这个属性，传入会报错。
	ProductI18nReqs struct {
		Language    string `json:"language"`    // 语言编码，en-美国
		ProductName string `json:"productName"` // 对应语言的商品标题
	} `json:"productI18nReqs"` // 多语言标题设置
	ProductName                string `json:"productName "` // 货品名称
	ProductCarouseVideoReqList []struct {
		Vid      string `json:"vid"`      // 视频 VID
		CoverUrl string `json:"coverUrl"` // 视频封面图
		VideoUrl string `json:"videoUrl"` // 视频 URL
		Width    int    `json:"width"`    // 视频宽度
		Height   int    `json:"height"`   // 视频高度
	} `json:"productCarouseVideoReqList"` // 商品主图视频 关于如何上传视频，请对接视频上传相关接口，获取图片相关参数即可用于此处入参 https://seller.kuajingmaihuo.com/sop/view/852640595329867111
	ProductCustomReq struct {
		GoodsLabelName   string `json:"goodsLabelName"`   // 货品关务标签名称
		IsRecommendedTag bool   `json:"isRecommendedTag"` // 是否使用推荐标签
	} `json:"productCustomReq"` // 货品关务标签
	CarouselImageUrls     []string `json:"carouselImageUrls"` // 货品轮播图
	CarouselImageI18nReqs []struct {
		ImgUrlList []string `json:"imgUrlList"` // 图片列表
		Language   string   `json:"language"`   // 语言
	} `json:"carouselImageI18nReqs"` // 货品 SPU 多语言轮播图，服饰类不传，非服饰必传
	ProductOuterPackageImageReqs []struct {
		ImageUrl string `json:"imageUrl"` // 图片链接，通过图片上传接口，imageBizType=1获取
	} `json:"productOuterPackageImageReqs"` // 外包装图片
	MaterialImgUrl      string `json:"materialImgUrl"` // 素材图
	ProductPropertyReqs []struct {
		TemplatePid      int64  `json:"templatePid"`      // 模板属性id
		PID              int64  `json:"pid"`              // 属性 id
		RefPid           int64  `json:"refPid"`           // 引用属性 id
		PropName         string `json:"propName"`         // 引用属性名
		Vid              int64  `json:"vid"`              // 基础属性值id，没有的情况传0
		PropValue        string `json:"propValue"`        // 基础属性值
		ValueUnit        string `json:"valueUnit"`        // 属性值单位，没有的情况传空字符串
		NumberInputValue string `json:"numberInputValue"` // 属性输入值，例如：65.66
		ValueExtendInfo  string `json:"valueExtendInfo"`  // 属性扩展信息，attrs.get返回
	} `json:"productPropertyReqs"` // 货品属性
	ProductSpecPropertyReqs []struct {
		TemplatePid      int64  `json:"templatePid"`      // 模板属性id
		PID              int64  `json:"pid"`              // 属性 id
		RefPid           int64  `json:"refPid"`           // 引用属性 id
		PropName         string `json:"propName"`         // 引用属性名
		Vid              int64  `json:"vid"`              // 基础属性值id，没有的情况传0
		PropValue        string `json:"propValue"`        // 基础属性值
		ParentSpecId     int64  `json:"parentSpecId"`     // 父规格id
		ParentSpecName   string `json:"parentSpecName"`   // 父规格名称
		SpecId           int64  `json:"specId"`           // 规格id
		SpecName         string `json:"specName"`         // 规格名称
		ValueGroupId     int    `json:"valueGroupId"`     // 属性值组id，没有的情况传0
		ValueGroupName   string `json:"valueGroupName"`   // 属性值组名称，没有的情况传空字符串
		ValueUnit        string `json:"valueUnit"`        // 属性值单位，没有的情况传空字符串
		NumberInputValue string `json:"numberInputValue"` // 属性输入值，例如：65.66
		ValueExtendInfo  string `json:"valueExtendInfo"`  // 属性组扩展信息（色板）
	} `json:"productSpecPropertyReqs"` // 货品规格属性
	ProductWhExtAttrReq struct {
		OuterGoodsUrl string `json:"outerGoodsUrl"` //  站外商品链接
		ProductOrigin struct {
			CountryShortName string `json:"countryShortName"` // 国家简称 (二字简码)
		} `json:"productOrigin"` // 货品产地 (灰度内必传)，请注意，日本站点发品必须传产地，否则会被拦截
	} `json:"productWhExtAttrReq"` // 货品仓配供应链侧扩展属性请求
	ProductSkcReqs []struct {
		PreviewImgUrls                  []string `json:"previewImgUrls"` // SKC 轮播图列表
		ProductSkcCarouselImageI18nReqs []struct {
			ImgUrlList []string `json:"imgUrlList"` // 图片列表
			Language   string   `json:"language"`   // 语言
		} `json:"productSkcCarouselImageI18nReqs"` // SKC多语言轮播图，服饰类必传，非服饰不传
		ColorImageUrl          string `json:"colorImageUrl"` // SKC 色块图，可通过（bg.colorimageurl.get）转换获取
		MainProductSkuSpecReqs []struct {
			ParentSpecId   int64  `json:"parentSpecId"`   // 父规格id
			ParentSpecName string `json:"parentSpecName"` // 父规格名称
			SpecId         int64  `json:"specId"`         // 规格id
			SpecName       string `json:"specName"`       // 规格名称
		} `json:"mainProductSkuSpecReqs"` //  主销售规格列表
		IsBasePlate    int `json:"isBasePlate"` // 是否底板
		ProductSkuReqs []struct {
			ThumbUrl                   string `json:"thumbUrl"` // 预览图
			ProductSkuThumbUrlI18nReqs []struct {
				ImgUrlList []string `json:"imgUrlList"` // 图片列表
				Language   string   `json:"language"`   // 语言
			} `json:"productSkuThumbUrlI18nReqs"` // SKU多语言预览图，服饰类不传，非服饰非必传 （英国英语、中东英语必传）
			CurrencyType       string `json:"currencyType"` // 币种 (CNY: 人民币, USD: 美元) (默认人民币)
			ProductSkuSpecReqs []struct {
				ParentSpecId   int64  `json:"parentSpecId"`   // 父规格id
				ParentSpecName string `json:"parentSpecName"` // 父规格名称
				SpecId         int64  `json:"specId"`         // 规格id
				SpecName       string `json:"specName"`       // 规格名称
			} `json:"productSkuSpecReqs"` // 货品sku规格列表
			SupplierPrice      int64 `json:"supplierPrice"` // 全托供货价 （单位：人民币-分/美元-美分），半托不传
			SiteSupplierPrices []struct {
				SiteId        int64 `json:"siteId"`        // 申报价格站点id
				SupplierPrice int64 `json:"supplierPrice"` // 站点申报价格，单位 人民币：分，美元：美分
			} `json:"siteSupplierPrices"` // 站点供货价列表，半托必传
			ProductSkuStockQuantityReq struct {
				WarehouseStockQuantityReqs []struct {
					TargetStockAvailable int    `json:"targetStockAvailable"` // sku目标库存值（覆盖写）
					WarehouseId          string `json:"warehouseId"`          // 仓库 ID
				} `json:"warehouseStockQuantityReqs"` // 仓库存库存请求列表
			} `json:"productSkuStockQuantityReq"` // 货品sku仓库库存，半托管发品必传
			ProductSkuMultiPackReq struct {
				NumberOfPieces          int `json:"numberOfPieces"` // 件数，单品默认是1
				ProductSkuNetContentReq struct {
					NetContentUnitCode int `json:"netContentUnitCode"` // 净含量单位，1：液体盎司，2：毫升，3：加仑，4：升，5：克，6：千克，7：常衡盎司，8：磅
					NetContentNumber   int `json:"netContentNumber"`   // 净含量数值
				} `json:"productSkuNetContentReq"` // 净含量请求，传空对象表示空，指定类目灰度管控
				SkuClassification int `json:"skuClassification"` // sku分类，1：单品，2：组合装，3：混合套装
				PieceUnitCode     int `json:"pieceUnitCode"`     // 单件单位，1：件，2：双，3：包
			} `json:"productSkuMultiPackReq"` // 货品多包规请求
			// 货品sku建议价格请求
			// 1. 建议零售价是制造商为产品设定的建议零售价或推荐零售价。建议零售价必须是市场上的真实销售价格，且符合任何可适用的法律法规的规定。如您的商品在欧盟市场上销售，则该产品必须有欧盟零售商以此价格进行真实的广告宣传和销售。如果您的产品没有符合这些标准的建议零售价，请勿填写建议零售价，而应该填写NA。当您所提供的建议零售价有所更新时，您需要确保对建议零售价进行更新。
			// 2. 通过输入建议零售价，您确认：
			//  a. - 您不是该产品在所销售的市场上唯一的卖家（因此在该市场上，建议零售价可以被用作比较价格）；并且
			//  b. - 您有证据表明您提供的建议零售价是该产品真实的一般销售价格，如您的商品在欧盟市场上销售，则该产品必须有欧盟零售商以此价格进行真实的广告宣传和销售，且该建议零售价是经由制造商审慎计算的。当Temu要求的时候，您必须向其提供此类证据。
			// 3. 如果得知或发现建议零售价不符合上述标准，Temu 有权自行决定删除任何建议零售价相关信息。
			ProductSkuSuggestedPriceReq struct {
			} `json:"productSkuSuggestedPriceReq"` // 货品sku建议价格请求
		} `json:"productSkuReqs"` // 货品 sku 列表
	} `json:"productSkcReqs"` // 货品 skc 列表
}

// Create 添加货品
func (s goodsService) Create(ctx context.Context, name string) error {
	return nil
}
