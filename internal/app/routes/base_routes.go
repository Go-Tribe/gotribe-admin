// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package routes

import (
	"context"
	"net/http"
	"time"

	"gotribe-admin/internal/app/controller"
	"gotribe-admin/internal/pkg/common"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// 注册基础路由
func InitBaseRoutes(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) gin.IRoutes {
	systemConfigController := controller.NewSystemConfigController()
	router := r.Group("/base")
	{
		// 登录登出刷新token无需鉴权
		router.POST("/login", authMiddleware.LoginHandler)
		router.POST("/logout", authMiddleware.LogoutHandler)
		router.POST("/refreshToken", authMiddleware.RefreshHandler)
		router.GET("/config", systemConfigController.GetSystemConfigInfo)

	}
	return r
}

// HealthResponse 健康检查响应结构
type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp int64             `json:"timestamp"`
	Checks    map[string]string `json:"checks"`
}

// InitHealthRoute 注册健康检查路由（不需要认证）
func InitHealthRoute(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()

		checks := make(map[string]string)
		status := "healthy"
		httpStatus := http.StatusOK

		// 检查数据库连接
		if err := common.CheckHealth(ctx); err != nil {
			checks["database"] = "unhealthy: " + err.Error()
			status = "unhealthy"
			httpStatus = http.StatusServiceUnavailable
		} else {
			checks["database"] = "healthy"
		}

		// 检查 Casbin
		if common.CasbinEnforcer == nil {
			checks["casbin"] = "unhealthy: not initialized"
			status = "unhealthy"
			httpStatus = http.StatusServiceUnavailable
		} else {
			checks["casbin"] = "healthy"
		}

		c.JSON(httpStatus, HealthResponse{
			Status:    status,
			Timestamp: time.Now().Unix(),
			Checks:    checks,
		})
	})
}
