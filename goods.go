package temu

import (
    "context"
    "errors"
    "fmt"
    validation "github.com/go-ozzo/ozzo-validation/v4"
    "github.com/hiscaler/temu-go/entity"
    "github.com/hiscaler/temu-go/helpers"
    "github.com/hiscaler/temu-go/normal"
    "github.com/hiscaler/temu-go/validators/is"
    "gopkg.in/guregu/null.v4"
    "slices"
    "time"
)

// 商品数据服务
type goodsService struct {
    service
    Barcode             goodsBarcodeService             // 条码数据
    Brand               goodsBrandService               // 商品品牌数据
    Category            goodsCategoryService            // 商品分类
    Certification       goodsCertificationService       // 资质
    LifeCycle           goodsLifeCycleService           // 商品生命周期数据
    Sales               goodsSalesService               // 销售数据
    SizeChartClass      goodsSizeChartClassService      // 尺码类
    SizeChart           goodsSizeChartService           // 尺码表
    SizeChartSetting    goodsSizeChartSettingService    // 尺码表设置
    SizeChartTemplate   goodsSizeChartTemplateService   // 尺码表模板
    TopSelling          goodsTopSellingService          // 畅销商品数据
    Warehouse           goodsWarehouseService           // 仓库数据
    Quantity            goodsQuantityService            // 虚拟库存
    ParentSpecification goodsParentSpecificationService // 父规格
    Specification       goodsSpecificationService       // 规格
    Price               goodsPriceService               // 供货价/核价/调价
}

type GoodsQueryParams struct {
    normal.ParameterWithPager
    Page                   int       `json:"page"`                             // 页码
    Cat1Id                 int64     `json:"cat1Id,omitempty"`                 // 一级分类 ID
    Cat2Id                 int64     `json:"cat2Id,omitempty"`                 // 二级分类 ID
    Cat3Id                 int64     `json:"cat3Id,omitempty"`                 // 三级分类 ID
    Cat4Id                 int64     `json:"cat4Id,omitempty"`                 // 四级分类 ID
    Cat5Id                 int64     `json:"cat5Id,omitempty"`                 // 五级分类 ID
    Cat6Id                 int64     `json:"cat6Id,omitempty"`                 // 六级分类 ID
    Cat7Id                 int64     `json:"cat7Id,omitempty"`                 // 七级分类 ID
    Cat8Id                 int64     `json:"cat8Id,omitempty"`                 // 八级分类 ID
    Cat9Id                 int64     `json:"cat9Id,omitempty"`                 // 九级分类 ID
    Cat10Id                int64     `json:"cat10Id,omitempty"`                // 十级分类 ID
    SkcExtCode             string    `json:"skcExtCode,omitempty"`             // 货品 SKC 外部编码
    SupportPersonalization int       `json:"supportPersonalization,omitempty"` // 是否支持定制品模板
    BindSiteIds            []int     `json:"bindSiteIds,omitempty"`            // 经营站点
    ProductName            string    `json:"productName,omitempty"`            // 货品名称
    ProductSkcIds          []int64   `json:"productSkcIds,omitempty"`          // SKC 列表
    SkuExtCodes            []string  `json:"skuExtCodes,omitempty"`            // SKU 货号列表
    QuickSellAgtSignStatus null.Int  `json:"quickSellAgtSignStatus,omitempty"` // 快速售卖协议签署状态 0-未签署 1-已签署
    MatchJitMode           null.Bool `json:"matchJitMode,omitempty"`           // 是否命中 JIT 模式
    SkcSiteStatus          null.Int  `json:"skcSiteStatus,omitempty"`          // skc 加站点状态 (0: 未加入站点, 1: 已加入站点)
    CreatedAtStart         string    `json:"createdAtStart,omitempty"`         // 创建时间开始（年-月-日 时:分:秒）
    CreatedAtEnd           string    `json:"createdAtEnd,omitempty"`           // 创建时间结束（年-月-日 时:分:秒）
}

func (m GoodsQueryParams) validate() error {
    return validation.ValidateStruct(&m,
        validation.Field(&m.BindSiteIds, validation.By(is.SiteIds(entity.SiteIds))),
        validation.Field(&m.CreatedAtStart,
            validation.When(m.CreatedAtStart != "" || m.CreatedAtEnd != "", validation.By(is.TimeRange(m.CreatedAtStart, m.CreatedAtEnd, time.DateTime))),
        ),
        validation.Field(&m.SkcSiteStatus,
            validation.When(m.SkcSiteStatus.Valid, validation.By(func(value interface{}) error {
                v, ok := value.(null.Int)
                if !ok {
                    return errors.New("无效的 SKC 加站点状态")
                }
                return validation.Validate(int(v.Int64), validation.In(entity.TrueNumber, entity.FalseNumber).Error("无效的 SKC 加站点状态"))
            })),
        ),
        validation.Field(&m.QuickSellAgtSignStatus,
            validation.When(m.QuickSellAgtSignStatus.Valid, validation.By(func(value interface{}) error {
                v, ok := value.(null.Int)
                if !ok {
                    return errors.New("无效的快速售卖协议签署状态")
                }
                return validation.Validate(int(v.Int64), validation.In(entity.TrueNumber, entity.FalseNumber).Error("无效的快速售卖协议签署状态"))
            })),
        ),
    )
}

// Query 货品列表查询
// https://seller.kuajingmaihuo.com/sop/view/750197804480663142#SjadVR
func (s goodsService) Query(ctx context.Context, params GoodsQueryParams) (items []entity.Goods, total, totalPages int, isLastPage bool, err error) {
    params.Page = params.TidyPager().Page
    params.OmitPage()
    if err = params.validate(); err != nil {
        err = invalidInput(err)
        return
    }

    if params.CreatedAtStart != "" && params.CreatedAtEnd != "" {
        if start, end, e := helpers.StrTime2UnixMilli(params.CreatedAtStart, params.CreatedAtEnd); e == nil {
            params.CreatedAtStart = start
            params.CreatedAtEnd = end
        }
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

// Detail 货品详情查询（bg.goods.detail.get）
// https://seller.kuajingmaihuo.com/sop/view/750197804480663142#VSGe8J
func (s goodsService) Detail(ctx context.Context, productId int64) (item entity.GoodsDetail, err error) {
    var result = struct {
        normal.Response
        Result entity.GoodsDetail `json:"result"`
    }{}
    resp, err := s.httpClient.R().
        SetContext(ctx).
        SetBody(map[string]int64{"productId": productId}).
        SetResult(&result).
        Post("bg.goods.detail.get")
    if err = recheckError(resp, result.Response, err); err != nil {
        return
    }

    return result.Result, nil
}

// 添加货品

type ProductImageUrl struct {
    ImgUrlList []string `json:"imgUrlList"` // 图片列表
    Language   string   `json:"language"`   // 语言
}

func (m ProductImageUrl) validate() error {
    return validation.ValidateStruct(&m,
        validation.Field(&m.Language, validation.Required.Error("语种不能为空")),
        validation.Field(&m.ImgUrlList,
            validation.Required.Error("图片列表不能为空"),
            validation.Each(validation.By(is.ImageUrl())),
        ),
    )
}

// GoodsCreateProductWarehouse 库存仓库配置对象
type GoodsCreateProductWarehouse struct {
    TargetRouteList []struct {
        SiteIdList  []int  `json:"siteIdList"`  // 站点列表
        WarehouseId string `json:"warehouseId"` // 仓库ID，使用goods.warehouse.list.get查询
    } `json:"targetRouteList"` // 商品的目标路由列表
}

// GoodsCreateProductI18n 多语言标题设置
type GoodsCreateProductI18n struct {
    Language    string `json:"language"`    // 语言编码，en-美国
    ProductName string `json:"productName"` // 对应语言的商品标题
}

func (m GoodsCreateProductI18n) validate() error {
    return validation.ValidateStruct(&m,
        validation.Field(&m.Language,
            validation.Required.Error("语言编码不能为空"),
            validation.In("en").Error("无效的语言编码"),
        ),
        validation.Field(&m.ProductName, validation.Required.Error("商品标题不能为空")),
    )
}

// GoodsCreateProductCarouseVideo 商品主图视频
type GoodsCreateProductCarouseVideo struct {
    Vid      string `json:"vid"`      // 视频 VID
    CoverUrl string `json:"coverUrl"` // 视频封面图
    VideoUrl string `json:"videoUrl"` // 视频 URL
    Width    int    `json:"width"`    // 视频宽度
    Height   int    `json:"height"`   // 视频高度
}

func (m GoodsCreateProductCarouseVideo) validate() error {
    return validation.ValidateStruct(&m,
        validation.Field(&m.Vid, validation.Required.Error("视频 ID 不能为空")),
        validation.Field(&m.CoverUrl,
            validation.Required.Error("视频封面图链接不能为空"),
            validation.By(is.ImageUrl()),
        ),
        validation.Field(&m.VideoUrl, validation.Required.Error("视频 URL 不能为空")),
    )
}

// GoodsCreateProductCustom 货品关务标签
type GoodsCreateProductCustom struct {
    GoodsLabelName   string `json:"goodsLabelName"`   // 货品关务标签名称
    IsRecommendedTag bool   `json:"isRecommendedTag"` // 是否使用推荐标签
}

func (m GoodsCreateProductCustom) validate() error {
    return validation.ValidateStruct(&m,
        validation.Field(&m.GoodsLabelName, validation.Required.Error("货品关务标签名称不能为空")),
    )
}

// GoodsCreateProductOuterPackageImage 外包装图片
type GoodsCreateProductOuterPackageImage struct {
    ImageUrl string `json:"imageUrl"` // 图片链接，通过图片上传接口，imageBizType=1获取
}

func (m GoodsCreateProductOuterPackageImage) validate() error {
    return validation.ValidateStruct(&m,
        validation.Field(&m.ImageUrl, validation.When(m.ImageUrl != "", validation.By(is.ImageUrl()))),
    )
}

// GoodsCreateProductProperty 货品属性
type GoodsCreateProductProperty struct {
    TemplatePid      int64  `json:"templatePid"`      // 模板属性 id
    Pid              int64  `json:"pid"`              // 属性 id
    RefPid           int64  `json:"refPid"`           // 引用属性 id
    PropName         string `json:"propName"`         // 引用属性名
    Vid              int64  `json:"vid"`              // 基础属性值 id，没有的情况传 0
    PropValue        string `json:"propValue"`        // 基础属性值
    ValueUnit        string `json:"valueUnit"`        // 属性值单位，没有的情况传空字符串
    NumberInputValue string `json:"numberInputValue"` // 属性输入值，例如：65.66
    ValueExtendInfo  string `json:"valueExtendInfo"`  // 属性扩展信息，attrs.get 返回
}

func (m GoodsCreateProductProperty) validate(attr entity.GoodsCategoryAttribute) error {
    templatePidIndex := -1
    return validation.ValidateStruct(&m,
        validation.Field(&m.TemplatePid,
            validation.Required.Error("模板属性 ID 不能为空"),
            validation.When(m.TemplatePid != 0, validation.By(func(value interface{}) error {
                id, ok := value.(int64)
                if !ok {
                    return fmt.Errorf("无效的模板属性 ID %v", value)
                }
                index := slices.IndexFunc(attr.Properties, func(prop entity.GoodsCategoryAttributeProperty) bool {
                    return id == prop.TemplatePid
                })
                if index == -1 {
                    return fmt.Errorf("模板属性 ID %d 在类目属性中不存在", id)
                }
                templatePidIndex = index
                return nil
            })),
        ),
        validation.Field(&m.Pid,
            validation.Required.Error("属性 ID 不能为空"),
            validation.When(m.Pid != 0, validation.By(func(value interface{}) error {
                id, ok := value.(int64)
                if !ok {
                    return fmt.Errorf("无效的基础属性 ID %v", value)
                }
                index := slices.IndexFunc(attr.Properties, func(prop entity.GoodsCategoryAttributeProperty) bool {
                    return id == prop.Pid
                })
                if index == -1 {
                    return fmt.Errorf("基础属性 ID %d 在类目属性中不存在", id)
                }
                return nil
            })),
        ),
        validation.Field(&m.RefPid,
            validation.Required.Error("引用属性 ID 不能为空"),
            validation.When(m.RefPid != 0, validation.By(func(value interface{}) error {
                id, ok := value.(int64)
                if !ok {
                    return fmt.Errorf("无效的引用属性 ID %v", value)
                }
                index := slices.IndexFunc(attr.Properties, func(prop entity.GoodsCategoryAttributeProperty) bool {
                    return id == prop.RefPid
                })
                if index == -1 {
                    return fmt.Errorf("引用属性 ID %d 在类目属性中不存在", id)
                }
                return nil
            })),
        ),
        validation.Field(&m.PropName,
            validation.Required.Error("引用属性名不能为空"),
            validation.When(m.PropName != "", validation.By(func(value interface{}) error {
                name, ok := value.(string)
                if !ok {
                    return fmt.Errorf("无效的引用属性名 %v", value)
                }
                index := slices.IndexFunc(attr.Properties, func(prop entity.GoodsCategoryAttributeProperty) bool {
                    return name == prop.Name
                })
                if index == -1 {
                    return fmt.Errorf("引用属性名 %s 在类目属性中不存在", name)
                }
                return nil
            })),
        ),
        validation.Field(&m.PropValue,
            validation.Required.Error("基础属性值不能为空"),
            validation.By(func(value interface{}) error {
                v, ok := value.(string)
                if !ok {
                    return fmt.Errorf("无效的基础属性值 %v", value)
                }
                index := slices.IndexFunc(attr.Properties[templatePidIndex].Values, func(e entity.GoodsCategoryAttributePropertyValue) bool {
                    return m.PropValue == e.Value
                })
                if index == -1 {
                    return fmt.Errorf("无效的基础属性值 %s", v)
                }
                return nil
            }),
        ),
        validation.Field(&m.ValueUnit,
            validation.By(func(value interface{}) error {
                v, ok := value.(string)
                if !ok {
                    return fmt.Errorf("无效的属性值单位 %v", value)
                }
                index := slices.Index(attr.Properties[templatePidIndex].ValueUnit, m.ValueUnit)
                if index == -1 {
                    return fmt.Errorf("无效的属性值单位 %s", v)
                }
                return nil
            }),
        ),
    )
}

// GoodsCreateProductSpecProperty 货品规格属性
type GoodsCreateProductSpecProperty struct {
    TemplatePid      int64  `json:"templatePid"`      // 模板属性id
    Pid              int64  `json:"pid"`              // 属性 id
    RefPid           int64  `json:"refPid"`           // 引用属性 id
    PropName         string `json:"propName"`         // 引用属性名
    Vid              int64  `json:"vid"`              // 基础属性值id，没有的情况传0
    PropValue        string `json:"propValue"`        // 基础属性值
    ParentSpecId     int    `json:"parentSpecId"`     // 父规格id
    ParentSpecName   string `json:"parentSpecName"`   // 父规格名称
    SpecId           int    `json:"specId"`           // 规格id
    SpecName         string `json:"specName"`         // 规格名称
    ValueGroupId     int    `json:"valueGroupId"`     // 属性值组id，没有的情况传0
    ValueGroupName   string `json:"valueGroupName"`   // 属性值组名称，没有的情况传空字符串
    ValueUnit        string `json:"valueUnit"`        // 属性值单位，没有的情况传空字符串
    NumberInputValue string `json:"numberInputValue"` // 属性输入值，例如：65.66
    ValueExtendInfo  string `json:"valueExtendInfo"`  // 属性组扩展信息（色板）
}

// GoodsCreateProductSaleExtAttr 货品销售类扩展属性请求
type GoodsCreateProductSaleExtAttr struct {
    ProductSecondHandReq struct {
        IsSecondHand    bool `json:"isSecondHand"`    // 是否二手货品，二手店铺传true，其他店铺不传值
        SecondHandLevel int  `json:"secondHandLevel"` // 成色定义，二手货品必传值，非二手货品不可传值，枚举值：（1：接近全新，2：状况极佳，3：状况良好，4：尚可接受）
    } `json:"productSecondHandReq"`              // 货品二手信息，二手店铺传值，其他店铺不传值
    InventoryRegion int `json:"inventoryRegion"` // 备货区域
    // 满足如下两个条件时，必填无充电器版本货品id
    // 引用属性：refPid=6919, propName="是否售卖不含充电器的同款商品"
    // 属性值：vid=147374, value="是"
    ProductNoChargerReq struct {
        NoChargerProductIds []int `json:"noChargerProductIds"` // 无充电器版本spuid，至少入参1个，至多入参3个
    } `json:"productNoChargerReq"` // 无充电器版本spuid请求
}

func (m GoodsCreateProductSaleExtAttr) validate() error {
    return nil
}

// GoodsCreateProductWhExtAttr 货品仓配供应链侧扩展属性请求
type GoodsCreateProductWhExtAttr struct {
    OuterGoodsUrl string `json:"outerGoodsUrl"` //  站外商品链接
    ProductOrigin struct {
        Region1ShortName string `json:"region1ShortName"` // 一级区域简称 (二字简码)
        // 枚举值：https://seller.kuajingmaihuo.com/sop/view/735407268315104853
        Region2Id string `json:"region2Id"` // 二级区域地址ID，当region1ShortName为CN时，必传
    } `json:"productOrigin"` // 货品产地 (灰度内必传)，请注意，日本站点发品必须传产地，否则会被拦截
}

// GoodsCreateProductSkuSiteSupplierPrice 站点供货价列表，半托必传
type GoodsCreateProductSkuSiteSupplierPrice struct {
    SiteId        int64 `json:"siteId"`        // 申报价格站点id
    SupplierPrice int64 `json:"supplierPrice"` // 站点申报价格，单位 人民币：分，美元：美分
}

// GoodsCreateProductSkuProductSkuStockQuantity 货品 sku 仓库库存，半托管发品必传
type GoodsCreateProductSkuProductSkuStockQuantity struct {
    WarehouseStockQuantityReqs []struct {
        TargetStockAvailable int    `json:"targetStockAvailable"` // sku目标库存值（覆盖写）
        WarehouseId          string `json:"warehouseId"`          // 仓库 ID
    } `json:"warehouseStockQuantityReqs"` // 仓库存库存请求列表
}

// GoodsCreateProductSkuProductSkuMultiPack 货品多包规请求
type GoodsCreateProductSkuProductSkuMultiPack struct {
    NumberOfPieces          int `json:"numberOfPieces"` // 件数，单品默认是1
    ProductSkuNetContentReq struct {
        NetContentUnitCode int `json:"netContentUnitCode"` // 净含量单位，1：液体盎司，2：毫升，3：加仑，4：升，5：克，6：千克，7：常衡盎司，8：磅
        NetContentNumber   int `json:"netContentNumber"`   // 净含量数值
    } `json:"productSkuNetContentReq"`               // 净含量请求，传空对象表示空，指定类目灰度管控
    SkuClassification int `json:"skuClassification"` // sku分类，1：单品，2：组合装，3：混合套装
    PieceUnitCode     int `json:"pieceUnitCode"`     // 单件单位，1：件，2：双，3：包
}

// GoodsCreateProductSkuSuggestedPrice 货品sku建议价格
type GoodsCreateProductSkuSuggestedPrice struct {
    // 特殊建议价格，用来标记商家是否有建议价格，传参规则如下：
    // - 传参为NA，则认为商家没有货品建议价格，即suggestedPrice和suggestedPriceCurrencyType这两个字段都不需要传；
    // - 不传该字段，则要求suggestedPrice和suggestedPriceCurrencyType字段必传，不传则会报错；
    // 示例：
    // - productSkuSuggestedPriceReq
    // -specialSuggestedPrice：NA
    // -----------------------------------------
    // - productSkuSuggestedPriceReq
    // -suggestedPriceCurrencyType：CNY
    // -suggestedPrice：10
    SpecialSuggestedPrice      string `json:"specialSuggestedPrice"`                //  特殊建议价格，用来标记商家是否有建议价格
    SuggestedPriceCurrencyType string `json:"suggestedPriceCurrencyType,omitempty"` // 建议价格币种（USD:美元,CNY:人民币,JPY:日元,CAD:加拿大元,GBP:英镑,AUD:澳大利亚元,NZD:新西兰元,EUR:欧元,MXN:墨西哥比索,PLN:波兰兹罗提,SEK:瑞典克朗,CHF:瑞士法郎,KRW:韩元,SAR:沙特里亚尔,SGD:新加坡元,AED:阿联酋迪拉姆,KWD:科威特第纳尔,NOK:挪威克朗,CLP:智利比索,MYR:马来西亚林吉特,PHP:菲律宾比索,TWD:新台湾元,THB:泰铢,QAR:卡塔尔里亚尔,JOD:约旦第纳尔,BRL:巴西雷亚尔,OMR:阿曼里亚尔,BHD:巴林第纳尔,ILS:以色列新锡克尔,ZAR:南非兰特,CZK:捷克克朗,HUF:匈牙利福林,DKK:丹麦克朗,RON:罗马尼亚列伊,BGN:保加利亚列瓦,HKD:港元,COP:哥伦比亚比索,GEL:格鲁吉亚拉里）
    // 建议价格，币种枚举值：
    // 备注：辅助单位分别为0、1、2、3分别对应前端录入信息时需要原值上
    // ×1、10、100、1000，再把转换后的数据传给后端
    //export declare enum Currency {
    //    /** 美元，辅助单位为 2 */
    //    USD = "USD",
    //    /** 人民币，辅助单位为 2 */
    //    CNY = "CNY",
    //    /** 日元，辅助单位为 0 */
    //    JPY = "JPY",
    //    /** 加拿大元, 辅助单位为 2 */
    //    CAD = "CAD",
    //    /** 英镑，辅助单位为 2 */
    //    GBP = "GBP",
    //    /** 澳大利亚，辅助单位为 2 */
    //    AUD = "AUD",
    //    /** 新西兰，辅助单位为 2 */
    //    NZD = "NZD",
    //    /**
    //     * 欧盟地区，统一用欧元，辅助单位为 2
    //     * 欧元区： 亚克罗提利与德凯利亚、 安道尔（AD）、 奥地利（AT）、 比利时（BE）、 赛普勒斯（CY）、 爱沙尼亚（EE）、
    //     * 芬兰（FI）、 法国（FR）、 德国（DE）、 希腊（GR）、 瓜德罗普（GP）、 爱尔兰（IE）、
    //     * 义大利（IT）、 科索沃、 拉脱维亚（LV）、 立陶宛（LT）、 卢森堡（LU）、 马尔他（MT）、 马提尼克（MQ）、
    //     * 马约特（YT）、 摩纳哥（MC）、 蒙特内哥罗（ME）、 荷兰（NL）、 葡萄牙（PT）、 留尼汪（RE）、 圣巴泰勒米（BL）、
    //     * 圣皮埃尔和密克隆（PM）、 圣马力诺（SM）、 斯洛伐克（SK）、 斯洛维尼亚（SI）、 西班牙（ES）、 梵蒂冈（VA）;
    //     */
    //    EUR = "EUR",
    //    /** 墨西哥，辅助单位为 2 */
    //    MXN = "MXN",
    //    /** 波兰，辅助单位为 2 */
    //    PLN = "PLN",
    //    /** 瑞典，辅助单位为 2 */
    //    SEK = "SEK",
    //    /** 瑞士，辅助单位为 2 */
    //    CHF = "CHF",
    //    /** 韩元，辅助单位为 0 */
    //    KRW = "KRW",
    //    /** 沙特, 辅助单位为 2 */
    //    SAR = "SAR",
    //    /** 新加坡, 辅助单位为 2 */
    //    SGD = "SGD",
    //    /** 阿联酋, 辅助单位为 2 */
    //    AED = "AED",
    //    /** 科威特，辅助单位为 3 */
    //    KWD = "KWD",
    //    /** 挪威, 辅助单位为 2 */
    //    NOK = "NOK",
    //    /** 智利, 辅助单位为 0 */
    //    CLP = "CLP",
    //    /** 马来西亚, 辅助单位为 2 */
    //    MYR = "MYR",
    //    /** 菲律宾, 辅助单位为 2 */
    //    PHP = "PHP",
    //    /** 台湾, 辅助单位为 2 */
    //    TWD = "TWD",
    //    /** 泰国, 辅助单位为 2 */
    //    THB = "THB",
    //    /** 卡塔尔, 辅助单位为 2 */
    //    QAR = "QAR",
    //    /** 约旦, 辅助单位为 3 */
    //    JOD = "JOD",
    //    /** 巴西, 辅助单位为 2 */
    //    BRL = "BRL",
    //    /** 阿曼, 辅助单位为 3 */
    //    OMR = "OMR",
    //    /** 巴林, 辅助单位为 3 */
    //    BHD = "BHD",
    //    /** 以色列, 辅助单位为 2 */
    //    ILS = "ILS",
    //    /** 南非, 辅助单位为 2 */
    //    ZAR = "ZAR",
    //    /** 捷克, 辅助单位为 2，但是输入的时候不能输入小数需特殊处理 */
    //    CZK = "CZK",
    //    /** 匈牙利, 辅助单位为 2，但是输入的时候不能输入小数需特殊处理 */
    //    HUF = "HUF",
    //    /** 丹麦, 辅助单位为 2 */
    //    DKK = "DKK",
    //    /** 罗马尼亚, 辅助单位为 2 */
    //    RON = "RON",
    //    /** 保加利亚, 辅助单位为 2 */
    //    BGN = "BGN",
    //    /** 香港, 辅助单位为 2 */
    //    HKD = "HKD",
    //    /** 哥伦比亚, 辅助单位为 2 */
    //    COP = "COP",
    //    /** 格鲁吉亚拉里, 辅助单位为 2 */
    //    GEL = "GEL"
    // }
    SuggestedPrice int `json:"suggestedPrice,omitempty"` // 建议价格
}

// GoodsCreateProductSkuWhExtAttrSensitiveLimitProductSkuWeight 货品sku重量
type GoodsCreateProductSkuWhExtAttrSensitiveLimitProductSkuWeight struct {
    Value int `json:"value"` // 重量值，单位mg
}

// GoodsCreateProductSkuWhExtAttrSameReferPrice 同款参考
type GoodsCreateProductSkuWhExtAttrSameReferPrice struct {
    Url string `json:"url"` // 站外同款商品售卖链接，有效链接规则，链接开头含：http:// 、 https:// 等
}

// GoodsCreateProductSkuWhExtAttrSensitiveLimit 货品sku敏感属性限制请求
type GoodsCreateProductSkuWhExtAttrSensitiveLimit struct {
    MaxBatteryCapacity   int `json:"maxBatteryCapacity"`   // 最大电池容量 (Wh)
    MaxBatteryCapacityHp int `json:"maxBatteryCapacityHp"` // 最大电池容量 (mWh)
    MaxLiquidCapacity    int `json:"maxLiquidCapacity"`    // 最大液体容量 (mL)
    MaxLiquidCapacityHp  int `json:"maxLiquidCapacityHp"`  // 最大液体容量 (μL)
    MaxKnifeLength       int `json:"maxKnifeLength"`       // 最大刀具长度 (mm)
    MaxKnifeLengthHp     int `json:"maxKnifeLengthHp"`     // 最大刀具长度 (μm)
    KnifeTipAngle        struct {
        Degrees int `json:"degrees"` //	度数
    } `json:"knifeTipAngle"` // 刀尖角度
}

// GoodsCreateProductSkuWhExtAttrVolume 货品sku体积
type GoodsCreateProductSkuWhExtAttrVolume struct {
    Len    int `json:"len"`    // 长，单位mm
    Width  int `json:"width"`  // 宽，单位mm
    Height int `json:"height"` // 高，单位mm
}

// GoodsCreateProductSkuWhExtAttrSensitiveAttr 货品 sku 敏感属性请求
type GoodsCreateProductSkuWhExtAttrSensitiveAttr struct {
    IsSensitive   int   `json:"isSensitive"`   // 是否敏感属性，0：非敏感，1：敏感
    SensitiveList []int `json:"sensitiveList"` // 敏感类型，        PURE_ELECTRIC(110001, "纯电"),    INTERNAL_ELECTRIC(120001, "内电"),    MAGNETISM(130001, "磁性"),    LIQUID(140001, "液体"),    POWDER(150001, "粉末"),    PASTE(160001, "膏体"),    CUTTER(170001, "刀具")
}

func (m GoodsCreateProductSkuWhExtAttrSensitiveAttr) validate() error {
    return validation.ValidateStruct(&m,
        validation.Field(&m.SensitiveList,
            validation.When(
                m.IsSensitive == 1,
                validation.Required.Error("敏感类型不能为空"),
                validation.In(110001, 120001, 130001, 140001, 150001, 160001, 170001).Error("无效的敏感类型"),
            ),
        ),
    )
}

type GoodsCreateProductSkuWhExtAttrBarCode struct {
    Code     string `json:"code"`     // 商品标准编码
    CodeType int    `json:"codeType"` // 条码类型 (1: EAN, 2: UPC, 3: ISBN)
}

func (m GoodsCreateProductSkuWhExtAttrBarCode) validate() error {
    return validation.ValidateStruct(&m,
        validation.Field(&m.Code, validation.Required.Error("商品标准编码不能为空")),
        validation.Field(&m.CodeType,
            validation.In(1, 2, 3).Error("无效的条码类型"),
        ),
    )
}

// GoodsCreateProductSkuWhExtAttr 货品sku扩展属性
type GoodsCreateProductSkuWhExtAttr struct {
    ProductSkuWeightReq         GoodsCreateProductSkuWhExtAttrSensitiveLimitProductSkuWeight `json:"productSkuWeightReq"`                   // 货品sku重量
    ProductSkuSameReferPriceReq GoodsCreateProductSkuWhExtAttrSameReferPrice                 `json:"productSkuSameReferPriceReq"`           // 同款参考
    ProductSkuSensitiveLimitReq *GoodsCreateProductSkuWhExtAttrSensitiveLimit                `json:"productSkuSensitiveLimitReq,omitempty"` // 货品sku敏感属性限制请求
    ProductSkuVolumeReq         GoodsCreateProductSkuWhExtAttrVolume                         `json:"productSkuVolumeReq"`                   // 货品sku体积
    ProductSkuSensitiveAttrReq  GoodsCreateProductSkuWhExtAttrSensitiveAttr                  `json:"productSkuSensitiveAttrReq"`            // 货品 sku 敏感属性请求
    ProductSkuBarCodeReqs       []GoodsCreateProductSkuWhExtAttrBarCode                      `json:"productSkuBarCodeReqs,omitempty"`
}

// GoodsCreateProductSku 货品 SKC 下的 SKU 信息
type GoodsCreateProductSku struct {
    ThumbUrl                   string                                        `json:"thumbUrl"`                             // 预览图
    ProductSkuThumbUrlI18nReqs []ProductImageUrl                             `json:"productSkuThumbUrlI18nReqs"`           // SKU多语言预览图，服饰类不传，非服饰非必传 （英国英语、中东英语必传）
    CurrencyType               string                                        `json:"currencyType"`                         // 币种 (CNY: 人民币, USD: 美元) (默认人民币)
    ProductSkuSpecReqs         []entity.Specification                        `json:"productSkuSpecReqs"`                   // 货品sku规格列表
    SupplierPrice              int64                                         `json:"supplierPrice"`                        // 全托供货价 （单位：人民币-分/美元-美分），半托不传
    SiteSupplierPrices         []GoodsCreateProductSkuSiteSupplierPrice      `json:"siteSupplierPrices"`                   // 站点供货价列表，半托必传
    ProductSkuStockQuantityReq *GoodsCreateProductSkuProductSkuStockQuantity `json:"productSkuStockQuantityReq,omitempty"` // 货品sku仓库库存，半托管发品必传
    ProductSkuMultiPackReq     *GoodsCreateProductSkuProductSkuMultiPack     `json:"productSkuMultiPackReq,omitempty"`     // 货品多包规请求
    // 货品sku建议价格请求
    // 1. 建议零售价是制造商为产品设定的建议零售价或推荐零售价。建议零售价必须是市场上的真实销售价格，且符合任何可适用的法律法规的规定。如您的商品在欧盟市场上销售，则该产品必须有欧盟零售商以此价格进行真实的广告宣传和销售。如果您的产品没有符合这些标准的建议零售价，请勿填写建议零售价，而应该填写NA。当您所提供的建议零售价有所更新时，您需要确保对建议零售价进行更新。
    // 2. 通过输入建议零售价，您确认：
    //  a. - 您不是该产品在所销售的市场上唯一的卖家（因此在该市场上，建议零售价可以被用作比较价格）；并且
    //  b. - 您有证据表明您提供的建议零售价是该产品真实的一般销售价格，如您的商品在欧盟市场上销售，则该产品必须有欧盟零售商以此价格进行真实的广告宣传和销售，且该建议零售价是经由制造商审慎计算的。当Temu要求的时候，您必须向其提供此类证据。
    // 3. 如果得知或发现建议零售价不符合上述标准，Temu 有权自行决定删除任何建议零售价相关信息。
    ProductSkuSuggestedPriceReq *GoodsCreateProductSkuSuggestedPrice `json:"productSkuSuggestedPriceReq,omitempty"` // 货品sku建议价格请求
    ProductSkuWhExtAttrReq      *GoodsCreateProductSkuWhExtAttr      `json:"productSkuWhExtAttrReq,omitempty"`      // 货品sku扩展属性
    ExtCode                     string                               `json:"extCode"`                               // 货品 skc 外部编码，没有的场景传空字符串
}

// GoodsCreateProductSkc 货品 SKC
type GoodsCreateProductSkc struct {
    PreviewImgUrls                  []string                `json:"previewImgUrls"`                  // SKC 轮播图列表
    ProductSkcCarouselImageI18nReqs []ProductImageUrl       `json:"productSkcCarouselImageI18nReqs"` // SKC多语言轮播图，服饰类必传，非服饰不传
    ColorImageUrl                   string                  `json:"colorImageUrl"`                   // SKC 色块图，可通过（bg.colorimageurl.get）转换获取
    MainProductSkuSpecReqs          []entity.Specification  `json:"mainProductSkuSpecReqs"`          // 主销售规格列表
    IsBasePlate                     int                     `json:"isBasePlate"`                     // 是否底板
    ProductSkuReqs                  []GoodsCreateProductSku `json:"productSkuReqs"`                  // 货品 sku 列表
    ExtCode                         string                  `json:"extCode"`                         //
}

// GoodsCreateModel 商品模特
type GoodsCreateModel struct {
    ModelProfileUrl string `json:"modelProfileUrl"` // 模特头像
    SizeSpecName    string `json:"sizeSpecName"`    // 试穿尺码规格名称
    ModelId         int    `json:"modelId"`         // 模特id，通过模特信息查询接口获取
    SizeSpecId      int    `json:"sizeSpecId"`      // 试穿尺码规格id
    ModelWaist      string `json:"modelWaist"`      // 模特腰围文本, modelType=2传空值
    ModelType       int    `json:"modelType"`       // 模特类型，1：成衣模特，2：鞋模
    ModelName       string `json:"modelName"`       // 模特名称
    ModelHeight     string `json:"modelHeight"`     // 模特身高文本modelType=2传空值
    ModelFeature    int    `json:"modelFeature"`    // 模特特性，1：真实模特
    ModelFootWidth  string `json:"modelFootWidth"`  // 模特脚宽文本modelType=1传空值
    ModelBust       string `json:"modelBust"`       // 模特胸围文本modelType=2传空值
    ModelFootLength string `json:"modelFootLength"` // 模特脚长文本modelType=1传空值
    TryOnResult     int    `json:"tryOnResult"`     // 试穿心得 TRUE_TO_SIZE(1, "舒适"),    TOO_SMALL(2, "紧身"),    TOO_LARGE(3, "宽松"),
    ModelHip        string `json:"modelHip"`        // 模特臀围文本modelType=2传空值
}

// GoodsCreateProductOuterPackage 货品外包装信息
type GoodsCreateProductOuterPackage struct {
    PackageShape int `json:"packageShape"` // 外包装形状0:不规则形状 1:长方体 2:圆柱体
    PackageType  int `json:"packageType"`  // 外包装类型0:硬包装 1:软包装+硬物 2:软包装+软物
}

// GoodsCreateProductGuideFile 说明书文件
type GoodsCreateProductGuideFile struct {
    FileName      string   `json:"fileName"`      // 文件名称
    PdfMaterialId int      `json:"pdfMaterialId"` // pdf文件id，通过file.upload上传返回得到
    Languages     []string `json:"languages"`     // 语言（zh-中文、en-英文）
}

// GoodsCreateGoodsLayerDecoration 商详装饰
type GoodsCreateGoodsLayerDecoration struct {
    FloorId     null.Int `json:"floorId"`  // 楼层id,null:新增,否则为更新
    GoodsId     int64    `json:"goodsId"`  // 商品 ID
    Lang        string   `json:"lang"`     // 语言类型
    Type        string   `json:"type"`     // 组件类型type,图片-image,文本-text 商详需要包含至少一个图片类型组件
    Priority    int      `json:"priority"` // 楼层排序
    Key         string   `json:"key"`      // 楼层类型的key,目前默认传'DecImage'
    ContentList []struct {
        ImgUrl            string `json:"imgUrl"` // 图片地址--通用，图片最大3M
        Width             int    `json:"width"`  // 图片宽度--通用，宽度最小480px
        Text              string `json:"text"`   // 文字信息--文字模块，文本-text必填，长度限制500字符内
        Height            int    `json:"height"` // 图片高度--通用，高度最小480px
        TextModuleDetails struct {
            BackgroundColor string `json:"backgroundColor"` // 背景颜色文本-text必填，六位值，例#ffffff
            FontFamily      int    `json:"fontFamily"`      // 字体类型文本-text不传
            FontSize        int    `json:"fontSize"`        // 文字模块字体大小文本-text必传12
            Align           string `json:"align"`           // 文字对齐方式，left--左对齐；right--右对齐；center--居中；justify--两端对齐文本-text必填
            FontColor       string `json:"fontColor"`       // 文字颜色文本-text必填，六位值，例#333333
        } `json:"textModuleDetails"` // 文字模块详情文本-text必填
    } `json:"contentList"` // 楼层内容
}

// GoodsCreateProductSemiManaged 半托管相关信息
type GoodsCreateProductSemiManaged struct {
    BindSiteIds []int `json:"bindSiteIds"` // 绑定站点列表
    // 半托管-素材语种策略，不传默认2
    // 1：仅站点本地语种素材，允许只上传站点本地语种的素材（多语言素材节点上传本地素材，英语素材也可使用本地语种素材填充）
    // 关联节点如下：
    // - 多语言标题（productI18nReqs）
    // - 多语言素材（materialMultiLanguages、carouselImageI18nReqs、productSkcCarouselImageI18nReqs、productSkuThumbUrlI18nReqs）
    //
    // 当前支持日本站、墨西哥站，使用语言如下：
    // 日本站：多语言标题和素材均使用ja
    // 墨西哥站：多语言标题语言传es-MX。多语言素材语言传es
    //
    // 2：英语以及其他语种
    SemiLanguageStrategy int `json:"emiLanguageStrategy"` // 半托管-素材语种策略
}

// GoodsCreateProductShipment 半托管货品配送信息请求
type GoodsCreateProductShipment struct {
    FreightTemplateId   string `json:"freightTemplateId"`   // 运费模板 id，使用bg.logistics.template.get查询，详见：https://seller.kuajingmaihuo.com/sop/view/867739977041685428#pa858C
    ShipmentLimitSecond int    `json:"shipmentLimitSecond"` // 承诺发货时间(单位:s)，可选值：86400，172800，259200（仅定制品可用）
}

func (m GoodsCreateProductShipment) validate() error {
    return validation.ValidateStruct(&m,
        validation.Field(&m.FreightTemplateId, validation.Required.Error("运费模板 ID 不能为空")),
        validation.Field(&m.ShipmentLimitSecond,
            validation.Required.Error("承诺发货时间不能为空"),
            validation.In(86400, 172800, 259200).Error("无效的承诺发货时间"),
        ),
    )
}

type GoodsCreateRequest struct {
    Cat1Id                       int64                                 `json:"cat1Id"`                             // 一级类目id
    Cat2Id                       int64                                 `json:"cat2Id"`                             // 二级类目id
    Cat3Id                       int64                                 `json:"cat3Id"`                             // 三级类目id
    Cat4Id                       int64                                 `json:"cat4Id"`                             // 四级类目id（没有的情况传 0）
    Cat5Id                       int64                                 `json:"cat5Id"`                             // 五级类目id（没有的情况传 0）
    Cat6Id                       int64                                 `json:"cat6Id"`                             // 六级类目id（没有的情况传 0）
    Cat7Id                       int64                                 `json:"cat7Id"`                             // 七级类目id（没有的情况传 0）
    Cat8Id                       int64                                 `json:"cat8Id"`                             // 八级类目id（没有的情况传 0）
    Cat9Id                       int64                                 `json:"cat9Id"`                             // 九级类目id（没有的情况传 0）
    Cat10Id                      int64                                 `json:"cat10Id"`                            // 十级类目id（没有的情况传 0）
    ProductWarehouseRouteReq     *GoodsCreateProductWarehouse          `json:"productWarehouseRouteReq,omitempty"` // 库存仓库配置对象。	半托管发品必传，全托管店铺不需要传这个属性，传入会报错。
    ProductI18nReqs              []GoodsCreateProductI18n              `json:"productI18nReqs,omitempty"`          // 多语言标题设置
    ProductName                  string                                `json:"productName"`                        // 货品名称
    ProductCarouseVideoReqList   []GoodsCreateProductCarouseVideo      `json:"productCarouseVideoReqList"`         // 商品主图视频 关于如何上传视频，请对接视频上传相关接口，获取图片相关参数即可用于此处入参 https://seller.kuajingmaihuo.com/sop/view/852640595329867111
    ProductCustomReq             *GoodsCreateProductCustom             `json:"productCustomReq,omitempty"`         // 货品关务标签
    CarouselImageUrls            []string                              `json:"carouselImageUrls"`                  // 货品轮播图
    CarouselImageI18nReqs        []ProductImageUrl                     `json:"carouselImageI18nReqs"`              // 货品 SPU 多语言轮播图，服饰类不传，非服饰必传
    ProductOuterPackageImageReqs []GoodsCreateProductOuterPackageImage `json:"productOuterPackageImageReqs"`       // 外包装图片
    MaterialImgUrl               string                                `json:"materialImgUrl"`                     // 素材图
    ProductPropertyReqs          []GoodsCreateProductProperty          `json:"productPropertyReqs"`                // 货品属性
    ProductSpecPropertyReqs      []GoodsCreateProductSpecProperty      `json:"productSpecPropertyReqs"`
    // 货品规格属性
    ProductSaleExtAttrReq    *GoodsCreateProductSaleExtAttr    `json:"productSaleExtAttrReq,omitempty"`  // 货品销售类扩展属性请求
    ProductWhExtAttrReq      GoodsCreateProductWhExtAttr       `json:"productWhExtAttrReq"`              // 货品仓配供应链侧扩展属性请求
    ProductSkcReqs           []GoodsCreateProductSkc           `json:"productSkcReqs"`                   // 货品 skc 列表
    SizeTemplateIds          []int                             `json:"sizeTemplateIds"`                  // 尺码表模板id（从sizecharts.template.create获取），无尺码表时传空数组[]
    GoodsModelReqs           []GoodsCreateModel                `json:"goodsModelReqs"`                   // 商品模特列表请求
    ShowSizeTemplateIds      []int64                           `json:"showSizeTemplateIds"`              // 套装尺码表展示，至多2个尺码表模板id入参
    ProductOuterPackageReq   *GoodsCreateProductOuterPackage   `json:"productOuterPackageReq,omitempty"` // 货品外包装信息
    ProductGuideFileReqs     []GoodsCreateProductGuideFile     `json:"productGuideFileReqs"`             // 说明书请求对象
    GoodsLayerDecorationReqs []GoodsCreateGoodsLayerDecoration `json:"goodsLayerDecorationReqs"`         // 商详装饰
    PersonalizationSwitch    int                               `json:"personalizationSwitch"`            // 是否定制品，API发品标记定制品后，请及时在卖家中心配置定制模版信息，否则无法正常加站点售卖 0：非定制品、1：定制品
    ProductSemiManagedReq    *GoodsCreateProductSemiManaged    `json:"productSemiManagedReq,omitempty"`  // 半托管相关信息
    ProductShipmentReq       *GoodsCreateProductShipment       `json:"productShipmentReq,omitempty"`     // 半托管货品配送信息请求
    AddProductChannelType    int                               `json:"addProductChannelType"`            // 发品渠道
    MaterialMultiLanguages   []string                          `json:"materialMultiLanguages"`           // 图片多语言列表
}

func (m GoodsCreateRequest) validate(ctx context.Context, s goodsService) error {
    return validation.ValidateStruct(&m,
        validation.Field(&m.Cat1Id, validation.Required.Error("一级类目不能为空")),
        validation.Field(&m.Cat2Id, validation.Required.Error("二级类目不能为空")),
        validation.Field(&m.Cat3Id, validation.Required.Error("三级类目不能为空")),
        validation.Field(&m.ProductName,
            validation.Required.Error("商品名称不能为空"),
            validation.Length(1, 60).Error("商品名称最多 {{.max}} 个字符"),
        ),
        validation.Field(&m.CarouselImageUrls,
            validation.Required.Error("货品轮播图不能为空"),
            validation.Min(5).Error("货品轮播图不能少于 {{.min}} 张"),
            validation.Each(validation.By(is.ImageUrl())),
        ),
        validation.Field(&m.MaterialImgUrl, validation.When(m.MaterialImgUrl != "", validation.By(is.ImageUrl()))),
        validation.Field(&m.ProductPropertyReqs,
            validation.Required.Error("货品属性不能为空"),
            // 判断是否有重复的数据
            validation.By(func(value interface{}) error {
                properties, ok := value.([]GoodsCreateProductProperty)
                if !ok {
                    return errors.New("无效的货品属性。")
                }
                templatePids := make([]int64, 0)
                for _, prop := range properties {
                    if slices.Index(templatePids, prop.TemplatePid) != -1 {
                        return fmt.Errorf("重复的模板属性 ID %d", prop.TemplatePid)
                    }
                    templatePids = append(templatePids, prop.TemplatePid)
                }

                return nil
            }),
            validation.Each(validation.By(func(value interface{}) error {
                v, ok := value.(GoodsCreateProductProperty)
                if !ok {
                    return errors.New("无效的货品属性。")
                }

                catIds := []int64{m.Cat10Id, m.Cat9Id, m.Cat8Id, m.Cat7Id, m.Cat6Id, m.Cat5Id, m.Cat4Id, m.Cat3Id, m.Cat2Id, m.Cat1Id}
                var catId int64 = -1
                if index := slices.IndexFunc(catIds, func(i int64) bool {
                    return i != 0
                }); index != -1 {
                    catId = catIds[index]
                }
                attr, err := s.Category.Attribute.Query(ctx, catId)
                if err != nil {
                    return nil
                }
                if attr == nil {
                    return fmt.Errorf("%d 类目属性查询失败", catId)
                }
                return v.validate(*attr)
            })),
        ),
        validation.Field(&m.AddProductChannelType, validation.Required.Error("发品渠道不能为空")),
    )
}

type GoodsCreateResult struct {
    ProductId      int64 `json:"productId"` // 货品 id
    ProductSkcList []struct {
        ProductSkcId int64 `json:"productSkcId"` // skc id
    } `json:"productSkcList"` //  skc列表
    ProductSkuList []struct {
        ProductSkcId int64                  `json:"productSkcId "` // SKC ID
        ProductSkuId int64                  `json:"productSkuId"`  // sku id
        ExtCode      string                 `json:"extCode"`       // sku 外部编码
        SkuSpecList  []entity.Specification `json:"skuSpecList"`   // sku 规格
    } `json:"productSkuList"` // sku 列表
}

// Create 添加货品
// https://seller.kuajingmaihuo.com/sop/view/750197804480663142#MwT6Ha
func (s goodsService) Create(ctx context.Context, request GoodsCreateRequest) (res GoodsCreateResult, err error) {
    if err = request.validate(ctx, s); err != nil {
        return res, invalidInput(err)
    }

    var result = struct {
        normal.Response
        Result GoodsCreateResult `json:"result"`
    }{}
    resp, err := s.httpClient.R().
        SetContext(ctx).
        SetBody(request).
        SetResult(&result).
        Post("bg.goods.add")
    if err = recheckError(resp, result.Response, err); err != nil {
        return
    }

    return result.Result, nil
}

// ImageUpload 上传货品图片（bg.goods.image.upload）
// https://seller.kuajingmaihuo.com/sop/view/338873192956832611#atWm1f

type GoodsImageUploadRequest struct {
    Image        string   `json:"image"`                  // 支持格式有：jpg/jpeg、png等图片格式，注意入参图片必须转码为base64编码
    ImageBizType null.Int `json:"imageBizType,omitempty"` // 枚举值：0、1，入参1返回的url用以货品发布时的外包装使用
    Options      struct {
        Boost              bool `json:"boost"`              // 是否 AI 清晰度提升
        CateId             int  `json:"cateId"`             // 叶子类目ID，按不同类型进行裁剪，当doIntelligenceCrop=true生效
        DoIntelligenceCrop bool `json:"doIntelligenceCrop"` // 是否AI智能裁剪，true-根据sizeMode返回一组智能裁剪图（1张原图+3张裁剪图）
        SizeMode           int  `json:"sizeMode"`           // 返回尺寸大小，0-原图大小，1-800*800（1:1），2-1350*1800（3:4）
    } `json:"options,omitempty"` //

}

func (m GoodsImageUploadRequest) validate() error {
    return nil
}

func (s goodsService) ImageUpload(ctx context.Context, request GoodsImageUploadRequest) (res entity.GoodsImageUploadResult, err error) {
    if err = request.validate(); err != nil {
        return res, invalidInput(err)
    }

    var result = struct {
        normal.Response
        Result entity.GoodsImageUploadResult `json:"result"`
    }{}
    resp, err := s.httpClient.R().
        SetContext(ctx).
        SetBody(request).
        SetResult(&result).
        Post("bg.goods.image.upload")
    if err = recheckError(resp, result.Response, err); err != nil {
        return
    }

    return result.Result, nil
}
