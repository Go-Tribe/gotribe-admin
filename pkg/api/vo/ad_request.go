// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package vo

// 创建推广位结构体
type CreateAdRequest struct {
	Title       string `form:"title" json:"title" validate:"required,min=2,max=50"`
	Description string `form:"description" json:"description" validate:"min=0,max=150"`
	URL         string `form:"url" json:"url" validate:"required,min=2,max=255"`
	URLType     uint   `form:"urlType" json:"urlType" validate:"required"`
	Image       string `form:"image" json:"image" validate:"required,min=2,max=255"`
	Sort        uint   `form:"sort" json:"sort" validate:"required"`
	Status      uint   `form:"status" json:"status" validate:"oneof=1 2"`
	SceneID     string `form:"sceneID" json:"sceneID" validate:"required"`
	Ext         string `form:"ext" json:"ext"`
}

// 获取推广位列表结构体
type AdListRequest struct {
	SceneID  string `form:"sceneID" json:"sceneID"`
	Title    string `form:"title" json:"title"`
	Status   uint   `form:"status" json:"status"`
	PageNum  uint   `json:"pageNum" form:"pageNum"`
	PageSize uint   `json:"pageSize" form:"pageSize"`
}

// 更新推广位内容
type UpdateAdRequest struct {
	Title       string `form:"title" json:"title" validate:"required,min=2,max=50"`
	Description string `form:"description" json:"description" validate:"min=0,max=150"`
	URL         string `form:"url" json:"url" validate:"required,min=2,max=255"`
	URLType     uint   `form:"urlType" json:"urlType" validate:"required"`
	Image       string `form:"image" json:"image" validate:"required,min=2,max=255"`
	Sort        uint   `form:"sort" json:"sort" validate:"required"`
	Status      uint   `form:"status" json:"status" validate:"oneof=1 2"`
	SceneID     string `form:"sceneID" json:"sceneID" validate:"required"`
	Ext         string `form:"ext" json:"ext"`
}

// 批量删除项目结构体
type DeleteAdsRequest struct {
	AdIds string `json:"adsIds" form:"adsIds"`
}
