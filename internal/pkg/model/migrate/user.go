// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package migrate

import (
	"gorm.io/gorm"
	"gotribe-admin/internal/pkg/model"
)

// 自动迁移表结构
func userMigrate(db *gorm.DB) {
	db.AutoMigrate(
		&model.User{},
	)
	// 添加字段
	//db.Migrator().AddColumn(&model.User{}, "Test")
	// 删除字段
	//db.Migrator().DropColumn(&model.User{}, "Test")
	// 添加索引
	//db.Migrator().CreateIndex(&model.User{}, "idx_name")
}
