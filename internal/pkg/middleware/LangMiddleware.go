package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// LangMiddleware parses request language and stores it into context
// Priority: query param `lang` > header `Accept-Language`
// Supported: zh, en; default: zh
func LangMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := c.Query("lang")
		if lang == "" {
			al := c.GetHeader("Accept-Language")
			if al != "" {
				// basic parse: take leading tag like "en" or "zh"
				al = strings.ToLower(al)
				// split by comma and semicolon, take first token
				first := al
				if idx := strings.IndexAny(al, ",; "); idx >= 0 {
					first = al[:idx]
				}
				if strings.HasPrefix(first, "en") {
					lang = "en"
				} else if strings.HasPrefix(first, "zh") {
					lang = "zh"
				}
			}
		}
		if lang != "en" { // default fallback to Chinese
			lang = "zh"
		}
		c.Set("lang", lang)
		c.Next()
	}
}
