// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package vo

// 创建系统配置结构体
type CreateSystemConfigRequest struct {
	Title   string `form:"title" json:"title" validate:"required,min=2,max=20"`
	Content string `form:"content" json:"content"`
	Logo    string `form:"logo" json:"logo" validate:"required"`
	Icon    string `form:"icon" json:"icon"`
	Footer  string `form:"footer" json:"footer"`
}
