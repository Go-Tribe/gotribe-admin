// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package vo

// 创建接口结构体
type CreateProductCategoryRequest struct {
	Title       string `json:"title" form:"title" validate:"required,min=1,max=50"`
	Icon        string `json:"icon" form:"icon" `
	Path        string `json:"path" form:"path"`
	Sort        uint   `json:"sort" form:"sort" validate:"gte=1,lte=999"`
	Status      uint   `json:"status" form:"status"`
	Hidden      uint   `json:"hidden" form:"hidden" validate:"oneof=1 2"`
	ProjectID   string `json:"projectID" form:"projectID"`
	Description string `json:"description" form:"description"`
	ParentID    *uint  `json:"parentID" form:"parentID"`
}

// 更新接口结构体
type UpdateProductCategoryRequest struct {
	Title       string `json:"title" form:"title" validate:"required,min=1,max=50"`
	Icon        string `json:"icon" form:"icon" `
	Path        string `json:"path" form:"path"`
	Sort        uint   `json:"sort" form:"sort" validate:"gte=1,lte=999"`
	Status      uint   `json:"status" form:"status"`
	Hidden      uint   `json:"hidden" form:"hidden" validate:"oneof=1 2"`
	ParentID    *uint  `json:"parentID" form:"parentID" validate:"required"`
	ProjectID   string `json:"projectID" form:"projectID" validate:"required"`
	Description string `json:"description" form:"description"`
}

// 删除接口结构体
type DeleteProductCategoryRequest struct {
	ProductCategoryIds string `json:"productCategoryIds" form:"productCategoryIds"`
}
