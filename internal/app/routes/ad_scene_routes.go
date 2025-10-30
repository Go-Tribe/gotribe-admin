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

// 注册推广位管理路由
func InitAdSceneRoutes(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) gin.IRoutes {
	adSceneController := controller.NewAdSceneController()
	router := r.Group("/ad/scene")
	// 开启jwt认证中间件
	router.Use(authMiddleware.MiddlewareFunc())
	// 开启casbin鉴权中间件
	router.Use(middleware.CasbinMiddleware())
	{
		router.GET("/:adSceneID", adSceneController.GetAdSceneInfo)
		router.GET("", adSceneController.GetAdScenes)
		router.POST("", adSceneController.CreateAdScene)
		router.PATCH("/:adSceneID", adSceneController.UpdateAdSceneByID)
		router.DELETE("", adSceneController.BatchDeleteAdSceneByIds)
	}
	return r
}
