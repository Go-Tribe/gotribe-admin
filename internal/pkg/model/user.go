// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

import (
	"gotribe-admin/pkg/util"
	"time"

	"github.com/dengmengmian/ghelper/gid"
	"gorm.io/gorm"
)

type User struct {
	Model
	UserID     string     `gorm:"type:char(10);not null;uniqueIndex;comment:字符ID，分布式 ID;" json:"user_id"`
	Username   string     `gorm:"type:varchar(30);not null;uniqueIndex;comment:用户名" json:"username"`
	ProjectID  string     `gorm:"type:char(10);not null;index;comment:项目ID;" json:"project_id"`
	Password   string     `gorm:"type:varchar(255);not null;comment:密码" json:"-"`
	Nickname   string     `gorm:"type:varchar(30);not null;comment:昵称" json:"nickname"`
	Email      string     `gorm:"type:varchar(30);not null;uniqueIndex;comment:邮箱" json:"email"`
	Phone      string     `gorm:"type:varchar(21);not null;uniqueIndex;comment:电话" json:"phone"`
	Sex        string     `gorm:"type:char(1);not null;default:M;comment:M:男 F:女" json:"sex"`
	Point      float64    `gorm:"-" json:"point"`
	Status     uint8      `gorm:"type:smallint;not null;default:1;comment:用户状态，1-正常；2-禁用" json:"status"`
	Birthday   *time.Time `gorm:"type:date;comment:'用户生日，格式为YYYY-MM-DD'" json:"birthday"`
	Background string     `gorm:"type:varchar(255);comment:个人中心背景" json:"background"`
	Ext        string     `gorm:"type:text;comment:扩展字段" json:"ext"`
	AvatarURL  string     `gorm:"type:varchar(255);comment:头像地址" json:"avatar_url"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.UserID = gid.GenShortID()
	// Encrypt the user password.
	u.Password, err = util.PasswordUtil.Encrypt(u.Password)
	if err != nil {
		return err
	}
	return nil
}
