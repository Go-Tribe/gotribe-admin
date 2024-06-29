// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package vo

// 创建配置结构体
type CreateConfigRequest struct {
	ProjectID   string `form:"projectID" json:"projectID" validate:"required"`
	Alias       string `form:"alias" json:"alias" validate:"required,min=2,max=20"`
	Type        uint   `form:"type" json:"type" validate:"required"`
	Title       string `form:"title" json:"title" validate:"required,min=2,max=20"`
	Description string `form:"description" json:"description" validate:"required,min=2,max=150"`
	Info        string `form:"info" json:"info" validate:"required,min=2,max=3000"`
	MDContent   string `form:"mdContent" json:"mdContent"`
}

// 获取配置列表结构体
type ConfigListRequest struct {
	ConfigID  string `form:"configID" json:"configID"`
	ProjectID string `form:"ProjectID" json:"ProjectID"`
	Title     string `form:"title" json:"title"`
	Type      uint   `form:"type" json:"type"`
	PageNum   uint   `json:"pageNum" form:"pageNum"`
	PageSize  uint   `json:"pageSize" form:"pageSize"`
}

// 更新配置内容
type UpdateConfigRequest struct {
	ProjectID   string `form:"projectID" json:"projectID" validate:"required"`
	Title       string `form:"title" json:"title" validate:"required,min=2,max=20"`
	Description string `form:"description" json:"description" validate:"required,min=2,max=150"`
	MDContent   string `form:"mdContent" json:"mdContent" `
	Info        string `form:"info" json:"info" validate:"required,min=2,max=3000"`
}

// 批量删除项目结构体
type DeleteConfigsRequest struct {
	ConfigIds string `json:"configIds" form:"configIds"`
}
