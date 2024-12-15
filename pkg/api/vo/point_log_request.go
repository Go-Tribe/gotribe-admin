// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package vo

// 创建积分结构体
type CreatePointLogRequest struct {
	ProjectID string  `form:"projectID" json:"projectID" validate:"required"`
	UserID    string  `form:"userID" json:"userID" validate:"required"`
	Point     float32 `form:"point" json:"point" validate:"required"`
}

// 获取积分列表结构体
type PointLogListRequest struct {
	UserID    string `form:"userID" json:"userID"`
	Nickname  string `form:"nickname" json:"nickname"`
	ProjectID string `form:"projectID" json:"projectID"`
	PageNum   uint   `json:"pageNum" form:"pageNum"`
	PageSize  uint   `json:"pageSize" form:"pageSize"`
}
