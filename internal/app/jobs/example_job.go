// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package jobs

import (
	"context"
	"time"

	"gotribe-admin/internal/pkg/common"
)

// ExampleJob 示例任务
type ExampleJob struct {
	*BaseJob
}

// NewExampleJob 创建示例任务
func NewExampleJob(config JobConfig) *ExampleJob {
	job := &ExampleJob{}
	job.BaseJob = NewBaseJob(config, job.execute)
	return job
}

// execute 执行示例任务
func (j *ExampleJob) execute(ctx context.Context) error {
	common.Log.Info("Starting example job")

	// 模拟一些工作
	select {
	case <-time.After(2 * time.Second):
		common.Log.Info("Example job completed")
		return nil
	case <-ctx.Done():
		common.Log.Info("Example job cancelled")
		return ctx.Err()
	}
}
