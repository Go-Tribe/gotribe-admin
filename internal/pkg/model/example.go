// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

import (
	"github.com/dengmengmian/ghelper/gid"
	"gorm.io/gorm"
)

type Example struct {
	Model
	ExampleID   string `gorm:"type:varchar(10);uniqueIndex;comment:唯一字符ID/分布式ID" json:"exampleID"`
	ProjectID   string `gorm:"type:varchar(10);not null;index;comment:项目ID;" json:"projectID"`
	Username    string `gorm:"type:varchar(30);not null;index:idx_username;comment:用户名" json:"username"`
	Title       string `gorm:"type:varchar(255);not null;comment:标题" json:"title"`
	Content     string `gorm:"not null;type:text;not null;comment:内容" json:"content"`
	Description string `gorm:"not null;size:300;not null;comment:描述" json:"description"`
	Status      uint8  `gorm:"type:smallint;not null;default:1;comment:状态，1-正常；2-禁用" json:"status,omitempty"`
}

func (e *Example) BeforeCreate(tx *gorm.DB) error {
	e.ExampleID = gid.GenShortID()

	return nil
}
