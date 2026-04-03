// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package vo

// 创建内容结构体
type CreatePostRequest struct {
	Title       string   `form:"title" json:"title" validate:"required,min=2,max=60"`
	Description string   `form:"description" json:"description" validate:"required,min=2,max=300"`
	CategoryID  uint     `form:"categoryID" json:"categoryID" validate:"required"`
	ProjectID   string   `form:"projectID" json:"projectID" validate:"required"`
	UserID      uint     `form:"userID" json:"userID" validate:"required"`
	Author      string   `form:"author" json:"author" validate:"required"`
	Content     string   `form:"content" json:"content" validate:"required"`
	HtmlContent string   `form:"htmlContent" json:"htmlContent" validate:"required"`
	ColumnID    uint     `form:"columnID" json:"columnID"`
	Tag         string   `form:"tag" json:"tag"`
	Ext         string   `form:"ext" json:"ext"`
	Icon        string   `form:"icon" json:"icon"`
	Type        uint     `form:"type" json:"type" validate:"required"`
	IsTop       uint     `form:"isTop" json:"isTop"`
	IsPasswd    uint     `form:"isPasswd" json:"isPasswd"`
	Password    string   `form:"password" json:"password"`
	Status      uint     `form:"status" json:"status"`
	Location    string   `form:"location" json:"location"`
	People      string   `form:"people" json:"people"`
	Time        string   `form:"time" json:"time"`
	Images      []string `form:"images" json:"images"`
	UnitPrice   float64  `form:"unitPrice" json:"unitPrice"`
	Video       string   `form:"video" json:"video"`
	ShowTime    string   `form:"showTime" json:"showTime"`
}

// 更新内容结构体
type UpdatePostRequest struct {
	Title       string   `form:"title" json:"title" validate:"required,min=2,max=60"`
	Description string   `form:"description" json:"description" validate:"required,min=2,max=300"`
	CategoryID  uint     `form:"categoryID" json:"categoryID" validate:"required"`
	ProjectID   string   `form:"projectID" json:"projectID" validate:"required"`
	UserID      uint     `form:"userID" json:"userID" validate:"required"`
	Author      string   `form:"author" json:"author" validate:"required"`
	Content     string   `form:"content" json:"content" validate:"required"`
	HtmlContent string   `form:"htmlContent" json:"htmlContent" validate:"required"`
	ColumnID    uint     `form:"columnID" json:"columnID"`
	Tag         string   `form:"tag" json:"tag"`
	Ext         string   `form:"ext" json:"ext"`
	Icon        string   `form:"icon" json:"icon"`
	Type        uint     `form:"type" json:"type" validate:"required"`
	IsTop       uint     `form:"isTop" json:"isTop"`
	IsPasswd    uint     `form:"isPasswd" json:"isPasswd"`
	Password    string   `form:"password" json:"password"`
	Status      uint     `form:"status" json:"status"`
	Location    string   `form:"location" json:"location"`
	People      string   `form:"people" json:"people"`
	Time        string   `form:"time" json:"time"`
	Images      []string `form:"images" json:"images"`
	UnitPrice   float64  `form:"unitPrice" json:"unitPrice"`
	ShowTime    string   `form:"showTime" json:"showTime"`
	Video       string   `form:"video" json:"video"`
}

// 获取内容列表结构体
type PostListRequest struct {
	PostID    string `form:"postID" json:"postID"`
	Title     string `form:"title" json:"title"`
	Status    uint   `form:"status" json:"status"`
	ProjectID string `form:"projectID" json:"projectID"`
	PageNum   uint   `json:"pageNum" form:"pageNum"`
	PageSize  uint   `json:"pageSize" form:"pageSize"`
	SortBy    string `json:"sortBy" form:"sortBy"`
	SortOrder string `json:"sortOrder" form:"sortOrder"`
}

// 批量删除内容结构体
type DeletePostsRequest struct {
	PostIds string `json:"postIds" form:"postIds"`
}
