// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package routes

import (
	"gotribe-admin/internal/app/controller"
	"gotribe-admin/internal/pkg/middleware"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// 注册用户路由
func InitAdminRoutes(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) gin.IRoutes {
	userController := controller.NewAdminController()
	router := r.Group("/admin")
	// 开启jwt认证中间件
	router.Use(authMiddleware.MiddlewareFunc())
	// 开启casbin鉴权中间件
	router.Use(middleware.CasbinMiddleware())
	{
		router.GET("/info", userController.GetAdminInfo)
		router.GET("/list", userController.GetAdmins)
		router.PUT("/changePwd", userController.ChangePwd)
		router.POST("/create", userController.CreateAdmin)
		router.PATCH("/update/:userID", userController.UpdateAdminByID)
		router.DELETE("/delete/batch", userController.BatchDeleteAdminByIds)
	}
	return r
}
