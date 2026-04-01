// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

type Tag struct {
	Model
	Title       string `gorm:"type:varchar(30);not null;uniqueIndex;comment:标题" json:"title"`
	Slug        string `gorm:"type:varchar(30);not null;uniqueIndex:idx_tag_slug;comment:URL别名" json:"slug"`
	Description string `gorm:"type:varchar(300);comment:描述" json:"description,omitempty"`
	Color       string `gorm:"type:varchar(20);comment:展示颜色" json:"color,omitempty"`
	Sort        uint   `gorm:"default:1;comment:排序，越大越靠前" json:"sort"`
	Count       uint   `gorm:"default:0;comment:引用次数" json:"count"`
	Status      uint8  `gorm:"type:smallint;not null;default:1;comment:状态，1-正常；2-禁用" json:"status,omitempty"`
}
