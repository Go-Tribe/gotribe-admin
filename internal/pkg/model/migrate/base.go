// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package migrate

import (
	"fmt"

	"gorm.io/gorm"
)

// BaseMigrator 基础迁移器
type BaseMigrator struct {
	name   string
	models []interface{}
}

// NewBaseMigrator 创建基础迁移器
func NewBaseMigrator(name string, models ...interface{}) *BaseMigrator {
	return &BaseMigrator{
		name:   name,
		models: models,
	}
}

// Migrate 执行迁移
func (m *BaseMigrator) Migrate(db *gorm.DB) error {
	if m == nil {
		return fmt.Errorf("migrator is nil")
	}
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}
	if m.models == nil {
		return fmt.Errorf("no models to migrate")
	}

	for _, model := range m.models {
		if model == nil {
			continue
		}
		result := db.AutoMigrate(model)
		if result != nil && result.Error != nil {
			return fmt.Errorf("failed to migrate model: %w", result.Error)
		}
	}
	return nil
}

// Name 返回迁移器名称
func (m *BaseMigrator) Name() string {
	return m.name
}
