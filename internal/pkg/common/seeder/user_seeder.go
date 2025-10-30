// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package seeder

import (
	"gotribe-admin/internal/pkg/model"

	"gorm.io/gorm"
)

// UserSeeder 用户种子
type UserSeeder struct {
	*BaseSeeder
}

// NewUserSeeder 创建用户种子
func NewUserSeeder() *UserSeeder {
	return &UserSeeder{
		BaseSeeder: NewBaseSeeder("user"),
	}
}

// Run 执行用户数据种子
func (s *UserSeeder) Run(db *gorm.DB) error {
	users := []*model.User{
		{
			Model:     model.Model{ID: 1},
			UserID:    "gotribe",
			Username:  "gotribe",
			Nickname:  "gotribe",
			ProjectID: "245eko",
		},
	}

	for _, user := range users {
		if err := createIfNotExists(db, user, user.ID); err != nil {
			return err
		}
	}

	return nil
}
