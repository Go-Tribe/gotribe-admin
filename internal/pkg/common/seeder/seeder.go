// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package seeder

import (
	"errors"
	"log"

	"gorm.io/gorm"
)

// Seeder 数据种子接口
type Seeder interface {
	Run(db *gorm.DB) error
	Name() string
}

// BaseSeeder 基础种子实现
type BaseSeeder struct {
	name string
}

// NewBaseSeeder 创建基础种子
func NewBaseSeeder(name string) *BaseSeeder {
	return &BaseSeeder{name: name}
}

// Name 返回种子名称
func (s *BaseSeeder) Name() string {
	return s.name
}

// SeedRegistry 种子注册表
type SeedRegistry struct {
	seeders []Seeder
}

// NewSeedRegistry 创建种子注册表
func NewSeedRegistry() *SeedRegistry {
	return &SeedRegistry{
		seeders: make([]Seeder, 0),
	}
}

// Register 注册种子
func (r *SeedRegistry) Register(seeder Seeder) {
	r.seeders = append(r.seeders, seeder)
}

// RunAll 运行所有种子
func (r *SeedRegistry) RunAll(db *gorm.DB) error {
	for _, seeder := range r.seeders {
		log.Printf("开始执行种子: %s", seeder.Name())
		if err := seeder.Run(db); err != nil {
			log.Printf("种子 %s 执行失败: %v", seeder.Name(), err)
			return err
		}
		log.Printf("种子 %s 执行完成", seeder.Name())
	}
	return nil
}

// 全局种子注册表
var GlobalRegistry = NewSeedRegistry()

// RegisterSeeder 注册种子的便捷函数
func RegisterSeeder(seeder Seeder) {
	GlobalRegistry.Register(seeder)
}

// RunSeeders 运行所有注册的种子
func RunSeeders(db *gorm.DB) error {
	return GlobalRegistry.RunAll(db)
}

// 通用创建函数
func createIfNotExists(db *gorm.DB, model interface{}, id uint) error {
	err := db.First(model, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return db.Create(model).Error
	}
	return err
}

// 批量创建函数
func createBatchIfNotExists(db *gorm.DB, models interface{}, getID func(interface{}, int) uint) error {
	// 这里可以根据具体需求实现批量检查逻辑
	return db.Create(models).Error
}
