// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package middleware

import (
	"gotribe-admin/config"
	"gotribe-admin/internal/app/repository"
	"gotribe-admin/internal/pkg/model"

	"github.com/gin-gonic/gin"

	"strings"
	"time"
)

// 操作日志channel
var OperationLogChan = make(chan *model.OperationLog, 30)

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
		apiDesc := getApiDescription(path, method)

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
			// 如果channel已满，可以选择记录错误或使用非阻塞方式处理
			go func() {
				OperationLogChan <- log
			}()
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
func getApiDescription(path, method string) string {
	apiRepository := repository.NewApiRepository()
	apiDesc, err := apiRepository.GetApiDescByPath(path, method)
	if err != nil {
		return ""
	}
	return apiDesc
}
