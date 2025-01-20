// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package vo

// 创建专栏结构体
type CreateColumnRequest struct {
	Title       string `form:"title" json:"title" validate:"required,min=2,max=20"`
	Description string `form:"description" json:"description" validate:"required,min=2,max=300"`
	Info        string `form:"info" json:"info"`
	Icon        string `form:"icon" json:"icon" validate:"required,min=2,max=300"`
	ProjectID   string `form:"projectID" json:"projectID" validate:"required,min=2,max=10"`
}

type UpdateColumnRequest struct {
	Title       string `form:"title" json:"title" validate:"required,min=2,max=20"`
	Description string `form:"description" json:"description" validate:"required,min=2,max=300"`
	Icon        string `form:"icon" json:"icon" validate:"required,min=2,max=300"`
	Info        string `form:"info" json:"info"`
}

// 获取专栏列表结构体
type ColumnListRequest struct {
	ColumnID  string `form:"columnID" json:"columnID"`
	ProjectID string `form:"projectID" json:"projectID"`
	Title     string `form:"title" json:"title"`
	PageNum   uint   `json:"pageNum" form:"pageNum"`
	PageSize  uint   `json:"pageSize" form:"pageSize"`
}

// 批量删除专栏结构体
type DeleteColumnsRequest struct {
	ColumnIds string `json:"columnIds" form:"columnIds"`
}
