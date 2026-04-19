// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

type Column struct {
	Model
	ProjectId   uint   `gorm:"not null;index;comment:项目ID;" json:"projectId"`
	Title       string `gorm:"type:varchar(30);not null;comment:标题" json:"title,omitempty"`
	Description string `gorm:"type:varchar(300);comment:描述" json:"description,omitempty"`
	Icon        string `gorm:"type:varchar(300);comment:图片" json:"icon,omitempty"`
	Info        string `gorm:"type:text;comment:内容" json:"info,omitempty"`
	Ext         string `gorm:"type:text;comment:扩展字段" json:"ext"`
	Status      uint8  `gorm:"type:smallint;not null;default:1;comment:状态，1-正常；2-禁用" json:"status,omitempty"`
}
