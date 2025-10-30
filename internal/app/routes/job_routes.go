// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package routes

import (
	"gotribe-admin/internal/app/controller"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// JobRoutes 任务路由
func JobRoutes(r *gin.Engine, authMiddleware *jwt.GinJWTMiddleware) {
	jobController := controller.NewJobController()

	// 任务管理API
	jobGroup := r.Group("/api/jobs")
	jobGroup.Use(authMiddleware.MiddlewareFunc())
	{
		// 获取任务列表
		jobGroup.GET("/", jobController.ListJobs)

		// 获取任务状态
		jobGroup.GET("/:name/status", jobController.GetJobStatus)

		// 获取任务历史
		jobGroup.GET("/:name/history", jobController.GetJobHistory)

		// 启用任务
		jobGroup.POST("/:name/enable", jobController.EnableJob)

		// 禁用任务
		jobGroup.POST("/:name/disable", jobController.DisableJob)
	}
}
