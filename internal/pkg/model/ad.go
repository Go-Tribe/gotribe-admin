// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

import (
	"github.com/dengmengmian/ghelper/gid"
	"gorm.io/gorm"
)

type Ad struct {
	Model
	AdID        string   `gorm:"type:char(10);uniqueIndex;comment:唯一字符ID/分布式ID" json:"AdID"`
	Title       string   `gorm:"type:varchar(255);not null;comment:标题" json:"title"`
	Description string   `gorm:"not null;size:300;not null;comment:描述" json:"description"`
	URL         string   `gorm:"type:varchar(255);not null;comment:广告链接" json:"url"`
	URLType     uint     `gorm:"type:tinyint;default:1;comment:1.链接，2.文章，3.商品" json:"urlType"`
	Sort        uint     `gorm:"type:tinyint;default:1;comment:排序" json:"sort"`
	Status      uint     `gorm:"type:tinyint;not null;default:1;comment:状态，1-未发布；2-发布" json:"status,omitempty"`
	SceneID     string   `gorm:"type:char(10);Index;comment:场景 ID" json:"sceneID"`
	Ext         string   `gorm:"type:text;comment:扩展字段" json:"ext"`
	Image       string   `gorm:"type:varchar(255);comment:图片地址" json:"image"`
	Scene       *AdScene `gorm:"-" json:"scene"`
}

func (r *Ad) BeforeCreate(tx *gorm.DB) error {
	r.AdID = gid.GenShortID()
	return nil
}
