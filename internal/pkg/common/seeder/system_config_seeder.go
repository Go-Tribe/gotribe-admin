// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package seeder

import (
	"gotribe-admin/internal/pkg/model"

	"gorm.io/gorm"
)

// SystemConfigSeeder 系统配置种子
type SystemConfigSeeder struct {
	*BaseSeeder
}

// NewSystemConfigSeeder 创建系统配置种子
func NewSystemConfigSeeder() *SystemConfigSeeder {
	return &SystemConfigSeeder{
		BaseSeeder: NewBaseSeeder("system_config"),
	}
}

// Run 执行系统配置数据种子
func (s *SystemConfigSeeder) Run(db *gorm.DB) error {
	configs := []*model.SystemConfig{
		{
			Model:          model.Model{ID: 1},
			SystemConfigID: "245eko",
			Title:          "GoTribe",
			Logo:           "https://raw.gitcode.com/Go-Tribe/gotribe/raw/5ae01df24c556094f74a9b23086f35c3929fe0f3/106083123.png",
			Icon:           "https://raw.gitcode.com/Go-Tribe/gotribe/raw/5ae01df24c556094f74a9b23086f35c3929fe0f3/106083123.png",
		},
	}

	for _, config := range configs {
		if err := createIfNotExists(db, config, config.ID); err != nil {
			return err
		}
	}

	return nil
}
