package helpers

import (
	"regexp"
	"slices"
	"strings"
)

// SensitiveWordChecker 敏感词检测
type SensitiveWordChecker struct {
	words    []string
	patterns map[string]string
}

func NewSensitiveWordChecker(words ...string) *SensitiveWordChecker {
	swc := &SensitiveWordChecker{
		words:    []string{},
		patterns: make(map[string]string),
	}
	return swc.SetWords(words)
}

func (swc *SensitiveWordChecker) LoadDefaultWords() *SensitiveWordChecker {
	swc.SetWords([]string{
		"sex", "fuck", "shit", "bitch", "cunt", "pussy", "dick", "asshole", "bastard", "fag", "crap", "dumbass", "twat", "bollocks", "arsehole",
	})
	return swc
}

func (swc *SensitiveWordChecker) SetWords(words []string) *SensitiveWordChecker {
	for _, word := range words {
		word = strings.TrimSpace(word)
		if word != "" && slices.IndexFunc(swc.words, func(w string) bool { return strings.EqualFold(w, word) }) == -1 {
			pattern := "(?i)\\b"
			for i, r := range word {
				if i > 0 {
					pattern += "[\\s\\-]*"
				}
				pattern += regexp.QuoteMeta(string(r))
			}
			pattern += "\\b"
			swc.patterns[word] = pattern
			swc.words = append(swc.words, word)
		}
	}
	return swc
}

// Execute 执行验证判断
// 返回值为：bool, []string, error
// 第 1 个返回值为：bool 类型（true：无敏感词，false：有敏感词）
// 第 2 个返回值为：[]string 类型，返回敏感词列表
// 第 3 个返回值为：error 类型，返回错误信息
func (swc *SensitiveWordChecker) Execute(text string) (valid bool, sensitiveWords []string, err error) {
	if text == "" {
		return true, nil, nil
	}

	sWords := make([]string, 0)
	for _, word := range swc.words {
		var matched bool
		if matched, err = regexp.MatchString(swc.patterns[word], text); err != nil {
			return
		}
		if matched {
			sWords = append(sWords, word)
		}
	}
	return len(sWords) == 0, sWords, nil
}
