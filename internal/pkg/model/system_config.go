// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

import (
	"github.com/dengmengmian/ghelper/gid"
	"gorm.io/gorm"
)

type SystemConfig struct {
	Model
	SystemConfigID string `gorm:"type:varchar(10);uniqueIndex;comment:唯一字符ID/分布式ID" json:"systemConfigID"`
	Title          string `gorm:"type:varchar(255);uniqueIndex;not null;comment:标题" json:"title"`
	Content        string `gorm:"type:text;not null;comment:内容" json:"content"`
	Logo           string `gorm:"type:varchar(255);comment:logo" json:"logo"`
	Icon           string `gorm:"type:varchar(255);comment:icon" json:"icon"`
	Footer         string `gorm:"type:varchar(255);comment:footer" json:"footer"`
}

func (t *SystemConfig) BeforeCreate(tx *gorm.DB) error {
	t.SystemConfigID = gid.GenShortID()
	return nil
}
