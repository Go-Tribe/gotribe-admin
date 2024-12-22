// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

import (
	"github.com/dengmengmian/ghelper/gid"
	"gorm.io/gorm"
)

type ProductSpecItem struct {
	Model
	ProductSpecItemID string `gorm:"type:char(10);uniqueIndex;comment:唯一字符ID/分布式ID" json:"productSpecItemID"`
	SpecID            string `gorm:"type:char(10);comment:唯一字符ID/分布式ID" json:"specID"`
	Title             string `gorm:"type:varchar(255);not null;comment:标题" json:"title"`
	Sort              uint   `gorm:"type:tinyint(4);not null;default:1;comment:排序" json:"sort"`
	Enabled           uint   `form:"enabled" json:"enabled" validate:"oneof=1 2"`
}

func (e *ProductSpecItem) BeforeCreateItem(tx *gorm.DB) error {
	e.ProductSpecItemID = gid.GenShortID()
	return nil
}
