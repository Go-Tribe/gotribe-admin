// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package model

type Role struct {
	Model
	Name    string   `gorm:"type:varchar(20);not null;unique" json:"name"`
	Keyword string   `gorm:"type:varchar(20);not null;unique" json:"keyword"`
	Desc    *string  `gorm:"type:varchar(100);" json:"desc"`
	Status  uint     `gorm:"type:smallint;default:1;comment:1正常, 2禁用" json:"status"`
	Sort    uint     `gorm:"type:integer;default:999;comment:角色排序(排序越大权限越低, 不能查看比自己序号小的角色, 不能编辑同序号用户权限, 排序为1表示超级管理员)" json:"sort"`
	Creator string   `gorm:"type:varchar(20);" json:"creator"`
	Admin   []*Admin `gorm:"many2many:admin_roles" json:"admins"`
	Menus   []*Menu  `gorm:"many2many:role_menus;" json:"menus"` // 角色菜单多对多关系
}
