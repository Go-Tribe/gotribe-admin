// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package migrate

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

// CustomMigrator 自定义迁移器示例
type CustomMigrator struct {
	name string
}

// NewCustomMigrator 创建自定义迁移器
func NewCustomMigrator(name string) *CustomMigrator {
	return &CustomMigrator{name: name}
}

// Migrate 执行自定义迁移逻辑
func (m *CustomMigrator) Migrate(db *gorm.DB) error {
	log.Printf("Executing custom migration: %s", m.name)

	// 示例：创建自定义索引
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_custom_field ON user (created_at)").Error; err != nil {
		return fmt.Errorf("failed to create custom index: %w", err)
	}

	// 示例：添加自定义列
	if err := db.Exec("ALTER TABLE user ADD COLUMN IF NOT EXISTS custom_field VARCHAR(255)").Error; err != nil {
		return fmt.Errorf("failed to add custom column: %w", err)
	}

	log.Printf("Custom migration %s completed successfully", m.name)
	return nil
}

// Name 返回迁移器名称
func (m *CustomMigrator) Name() string {
	return m.name
}

// BatchMigrator 批量迁移器示例
type BatchMigrator struct {
	name   string
	models []interface{}
}

// NewBatchMigrator 创建批量迁移器
func NewBatchMigrator(name string, models ...interface{}) *BatchMigrator {
	return &BatchMigrator{
		name:   name,
		models: models,
	}
}

// Migrate 批量迁移多个模型
func (m *BatchMigrator) Migrate(db *gorm.DB) error {
	log.Printf("Starting batch migration: %s (%d models)", m.name, len(m.models))

	// 使用事务确保批量迁移的原子性
	return db.Transaction(func(tx *gorm.DB) error {
		for i, model := range m.models {
			if err := tx.AutoMigrate(model).Error; err != nil {
				return fmt.Errorf("failed to migrate model %d: %w", i, err)
			}
		}
		return nil
	})
}

// Name 返回迁移器名称
func (m *BatchMigrator) Name() string {
	return m.name
}
