// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package seeder

import (
	"gotribe-admin/internal/pkg/model"

	"gorm.io/gorm"
)

// TagSeeder 标签种子
type TagSeeder struct {
	*BaseSeeder
}

// NewTagSeeder 创建标签种子
func NewTagSeeder() *TagSeeder {
	return &TagSeeder{
		BaseSeeder: NewBaseSeeder("tag"),
	}
}

// Run 执行标签数据种子
func (s *TagSeeder) Run(db *gorm.DB) error {
	tags := []*model.Tag{
		{
			Model:       model.Model{ID: 1},
			TagID:       "default",
			Title:       "默认标签",
			Description: "默认标签",
		},
	}

	for _, tag := range tags {
		if err := createIfNotExists(db, tag, tag.ID); err != nil {
			return err
		}
	}

	return nil
}
