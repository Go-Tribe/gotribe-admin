// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package seeder

import (
	"gotribe-admin/internal/pkg/model"

	"gorm.io/gorm"
)

// ProjectSeeder 项目种子
type ProjectSeeder struct {
	*BaseSeeder
}

// NewProjectSeeder 创建项目种子
func NewProjectSeeder() *ProjectSeeder {
	return &ProjectSeeder{
		BaseSeeder: NewBaseSeeder("project"),
	}
}

// Run 执行项目数据种子
func (s *ProjectSeeder) Run(db *gorm.DB) error {
	projects := []*model.Project{
		{
			Model:       model.Model{ID: 1},
			Name:        "default",
			Title:       "默认项目",
			ProjectID:   "245eko",
			Description: "默认项目",
		},
	}

	for _, project := range projects {
		if err := createIfNotExists(db, project, project.ID); err != nil {
			return err
		}
	}

	return nil
}
