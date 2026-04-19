// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

import (
	"time"
)

type User struct {
	Model
	Username   string     `gorm:"type:varchar(30);not null;uniqueIndex:idx_user_project_username,priority:2;comment:用户名" json:"username"`
	ProjectId  uint       `gorm:"not null;index;uniqueIndex:idx_user_project_username,priority:1;uniqueIndex:idx_user_project_email,priority:1;uniqueIndex:idx_user_project_phone,priority:1;comment:项目ID" json:"projectId"`
	Password   string     `gorm:"type:varchar(255);not null;comment:密码" json:"-"`
	Nickname   string     `gorm:"type:varchar(30);not null;comment:昵称" json:"nickname"`
	Email      *string    `gorm:"type:varchar(254);default:null;uniqueIndex:idx_user_project_email,priority:2;comment:邮箱" json:"email,omitempty"`
	Phone      *string    `gorm:"type:varchar(32);default:null;uniqueIndex:idx_user_project_phone,priority:2;comment:电话" json:"phone,omitempty"`
	Sex        string     `gorm:"type:char(1);not null;default:M;comment:M:男 F:女" json:"sex,omitempty"`
	Point      float64    `gorm:"-" json:"point"`
	Status     uint8      `gorm:"type:smallint;not null;default:1;comment:用户状态，1-正常；2-禁用" json:"status"`
	Birthday   *time.Time `gorm:"type:date;comment:用户生日，格式为YYYY-MM-DD" json:"birthday,omitempty"`
	Background string     `gorm:"type:varchar(255);comment:个人中心背景" json:"background,omitempty"`
	Ext        string     `gorm:"type:text;comment:扩展字段" json:"ext,omitempty"`
	AvatarURL  string     `gorm:"type:varchar(255);comment:头像地址" json:"avatar_url,omitempty"`
}
