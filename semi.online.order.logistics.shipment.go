package temu

import (
	"context"
	"crypto/md5"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-resty/resty/v2"
	"github.com/hiscaler/gox/filex"
	"github.com/hiscaler/gox/randx"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
	"gopkg.in/guregu/null.v4"
)

// 物流发货服务
type semiOnlineOrderLogisticsShipmentService service

// SemiOnlineOrderLogisticsShipmentPackageItem 半托在线下单物流发货包裹数据
type SemiOnlineOrderLogisticsShipmentPackageItem struct {
	ParentOrderSn string `json:"parentOrderSn"` // 父订单号
	OrderSn       string `json:"orderSn"`       // 订单号
	GoodsId       int64  `json:"goodsId"`       // 商品 goodsId
	SkuId         int64  `json:"skuId"`         // 商品 skuId
	Quantity      int    `json:"quantity"`      // 发货数量
}

func (m SemiOnlineOrderLogisticsShipmentPackageItem) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.ParentOrderSn, validation.Required.Error("父订单号不能为空")),
		validation.Field(&m.OrderSn, validation.Required.Error("订单号不能为空")),
		validation.Field(&m.GoodsId, validation.Required.Error("商品 ID 不能为空")),
		validation.Field(&m.SkuId, validation.Required.Error("SKU ID 不能为空")),
		validation.Field(&m.Quantity, validation.Min(1).Error("发货数量不能小于 {min}")),
	)
}

// SemiOnlineOrderLogisticsShipmentAdditionalPackage 附加包裹
// 单件 sku 多包裹场景，比如如下场景：
// 一个 sku 包含一桌四椅。但是发货的时候，需要多个包裹才能发出·
type SemiOnlineOrderLogisticsShipmentAdditionalPackage struct {
	ExtendWeightUnit string   `json:"extendWeightUnit,omitempty"` // 扩展重量单位
	ExtendWeight     string   `json:"extendWeight,omitempty"`     // 扩展重量
	WeightUnit       string   `json:"weightUnit"`                 // 重量单位
	DimensionUnit    string   `json:"dimensionUnit"`              // 尺寸单位
	Weight           string   `json:"weight"`                     // 包裹重量（默认 2 位小数）
	Length           string   `json:"length"`                     // 包裹长度（默认 2 位小数）
	Height           string   `json:"height"`                     // 包裹高度（默认 2 位小数）
	Width            string   `json:"width"`                      // 包裹宽度（默认 2 位小数）
	WarehouseId      string   `json:"warehouseId"`                // 仓库 ID
	ShipCompanyId    string   `json:"shipCompanyId"`              // 物流公司 ID
	ChannelId        int64    `json:"channelId"`                  // 物流渠道 ID
	SignServiceId    null.Int `json:"signServiceId,omitempty"`    // 想使用的签收服务 ID
}

func (m SemiOnlineOrderLogisticsShipmentAdditionalPackage) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.ExtendWeightUnit, validation.Required.Error("扩展重量单位不能为空")),
		// todo 更严格的验证
	)
}

type SemiOnlineOrderLogisticsShipmentCreateSendItem struct {
	ShipCompanyId      int64                                               `json:"shipCompanyId"`                // 物流公司 id
	TrackingNumber     null.String                                         `json:"trackingNumber,omitempty"`     // 运单号
	OrderSendInfoList  []SemiOnlineOrderLogisticsShipmentPackageItem       `json:"orderSendInfoList"`            // 发货商品信息
	WarehouseId        string                                              `json:"warehouseId"`                  // 仓库id
	Weight             string                                              `json:"weight"`                       // 重量（默认 2 位小数）
	WeightUnit         string                                              `json:"weightUnit"`                   // 重量单位，美国为 lb（磅），其他国家为 kg（千克）
	Length             string                                              `json:"length"`                       // 包裹长度（默认 2 位小数）
	Width              string                                              `json:"width"`                        // 包裹宽度（默认 2 位小数）
	Height             string                                              `json:"height"`                       // 包裹高度（默认 2 位小数）
	DimensionUnit      string                                              `json:"dimensionUnit"`                // 尺寸单位高度，美国为in（英寸）其他国家为cm（厘米）
	ChannelId          int64                                               `json:"channelId"`                    // 渠道id，取自 shipservice.get
	PickupStartTime    null.Int                                            `json:"pickupStartTime,omitempty"`    // 预约上门取件开始时间（当渠道为需要下 call 同时入参预约时间渠道时，需入参。剩余渠道无需入参。）
	PickupEndTime      null.Int                                            `json:"pickupEndTime,omitempty"`      // 预约上门取件结束时间（当渠道为需要下 call 同时入参预约时间渠道时，需入参。剩余渠道无需入参。）
	SignServiceId      null.Int                                            `json:"signServiceId,omitempty"`      // 想使用的签收服务 ID
	SplitSubPackage    bool                                                `json:"splitSubPackage"`              // 是否为单件 SKU 拆多包裹（TRUE：是单件SKU多包裹场景 FALSE/不填：不是单件SKU多包裹场景）
	SendSubRequestList []SemiOnlineOrderLogisticsShipmentAdditionalPackage `json:"sendSubRequestList,omitempty"` // 单件 sku 多包裹场景，附属包裹入参
}

func (m SemiOnlineOrderLogisticsShipmentCreateSendItem) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.SendSubRequestList,
			validation.When(m.SplitSubPackage,
				validation.Required.Error("单件 sku 多包裹数据不能为空"),
				validation.Each(validation.By(func(value interface{}) error {
					v, ok := value.(SemiOnlineOrderLogisticsShipmentAdditionalPackage)
					if !ok {
						return errors.New("无效的单件 sku 多包裹数据")
					}
					return v.validate()
				})),
			),
		),
	)
}

type SemiOnlineOrderLogisticsShipmentCreateRequest struct {
	SendType           int                                              `json:"sendType"`                     // 发货类型：0-单个运单发货 1-拆成多个运单发货 2-合并发货
	ShipLater          bool                                             `json:"shipLater"`                    // 下 call 成功后是否延迟发货（TRUE：下call成功之后延迟发货，FALSE/不填：下call成功订单自动流转为已发货）
	ShipLaterLimitTime null.Int                                         `json:"shipLaterLimitTime,omitempty"` // 稍后发货兜底配置时间（单位:h），枚举：24, 48, 72, 96
	SendRequestList    []SemiOnlineOrderLogisticsShipmentCreateSendItem `json:"sendRequestList"`              // 包裹信息
}

func (m SemiOnlineOrderLogisticsShipmentCreateRequest) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.SendType, validation.In(
			entity.SemiShippingTypeSingle,
			entity.SemiShippingTypeSplit,
			entity.SemiShippingTypeMerge,
		).Error("无效的发货类型")),
		validation.Field(&m.SendRequestList,
			validation.Required.Error("包裹信息不能为空"),
			validation.Each(validation.By(func(value interface{}) error {
				v, ok := value.(SemiOnlineOrderLogisticsShipmentCreateSendItem)
				if !ok {
					return errors.New("无效的包裹信息")
				}
				return v.validate()
			})),
		),
	)
}

// Create 物流在线发货下单接口（bg.logistics.shipment.create）
// https://seller.kuajingmaihuo.com/sop/view/144659541206936016#Tf6UNY
func (s semiOnlineOrderLogisticsShipmentService) Create(ctx context.Context, request SemiOnlineOrderLogisticsShipmentCreateRequest) (packageNumbers []string, limitTime null.String, err error) {
	if err = request.validate(); err != nil {
		return
	}

	var result = struct {
		normal.Response
		Result struct {
			PackageSnList      []string    `json:"packageSnList"`      // 可使用的渠道列表
			ShipLaterLimitTime null.String `json:"shipLaterLimitTime"` // 稍后发货兜底配置时间，如下 call 时有则返回
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(request).
		SetResult(&result).
		Post("bg.logistics.shipment.create")
	if err = recheckError(resp, result.Response, err); err != nil {
		return
	}

	return result.Result.PackageSnList, result.Result.ShipLaterLimitTime, nil
}

// Query 物流在线发货下单查询接口
// https://seller.kuajingmaihuo.com/sop/view/144659541206936016#S8m7N3
func (s semiOnlineOrderLogisticsShipmentService) Query(ctx context.Context, packageNumbers ...string) ([]entity.SemiOnlineOrderLogisticsShipmentPackageResult, error) {
	if len(packageNumbers) == 0 {
		return nil, ErrInvalidParameters
	}

	var result = struct {
		normal.Response
		Result struct {
			PackageInfoResultList []entity.SemiOnlineOrderLogisticsShipmentPackageResult `json:"packageInfoResultList"` // 包裹下单结果
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(map[string][]string{"packageSnList": packageNumbers}).
		SetResult(&result).
		Post("bg.logistics.shipment.result.get")
	if err = recheckError(resp, result.Response, err); err != nil {
		return nil, err
	}

	return result.Result.PackageInfoResultList, nil
}

// 重新下单

// SemiOnlineOrderLogisticsShipmentUpdatePackage 包裹信息
type SemiOnlineOrderLogisticsShipmentUpdatePackage struct {
	PackageSn         string                                        `json:"packageSn"`         // 包裹号
	PickupStartTime   int64                                         `json:"pickupStartTime"`   // 预约上门取件的开始时间 秒级时间戳
	PickupEndTime     int64                                         `json:"pickupEndTime"`     // 预约上门取件的结束时间 秒级时间戳
	SignServiceId     int64                                         `json:"signServiceId"`     // 签收服务 ID
	ChannelId         int64                                         `json:"channelId"`         // 渠道 ID
	ShipCompanyId     int64                                         `json:"shipCompanyId"`     // 物流公司 ID
	OrderSendInfoList []SemiOnlineOrderLogisticsShipmentPackageItem `json:"orderSendInfoList"` // 发货商品信息
	// TRUE：是单件SKU多包裹场景
	// FALSE/不填：不是单件SKU多包裹场景
	SplitSubPackage    bool                                                `json:"splitSubPackage"`    // 是否为单件SKU拆多包裹
	SendSubRequestList []SemiOnlineOrderLogisticsShipmentAdditionalPackage `json:"sendSubRequestList"` // 单件sku多包裹场景，附属包裹入参
	// 具体确认场景，目前存在枚举为：
	// SUCCESSFUL_RETRY//确认是下call成功之后再次call
	// NO_DELIVERY_ON_SATURDAY//确认允许周六不上门派送】强制发货
	// DENY_CANCELLATION//确认驳回取消待确认请求，强制发货
	// DENY_ADDRESS_CHANGE://确认驳回改地址待确认请求，强制发货
	// DENY_PARENT_RISK_WARNING//确认驳回风控，强制发货
	ConfirmAcceptance []string `json:"confirmAcceptance"` // 确认场景
	WarehouseId       int64    `json:"warehouseId"`       // 仓库 id
	Weight            string   `json:"weight"`            // 包裹重量（默认 2 位小数）
	WeightUnit        string   `json:"weightUnit"`        // 重量单位
	Height            string   `json:"height"`            // 包包裹高度（默认 2 位小数）
	Length            string   `json:"length"`            // 包裹长度（默认 2 位小数）
	Width             string   `json:"width"`             // 包裹宽度（默认 2 位小数）
	DimensionUnit     string   `json:"dimensionUnit"`     // 尺寸单位高度
}

func (m SemiOnlineOrderLogisticsShipmentUpdatePackage) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.ChannelId, validation.Required.Error("物流渠道不能为空")),
		validation.Field(&m.OrderSendInfoList, validation.Required.Error("包裹发货商品不能为空")),
	)
}

type SemiOnlineOrderLogisticsShipmentUpdateRequest struct {
	RetrySendPackageRequestList []SemiOnlineOrderLogisticsShipmentUpdatePackage `json:"retrySendPackageRequestList"` // 包裹信息
}

func (m SemiOnlineOrderLogisticsShipmentUpdateRequest) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.RetrySendPackageRequestList,
			validation.Required.Error("包裹列表不能为空"),
			validation.By(func(value interface{}) error {
				packages, ok := value.([]SemiOnlineOrderLogisticsShipmentUpdatePackage)
				if !ok {
					return errors.New("无效的发货包裹")
				}
				for _, pkg := range packages {
					if err := pkg.validate(); err != nil {
						return err
					}
				}
				return nil
			}),
		),
	)
}

// Update 物流在线发货重新下单接口
// https://seller.kuajingmaihuo.com/sop/view/144659541206936016#Ff9JoY
func (s semiOnlineOrderLogisticsShipmentService) Update(ctx context.Context, request SemiOnlineOrderLogisticsShipmentUpdateRequest) (bool, error) {
	if err := request.validate(); err != nil {
		return false, invalidInput(err)
	}

	var result = struct {
		normal.Response
		Result bool `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(request).
		SetResult(&result).
		Post("bg.logistics.shipment.update")
	if err = recheckError(resp, result.Response, err); err != nil {
		return false, err
	}

	return result.Result, nil
}

// 物流在线发货修改物流接口

type EditPackageRequestItem struct {
	PackageSn      string `json:"packageSn"`      // 包裹号
	TrackingNumber string `json:"trackingNumber"` // 运单号
	ShipCompanyId  int64  `json:"shipCompanyId"`  // 物流公司 id
}

func (m EditPackageRequestItem) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.PackageSn, validation.Required.Error("包裹号不能为空")),
		validation.Field(&m.TrackingNumber, validation.Required.Error("运单号不能为空")),
		validation.Field(&m.ShipCompanyId, validation.Required.Error("物流公司不能为空")),
	)
}

type SemiOnlineOrderLogisticsShipmentUpdateShippingTypeRequest struct {
	EditPackageRequestList []EditPackageRequestItem `json:"editPackageRequestList"` // 编辑请求列表
}

func (m SemiOnlineOrderLogisticsShipmentUpdateShippingTypeRequest) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.EditPackageRequestList,
			validation.Required.Error("编辑请求列表不能为空"),
			validation.Each(validation.By(func(value interface{}) error {
				v, ok := value.(EditPackageRequestItem)
				if !ok {
					return errors.New("无效的编辑请求项")
				}
				return v.validate()
			})),
		),
	)
}

// UpdateShippingType 物流在线发货修改物流接口
// https://seller.kuajingmaihuo.com/sop/view/144659541206936016#bYZCmU
func (s semiOnlineOrderLogisticsShipmentService) UpdateShippingType(ctx context.Context, request SemiOnlineOrderLogisticsShipmentUpdateShippingTypeRequest) (bool, error) {
	if err := request.validate(); err != nil {
		return false, invalidInput(err)
	}

	var result = struct {
		normal.Response
		Result bool `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(request).
		SetResult(&result).
		Post("bg.logistics.shipment.shippingtype.update")
	if err = recheckError(resp, result.Response, err); err != nil {
		return false, err
	}

	return result.Result, nil
}

// 物流在线发货打印面单接口（bg.logistics.shipment.document.get）

type SemiOnlineOrderLogisticsShipmentDocumentRequest struct {
	// - SHIPPING_LABEL_PDF:入参此参数，返回的 URL 加签后只返回 PDF 格式的面单文件
	// - 不入参，按照旧有逻辑返回面单文件，即按物流商的面单文件返回确定图片格式或 PDF 格式；
	// - 入不合法的参数值：接口报错，报错文案：Document type is invalid.
	DocumentType  string   `json:"documentType"`  // 文件类型
	PackageSnList []string `json:"packageSnList"` // 需要打印面单的包裹号列表
	// 自行添加，非接口字段，用于下载面单文件
	Download      bool `json:"download"`      // 是否下载面单文件
	RetryDownload bool `json:"retryDownload"` // 面单文件已经存在的情况下是否重新下载，为 false 的情况下直接返回，不会去再次请求
}

func (m SemiOnlineOrderLogisticsShipmentDocumentRequest) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.DocumentType,
			validation.Required.Error("面单文件类型不能为空"),
			validation.In(entity.LogisticsShipmentDocumentPdfFile, entity.LogisticsShipmentDocumentImageFile).Error("无效的面单文件类型"),
		),
		validation.Field(&m.PackageSnList, validation.Required.Error("需要打印面单的包裹号列表不能为空")),
	)
}

// urlJoin 拼接 URL 地址
func urlJoin(prefix, file string) string {
	if strings.HasSuffix(prefix, "/") {
		prefix = prefix[0 : len(prefix)-1]
	}

	file = path.Clean(file)
	if !strings.HasPrefix(file, "/") {
		file = "/" + file
	}
	return prefix + file
}

// Document 物流在线发货打印面单接口
// https://seller.kuajingmaihuo.com/sop/view/144659541206936016#IYqSks
func (s semiOnlineOrderLogisticsShipmentService) Document(ctx context.Context, request SemiOnlineOrderLogisticsShipmentDocumentRequest) ([]entity.SemiOnlineOrderLogisticsShipmentDocument, error) {
	if err := request.validate(); err != nil {
		return nil, invalidInput(err)
	}

	var result = struct {
		normal.Response
		Result struct {
			ShippingLabelUrlList []entity.SemiOnlineOrderLogisticsShipmentDocument `json:"shippingLabelUrlList"` // 包裹对应的面单文件 url（PDF 或图片）
		} `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(request).
		SetResult(&result).
		Post("bg.logistics.shipment.document.get")
	if err = recheckError(resp, result.Response, err); err != nil {
		return nil, err
	}

	documents := result.Result.ShippingLabelUrlList
	if !request.Download || len(documents) == 0 {
		return documents, nil
	}

	keys := []string{
		"toa-access-token",
		"toa-app-key",
		"toa-random",
		"toa-timestamp",
	}
	expireTime := time.Now().Add(10 * time.Minute).Unix() // 10 分钟后过期
	dir := "./static_files/temu/logistics-labels"
	sb := strings.Builder{}
	headers := map[string]string{
		"toa-app-key":      s.config.AppKey,
		"toa-access-token": s.config.AccessToken,
	}
	httpClient := resty.New().
		SetDebug(s.debug).
		SetHeaders(map[string]string{
			"Content-Type": "application/json",
			"User-Agent":   UserAgent,
		}).
		SetAllowGetMethodPayload(true).
		SetTimeout(s.config.Timeout * time.Second).
		SetTransport(&http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: !s.config.VerifySSL},
			DialContext: (&net.Dialer{
				Timeout: s.config.Timeout * time.Second,
			}).DialContext,
		})
	if s.debug {
		httpClient.EnableTrace()
	}
	for i, doc := range documents {
		documents[i].ExpireTime = expireTime
		if doc.Url == "" {
			documents[i].Error = null.StringFrom("URL is empty")
			continue
		}

		filename := path.Base(doc.Url)
		if index := strings.Index(filename, "?"); index != -1 {
			filename = filename[0:index]
		}
		filename = strings.ToLower(fmt.Sprintf("%s%s", doc.PackageSn, path.Ext(filename)))
		savePath := filepath.Join(dir, filename)
		if !request.RetryDownload && filex.Exists(savePath) {
			documents[i].Path = null.StringFrom(urlJoin(s.config.StaticFileServer, savePath))
			continue
		}

		headers["toa-random"] = randx.Letter(32, true)
		headers["toa-timestamp"] = strconv.FormatInt(time.Now().Unix(), 10)
		sb.Reset()
		sb.WriteString(s.config.AppSecret)
		for _, key := range keys {
			sb.WriteString(key)
			sb.WriteString(headers[key])
		}
		sb.WriteString(s.config.AppSecret)
		headers["toa-sign"] = strings.ToUpper(fmt.Sprintf("%x", md5.Sum([]byte(sb.String()))))
		resp, err = httpClient.
			SetOutputDirectory(dir).
			R().
			SetHeaders(headers).
			SetOutput(filename).
			Get(doc.Url)
		if err != nil {
			documents[i].Error = null.StringFrom(err.Error())
		} else {
			if resp.IsError() {
				documents[i].Error = null.StringFrom(resp.String())
			} else if resp.IsSuccess() {
				documents[i].Path = null.StringFrom(urlJoin(s.config.StaticFileServer, path.Join(dir, filename)))
			} else {
				documents[i].Error = null.StringFrom(resp.String())
			}
		}
	}

	return documents, nil
}
