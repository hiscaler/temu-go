package downloadableurl

import (
	"crypto/md5"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/hiscaler/gox/filex"
	"github.com/hiscaler/gox/randx"
	"gopkg.in/guregu/null.v4"
	"net"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// DownloadableUrl 可下载 URL
type DownloadableUrl string

type Option struct {
	Debug            bool
	UserAgent        string
	AppKey           string
	AppSecret        string
	AccessToken      string
	Timeout          time.Duration
	VerifySSL        bool
	StaticFileServer string
	Dir              string
}

// Download 下载信息
type Download struct {
	Url        string      `json:"url"`         // 下载地址
	ExpireTime null.Time   `json:"expire_time"` // 失效时间
	Error      null.String `json:"error"`       // 错误消息
}

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

func (dau DownloadableUrl) Download(opt Option) (Download, error) {
	var d Download
	var err error
	originalUrl := string(dau)
	if originalUrl == "" {
		return d, errors.New("url is empty")
	}
	keys := []string{
		"toa-access-token",
		"toa-app-key",
		"toa-random",
		"toa-timestamp",
	}
	d.ExpireTime = null.TimeFrom(time.Now().Add(10 * time.Minute)) // 10 分钟后过期
	dir := opt.Dir
	sb := strings.Builder{}
	headers := map[string]string{
		"toa-app-key":      opt.AppKey,
		"toa-access-token": opt.AccessToken,
	}
	httpClient := resty.New().
		SetDebug(opt.Debug).
		SetHeaders(map[string]string{
			"Content-Type": "application/json",
			"User-Agent":   opt.UserAgent,
		}).
		SetAllowGetMethodPayload(true).
		SetTimeout(opt.Timeout * time.Second).
		SetTransport(&http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: !opt.VerifySSL},
			DialContext: (&net.Dialer{
				Timeout: opt.Timeout * time.Second,
			}).DialContext,
		})
	if opt.Debug {
		httpClient.EnableTrace()
	}

	var parsedURL *url.URL
	parsedURL, err = url.Parse(originalUrl)
	if err != nil {
		return d, err
	}
	filename := strings.ToLower(filepath.Base(parsedURL.Path))
	if filename == "" {
		return d, errors.New("无法获取文件名")
	}
	savePath := filepath.Join(dir, filename)
	if !filex.Exists(savePath) {
		d.Url = urlJoin(opt.StaticFileServer, savePath)
		return d, nil
	}

	headers["toa-random"] = randx.Letter(32, true)
	headers["toa-timestamp"] = strconv.FormatInt(time.Now().Unix(), 10)
	sb.Reset()
	sb.WriteString(opt.AppSecret)
	for _, key := range keys {
		sb.WriteString(key)
		sb.WriteString(headers[key])
	}
	sb.WriteString(opt.AppSecret)
	headers["toa-sign"] = strings.ToUpper(fmt.Sprintf("%x", md5.Sum([]byte(sb.String()))))
	resp, err := httpClient.
		SetOutputDirectory(dir).
		R().
		SetHeaders(headers).
		SetOutput(filename).
		Get(originalUrl)
	if err != nil {
		return d, err
	}
	if !resp.IsSuccess() {
		return d, errors.New(resp.String())
	}
	d.Url = urlJoin(opt.StaticFileServer, path.Join(dir, filename))
	return d, nil
}
