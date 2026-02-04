// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

import (
	"github.com/dengmengmian/ghelper/gid"
	"gorm.io/gorm"
)

type Tag struct {
	Model
	TagID       string `gorm:"type:varchar(10);uniqueIndex;comment:唯一字符ID/分布式ID" json:"tagID"`
	Title       string `gorm:"type:varchar(255);uniqueIndex;not null;comment:标题" json:"title"`
	Description string `gorm:"not null;size:300;comment:描述" json:"description"`
	Color       string `gorm:"type:varchar(20);comment:颜色" json:"color"`
}

func (t *Tag) BeforeCreate(tx *gorm.DB) error {
	t.TagID = gid.GenShortID()
	return nil
}
