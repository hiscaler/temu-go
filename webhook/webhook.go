package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"slices"
	"strings"
)

// Webhook
// https://partner-us.temu.com/documentation?menu_code=38e79b35d2cb463d85619c1c786dd303&sub_menu_code=b996f4717c5c44da8410069a63da6723#callback_message_content
type Webhook struct {
	secretKey string
	header    http.Header
	rawBody   string
	body      string
}

func NewWebhook(secretKey string, header http.Header, eventData string) *Webhook {
	return &Webhook{
		secretKey: secretKey,
		header:    header,
		rawBody:   eventData,
	}
}

func (w *Webhook) Valid() bool {
	if len(w.header) == 0 || len(w.rawBody) == 0 {
		return false
	}
	headerKeys := []string{
		"x-tm-app-key",
		"x-tm-event-code",
		"x-tm-timestamp",
		"x-tm-signature",
		"x-tm-ext-param",
	}
	sb := strings.Builder{}
	sb.WriteString(w.rawBody)
	for k, v := range w.header {
		if !slices.Contains(headerKeys, k) || len(v) == 0 || v[0] == "" {
			return false
		}
		sb.WriteString(k)
		sb.WriteString("=")
		sb.WriteString(v[0])
	}

	mac := hmac.New(sha256.New, []byte(w.secretKey))
	_, err := mac.Write([]byte(sb.String()))
	if err != nil {
		return false
	}
	return string(mac.Sum(nil)) == w.rawBody
}

func (w *Webhook) Decrypt() (string, error) {
	sig, err := hex.DecodeString(w.rawBody)
	if err != nil {
		return "", err
	}

	mac := hmac.New(sha256.New, []byte(w.secretKey))
	mac.Write([]byte(w.rawBody))
	if !hmac.Equal(sig, mac.Sum(nil)) {
		return "", errors.New("invalid webhook body")
	}
	return "", nil
}
