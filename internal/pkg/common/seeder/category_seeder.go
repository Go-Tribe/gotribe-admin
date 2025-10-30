// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package seeder

import (
	"gotribe-admin/internal/pkg/model"

	"gorm.io/gorm"
)

// CategorySeeder 分类种子
type CategorySeeder struct {
	*BaseSeeder
}

// NewCategorySeeder 创建分类种子
func NewCategorySeeder() *CategorySeeder {
	return &CategorySeeder{
		BaseSeeder: NewBaseSeeder("category"),
	}
}

// Run 执行分类数据种子
func (s *CategorySeeder) Run(db *gorm.DB) error {
	categories := []*model.Category{
		{
			Model:       model.Model{ID: 1},
			Title:       "默认分类",
			Description: "默认分类",
			CategoryID:  "24ejga",
			Status:      1,
		},
	}

	for _, category := range categories {
		if err := createIfNotExists(db, category, category.ID); err != nil {
			return err
		}
	}

	return nil
}
