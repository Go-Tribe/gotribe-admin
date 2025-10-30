// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package jobs

import (
	"context"
	"time"
)

// JobStatus 任务状态
type JobStatus string

const (
	StatusPending   JobStatus = "pending"   // 等待执行
	StatusRunning   JobStatus = "running"   // 正在执行
	StatusCompleted JobStatus = "completed" // 执行完成
	StatusFailed    JobStatus = "failed"    // 执行失败
	StatusSkipped   JobStatus = "skipped"   // 跳过执行
)

// JobResult 任务执行结果
type JobResult struct {
	JobID      string        `json:"job_id"`
	Status     JobStatus     `json:"status"`
	StartTime  time.Time     `json:"start_time"`
	EndTime    *time.Time    `json:"end_time,omitempty"`
	Duration   time.Duration `json:"duration"`
	Error      error         `json:"error,omitempty"`
	Message    string        `json:"message,omitempty"`
	RetryCount int           `json:"retry_count"`
}

// Job 任务接口
type Job interface {
	// Name 返回任务名称
	Name() string

	// Description 返回任务描述
	Description() string

	// Schedule 返回任务调度表达式
	Schedule() string

	// Execute 执行任务
	Execute(ctx context.Context) error

	// Timeout 返回任务超时时间
	Timeout() time.Duration

	// RetryCount 返回重试次数
	RetryCount() int

	// IsEnabled 返回任务是否启用
	IsEnabled() bool
}

// JobConfig 任务配置
type JobConfig struct {
	Name        string        `yaml:"name" json:"name"`
	Description string        `yaml:"description" json:"description"`
	Schedule    string        `yaml:"schedule" json:"schedule"`
	Enabled     bool          `yaml:"enabled" json:"enabled"`
	Timeout     time.Duration `yaml:"timeout" json:"timeout"`
	RetryCount  int           `yaml:"retry_count" json:"retry_count"`
}

// JobManager 任务管理器接口
type JobManager interface {
	// Register 注册任务
	Register(job Job) error

	// Start 启动任务调度器
	Start() error

	// Stop 停止任务调度器
	Stop() error

	// GetJobStatus 获取任务状态
	GetJobStatus(jobName string) (*JobResult, error)

	// GetJobHistory 获取任务执行历史
	GetJobHistory(jobName string, limit int) ([]*JobResult, error)

	// EnableJob 启用任务
	EnableJob(jobName string) error

	// DisableJob 禁用任务
	DisableJob(jobName string) error
}
