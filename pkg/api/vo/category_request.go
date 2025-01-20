// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package vo

// 创建接口结构体
type CreateCategoryRequest struct {
	Title       string `json:"title" form:"title" validate:"required,min=1,max=50"`
	Icon        string `json:"icon" form:"icon" `
	Path        string `json:"path" form:"path"`
	Sort        uint   `json:"sort" form:"sort" validate:"gte=1,lte=999"`
	Status      uint   `json:"status" form:"status"`
	Hidden      uint   `json:"hidden" form:"hidden" validate:"oneof=1 2"`
	Description string `json:"description" form:"description"`
	ParentID    *uint  `json:"parentID" form:"parentID"`
}

// 更新接口结构体
type UpdateCategoryRequest struct {
	Title       string `json:"title" form:"title" validate:"required,min=1,max=50"`
	Icon        string `json:"icon" form:"icon" `
	Path        string `json:"path" form:"path"`
	Sort        uint   `json:"sort" form:"sort" validate:"gte=1,lte=999"`
	Status      uint   `json:"status" form:"status"`
	Hidden      uint   `json:"hidden" form:"hidden" validate:"oneof=1 2"`
	ParentID    *uint  `json:"parentID" form:"parentID" validate:"required"`
	Description string `json:"description" form:"description"`
}

// 删除接口结构体
type DeleteCategoryRequest struct {
	CategoryIds string `json:"categoryIds" form:"categoryIds"`
}
