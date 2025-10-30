// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package jobs

import (
	"fmt"
	"sync"

	"gotribe-admin/internal/pkg/common"
)

// JobRegistry 任务注册表
type JobRegistry struct {
	manager JobManager
	jobs    map[string]Job
	mu      sync.RWMutex
}

// NewJobRegistry 创建任务注册表
func NewJobRegistry(manager JobManager) *JobRegistry {
	return &JobRegistry{
		manager: manager,
		jobs:    make(map[string]Job),
	}
}

// Register 注册任务
func (r *JobRegistry) Register(job Job) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if job == nil {
		return fmt.Errorf("job cannot be nil")
	}

	jobName := job.Name()
	if jobName == "" {
		return fmt.Errorf("job name cannot be empty")
	}

	// 检查任务是否已存在
	if _, exists := r.jobs[jobName]; exists {
		return fmt.Errorf("job %s already registered", jobName)
	}

	// 注册到管理器
	if err := r.manager.Register(job); err != nil {
		return fmt.Errorf("failed to register job %s: %w", jobName, err)
	}

	r.jobs[jobName] = job
	common.Log.Infof("Job %s registered successfully", jobName)
	return nil
}

// GetJob 获取任务
func (r *JobRegistry) GetJob(jobName string) (Job, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	job, exists := r.jobs[jobName]
	return job, exists
}

// ListJobs 列出所有任务
func (r *JobRegistry) ListJobs() []Job {
	r.mu.RLock()
	defer r.mu.RUnlock()

	jobs := make([]Job, 0, len(r.jobs))
	for _, job := range r.jobs {
		jobs = append(jobs, job)
	}
	return jobs
}

// Start 启动所有任务
func (r *JobRegistry) Start() error {
	return r.manager.Start()
}

// Stop 停止所有任务
func (r *JobRegistry) Stop() error {
	return r.manager.Stop()
}

// GetJobStatus 获取任务状态
func (r *JobRegistry) GetJobStatus(jobName string) (*JobResult, error) {
	return r.manager.GetJobStatus(jobName)
}

// GetJobHistory 获取任务历史
func (r *JobRegistry) GetJobHistory(jobName string, limit int) ([]*JobResult, error) {
	return r.manager.GetJobHistory(jobName, limit)
}

// EnableJob 启用任务
func (r *JobRegistry) EnableJob(jobName string) error {
	return r.manager.EnableJob(jobName)
}

// DisableJob 禁用任务
func (r *JobRegistry) DisableJob(jobName string) error {
	return r.manager.DisableJob(jobName)
}

// 全局任务注册表
var (
	globalRegistry *JobRegistry
	registryOnce   sync.Once
)

// GetGlobalRegistry 获取全局任务注册表
func GetGlobalRegistry() *JobRegistry {
	registryOnce.Do(func() {
		manager := NewCronJobManager()
		globalRegistry = NewJobRegistry(manager)
	})
	return globalRegistry
}

// RegisterJob 注册任务到全局注册表
func RegisterJob(job Job) error {
	return GetGlobalRegistry().Register(job)
}

// StartJobs 启动所有任务
func StartJobs() error {
	return GetGlobalRegistry().Start()
}

// StopJobs 停止所有任务
func StopJobs() error {
	return GetGlobalRegistry().Stop()
}
