// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package vo

// 创建标签结构体
type CreateTagRequest struct {
	Title       string `form:"title" json:"title" validate:"required,min=2,max=20"`
	Slug        string `form:"slug" json:"slug" validate:"required,min=2,max=30"`
	Description string `form:"description" json:"description"`
	Color       string `form:"color" json:"color"`
	Sort        uint   `form:"sort" json:"sort"`
	Status      uint8  `form:"status" json:"status"`
}

// 获取标签列表结构体
type TagListRequest struct {
	ID        uint   `form:"id" json:"id"`
	Title     string `form:"title" json:"title"`
	PageNum   uint   `json:"pageNum" form:"pageNum"`
	PageSize  uint   `json:"pageSize" form:"pageSize"`
	SortBy    string `json:"sortBy" form:"sortBy"`
	SortOrder string `json:"sortOrder" form:"sortOrder"`
}

// 批量删除标签结构体
type DeleteTagsRequest struct {
	Ids []uint `json:"ids" form:"ids"`
}
