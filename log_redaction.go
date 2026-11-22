package temu

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"unicode"
)

const redacted = "[REDACTED]"

var sensitiveLogQuery = regexp.MustCompile(`(?i)([?&](?:access[_-]?token|app[_-]?secret|password|authorization|phone|mobile|email|address|receiver|recipient|consignee|contact)[^=&\s]*=)[^&\s]*`)

// redactLogJSON returns copyable JSON with sensitive fields and configured secrets removed.
func redactLogJSON(value any, secrets ...string) string {
	var data any
	var encoded []byte
	var err error
	if raw, ok := value.([]byte); ok {
		encoded = raw
	} else {
		encoded, err = json.Marshal(value)
		if err != nil {
			return `"[unavailable: cannot encode JSON]"`
		}
	}
	decoder := json.NewDecoder(bytes.NewReader(encoded))
	decoder.UseNumber()
	if decoder.Decode(&data) != nil || !json.Valid(encoded) {
		return `"[unavailable: non-JSON response]"`
	}
	clean, err := json.Marshal(redactValue(data, secrets))
	if err != nil {
		return `"[unavailable: cannot encode JSON]"`
	}
	var formatted bytes.Buffer
	if err := json.Indent(&formatted, clean, "", "   "); err != nil {
		return string(clean)
	}
	return formatted.String()
}

func sensitiveLogKey(key string) bool {
	key = strings.ToLower(strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return r
		}
		return -1
	}, key))
	for _, part := range []string{
		"token", "secret", "password", "passwd", "authorization", "credential", "cookie", "signature",
		"phone", "mobile", "email", "address", "recipient", "receiver", "consignee", "contact",
		"username", "buyername", "fullname", "realname", "postcode", "zipcode", "street",
	} {
		if strings.Contains(key, part) {
			return true
		}
	}
	switch key {
	case "apikey", "privatekey", "toasign", "name", "firstname", "lastname":
		return true
	}
	return false
}

func redactLogHeaders(headers http.Header, secrets ...string) {
	for key, values := range headers {
		if sensitiveLogKey(key) {
			headers[key] = []string{redacted}
			continue
		}
		for i, value := range values {
			for _, secret := range secrets {
				if secret != "" {
					value = strings.ReplaceAll(value, secret, redacted)
				}
			}
			values[i] = value
		}
	}
}

func redactValue(value any, secrets []string) any {
	switch v := value.(type) {
	case map[string]any:
		result := make(map[string]any, len(v))
		for key, item := range v {
			if sensitiveLogKey(key) {
				result[key] = redacted
			} else {
				result[key] = redactValue(item, secrets)
			}
		}
		return result
	case []any:
		result := make([]any, len(v))
		for i, item := range v {
			result[i] = redactValue(item, secrets)
		}
		return result
	case string:
		// Some APIs embed JSON as a string; inspect its fields without changing its type.
		var nested any
		trimmed := strings.TrimSpace(v)
		decoder := json.NewDecoder(strings.NewReader(trimmed))
		decoder.UseNumber()
		if len(trimmed) > 1 && (trimmed[0] == '{' || trimmed[0] == '[') && json.Valid([]byte(trimmed)) && decoder.Decode(&nested) == nil {
			if encoded, err := json.Marshal(redactValue(nested, secrets)); err == nil {
				return string(encoded)
			}
			return redacted
		}
		for _, secret := range secrets {
			if secret != "" {
				v = strings.ReplaceAll(v, secret, redacted)
			}
		}
		if parsed, err := url.Parse(v); err == nil && (parsed.Scheme == "http" || parsed.Scheme == "https") && parsed.Host != "" {
			if parsed.RawQuery != "" {
				parsed.RawQuery = "REDACTED"
			}
			parsed.Fragment = ""
			return parsed.String()
		}
		return v
	default:
		return value
	}
}
