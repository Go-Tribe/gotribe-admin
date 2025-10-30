// Copyright 2023 Innkeeper gotribe <info@gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package common

import (
	"gotribe-admin/config"
	"gotribe-admin/internal/pkg/common/seeder"
)

// 初始化数据库数据
func InitData() {
	// 是否初始化数据
	if !config.Conf.System.InitData {
		return
	}

	// 注册所有种子
	registerAllSeeders()

	// 运行所有种子
	if err := seeder.RunSeeders(DB); err != nil {
		Log.Errorf("数据库种子执行失败: %v", err)
	}
}

// registerAllSeeders 注册所有种子
func registerAllSeeders() {
	// 基础数据种子
	seeder.RegisterSeeder(seeder.NewRoleSeeder())
	seeder.RegisterSeeder(seeder.NewAdminSeeder())
	seeder.RegisterSeeder(seeder.NewMenuSeeder())
	seeder.RegisterSeeder(seeder.NewApiSeeder())
	seeder.RegisterSeeder(seeder.NewSystemConfigSeeder())

	// 内容管理种子
	seeder.RegisterSeeder(seeder.NewCategorySeeder())
	seeder.RegisterSeeder(seeder.NewTagSeeder())
	seeder.RegisterSeeder(seeder.NewPostSeeder())

	// 项目管理种子
	seeder.RegisterSeeder(seeder.NewProjectSeeder())
	seeder.RegisterSeeder(seeder.NewUserSeeder())

	// 可以继续添加其他种子...
}
