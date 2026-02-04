// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

import (
	"github.com/dengmengmian/ghelper/gid"
	"gorm.io/gorm"
)

type ProductSku struct {
	Model
	SkuID         string `gorm:"type:varchar(10);uniqueIndex;comment:唯一字符ID/分布式ID" json:"skuID"`
	Title         string `gorm:"type:varchar(255);not null;comment:标题" json:"title"`
	ProjectID     string `gorm:"type:varchar(10);Index;comment:项目 ID" json:"projectID"`
	ProductID     string `gorm:"type:varchar(10);Index;comment:产品ID" json:"productID"`
	Image         string `gorm:"type:varchar(255);not null;comment:产品主图" json:"image"`
	Video         string `gorm:"type:varchar(255);not null;comment:产品视频" json:"video"`
	CostPrice     int64  `gorm:"type:bigint;not null;comment:成本价(分)" json:"costPrice"`
	UnitPrice     int64  `gorm:"type:bigint;not null;comment:商品价格(分)" json:"unitPrice"`
	MarketPrice   int64  `gorm:"type:bigint;not null;comment:市场价格(分)" json:"marketPrice"`
	Quantity      uint   `gorm:"type:integer;not null;comment:库存" json:"quantity"`
	UnitPoint     int64  `gorm:"type:bigint;not null;comment:积分数值(分)" json:"unitPoint"`
	EnableDefault uint   `gorm:"type:smallint;not null;default:1;comment:是否启用：1-正常；2-默认" json:"enableDefault"`
}

func (e *ProductSku) BeforeCreate(tx *gorm.DB) error {
	e.SkuID = gid.GenShortID()
	return nil
}
