package temu

import (
	"crypto/md5"
	"crypto/tls"
	"embed"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"log/slog"
	"net"
	"net/http"
	"os"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-resty/resty/v2"
	"github.com/goccy/go-json"
	"github.com/hiscaler/gox/stringx"
	"github.com/hiscaler/temu-go/config"
	"github.com/hiscaler/temu-go/entity"
	"github.com/hiscaler/temu-go/normal"
)

const (
	Version   = "0.0.1"
	UserAgent = "Temu API Client-Golang/" + Version
)

const (
	devEnv  = "dev"  // 开发环境
	testEnv = "test" // 测试环境
	prodEnv = "prod" // 生产环境
)

const (
	BadRequestError              = 400     // 错误的请求
	UnauthorizedError            = 401     // 身份验证或权限错误
	NotFoundError                = 404     // 访问资源不存在
	InternalServerError          = 500     // 服务器内部错误
	MethodNotImplementedError    = 501     // 方法未实现
	SystemExceptionError         = 200000  // 系统异常
	InvalidSignError             = 7000015 // 签名无效
	NoAppKeyError                = 7000002 // 未设置 App Key
	NoAccessTokenError           = 7000003 // 未设置 Access Token
	InvalidAccessTokenError      = 7000018 // 无效的 Access Token
	AccessTokenKeyUnmatchedError = 7000006 // Access Token 和 Key 不匹配
	TypeIsNotExistsError         = 3000003 // 接口不存在
)

var ErrNotFound = validation.ErrorObject{}.SetCode("NotFound").SetMessage("数据不存在")
var ErrInvalidSign = validation.ErrorObject{}.SetCode("InvalidSign").SetMessage("无效的签名")
var ErrInvalidParameters = validation.ErrorObject{}.SetCode("InvalidParameters").SetMessage("无效的参数")

type service struct {
	debug      bool          // Is debug mode
	logger     *slog.Logger  // Logger
	config     config.Config // Config
	language   *string       // Message language
	httpClient *resty.Client // HTTP client
}

type services struct {
	PurchaseOrder purchaseOrderService
	ShipOrder     shipOrderService
	Logistics     logisticsService
	Goods         goodsService
	Mall          mallService
	Jit           jitService
	SemiManaged   semiManagedService
	Picture       pictureService
}

type Client struct {
	language     string         // 消息语种
	Env          string         // 环境
	Debug        bool           // 是否为 Debug 模式
	Region       string         // 接口所在区域
	Logger       *slog.Logger   // Log
	Services     services       // API services
	TimeLocation *time.Location // 时区
}

// url 获取方法对应的 URL
func url(typ, region, env string, proxies config.RegionEnvUrls) string {
	// key 为 type 值，value 为对应的区域，为空表示根据 region 确定 baseUrl，
	// 不为空的情况下表示无论传入的 region 为何值，均取 value 作为 region 值去获取 baseUrl
	//
	// 注意：以下列出的 type 均为半托请求，非半托请求不要添加进来
	semiTypes := map[string]string{
		"bg.logistics.shippingservices.get":         "",
		"bg.logistics.shipment.document.get":        "",
		"bg.logistics.shipment.create":              "",
		"bg.logistics.shipment.result.get":          "",
		"bg.logistics.shipment.update":              "",
		"bg.logistics.shipment.shippingtype.update": "",
		"bg.logistics.warehouse.list.get":           "",
		"bg.logistics.shipped.package.confirm":      "",
		"bg.order.unshipped.package.get":            "",
		"bg.order.list.v2.get":                      "",
		"bg.order.list.get":                         "",
		"bg.logistics.shipment.v2.get":              entity.AmericanRegion,
		"bg.logistics.shipment.get":                 entity.AmericanRegion,
		"bg.goods.quantity.get":                     "",
		"bg.goods.quantity.update":                  "",
		"bg.logistics.companies.get":                "",
		"bg.order.shippinginfo.get":                 "",
		"bg.logistics.shipment.confirm":             entity.AmericanRegion,
	}
	if v, ok := semiTypes[typ]; ok {
		if v != "" {
			region = v
		}
	} else {
		region = entity.ChinaRegion
	}

	urls := config.RegionEnvUrls{
		entity.ChinaRegion: {
			Prod: "https://openapi.kuajingmaihuo.com/openapi/router",
			Test: "https://openapi.kuajingmaihuo.com/openapi/router",
		},
		entity.AmericanRegion: {
			Prod: "https://openapi-b-us.temu.com/openapi/router",
			Test: "http://openapi-b-us.temudemo.com/openapi/router",
		},
		entity.EuropeanUnionRegion: {
			Prod: "https://openapi-b-eu.temu.com/openapi/router",
			Test: "http://openapi-b-eu.temudemo.com/openapi/router",
		},
	}
	// 如果有设置请求代理的话，则使用代理的地址替换 Temu 平台地址
	if proxies != nil {
		for k, proxy := range proxies {
			v, ok := urls[k]
			if !ok || (v.Prod == "" && v.Test == "") {
				continue
			}

			if proxy.Prod != "" {
				v.Prod = proxy.Prod
			}
			if proxy.Test != "" {
				v.Test = proxy.Test
			}
			urls[k] = v
		}
	}

	envUrl, _ := urls[region]
	if env == prodEnv {
		return envUrl.Prod
	} else {
		return envUrl.Test
	}
}

// generate sign string
func generateSign(values map[string]any, appSecret string) map[string]any {
	delete(values, "sign")
	values["timestamp"] = time.Now().Unix()
	size := len(values)
	keys := make([]string, size)
	for k := range values {
		size--
		keys[size] = k
	}
	sort.Strings(keys)
	sb := strings.Builder{}
	sb.WriteString(appSecret)
	for _, key := range keys {
		value := strings.TrimSpace(stringx.String(values[key]))
		if value == "" {
			delete(values, key)
			continue
		}
		sb.WriteString(key)
		sb.WriteString(value)
	}
	sb.WriteString(appSecret)
	values["sign"] = strings.ToUpper(fmt.Sprintf("%x", md5.Sum([]byte(sb.String()))))
	return values
}

var loc *time.Location
var lang *string
var i18nBundle *i18n.Bundle
var i18nLocalizer *i18n.Localizer
var versionPattern = regexp.MustCompile(`\.(v[1-9]+)`)

//go:embed locales/*.toml
var localeFS embed.FS

func init() {
	var err error
	loc, err = time.LoadLocation("Asia/Shanghai")
	if err != nil {
		loc = time.FixedZone("CST", 8*3600)
	}
	time.Local = loc

	i18nBundle = i18n.NewBundle(language.SimplifiedChinese)
	i18nBundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	l := language.SimplifiedChinese.String()
	lang = &l
	_, _ = i18nBundle.LoadMessageFileFS(localeFS, fmt.Sprintf("locales/%s.toml", l))
	i18nLocalizer = i18n.NewLocalizer(i18nBundle, l)
}

type simpleResponse struct {
	Success   bool   `json:"success"`
	ErrorCode int    `json:"errorCode"`
	ErrorMsg  string `json:"errorMsg"`
}

// retry 请求是否可重试
func (r simpleResponse) retry() bool {
	return !r.Success && r.ErrorCode == 4000000 && strings.ToLower(r.ErrorMsg) == "system_exception"
}

// 默认中国区
func parseRegion(region string) string {
	region = strings.ToUpper(region)
	if !slices.Contains([]string{entity.ChinaRegion, entity.AmericanRegion, entity.EuropeanUnionRegion}, region) {
		region = entity.ChinaRegion
	}
	return region
}

// getVersion 从 typ 中获取 API 版本
func getVersion(typ string) string {
	// todo 不是所有情况下 type 值都是 x.v2.y 的形式，存在 xv2.y 这样的 type，比如 bg.purchaseorderv2.get
	vs := versionPattern.FindStringSubmatch(typ)
	if len(vs) <= 1 {
		return "V1"
	}

	return strings.ToUpper(vs[1])
}

func NewClient(cfg config.Config) *Client {
	var l *slog.Logger
	debug := cfg.Debug
	if cfg.Logger != nil {
		l = cfg.Logger
	} else {
		if debug {
			l = slog.New(slog.NewTextHandler(os.Stdout, nil))
		} else {
			l = slog.New(slog.NewJSONHandler(os.Stdout, nil))
		}
	}

	env := strings.ToLower(cfg.Env)
	if env == "" {
		env = prodEnv
	} else if env != prodEnv {
		env = testEnv
	}
	region := parseRegion(cfg.Region)
	client := &Client{
		language:     *lang,
		Env:          env,
		Debug:        debug,
		Region:       region,
		Logger:       l,
		TimeLocation: loc,
	}
	httpClient := resty.New().
		SetDebug(debug).
		EnableTrace().
		SetHeaders(map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
			"User-Agent":   UserAgent,
		}).
		SetAllowGetMethodPayload(true).
		SetTimeout(cfg.Timeout * time.Second).
		SetTransport(&http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: !cfg.VerifySSL},
			DialContext: (&net.Dialer{
				Timeout: cfg.Timeout * time.Second,
			}).DialContext,
		}).
		OnBeforeRequest(func(c *resty.Client, request *resty.Request) error {
			values := make(map[string]any)
			if request.Body != nil {
				b, e := json.Marshal(request.Body)
				if e != nil {
					return e
				}

				e = json.Unmarshal(b, &values)
				if e != nil {
					return e
				}
			}
			values["app_key"] = cfg.AppKey
			values["access_token"] = cfg.AccessToken
			values["data_type"] = "JSON"
			values["timestamp"] = time.Now().Unix()
			typ := request.URL
			values["version"] = getVersion(typ)
			values["type"] = typ
			request.URL = ""
			request.SetBody(generateSign(values, cfg.AppSecret))
			c.SetBaseURL(url(typ, client.Region, client.Env, cfg.Proxies))
			return nil
		}).
		OnAfterResponse(func(client *resty.Client, response *resty.Response) error {
			if debug || response.IsError() {
				params := response.Request.Body
				endpoint := ""
				if v, ok := params.(map[string]any); ok {
					var vv any
					if vv, ok = v["type"]; ok {
						endpoint = vv.(string)
					}
				}
				l.Error(fmt.Sprintf(`%s %s
   ENDPOINT: %s
 PARAMETERS: %s
	 STATUS: %s
       BODY: %v`,
					response.Request.Method,
					response.Request.URL,
					endpoint,
					params,
					response.Status(),
					response.String(),
				))
			}
			return nil
		}).
		SetRetryCount(3).
		SetRetryWaitTime(time.Duration(500) * time.Millisecond).
		SetRetryMaxWaitTime(time.Duration(1) * time.Second).
		AddRetryCondition(func(response *resty.Response, err error) bool {
			if response == nil {
				return false
			}

			retry := response.StatusCode() == http.StatusTooManyRequests
			if !retry {
				var r simpleResponse
				retry = json.Unmarshal(response.Body(), &r) == nil && r.retry()
			}
			if retry {
				body := response.Request.Body
				endpoint := ""
				if body != nil {
					var values map[string]any
					var b []byte
					var e error
					if b, e = json.Marshal(body); e == nil {
						if e = json.Unmarshal(b, &values); e == nil {
							if v, ok := values["type"]; ok {
								endpoint = v.(string)
							}
							response.Request.SetBody(generateSign(values, cfg.AppSecret))
						}
					}
					retry = e == nil
				}
				if retry && debug {
					messages := make([]string, 0)
					messages = append(messages, "URL: "+response.Request.URL)
					if endpoint != "" {
						messages = append(messages, "Type: "+endpoint)
					}
					l.Info("Retry " + strings.Join(messages, " "))
				}
			}
			return retry
		}).
		SetRetryAfter(func(client *resty.Client, response *resty.Response) (time.Duration, error) {
			var milliseconds int64 = 0
			if response != nil {
				retry := response.StatusCode() == http.StatusTooManyRequests
				if !retry {
					var r simpleResponse
					retry = json.Unmarshal(response.Body(), &r) == nil && r.retry()
				}
				if retry {
					milliseconds = 1000 - time.Now().UnixMilli()%1000 // 最多等待下一秒钟到目前的毫秒数
				}
			}
			if milliseconds == 0 {
				return 0, nil
			}
			l.Debug(fmt.Sprintf("Retry waiting %d milliseconds...", milliseconds))
			return time.Duration(milliseconds) * time.Millisecond, nil
		})
	httpClient.JSONMarshal = json.Marshal
	httpClient.JSONUnmarshal = json.Unmarshal
	xService := service{
		debug:      debug,
		logger:     l,
		httpClient: httpClient,
		language:   lang,
		config:     cfg,
	}
	client.Services = services{
		PurchaseOrder: (purchaseOrderService)(xService),
		ShipOrder: shipOrderService{
			service:        xService,
			Package:        (shipOrderPackageService)(xService),
			Packing:        (shipOrderPackingService)(xService),
			ReceiveAddress: (shipOrderReceiveAddressService)(xService),
			Staging:        (shipOrderStagingService)(xService),
			Logistics:      (shipOrderLogisticsService)(xService),
		},
		Logistics: (logisticsService)(xService),
		Goods: goodsService{
			service: xService,
			Barcode: (goodsBarcodeService)(xService),
			Brand:   (goodsBrandService)(xService),
			Category: goodsCategoryService{
				service:   xService,
				Attribute: (goodsCategoryAttributeService)(xService),
			},
			Certification:       (goodsCertificationService)(xService),
			LifeCycle:           (goodsLifeCycleService)(xService),
			Sales:               (goodsSalesService)(xService),
			SizeChartClass:      (goodsSizeChartClassService)(xService),
			SizeChart:           (goodsSizeChartService)(xService),
			SizeChartSetting:    (goodsSizeChartSettingService)(xService),
			SizeChartTemplate:   (goodsSizeChartTemplateService)(xService),
			TopSelling:          (goodsTopSellingService)(xService),
			Warehouse:           (goodsWarehouseService)(xService),
			Quantity:            (goodsQuantityService)(xService),
			ParentSpecification: (goodsParentSpecificationService)(xService),
			Specification:       (goodsSpecificationService)(xService),
			Price: goodsPriceService{
				service:        xService,
				Review:         (goodsPriceReviewService)(xService),
				FullAdjustment: (goodsPriceFullAdjustmentService)(xService),
			},
		},
		Mall: mallService{
			service:         xService,
			DeliveryAddress: (mallDeliveryAddressService)(xService),
		},
		Jit: jitService{
			service:          xService,
			PresaleRule:      (jitPresaleRuleService)(xService),
			VirtualInventory: (jitVirtualInventoryService)(xService),
		},
		SemiManaged: semiManagedService{
			Order: (semiOrderService)(xService),
			OnlineOrder: semiOnlineOrderService{
				Logistics: semiOnlineOrderLogisticsService{
					ServiceProvider: (semiOnlineOrderLogisticsServiceProviderService)(xService),
					Shipment:        (semiOnlineOrderLogisticsShipmentService)(xService),
					Warehouse:       (semiOnlineOrderLogisticsWarehouseService)(xService),
				},
				Package: semiOnlineOrderPackageService{
					Unshipped: (semiOnlineOrderUnshippedPackageService)(xService),
					Shipped:   (semiOnlineOrderShippedPackageService)(xService),
				},
			},
			VirtualInventory:       (semiVirtualInventoryService)(xService),
			OrderLogisticsShipment: (semiOrderLogisticsShipmentService)(xService),
			Logistics:              (semiLogisticsService)(xService),
		},
		Picture: (pictureService)(xService),
	}

	return client
}

func (c *Client) SetRegion(region string) *Client {
	c.Region = parseRegion(region)
	return c
}

// SetLanguage 设置返回消息语种
// 设置有误的情况下（比如语种文件不存在等）默认为英文
func (c *Client) SetLanguage(l language.Tag) *Client {
	if l == language.Chinese {
		l = language.SimplifiedChinese
	}
	langString := l.String()
	if c.language == langString {
		return c
	}

	for _, v := range []language.Tag{l, language.English} {
		_, err := i18nBundle.LoadMessageFileFS(localeFS, fmt.Sprintf("locales/%s.toml", v.String()))
		if err == nil {
			i18nLocalizer = i18n.NewLocalizer(i18nBundle, v.String())
			break
		}
		slog.Error(fmt.Sprintf(`SetLanguage("%s") error: %s`, v.String(), err.Error()), slog.String("lang", v.String()))
	}
	c.language = langString
	lang = &langString
	return c
}

func parseResponseTotal(currentPage, pageSize, total int) (n, totalPages int, isLastPage bool) {
	if currentPage == 0 {
		currentPage = 1
	}

	totalPages = (total / pageSize) + 1
	return total, totalPages, currentPage >= totalPages
}

func invalidInput(e error) error {
	var errs validation.Errors
	if !errors.As(e, &errs) {
		return e
	}

	if len(errs) == 0 {
		return nil
	}

	fields := make([]string, 0)
	messages := make([]string, 0)
	for field := range errs {
		fields = append(fields, field)
	}
	sort.Strings(fields)

	localizeConfig := &i18n.LocalizeConfig{}
	for _, field := range fields {
		e1 := errs[field]
		if e1 == nil {
			continue
		}

		var errObj validation.ErrorObject
		if errors.As(e1, &errObj) {
			e1 = errors.New(errObj.Code())
			localizeConfig.MessageID = errObj.Code()
			localizeConfig.TemplateData = errObj.Params()
			msg, err := i18nLocalizer.Localize(localizeConfig)
			if err != nil {
				e1 = errors.New(errObj.Error())
			} else {
				e1 = errors.New(msg)
			}
		} else {
			var errs1 validation.Errors
			if errors.As(e1, &errs1) {
				e1 = invalidInput(errs1)
				if e1 == nil {
					continue
				}
			}
		}

		messages = append(messages, e1.Error())
	}
	return errors.New(strings.Join(messages, "; "))
}

func recheckError(resp *resty.Response, result normal.Response, e error) error {
	if e != nil {
		if errors.Is(e, http.ErrHandlerTimeout) {
			return errorWrap(http.StatusRequestTimeout, e.Error())
		}
		return e
	}

	if resp.IsError() {
		return errorWrap(resp.StatusCode(), resp.Error().(string))
	}

	if !result.Success {
		return errorWrap(result.ErrorCode, result.ErrorMessage)
	}
	return nil
}

// errorWrap wrap an error, if status code is 200, return nil, otherwise return an error
func errorWrap(code int, message string) error {
	if code == http.StatusOK {
		return nil
	}

	if code == NotFoundError {
		return ErrNotFound
	}

	// message not found in translate file if err not equal nil
	msg, err := i18nLocalizer.Localize(&i18n.LocalizeConfig{
		MessageID: strconv.Itoa(code),
	})
	if err == nil {
		return errors.New(msg)
	}

	switch code {
	case BadRequestError:
		message = "请求错误"
	case UnauthorizedError:
		message = "认证失败，无法访问系统资源"
	case InternalServerError:
		message = "服务器内部错误"
	case MethodNotImplementedError:
		message = "请求方法未实现"
	case SystemExceptionError, 4000000:
		message = "Temu 平台异常"
	case InvalidSignError:
		return ErrInvalidSign
	case NoAppKeyError:
		message = "未设置 App Key"
	case NoAccessTokenError:
		message = "未设置 Access Token"
	case InvalidAccessTokenError, 7000020:
		message = "无效的 Access Token"
	case AccessTokenKeyUnmatchedError:
		message = "Access Token 和 Key 不匹配"
	case TypeIsNotExistsError:
		return errors.New("接口不存在")
	case 7000016:
		message = "无效的请求地址"
	case 2000000, 2000090, 3000000:
	case 7000007:
		message = "Access Token 已过期，请联系卖家重新授权并与您共享新的 Access Token"
	case 2000060:
		message = "店铺类型不符合预期，不允许查询或变更库存操作"
	default:
		message = fmt.Sprintf("%d: %s", code, strings.TrimSpace(message))
	}

	return errors.New(message)
}
