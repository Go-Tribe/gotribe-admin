// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

import (
	"github.com/dengmengmian/ghelper/gid"
	"gorm.io/gorm"
)

type ProductSpec struct {
	Model
	ProductSpecID string             `gorm:"type:varchar(10);uniqueIndex;comment:唯一字符ID/分布式ID" json:"productSpecID"`
	Title         string             `gorm:"type:varchar(255);not null;comment:标题" json:"title"`
	Remark        string             `gorm:"type:varchar(50);not null;comment:备注" json:"remark"`
	Format        uint               `gorm:"type:smallint;not null;default:1;comment:规格类型:1-文字,2-图片" json:"format"`
	Image         string             `gorm:"type:varchar(255);comment:图片" json:"image"`
	Sort          uint               `gorm:"type:smallint;not null;default:1;comment:排序" json:"sort"`
	Items         []*ProductSpecItem `gorm:"-" json:"items"`
}

func (e *ProductSpec) BeforeCreate(tx *gorm.DB) error {
	e.ProductSpecID = gid.GenShortID()
	return nil
}
