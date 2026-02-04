// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

type ThirdPartyAccounts struct {
	Model
	UserID   string `gorm:"type:varchar(10);Index;comment:用户ID" json:"userID"`
	Platform string `gorm:"type:varchar(50);not null;comment:平台" json:"platform"`
	BindFlag uint   `gorm:"type:smallint;default:1;comment:是否绑定,2绑定" json:"bindFlag"`
	OpenID   string `gorm:"type:varchar(255);uniqueIndex;not null;comment:openID" json:"openID"`
}
