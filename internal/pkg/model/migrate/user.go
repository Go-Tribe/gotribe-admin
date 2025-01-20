// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package migrate

import (
	"fmt"
	"gorm.io/gorm"
	"gotribe-admin/internal/pkg/model"
	"log"
)

// 自动迁移表结构
func userMigrate(db *gorm.DB) error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 自动迁移表结构
	if err := tx.AutoMigrate(&model.User{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("auto migrate failed: %w", err)
	}

	// 添加生日字段
	if err := tx.Migrator().AddColumn(&model.User{}, "birthday").Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("add column birthday failed: %w", err)
	}
	// 增加用户头像
	if err := tx.Migrator().AddColumn(&model.User{}, "avatar_url").Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("add column avatar_url failed: %w", err)
	}

	// 删除字段
	// if err := tx.Migrator().DropColumn(&model.User{}, "Test").Error; err != nil {
	//     tx.Rollback()
	//     return fmt.Errorf("drop column Test failed: %w", err)
	// }

	//手机号&邮箱唯一索引
	if err := tx.Migrator().CreateIndex(&model.User{}, "idx_phone").Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("create index idx_phone failed: %w", err)
	}
	if err := tx.Migrator().CreateIndex(&model.User{}, "idx_email").Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("create index idx_phone failed: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("commit transaction failed: %w", err)
	}

	log.Println("User migration completed successfully")
	return nil
}
