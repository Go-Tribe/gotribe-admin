// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

import (
	"github.com/dengmengmian/ghelper/gid"
	"gorm.io/gorm"
)

type Category struct {
	Model
	CategoryID  string      `gorm:"type:char(10);uniqueIndex;comment:唯一字符ID/分布式ID" json:"categoryID"`
	ParentID    *uint       `gorm:"default:0;comment:父菜单编号(编号为0时表示根菜单)" json:"parentID"`
	Sort        uint        `gorm:"default:1;comment:排序" json:"sort"`
	Icon        string      `gorm:"type:varchar(255);comment:图标" json:"icon"`
	Title       string      `gorm:"type:varchar(255);not null;comment:'标题'" json:"title"`
	Path        string      `gorm:"type:varchar(100);comment:url" json:"path"`
	Hidden      uint        `gorm:"type:smallint;default:1;comment:1显示，2隐藏" json:"hidden"`
	Description string      `gorm:"type:varchar(300);not null;comment:描述" json:"description"`
	Ext         string      `gorm:"type:text;comment:扩展字段" json:"ext"`
	Status      uint        `gorm:"type:smallint;not null;default:1;comment:状态，1-正常；2-禁用" json:"status,omitempty"`
	Children    []*Category `gorm:"-" json:"children"`
}

func (c *Category) BeforeCreate(tx *gorm.DB) error {
	c.CategoryID = gid.GenShortID()

	return nil
}
