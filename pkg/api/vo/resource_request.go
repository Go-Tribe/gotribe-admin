// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package vo

// 创建资源结构体
type CreateResourceRequest struct {
	Title       string `form:"title" json:"title" validate:"required,min=2,max=20"`
	Description string `form:"description" json:"description" validate:"required,min=2,max=150"`
}

// 获取资源列表结构体
type ResourceListRequest struct {
	ResourceID string `form:"resourceID" json:"resourceID"`
	Type       uint   `form:"type" json:"type"`
	PageNum    uint   `json:"pageNum" form:"pageNum"`
	PageSize   uint   `json:"pageSize" form:"pageSize"`
}

// 批量删除资源结构体
type DeleteResourcesRequest struct {
	ResourceID string `json:"resourceID" form:"resourceID"`
}
