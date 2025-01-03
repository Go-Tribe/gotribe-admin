// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

import (
	"github.com/dengmengmian/ghelper/gid"
	"gorm.io/gorm"
)

type ProductType struct {
	Model
	ProductTypeID     string `gorm:"type:char(10);uniqueIndex;comment:唯一字符ID/分布式ID" json:"productTypeID"`
	Title             string `gorm:"type:varchar(255);not null;comment:标题" json:"title"`
	Remark            string `gorm:"type:varchar(50);not null;comment:备注" json:"remark"`
	ProductCategoryID string `gorm:"type:char(10);not null;index:idx_product_category_id;comment:分类ID;" json:"productCategoryID"`
	SpecIds           string `gorm:"type:varchar(255);not null;index:idx_spec_id;comment:规格编号;" json:"specIds"`
}

func (e *ProductType) BeforeCreate(tx *gorm.DB) error {
	e.ProductTypeID = gid.GenShortID()
	return nil
}
