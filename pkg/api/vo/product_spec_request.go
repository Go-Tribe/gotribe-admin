// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package vo

// 创建商品规格值结构体
type CreateProductSpecRequest struct {
	Title  string `form:"title" json:"title" validate:"required,min=2,max=20"`
	Sort   uint   `form:"sort" json:"sort" validate:"gte=1,lte=999"`
	Image  string `form:"image" json:"image"`
	Remark string `gorm:"type:varchar(50);not null;comment:备注" json:"remark"`
	Format uint   `gorm:"type:tinyint(4);not null;default:1;comment:规格类型:1-文字,2-图片" json:"format"`
}

// 获取商品规格值列表结构体
type ProductSpecListRequest struct {
	Title    string `form:"title" json:"title"`
	PageNum  uint   `json:"pageNum" form:"pageNum"`
	PageSize uint   `json:"pageSize" form:"pageSize"`
}

// 批量删除商品规格值结构体
type DeleteProductSpecRequest struct {
	ProductSpecIds string `json:"productSpecIds" form:"productSpecIds"`
}
