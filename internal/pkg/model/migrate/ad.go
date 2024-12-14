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
func adMigrate(db *gorm.DB) {
	db.AutoMigrate(
		&model.Ad{},
	)
}
