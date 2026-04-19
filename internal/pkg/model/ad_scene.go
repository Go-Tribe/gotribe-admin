// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

type AdScene struct {
	Model
	Title       string   `gorm:"type:varchar(255);not null;comment:标题" json:"title"`
	Description string   `gorm:"not null;size:300;not null;comment:描述" json:"description"`
	ProjectId   uint     `gorm:"not null;index;comment:项目ID;" json:"projectId"`
	Project     *Project `gorm:"-" json:"project"`
}
