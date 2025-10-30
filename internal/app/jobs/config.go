// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package jobs

import (
	"time"
)

// JobsConfig 任务配置
type JobsConfig struct {
	Enabled bool                 `yaml:"enabled" json:"enabled"`
	Jobs    map[string]JobConfig `yaml:"jobs" json:"jobs"`
}

// DefaultJobsConfig 默认任务配置
func DefaultJobsConfig() *JobsConfig {
	return &JobsConfig{
		Enabled: true,
		Jobs: map[string]JobConfig{
			"sitemap": {
				Name:        "sitemap",
				Description: "生成站点地图",
				Schedule:    "@every 1m", // 改为24小时执行一次
				Enabled:     true,
				Timeout:     5 * time.Minute,
				RetryCount:  3,
			},
			"example": {
				Name:        "example",
				Description: "示例任务",
				Schedule:    "@every 30s",
				Enabled:     false, // 默认禁用
				Timeout:     1 * time.Minute,
				RetryCount:  1,
			},
		},
	}
}

// GetJobConfig 获取任务配置
func (c *JobsConfig) GetJobConfig(jobName string) (JobConfig, bool) {
	config, exists := c.Jobs[jobName]
	return config, exists
}

// IsJobEnabled 检查任务是否启用
func (c *JobsConfig) IsJobEnabled(jobName string) bool {
	if !c.Enabled {
		return false
	}

	config, exists := c.GetJobConfig(jobName)
	return exists && config.Enabled
}
