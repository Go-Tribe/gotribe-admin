// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package seeder

import (
	"gotribe-admin/internal/pkg/model"

	"gorm.io/gorm"
)

// RoleSeeder 角色种子
type RoleSeeder struct {
	*BaseSeeder
}

// NewRoleSeeder 创建角色种子
func NewRoleSeeder() *RoleSeeder {
	return &RoleSeeder{
		BaseSeeder: NewBaseSeeder("role"),
	}
}

// Run 执行角色数据种子
func (s *RoleSeeder) Run(db *gorm.DB) error {
	roles := []*model.Role{
		{
			Model:   model.Model{ID: 1},
			Name:    "管理员",
			Keyword: "admin",
			Desc:    new(string),
			Sort:    1,
			Status:  1,
			Creator: "系统",
		},
		{
			Model:   model.Model{ID: 2},
			Name:    "普通管理员",
			Keyword: "user",
			Desc:    new(string),
			Sort:    3,
			Status:  1,
			Creator: "系统",
		},
		{
			Model:   model.Model{ID: 3},
			Name:    "访客",
			Keyword: "guest",
			Desc:    new(string),
			Sort:    5,
			Status:  1,
			Creator: "系统",
		},
	}

	for _, role := range roles {
		if err := createIfNotExists(db, role, role.ID); err != nil {
			return err
		}
	}

	return nil
}
