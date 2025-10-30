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

// 注册标签管理路由
func InitTagRoutes(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) gin.IRoutes {
	tagController := controller.NewTagController()
	router := r.Group("/tag")
	// 开启jwt认证中间件
	router.Use(authMiddleware.MiddlewareFunc())
	// 开启casbin鉴权中间件
	router.Use(middleware.CasbinMiddleware())
	{
		router.GET("/:tagID", tagController.GetTagInfo)
		router.GET("", tagController.GetTags)
		router.POST("", tagController.CreateTag)
		router.PATCH("/:tagID", tagController.UpdateTagByID)
		router.DELETE("", tagController.BatchDeleteTagByIds)
	}
	return r
}
