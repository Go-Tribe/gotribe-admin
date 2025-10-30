// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package migrate

import (
	"fmt"

	"gotribe-admin/internal/pkg/model"

	"gorm.io/gorm"
)

// 自动迁移表结构
func DBAutoMigrate(db *gorm.DB) {
	// 创建新的迁移注册表
	registry := NewMigrationRegistry()

	// 注册所有迁移
	registerAllMigrationsToRegistry(registry)

	// 运行所有迁移
	if err := registry.RunAll(db); err != nil {
		panic(fmt.Sprintf("database migration failed: %v", err))
	}
}

// registerAllMigrationsToRegistry 向指定注册表注册所有迁移
func registerAllMigrationsToRegistry(registry *MigrationRegistry) {
	// 用户相关表
	registry.Register(NewBaseMigrator("user", &model.User{}))
	registry.Register(NewBaseMigrator("admin", &model.Admin{}))
	registry.Register(NewBaseMigrator("userEvent", &model.UserEvent{}))

	// 权限相关表
	registry.Register(NewBaseMigrator("role", &model.Role{}))
	registry.Register(NewBaseMigrator("menu", &model.Menu{}))
	registry.Register(NewBaseMigrator("api", &model.Api{}))

	// 内容相关表
	registry.Register(NewBaseMigrator("post", &model.Post{}))
	registry.Register(NewBaseMigrator("tag", &model.Tag{}))
	registry.Register(NewBaseMigrator("category", &model.Category{}))
	registry.Register(NewBaseMigrator("column", &model.Column{}))
	registry.Register(NewBaseMigrator("comment", &model.Comment{}))

	// 项目相关表
	registry.Register(NewBaseMigrator("project", &model.Project{}))
	registry.Register(NewBaseMigrator("config", &model.Config{}))

	// 资源相关表
	registry.Register(NewBaseMigrator("resource", &model.Resource{}))

	// 广告相关表
	registry.Register(NewBaseMigrator("adScene", &model.AdScene{}))
	registry.Register(NewBaseMigrator("ad", &model.Ad{}))

	// 积分相关表
	registry.Register(NewBaseMigrator("pointLog", &model.PointLog{}))
	registry.Register(NewBaseMigrator("pointDeduction", &model.PointDeduction{}))
	registry.Register(NewBaseMigrator("pointAvailable", &model.PointAvailable{}))

	// 商品相关表
	registry.Register(NewBaseMigrator("productCategory", &model.ProductCategory{}))
	registry.Register(NewBaseMigrator("productType", &model.ProductType{}))
	registry.Register(NewBaseMigrator("productSpec", &model.ProductSpec{}))
	registry.Register(NewBaseMigrator("productSpecItem", &model.ProductSpecItem{}))
	registry.Register(NewBaseMigrator("product", &model.Product{}))
	registry.Register(NewBaseMigrator("productSku", &model.ProductSku{}))

	// 订单相关表
	registry.Register(NewBaseMigrator("order", &model.Order{}))
	registry.Register(NewBaseMigrator("orderLog", &model.OrderLog{}))

	// 系统相关表
	registry.Register(NewBaseMigrator("systemConfig", &model.SystemConfig{}))
	registry.Register(NewBaseMigrator("operationLog", &model.OperationLog{}))
	registry.Register(NewBaseMigrator("thirdPartyAccounts", &model.ThirdPartyAccounts{}))
	registry.Register(NewBaseMigrator("feedback", &model.Feedback{}))
}

// registerAllMigrations 注册所有迁移（保持向后兼容）
func registerAllMigrations() {
	// 用户相关表
	RegisterMigration(NewBaseMigrator("user", &model.User{}))
	RegisterMigration(NewBaseMigrator("admin", &model.Admin{}))
	RegisterMigration(NewBaseMigrator("userEvent", &model.UserEvent{}))

	// 权限相关表
	RegisterMigration(NewBaseMigrator("role", &model.Role{}))
	RegisterMigration(NewBaseMigrator("menu", &model.Menu{}))
	RegisterMigration(NewBaseMigrator("api", &model.Api{}))

	// 内容相关表
	RegisterMigration(NewBaseMigrator("post", &model.Post{}))
	RegisterMigration(NewBaseMigrator("tag", &model.Tag{}))
	RegisterMigration(NewBaseMigrator("category", &model.Category{}))
	RegisterMigration(NewBaseMigrator("column", &model.Column{}))
	RegisterMigration(NewBaseMigrator("comment", &model.Comment{}))

	// 项目相关表
	RegisterMigration(NewBaseMigrator("project", &model.Project{}))
	RegisterMigration(NewBaseMigrator("config", &model.Config{}))

	// 资源相关表
	RegisterMigration(NewBaseMigrator("resource", &model.Resource{}))

	// 广告相关表
	RegisterMigration(NewBaseMigrator("adScene", &model.AdScene{}))
	RegisterMigration(NewBaseMigrator("ad", &model.Ad{}))

	// 积分相关表
	RegisterMigration(NewBaseMigrator("pointLog", &model.PointLog{}))
	RegisterMigration(NewBaseMigrator("pointDeduction", &model.PointDeduction{}))
	RegisterMigration(NewBaseMigrator("pointAvailable", &model.PointAvailable{}))

	// 商品相关表
	RegisterMigration(NewBaseMigrator("productCategory", &model.ProductCategory{}))
	RegisterMigration(NewBaseMigrator("productType", &model.ProductType{}))
	RegisterMigration(NewBaseMigrator("productSpec", &model.ProductSpec{}))
	RegisterMigration(NewBaseMigrator("productSpecItem", &model.ProductSpecItem{}))
	RegisterMigration(NewBaseMigrator("product", &model.Product{}))
	RegisterMigration(NewBaseMigrator("productSku", &model.ProductSku{}))

	// 订单相关表
	RegisterMigration(NewBaseMigrator("order", &model.Order{}))
	RegisterMigration(NewBaseMigrator("orderLog", &model.OrderLog{}))

	// 系统相关表
	RegisterMigration(NewBaseMigrator("systemConfig", &model.SystemConfig{}))
	RegisterMigration(NewBaseMigrator("operationLog", &model.OperationLog{}))
	RegisterMigration(NewBaseMigrator("thirdPartyAccounts", &model.ThirdPartyAccounts{}))
	RegisterMigration(NewBaseMigrator("feedback", &model.Feedback{}))
}
