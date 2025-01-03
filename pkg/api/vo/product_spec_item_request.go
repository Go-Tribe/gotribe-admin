// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package vo

// 创建商品规格值结构体
type CreateProductSpecItemRequest struct {
	Title   string `form:"title" json:"title" validate:"required,min=2,max=20"`
	SpecID  string `form:"specID" json:"specID" validate:"required"`
	Sort    uint   `form:"sort" json:"sort" validate:"gte=1,lte=999"`
	Enabled uint   `form:"enabled" json:"enabled" validate:"oneof=1 2"`
}

// 获取商品规格值列表结构体
type ProductSpecItemListRequest struct {
	SpecID   string `form:"specID" json:"specID"`
	PageNum  uint   `json:"pageNum" form:"pageNum"`
	PageSize uint   `json:"pageSize" form:"pageSize"`
}

// 批量删除商品规格值结构体
type DeleteProductSpecItemsRequest struct {
	ProductSpecItemIds string `json:"productSpecItemIds" form:"productSpecItemIds"`
}
