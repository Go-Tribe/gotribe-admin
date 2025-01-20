// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

type Feedback struct {
	Model
	Title     string   `gorm:"type:varchar(255);uniqueIndex;not null;comment:标题" json:"title"`
	Content   string   `gorm:"type:longtext;comment:内容" json:"content"`
	Phone     string   `gorm:"type:varchar(20);comment:电话" json:"phone"`
	UserID    string   `gorm:"type:char(10);Index;comment:用户ID" json:"userID"`
	ProjectID string   `gorm:"type:char(10);Index;comment:项目 ID" json:"projectID"`
	User      *User    `gorm:"-" json:"user"`
	Project   *Project `gorm:"-" json:"project"`
}
