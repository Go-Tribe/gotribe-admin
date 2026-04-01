// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

type Category struct {
	Model
	ParentID    uint        `gorm:"default:0;index:idx_category_parent_id;comment:父分类ID，0表示根分类" json:"parentID"`
	Sort        uint        `gorm:"default:1;index:idx_category_sort;comment:排序" json:"sort"`
	Icon        string      `gorm:"type:varchar(255);comment:图标" json:"icon"`
	Title       string      `gorm:"type:varchar(30);not null;comment:标题" json:"title"`
	Slug        string      `gorm:"type:varchar(30);not null;uniqueIndex:idx_category_slug;comment:URL别名" json:"slug"`
	Path        string      `gorm:"type:varchar(255);comment:自定义url路径" json:"path"`
	Hidden      uint8       `gorm:"type:smallint;default:1;comment:1显示，2隐藏" json:"hidden"`
	Description string      `gorm:"type:varchar(300);comment:描述" json:"description,omitempty"`
	Ext         string      `gorm:"type:text;comment:扩展字段" json:"ext"`
	Status      uint8       `gorm:"type:smallint;not null;default:1;comment:状态，1-正常；2-禁用" json:"status,omitempty"`
	Count       uint        `gorm:"default:0;comment:内容数量" json:"count"`
	Children    []*Category `gorm:"-" json:"children"`
}
