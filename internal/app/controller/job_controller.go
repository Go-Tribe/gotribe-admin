// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package controller

import (
	"strconv"

	"gotribe-admin/internal/app/jobs"
	"gotribe-admin/internal/pkg/common"
	"gotribe-admin/pkg/api/response"
	"gotribe-admin/pkg/api/vo"

	"github.com/gin-gonic/gin"
)

// JobController 任务控制器
type JobController struct{}

// NewJobController 创建任务控制器
func NewJobController() *JobController {
	return &JobController{}
}

// ListJobs 列出所有任务
func (c *JobController) ListJobs(ctx *gin.Context) {
	registry := jobs.GetGlobalRegistry()
	jobList := registry.ListJobs()

	var jobsVO []vo.JobVO
	for _, job := range jobList {
		jobVO := vo.JobVO{
			Name:        job.Name(),
			Description: job.Description(),
			Schedule:    job.Schedule(),
			Enabled:     job.IsEnabled(),
			Timeout:     job.Timeout().String(),
			RetryCount:  job.RetryCount(),
		}
		jobsVO = append(jobsVO, jobVO)
	}

	response.Success(ctx, gin.H{"jobs": jobsVO}, common.Msg(ctx, common.MsgListSuccess))
}

// GetJobStatus 获取任务状态
func (c *JobController) GetJobStatus(ctx *gin.Context) {
	jobName := ctx.Param("name")
	if jobName == "" {
		response.Fail(ctx, gin.H{}, "任务名称不能为空")
		return
	}

	registry := jobs.GetGlobalRegistry()
	status, err := registry.GetJobStatus(jobName)
	if err != nil {
		response.Fail(ctx, gin.H{}, "任务状态获取失败: "+err.Error())
		return
	}

	response.Success(ctx, gin.H{"status": status}, "获取任务状态成功")
}

// GetJobHistory 获取任务执行历史
func (c *JobController) GetJobHistory(ctx *gin.Context) {
	jobName := ctx.Param("name")
	if jobName == "" {
		response.Fail(ctx, gin.H{}, "任务名称不能为空")
		return
	}

	limitStr := ctx.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	registry := jobs.GetGlobalRegistry()
	history, err := registry.GetJobHistory(jobName, limit)
	if err != nil {
		response.Fail(ctx, gin.H{}, "任务历史获取失败: "+err.Error())
		return
	}

	response.Success(ctx, gin.H{"history": history}, "获取任务历史成功")
}

// EnableJob 启用任务
func (c *JobController) EnableJob(ctx *gin.Context) {
	jobName := ctx.Param("name")
	if jobName == "" {
		response.Fail(ctx, gin.H{}, "任务名称不能为空")
		return
	}

	registry := jobs.GetGlobalRegistry()
	if err := registry.EnableJob(jobName); err != nil {
		response.Fail(ctx, gin.H{}, "启用任务失败: "+err.Error())
		return
	}

	response.Success(ctx, gin.H{}, "任务启用成功")
}

// DisableJob 禁用任务
func (c *JobController) DisableJob(ctx *gin.Context) {
	jobName := ctx.Param("name")
	if jobName == "" {
		response.Fail(ctx, gin.H{}, "任务名称不能为空")
		return
	}

	registry := jobs.GetGlobalRegistry()
	if err := registry.DisableJob(jobName); err != nil {
		response.Fail(ctx, gin.H{}, "禁用任务失败: "+err.Error())
		return
	}

	response.Success(ctx, gin.H{}, "任务禁用成功")
}
