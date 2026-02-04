// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package jobs

import (
	"context"
	"gotribe-admin/config"
	"time"
)

// BaseJob 基础任务实现
type BaseJob struct {
	name        string
	description string
	schedule    string
	enabled     bool
	timeout     time.Duration
	retryCount  int
	executor    func(ctx context.Context) error
}

// NewBaseJob 创建基础任务
func NewBaseJob(jobConfig config.JobConfig, executor func(ctx context.Context) error) *BaseJob {
	return &BaseJob{
		name:        jobConfig.Name,
		description: jobConfig.Description,
		schedule:    jobConfig.Schedule,
		enabled:     jobConfig.Enabled,
		timeout:     jobConfig.Timeout,
		retryCount:  jobConfig.RetryCount,
		executor:    executor,
	}
}

// Name 返回任务名称
func (j *BaseJob) Name() string {
	return j.name
}

// Description 返回任务描述
func (j *BaseJob) Description() string {
	return j.description
}

// Schedule 返回任务调度表达式
func (j *BaseJob) Schedule() string {
	return j.schedule
}

// Execute 执行任务
func (j *BaseJob) Execute(ctx context.Context) error {
	if j.executor == nil {
		return nil
	}
	return j.executor(ctx)
}

// Timeout 返回任务超时时间
func (j *BaseJob) Timeout() time.Duration {
	if j.timeout <= 0 {
		return 5 * time.Minute // 默认5分钟超时
	}
	return j.timeout
}

// RetryCount 返回重试次数
func (j *BaseJob) RetryCount() int {
	if j.retryCount < 0 {
		return 0
	}
	return j.retryCount
}

// IsEnabled 返回任务是否启用
func (j *BaseJob) IsEnabled() bool {
	return j.enabled
}

// SetEnabled 设置任务启用状态
func (j *BaseJob) SetEnabled(enabled bool) {
	j.enabled = enabled
}
