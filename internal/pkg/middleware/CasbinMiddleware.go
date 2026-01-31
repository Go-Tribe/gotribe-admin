// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package middleware

import (
	"gotribe-admin/config"
	"gotribe-admin/internal/app/repository"
	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/pkg/api/response"

	"github.com/gin-gonic/gin"

	"strings"
	"sync"
)

// 使用读写锁提升并发性能，权限检查是读操作，可以并发执行
var checkLock sync.RWMutex

// Casbin中间件, 基于RBAC的权限访问控制模型
func CasbinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ur := repository.NewAdminRepository()
		admin, err := ur.GetCurrentAdmin(c)
		if err != nil {
			response.Unauthorized(c, "用户未登录")
			c.Abort()
			return
		}
		if admin.Status != 1 {
			response.UserDisabled(c, "当前用户已被禁用")
			c.Abort()
			return
		}

		// 获得用户的全部角色
		roles := admin.Roles
		// 检查是否为超级管理员（拥有排序为1的角色）
		isSuperAdmin := false
		var subs []string
		for _, role := range roles {
			if role.Status == 1 { // 角色状态正常
				subs = append(subs, role.Keyword)
				// 超级管理员判断：排序为1
				if role.Sort == 1 {
					isSuperAdmin = true
				}
			}
		}

		// 超级管理员跳过权限检查
		if isSuperAdmin {
			c.Next()
			return
		}

		// 获得请求路径URL
		obj := strings.TrimPrefix(c.FullPath(), "/"+config.Conf.System.UrlPathPrefix)
		// 获取请求方式
		act := c.Request.Method

		isPass := check(subs, obj, act)
		if !isPass {
			response.PermissionDenied(c, "权限不足")
			c.Abort()
			return
		}

		c.Next()
	}
}

func check(subs []string, obj string, act string) bool {
	// 使用读锁允许并发权限检查，提升性能
	checkLock.RLock()
	defer checkLock.RUnlock()

	// 遍历用户的所有角色，只要有一个角色有权限就通过
	for _, sub := range subs {
		if pass, _ := common.CasbinEnforcer.Enforce(sub, obj, act); pass {
			return true
		}
	}
	return false
}
