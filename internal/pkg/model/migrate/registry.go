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

// Migrator 迁移接口
type Migrator interface {
	Migrate(db *gorm.DB) error
	Name() string
}

// MigrationRegistry 迁移注册表
type MigrationRegistry struct {
	migrators []Migrator
}

// NewMigrationRegistry 创建新的迁移注册表
func NewMigrationRegistry() *MigrationRegistry {
	return &MigrationRegistry{
		migrators: make([]Migrator, 0),
	}
}

// Register 注册迁移器
func (r *MigrationRegistry) Register(migrator Migrator) {
	if r == nil {
		return
	}
	if migrator == nil {
		log.Printf("Warning: attempting to register nil migrator")
		return
	}
	r.migrators = append(r.migrators, migrator)
}

// RunAll 运行所有迁移
func (r *MigrationRegistry) RunAll(db *gorm.DB) error {
	if r == nil {
		return fmt.Errorf("migration registry is nil")
	}
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}
	for _, migrator := range r.migrators {
		if migrator == nil {
			log.Printf("Skipping nil migrator")
			continue
		}
		log.Printf("Starting migration: %s", migrator.Name())
		if err := migrator.Migrate(db); err != nil {
			return fmt.Errorf("migration %s failed: %w", migrator.Name(), err)
		}
		log.Printf("Migration %s completed successfully", migrator.Name())
	}
	return nil
}

// 全局迁移注册表
var GlobalRegistry = NewMigrationRegistry()

// RegisterMigration 注册迁移的便捷函数
func RegisterMigration(migrator Migrator) {
	if GlobalRegistry == nil {
		GlobalRegistry = NewMigrationRegistry()
	}
	GlobalRegistry.Register(migrator)
}

// RunMigrations 运行所有注册的迁移
func RunMigrations(db *gorm.DB) error {
	if GlobalRegistry == nil {
		GlobalRegistry = NewMigrationRegistry()
	}
	return GlobalRegistry.RunAll(db)
}
