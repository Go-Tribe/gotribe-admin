// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package jobs

import (
	"gotribe-admin/internal/pkg/common"
)

// InitJobs 初始化所有任务
func InitJobs() error {
	common.Log.Info("Initializing jobs...")

	// 获取默认配置
	config := DefaultJobsConfig()

	// 注册站点地图任务
	if config.IsJobEnabled("sitemap") {
		sitemapConfig, _ := config.GetJobConfig("sitemap")
		sitemapJob := NewSitemapJob(sitemapConfig)
		if err := RegisterJob(sitemapJob); err != nil {
			common.Log.Errorf("Failed to register sitemap job: %v", err)
			return err
		}
	}

	// 注册示例任务
	if config.IsJobEnabled("example") {
		exampleConfig, _ := config.GetJobConfig("example")
		exampleJob := NewExampleJob(exampleConfig)
		if err := RegisterJob(exampleJob); err != nil {
			common.Log.Errorf("Failed to register example job: %v", err)
			return err
		}
	}

	common.Log.Info("Jobs initialized successfully")
	return nil
}

// StartAllJobs 启动所有任务
func StartAllJobs() error {
	common.Log.Info("Starting all jobs...")

	if err := StartJobs(); err != nil {
		common.Log.Errorf("Failed to start jobs: %v", err)
		return err
	}

	common.Log.Info("All jobs started successfully")
	return nil
}

// StopAllJobs 停止所有任务
func StopAllJobs() error {
	common.Log.Info("Stopping all jobs...")

	if err := StopJobs(); err != nil {
		common.Log.Errorf("Failed to stop jobs: %v", err)
		return err
	}

	common.Log.Info("All jobs stopped successfully")
	return nil
}
