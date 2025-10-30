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

// 注册用户管理路由
func InitUserRoutes(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) gin.IRoutes {
	userController := controller.NewUserController()
	router := r.Group("/user")
	// 开启jwt认证中间件
	router.Use(authMiddleware.MiddlewareFunc())
	// 开启casbin鉴权中间件
	router.Use(middleware.CasbinMiddleware())
	{
		router.GET("/:userID", userController.GetUserInfo)
		router.GET("", userController.GetUsers)
		router.POST("", userController.CreateUser)
		router.PATCH("/:userID", userController.UpdateUserByID)
		router.DELETE("", userController.BatchDeleteUserByIds)
		router.GET("/search", userController.SearchUserByUsername)
	}
	return r
}
