package config

import (
	"log/slog"
	"time"
)

type Config struct {
	Debug       bool          `json:"debug"`        // 是否为调试模式
	AppKey      string        `json:"app_key"`      // App Key
	AppSecret   string        `json:"app_secret"`   // App 秘钥
	AccessToken string        `json:"access_token"` // Access Token
	Timeout     time.Duration `json:"timeout"`      // 超时时间（秒）
	VerifySSL   bool          `json:"verify_ssl"`   // 是否验证 SSL
	Logger      *slog.Logger  `json:"-"`            // 日志
}
