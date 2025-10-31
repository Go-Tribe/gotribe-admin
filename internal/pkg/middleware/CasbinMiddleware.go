// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package middleware

import (
	"gotribe-admin/config"
	"gotribe-admin/internal/app/repository"
	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/pkg/api/known"
	"gotribe-admin/pkg/api/response"

	"github.com/gin-gonic/gin"

	"strings"
	"sync"
)

var checkLock sync.Mutex

// Casbin中间件, 基于RBAC的权限访问控制模型
func CasbinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ur := repository.NewAdminRepository()
		admin, err := ur.GetCurrentAdmin(c)
		if err != nil {
			response.ResponseFunc(c, 401, 401, nil, "用户未登录")
			c.Abort()
			return
		}
		if admin.Status != 1 {
			response.ResponseFunc(c, 401, 401, nil, "当前用户已被禁用")
			c.Abort()
			return
		}
		// 增加超级管理员账号
		if admin.ID == known.DEFAULT_ID {
			return
		}
		// 获得用户的全部角色
		roles := admin.Roles
		// 获得用户全部未被禁用的角色的Keyword
		var subs []string
		for _, role := range roles {
			if role.Status == known.DEFAULT_ID {
				subs = append(subs, role.Keyword)
			}
		}
		// 获得请求路径URL
		//obj := strings.Replace(c.Request.URL.Path, "/"+config.Conf.System.UrlPathPrefix, "", 1)
		obj := strings.TrimPrefix(c.FullPath(), "/"+config.Conf.System.UrlPathPrefix)
		// 获取请求方式
		act := c.Request.Method

		isPass := check(subs, obj, act)
		if !isPass {
			response.ResponseFunc(c, 401, 401, nil, "没有权限")
			c.Abort()
			return
		}

		c.Next()
	}
}

func check(subs []string, obj string, act string) bool {
	// 同一时间只允许一个请求执行校验, 否则可能会校验失败
	checkLock.Lock()
	defer checkLock.Unlock()
	isPass := false
	for _, sub := range subs {
		pass, _ := common.CasbinEnforcer.Enforce(sub, obj, act)
		if pass {
			isPass = true
			break
		}
	}
	return isPass
}
