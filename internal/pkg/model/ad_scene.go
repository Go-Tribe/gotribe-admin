// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

import (
	"github.com/dengmengmian/ghelper/gid"
	"gorm.io/gorm"
)

type AdScene struct {
	Model
	AdSceneID   string   `gorm:"type:char(10);uniqueIndex;comment:唯一字符ID/分布式ID" json:"AdSceneID"`
	Title       string   `gorm:"type:varchar(255);not null;comment:标题" json:"title"`
	Description string   `gorm:"not null;size:300;not null;comment:描述" json:"description"`
	ProjectID   string   `gorm:"type:char(10);not null;index;comment:项目ID;" json:"project_id"`
	Project     *Project `gorm:"-" json:"project"`
}

func (r *AdScene) BeforeCreate(tx *gorm.DB) error {
	r.AdSceneID = gid.GenShortID()
	return nil
}
