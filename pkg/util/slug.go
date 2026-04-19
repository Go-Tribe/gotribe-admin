package util

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/mozillazg/go-pinyin"
)

// GenerateSlug 根据标题生成 URL 友好的 slug
// 英文标题直接转小写+连字符
// 中文标题转为拼音
func GenerateSlug(title string) string {
	title = strings.TrimSpace(title)
	if title == "" {
		return ""
	}

	// 检测是否包含中文字符
	hasChinese := false
	for _, r := range title {
		if unicode.Is(unicode.Han, r) {
			hasChinese = true
			break
		}
	}

	var result string
	if hasChinese {
		// 中文转为拼音
		args := pinyin.NewArgs()
		args.Style = pinyin.NORMAL
		py := pinyin.LazyPinyin(title, args)
		result = strings.Join(py, "")
	} else {
		result = title
	}

	// 统一处理：转小写、替换非字母数字为空格、去重连字符
	result = strings.ToLower(result)
	result = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(result, "-")
	result = strings.Trim(result, "-")
	result = regexp.MustCompile(`-+-`).ReplaceAllString(result, "-")

	// 截断到 255 字符
	if len(result) > 255 {
		result = result[:255]
	}

	return result
}
