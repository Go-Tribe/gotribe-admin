// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package jobs

import (
	"context"
	"fmt"
	"sync"
	"time"

	"gotribe-admin/internal/pkg/common"

	"github.com/robfig/cron/v3"
)

// CronJobManager 基于cron的任务管理器
type CronJobManager struct {
	cron       *cron.Cron
	jobs       map[string]Job
	jobResults map[string][]*JobResult
	mu         sync.RWMutex
	entryIDs   map[string]cron.EntryID
}

// NewCronJobManager 创建新的任务管理器
func NewCronJobManager() *CronJobManager {
	secondParser := cron.NewParser(
		cron.SecondOptional | cron.Minute | cron.Hour |
			cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
	)

	return &CronJobManager{
		cron:       cron.New(cron.WithParser(secondParser), cron.WithChain()),
		jobs:       make(map[string]Job),
		jobResults: make(map[string][]*JobResult),
		entryIDs:   make(map[string]cron.EntryID),
	}
}

// Register 注册任务
func (m *CronJobManager) Register(job Job) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if job == nil {
		return fmt.Errorf("job cannot be nil")
	}

	jobName := job.Name()
	if jobName == "" {
		return fmt.Errorf("job name cannot be empty")
	}

	// 如果任务已存在，先移除
	if _, exists := m.jobs[jobName]; exists {
		if entryID, ok := m.entryIDs[jobName]; ok {
			m.cron.Remove(entryID)
			delete(m.entryIDs, jobName)
		}
	}

	m.jobs[jobName] = job

	// 如果任务启用，添加到cron调度器
	if job.IsEnabled() {
		entryID, err := m.cron.AddFunc(job.Schedule(), func() {
			m.executeJob(jobName)
		})
		if err != nil {
			return fmt.Errorf("failed to add job %s to cron: %w", jobName, err)
		}
		m.entryIDs[jobName] = entryID
	}

	common.Log.Infof("Job %s registered successfully", jobName)
	return nil
}

// Start 启动任务调度器
func (m *CronJobManager) Start() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.cron == nil {
		return fmt.Errorf("cron scheduler is not initialized")
	}

	m.cron.Start()
	common.Log.Info("Job scheduler started")
	return nil
}

// Stop 停止任务调度器
func (m *CronJobManager) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.cron == nil {
		return nil
	}

	ctx := m.cron.Stop()
	<-ctx.Done()
	common.Log.Info("Job scheduler stopped")
	return nil
}

// GetJobStatus 获取任务状态
func (m *CronJobManager) GetJobStatus(jobName string) (*JobResult, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	results, exists := m.jobResults[jobName]
	if !exists || len(results) == 0 {
		return nil, fmt.Errorf("no execution history found for job %s", jobName)
	}

	// 返回最新的执行结果
	return results[len(results)-1], nil
}

// GetJobHistory 获取任务执行历史
func (m *CronJobManager) GetJobHistory(jobName string, limit int) ([]*JobResult, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	results, exists := m.jobResults[jobName]
	if !exists {
		return nil, fmt.Errorf("no execution history found for job %s", jobName)
	}

	if limit <= 0 || limit > len(results) {
		limit = len(results)
	}

	// 返回最新的limit个结果
	start := len(results) - limit
	if start < 0 {
		start = 0
	}

	return results[start:], nil
}

// EnableJob 启用任务
func (m *CronJobManager) EnableJob(jobName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	job, exists := m.jobs[jobName]
	if !exists {
		return fmt.Errorf("job %s not found", jobName)
	}

	// 如果任务已经在运行，先移除
	if entryID, ok := m.entryIDs[jobName]; ok {
		m.cron.Remove(entryID)
		delete(m.entryIDs, jobName)
	}

	// 启用任务
	if baseJob, ok := job.(*BaseJob); ok {
		baseJob.SetEnabled(true)
	}

	// 添加到cron调度器
	entryID, err := m.cron.AddFunc(job.Schedule(), func() {
		m.executeJob(jobName)
	})
	if err != nil {
		return fmt.Errorf("failed to enable job %s: %w", jobName, err)
	}

	m.entryIDs[jobName] = entryID
	common.Log.Infof("Job %s enabled", jobName)
	return nil
}

// DisableJob 禁用任务
func (m *CronJobManager) DisableJob(jobName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	job, exists := m.jobs[jobName]
	if !exists {
		return fmt.Errorf("job %s not found", jobName)
	}

	// 从cron调度器中移除
	if entryID, ok := m.entryIDs[jobName]; ok {
		m.cron.Remove(entryID)
		delete(m.entryIDs, jobName)
	}

	// 禁用任务
	if baseJob, ok := job.(*BaseJob); ok {
		baseJob.SetEnabled(false)
	}

	common.Log.Infof("Job %s disabled", jobName)
	return nil
}

// executeJob 执行任务
func (m *CronJobManager) executeJob(jobName string) {
	m.mu.RLock()
	job, exists := m.jobs[jobName]
	m.mu.RUnlock()

	if !exists {
		common.Log.Errorf("Job %s not found", jobName)
		return
	}

	// 创建执行结果
	result := &JobResult{
		JobID:      fmt.Sprintf("%s_%d", jobName, time.Now().UnixNano()),
		Status:     StatusRunning,
		StartTime:  time.Now(),
		RetryCount: 0,
	}

	// 记录开始执行
	m.recordJobResult(jobName, result)

	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), job.Timeout())
	defer cancel()

	// 执行任务
	err := m.executeWithRetry(ctx, job, result)

	// 记录执行结果
	endTime := time.Now()
	result.EndTime = &endTime
	result.Duration = endTime.Sub(result.StartTime)

	if err != nil {
		result.Status = StatusFailed
		result.Error = err
		result.Message = err.Error()
		common.Log.Errorf("Job %s failed: %v", jobName, err)
	} else {
		result.Status = StatusCompleted
		result.Message = "Job completed successfully"
		common.Log.Infof("Job %s completed in %v", jobName, result.Duration)
	}

	m.recordJobResult(jobName, result)
}

// executeWithRetry 带重试的任务执行
func (m *CronJobManager) executeWithRetry(ctx context.Context, job Job, result *JobResult) error {
	var lastErr error

	for i := 0; i <= job.RetryCount(); i++ {
		result.RetryCount = i

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err := job.Execute(ctx); err != nil {
			lastErr = err
			if i < job.RetryCount() {
				common.Log.Warnf("Job %s failed (attempt %d/%d), retrying: %v",
					job.Name(), i+1, job.RetryCount()+1, err)
				time.Sleep(time.Duration(i+1) * time.Second) // 指数退避
				continue
			}
		} else {
			return nil
		}
	}

	return lastErr
}

// recordJobResult 记录任务执行结果
func (m *CronJobManager) recordJobResult(jobName string, result *JobResult) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.jobResults[jobName] == nil {
		m.jobResults[jobName] = make([]*JobResult, 0)
	}

	// 保持最近100条记录
	results := m.jobResults[jobName]
	if len(results) >= 100 {
		results = results[1:]
	}

	results = append(results, result)
	m.jobResults[jobName] = results
}
