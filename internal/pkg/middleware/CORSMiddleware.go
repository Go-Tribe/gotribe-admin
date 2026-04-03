// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package middleware

import (
	"github.com/gin-gonic/gin"
	"gotribe-admin/config"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// CORS跨域中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			if isAllowedOrigin(origin) {
				c.Header("Vary", "Origin")
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token, Session, Content-Type, Accept-Language")
				c.Header("Access-Control-Expose-Headers", "Content-Length")
				c.Header("Access-Control-Max-Age", formatMaxAge(config.Conf.CORS.MaxAge))
				if config.Conf.CORS.AllowCredentials {
					c.Header("Access-Control-Allow-Credentials", "true")
				}
			} else if method == http.MethodOptions {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
		}

		if method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func isAllowedOrigin(origin string) bool {
	if config.Conf.CORS == nil {
		return isDefaultDevOrigin(origin)
	}

	allowedOrigins := config.Conf.CORS.AllowedOrigins
	for _, allowed := range allowedOrigins {
		if strings.EqualFold(strings.TrimSpace(allowed), origin) {
			return true
		}
	}

	if len(allowedOrigins) == 0 {
		return isDefaultDevOrigin(origin)
	}

	return false
}

func isDefaultDevOrigin(origin string) bool {
	if config.Conf.System != nil && strings.EqualFold(config.Conf.System.Mode, gin.ReleaseMode) {
		return false
	}

	u, err := url.Parse(origin)
	if err != nil {
		return false
	}

	host := u.Hostname()
	return host == "localhost" || host == "127.0.0.1" || host == "::1"
}

func formatMaxAge(maxAge int) string {
	if maxAge <= 0 {
		return "600"
	}
	return strconv.Itoa(maxAge)
}
