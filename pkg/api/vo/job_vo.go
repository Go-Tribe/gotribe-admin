// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package vo

// JobVO 任务视图对象
type JobVO struct {
	Name        string `json:"name"`        // 任务名称
	Description string `json:"description"` // 任务描述
	Schedule    string `json:"schedule"`    // 调度表达式
	Enabled     bool   `json:"enabled"`     // 是否启用
	Timeout     string `json:"timeout"`     // 超时时间
	RetryCount  int    `json:"retry_count"` // 重试次数
}
