// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package vo

// 创建推广位结构体
type CreateAdSceneRequest struct {
	ProjectID   string `form:"projectID" json:"projectID" validate:"required"`
	Title       string `form:"title" json:"title" validate:"required,min=2,max=50"`
	Description string `form:"description" json:"description" validate:"min=0,max=150"`
}

// 获取推广位列表结构体
type AdSceneListRequest struct {
	ProjectID string `form:"ProjectID" json:"ProjectID"`
	PageNum   uint   `json:"pageNum" form:"pageNum"`
	PageSize  uint   `json:"pageSize" form:"pageSize"`
}

// 更新推广位内容
type UpdateAdSceneRequest struct {
	Title       string `form:"title" json:"title" validate:"required,min=2,max=50"`
	Description string `form:"description" json:"description" validate:"required,min=2,max=150"`
}

// 批量删除项目结构体
type DeleteAdScenesRequest struct {
	AdSceneIds string `json:"adScenesIds" form:"adScenesIds"`
}
