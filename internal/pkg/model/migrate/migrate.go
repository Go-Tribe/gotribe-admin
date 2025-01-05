// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package migrate

import (
	"gorm.io/gorm"
)

// 自动迁移表结构
func DBAutoMigrate(db *gorm.DB) {
	// user表
	userMigrate(db)
	// admin表
	adminMigrate(db)
	// role表
	roleMigrate(db)
	// menu表
	menuMigrate(db)
	// api表
	apiMigrate(db)
	// operationLog表
	operationLogMigrate(db)
	// post
	postMigrate(db)
	// example 示例表
	exampleMigrate(db)
	// config表,用于自定义配置
	configMigrate(db)
	// tag表
	tagMigrate(db)
	// resource 资源表
	resourceMigrate(db)
	// 分类表
	categoryMigrate(db)
	// 项目表
	projectMigrate(db)
	// 专栏表
	columnMigrate(db)
	// 推广场景
	adSceneMigrate(db)
	// 广告内容
	adMigrate(db)
	// 评论表
	commentMigrate(db)
	// 积分相关表
	pointMigrate(db)
	// 商品分类
	productCategoryMigrate(db)
	// 商品类型（规格组）
	productTypeMigrate(db)
	// 商品规格
	productSpecMigrate(db)
	// 商品规格属性
	productSpecItemMigrate(db)
	// SPU表
	productMigrate(db)
	// SKU表
	productSKUMigrate(db)
	// 订单表
	orderMigrate(db)
	// 订单日志表
	orderLogMigrate(db)
	// 系统设置
	systemConfigMigrate(db)
	// 第三方帐号表
	thirdPartyAccountsMigrate(db)
}
