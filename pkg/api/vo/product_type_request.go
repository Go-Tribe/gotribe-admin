// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package vo

// 创建商品类型结构体
type CreateProductTypeRequest struct {
	Title      string `form:"title" json:"title" validate:"required,min=2,max=20"`
	CategoryID string `form:"categoryID" json:"categoryID" validate:"required"`
	SpecIds    string `form:"specIDs" json:"specIDs" validate:"required"`
	Remark     string `form:"remark" json:"remark"`
}

// 获取商品类型列表结构体
type ProductTypeListRequest struct {
	ProductTypeID string `form:"productTypeID" json:"productTypeID"`
	Title         string `form:"title" json:"title"`
	PageNum       uint   `json:"pageNum" form:"pageNum"`
	PageSize      uint   `json:"pageSize" form:"pageSize"`
}

// 批量删除商品类型结构体
type DeleteProductTypesRequest struct {
	ProductTypeIds string `json:"productTypeIds" form:"productTypeIds"`
}
