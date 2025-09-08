package redownloadurl

import (
	"crypto/md5"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/hiscaler/gox/filex"
	"github.com/hiscaler/gox/randx"
	"github.com/hiscaler/temu-go/config"
	"gopkg.in/guregu/null.v4"
)

// RedownloadUrl 可下载 URL
// 用来简化 Temu 文件下载处理
// Temu 提供的文件下载地址是需要认证后才能查看的，且有时间的限制，该类型旨在简化繁琐的认证处理步骤，将文件下载到服务器地址路径，然后由服务器提供 HTTP 服务访问实现下载
type RedownloadUrl string

// File 下载后的文件信息
type File struct {
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

// Download 下载 Temu 资源文件
func (rdu RedownloadUrl) Download(cfg config.Config, saveDir string) (File, error) {
	var f File
	var err error
	originalUrl := string(rdu)
	if originalUrl == "" {
		return f, errors.New("url is empty")
	}
	keys := []string{
		"toa-access-token",
		"toa-app-key",
		"toa-random",
		"toa-timestamp",
	}
	f.ExpireTime = null.TimeFrom(time.Now().Add(10 * time.Minute)) // 10 分钟后过期
	sb := strings.Builder{}
	headers := map[string]string{
		"toa-app-key":      cfg.AppKey,
		"toa-access-token": cfg.AccessToken,
	}
	httpClient := resty.New().
		SetDebug(cfg.Debug).
		SetHeaders(map[string]string{
			"Content-Type": "application/json",
		}).
		SetAllowGetMethodPayload(true).
		SetTimeout(cfg.Timeout * time.Second).
		SetTransport(&http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: !cfg.VerifySSL},
			DialContext: (&net.Dialer{
				Timeout: cfg.Timeout * time.Second,
			}).DialContext,
		})
	if cfg.Debug {
		httpClient.EnableTrace()
	}

	var parsedURL *url.URL
	parsedURL, err = url.Parse(originalUrl)
	if err != nil {
		return f, err
	}
	filename := strings.ToLower(filepath.Base(parsedURL.Path))
	if filename == "" {
		return f, errors.New("无法获取文件名")
	}
	savePath := filepath.Join(saveDir, filename)
	if !filex.Exists(savePath) {
		f.Url = urlJoin(cfg.StaticFileServer, savePath)
		return f, nil
	}

	headers["toa-random"] = randx.Letter(32, true)
	headers["toa-timestamp"] = strconv.FormatInt(time.Now().Unix(), 10)
	sb.Reset()
	sb.WriteString(cfg.AppSecret)
	for _, key := range keys {
		sb.WriteString(key)
		sb.WriteString(headers[key])
	}
	sb.WriteString(cfg.AppSecret)
	headers["toa-sign"] = strings.ToUpper(fmt.Sprintf("%x", md5.Sum([]byte(sb.String()))))
	resp, err := httpClient.
		SetOutputDirectory(saveDir).
		R().
		SetHeaders(headers).
		SetOutput(filename).
		Get(originalUrl)
	if err != nil {
		return f, err
	}
	if !resp.IsSuccess() {
		return f, errors.New(resp.String())
	}
	f.Url = urlJoin(cfg.StaticFileServer, path.Join(saveDir, filename))
	return f, nil
}
