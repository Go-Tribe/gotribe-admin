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

func InitCategoryRoutes(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) gin.IRoutes {
	categoryController := controller.NewCategoryController()
	router := r.Group("/category")
	// 开启jwt认证中间件
	router.Use(authMiddleware.MiddlewareFunc())
	// 开启casbin鉴权中间件
	router.Use(middleware.CasbinMiddleware())
	{
		router.GET("/tree", categoryController.GetCategoryTree)
		router.GET("", categoryController.GetCategorys)
		router.POST("", categoryController.CreateCategory)
		router.PATCH("/:categoryID", categoryController.UpdateCategoryByID)
		router.DELETE("", categoryController.BatchDeleteCategoryByIds)
		router.GET("/:categoryID", categoryController.GetCategoryInfo)
	}

	return r
}
