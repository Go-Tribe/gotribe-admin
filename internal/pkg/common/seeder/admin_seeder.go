// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package seeder

import (
	"gotribe-admin/internal/pkg/model"
	"gotribe-admin/pkg/util"

	"gorm.io/gorm"
)

// AdminSeeder 管理员种子
type AdminSeeder struct {
	*BaseSeeder
}

// NewAdminSeeder 创建管理员种子
func NewAdminSeeder() *AdminSeeder {
	return &AdminSeeder{
		BaseSeeder: NewBaseSeeder("admin"),
	}
}

// Run 执行管理员数据种子
func (s *AdminSeeder) Run(db *gorm.DB) error {
	// 获取角色
	var adminRole model.Role
	if err := db.First(&adminRole, 1).Error; err != nil {
		return err
	}

	password, err := util.PasswordUtil.GenPasswd("123456")
	if err != nil {
		return err
	}

	admin := &model.Admin{
		Model:        model.Model{ID: 1},
		Username:     "admin",
		Password:     password,
		Mobile:       "18888888888",
		Avatar:       "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
		Nickname:     new(string),
		Introduction: new(string),
		Status:       1,
		Creator:      "系统",
		Roles:        []*model.Role{&adminRole},
	}

	return createIfNotExists(db, admin, admin.ID)
}
