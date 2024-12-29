// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package routes

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"gotribe-admin/internal/app/controller"
	"gotribe-admin/internal/pkg/middleware"
)

// 注册系统配置管理路由
func InitSystemConfigRoutes(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) gin.IRoutes {
	systemConfigController := controller.NewSystemConfigController()
	router := r.Group("/system")
	// 开启casbin鉴权中间件
	router.Use(middleware.CasbinMiddleware())
	{
		router.PATCH(":systemConfigID", systemConfigController.UpdateSystemConfigByID).Use(authMiddleware.MiddlewareFunc())
	}
	return r
}
