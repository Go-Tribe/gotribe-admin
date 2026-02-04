// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package jobs

import (
	"time"

	"gotribe-admin/config"
)

// DefaultJobsConfig 默认任务配置
func DefaultJobsConfig() *config.JobsConfig {
	return &config.JobsConfig{
		Enabled: false,
		List: map[string]config.JobConfig{
			"sitemap": {
				Name:        "sitemap",
				Description: "生成站点地图",
				Schedule:    "@every 1h",
				Enabled:     false,
				Timeout:     5 * time.Minute,
				RetryCount:  3,
			},
			"example": {
				Name:        "example",
				Description: "示例任务",
				Schedule:    "@every 30s",
				Enabled:     false,
				Timeout:     1 * time.Minute,
				RetryCount:  1,
			},
		},
	}
}

// GetJobConfig 获取任务配置
func GetJobConfig(c *config.JobsConfig, jobName string) (config.JobConfig, bool) {
	if c == nil || c.List == nil {
		return config.JobConfig{}, false
	}
	conf, exists := c.List[jobName]
	return conf, exists
}

// IsJobEnabled 检查任务是否启用
func IsJobEnabled(c *config.JobsConfig, jobName string) bool {
	if c == nil || !c.Enabled {
		return false
	}

	conf, exists := GetJobConfig(c, jobName)
	return exists && conf.Enabled
}
