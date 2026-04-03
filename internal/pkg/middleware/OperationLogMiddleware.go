// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package middleware

import (
	"gotribe-admin/config"
	"gotribe-admin/internal/app/repository"
	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/internal/pkg/model"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"

	"strings"
	"time"
)

// 操作日志channel
var OperationLogChan = make(chan *model.OperationLog, 256)
var apiDescCache = cache.New(10*time.Minute, 20*time.Minute)

// 定义静态资源路径前缀
var skipPaths = []string{
	"/static/",
	"/assets/",
	"/images/",
	"/favicon.ico",
	"/swagger/",
}

func OperationLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取实际请求路径
		requestPath := c.Request.URL.Path

		// 如果请求的是 Swagger 相关路径，直接跳过
		if strings.HasPrefix(requestPath, "/swagger/") {
			c.Next()
			return
		}

		// 获取访问路径
		path := strings.TrimPrefix(c.FullPath(), "/"+config.Conf.System.UrlPathPrefix)

		// 如果是空路径或静态资源，直接返回
		if shouldSkipLog(path) {
			c.Next()
			return
		}

		startTime := time.Now()
		c.Next()
		timeCost := time.Since(startTime).Milliseconds()

		username := getUsername(c)
		method := c.Request.Method

		// 获取接口描述
		apiDesc := getApiDescription(path, method, c)

		log := &model.OperationLog{
			Username:  username,
			Ip:        c.ClientIP(),
			Method:    method,
			Path:      path,
			Desc:      apiDesc,
			Status:    c.Writer.Status(),
			StartTime: startTime,
			TimeCost:  timeCost,
		}

		// 异步写入日志
		select {
		case OperationLogChan <- log:
		default:
			// 日志写入不应阻塞主请求链路，队列满时丢弃并记录告警。
			if requestPath != "/health" {
				repositoryLogDropWarn(path, method)
			}
		}
	}
}

// 判断是否需要跳过日志记录
func shouldSkipLog(path string) bool {
	if path == "" {
		return true
	}

	for _, prefix := range skipPaths {
		if strings.HasPrefix(path, prefix) || prefix == path {
			return true
		}
	}
	return false
}

// 获取用户名
func getUsername(c *gin.Context) string {
	ctxUser, exists := c.Get("user")
	if !exists {
		return "未登录"
	}

	user, ok := ctxUser.(model.Admin)
	if !ok {
		return "未登录"
	}

	return user.Username
}

// 获取API描述
func getApiDescription(path, method string, c *gin.Context) string {
	cacheKey := method + ":" + path
	if cachedDesc, ok := apiDescCache.Get(cacheKey); ok {
		return cachedDesc.(string)
	}

	apiRepository := repository.NewApiRepository()
	apiDesc, err := apiRepository.GetApiDescByPath(c.Request.Context(), path, method)
	if err != nil {
		apiDescCache.Set(cacheKey, "", cache.DefaultExpiration)
		return ""
	}
	apiDescCache.Set(cacheKey, apiDesc, cache.DefaultExpiration)
	return apiDesc
}

func repositoryLogDropWarn(path, method string) {
	// 使用固定文案降低告警噪音，避免在高压场景下再次放大日志量。
	if path == "" {
		common.Log.Warn("operation log queue is full, dropping request log")
		return
	}
	common.Log.Warnf("operation log queue is full, dropping request log for %s %s", method, path)
}
