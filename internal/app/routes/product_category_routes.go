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

func InitProductCategoryRoutes(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) gin.IRoutes {
	productCategoryController := controller.NewProductCategoryController()
	router := r.Group("/product/category")
	// 开启jwt认证中间件
	router.Use(authMiddleware.MiddlewareFunc())
	// 开启casbin鉴权中间件
	router.Use(middleware.CasbinMiddleware())
	{
		router.GET("/tree", productCategoryController.GetProductCategoryTree)
		router.GET("", productCategoryController.GetProductCategorys)
		router.POST("", productCategoryController.CreateProductCategory)
		router.PATCH("/:productCategoryID", productCategoryController.UpdateProductCategoryByID)
		router.DELETE("", productCategoryController.BatchDeleteProductCategoryByIds)
		router.GET("/:productCategoryID", productCategoryController.GetProductCategoryInfo)
	}

	return r
}
