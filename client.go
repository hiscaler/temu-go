package temu

import (
	"crypto/md5"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"log/slog"
	"net"
	"net/http"
	"os"
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
	language   *string       // Language
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
		value := stringx.String(values[key])
		value = strings.TrimSpace(value)
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
	_, _ = i18nBundle.LoadMessageFile(fmt.Sprintf("./locales/%s.toml", l))
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
	if !slices.Contains([]string{
		entity.ChinaRegion,
		entity.AmericanRegion,
		entity.EuropeanUnionRegion,
	}, region) {
		region = entity.ChinaRegion
	}
	return region
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

	urls := map[string]config.URLPair{
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

	if cfg.OverwriteUrls != nil {
		for region, overwriteURL := range cfg.OverwriteUrls {
			if _, exists := urls[region]; exists {
				if overwriteURL.Prod != "" {
					urls[region] = config.URLPair{Prod: overwriteURL.Prod, Test: urls[region].Test}
				}
				if overwriteURL.Test != "" {
					urls[region] = config.URLPair{Prod: urls[region].Prod, Test: overwriteURL.Test}
				}
			}
		}
	}

	env := strings.ToLower(cfg.Env)
	if env == "" {
		env = prodEnv
	} else if env != prodEnv {
		env = testEnv
	}
	region := parseRegion(cfg.Region)
	url := ""
	if v, ok := urls[region]; ok {
		if env == prodEnv {
			url = v.Prod
		} else {
			url = v.Test
		}
	}
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
		SetBaseURL(url).
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
		OnBeforeRequest(func(client *resty.Client, request *resty.Request) error {
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
			values["version"] = "V1"
			values["timestamp"] = time.Now().Unix()
			values["type"] = request.URL
			request.URL = ""
			request.SetBody(generateSign(values, cfg.AppSecret))
			return nil
		}).
		OnAfterResponse(func(client *resty.Client, response *resty.Response) error {
			if debug || response.IsError() {
				params := response.Request.Body
				endpoint := ""
				if v, ok := params.(map[string]any); ok {
					if vv, ok := v["type"]; ok {
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
		},
		Mall: mallService{
			service: xService,
			Address: (mallAddressService)(xService),
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
		},
	}

	return client
}

func (c *Client) SetRegion(region string) *Client {
	c.Region = parseRegion(region)
	return c
}

func (c *Client) SetLanguage(l string) *Client {
	if c.language != l {
		_, err := i18nBundle.LoadMessageFile(fmt.Sprintf("./locales/%s.toml", l))
		if err != nil {
			slog.Error(fmt.Sprintf(`SetLanguage("%s") error: %s`, l, err.Error()), slog.String("lang", l))
		} else {
			i18nLocalizer = i18n.NewLocalizer(i18nBundle, l)
		}
	}
	c.language = l
	lang = &l
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

	message = strings.TrimSpace(message)
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
	default:
		message = fmt.Sprintf("%d: %s", code, message)
	}

	// msg not found in translate file if err not equal nil
	msg, err := i18nLocalizer.Localize(&i18n.LocalizeConfig{
		MessageID: strconv.Itoa(code),
	})
	if err == nil {
		message = msg
	}

	return errors.New(message)
}
