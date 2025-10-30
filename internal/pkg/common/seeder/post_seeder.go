// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package seeder

import (
	"gotribe-admin/internal/pkg/model"

	"gorm.io/gorm"
)

// PostSeeder 文章种子
type PostSeeder struct {
	*BaseSeeder
}

// NewPostSeeder 创建文章种子
func NewPostSeeder() *PostSeeder {
	return &PostSeeder{
		BaseSeeder: NewBaseSeeder("post"),
	}
}

// Run 执行文章数据种子
func (s *PostSeeder) Run(db *gorm.DB) error {
	posts := []*model.Post{
		{
			Model:       model.Model{ID: 1},
			PostID:      "243x9",
			Title:       "欢迎使用GoTribe",
			Description: "这是一篇示例文章",
			Content:     "# 这是一篇示例文章",
			Icon:        "https://cdn.dengmengmian.com/20240528/1716909013037462047.jpg",
			HtmlContent: "<h1>这是一篇示例文章</h1>",
			UserID:      "245eko",
			CategoryID:  "24ejga",
			Author:      "GoTribe",
			ProjectID:   "245eko",
		},
	}

	for _, post := range posts {
		if err := createIfNotExists(db, post, post.ID); err != nil {
			return err
		}
	}

	return nil
}
